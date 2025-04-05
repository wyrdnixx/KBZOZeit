package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
)

// handleEchoTest processes echoTest type messages
func handleEchoTest(content interface{}) ([]byte, error) {
	// Convert content to string (it's expected to be a string)
	contentStr, ok := content.(string)
	if !ok {
		return generateResponse("handleEchoTestResponse", true, "Invalid content for echoTest")
	}

	// Respond with the same content
	response := map[string]string{
		"type":    "echoResponse",
		"content": contentStr,
	}
	return json.Marshal(response)
}

// handleTimeBooking processes timebooking type messages
func handleTimeBooking(content interface{}, user User) ([]byte, error) {
	// First, convert content to a map
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return generateResponse("handleTimeBookingResponse", true, "Invalid content format for timebooking")
	}

	// Extract the "from" and "to" fields
	fromStr, okFrom := contentMap["from"].(string)
	toStr, okTo := contentMap["to"].(string)
	/* 	if !okFrom || !okTo {
		return generateResponse("handleTimeBookingResponse", true, "Missing 'from' or 'to' in timebooking content")
	} */

	if okFrom && !okTo { /////////////////////////////////////////// only clock in

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) {
			log.Printf("ERROR: From string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From string not in valid DateTime format.")
		}
		// insert from to new booking
		log.Printf(`got "from" booking: %s`, fromStr)

		// Check if user has alrady open bookings
		//ToDO - Fehler - wird immer erstellt
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if openBooking != zeroBooking { // Error - user has open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		// insert booking
		errInsert := insertBooking(user.Id.(int64), fromStr, "", "")
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}

		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if !okFrom && okTo { ///////////////////////////////////////// only clock out

		// check for valid dateTime Format
		if !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: To string not in valid DateTime format.")
		}

		// insert "to" to existing booking
		log.Printf(`got "to" booking: %s`, toStr)

		/* 	// Check if user has alrady open bookings
		res, err := getOpenBookings(user)
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}
		if !res { // Error - user has no open bookings
			return generateResponse("handleTimeBookingResponse", true, "Error: User has no open booking")
		} */

		// Check if user has alrady open bookings
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}

		if openBooking == zeroBooking {
			log.Printf("Error: No open booking found. ")
			return generateResponse("handleTimeBookingResponse", true, "Error: User has no open booking")
		}

		duration, err := calcDuration(openBooking.From, toStr)
		if err != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: Error parsing dates for duration calculation")
		}
		// ToDo: calculate duration
		errInsert := insertBooking(user.Id.(int64), "", toStr, duration)
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	} else if okFrom && okTo { ///////////////////////////////////////// full clocking

		// check for valid dateTime Format
		if !checkDateTimeFormat(fromStr) || !checkDateTimeFormat(toStr) {
			log.Printf("ERROR: From or To string not in valid DateTime format.")
			return generateResponse("handleTimeBookingResponse", true, "ERROR: From or To string not in valid DateTime format.")
		}

		// insert from and to
		log.Printf(`got full booking: %s - %s `, fromStr, toStr)

		// Check if user has alrady open bookings
		openBooking, err := getOpenBookings(user)
		var zeroBooking Booking
		if err != nil {
			log.Printf("Error db check for open bookings: %s", err)
			return generateResponse("handleTimeBookingResponse", true, "Error: cannot check fo open bookings")
		}

		if openBooking != zeroBooking {
			log.Printf("Error: found open booking from: %s", openBooking.From)
			return generateResponse("handleTimeBookingResponse", true, "Error: User has already open booking")
		}

		duration, err := calcDuration(fromStr, toStr)
		if err != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: Error parsing dates for duration calculation")
		}

		errInsert := insertBooking(user.Id.(int64), fromStr, toStr, duration)
		if errInsert != nil {
			return generateResponse("handleTimeBookingResponse", true, "Error: DB error: "+errInsert.Error())
		}
		return generateResponse("handleTimeBookingResponse", false, "booking processed")

	}

	// error - got no values or other error
	return generateResponse("handleTimeBookingResponse", true, "Error: Missing / error timeBooking values.")

}

func handleGetBookings(content interface{}, user User) ([]byte, error) {
	bookings, err := dbGetBookings(user)
	if err != nil {
		return generateResponse("handleGetBookingsResponse", true, "Error: Error getting Bookings: "+err.Error())
	}

	return generateResponse("handleGetBookingsResponse", false, bookings)
}

func handleGetOpenBookings(content interface{}, user User) ([]byte, error) {
	bookings, err := getOpenBookings(user)
	if err != nil {
		return generateResponse("handleGetOpenBookingsResponse", true, "Error: Error getting openBookings: "+err.Error())
	}

	return generateResponse("handleGetOpenBookingsResponse", false, bookings)
}

func handleStartRecalc(content interface{}, user User) ([]byte, error) {

	/* 	// ToDo: Fehler - summe aus SQL Select sind nicht minuten (60er Basis)
	   	var monthSollZeit float64 = 10
	   	var monthCount float64 = 20
	   	log.Printf("Sollzeit: %f", monthSollZeit)
	   	var sumSoll = monthSollZeit * monthCount
	   	log.Printf("Sollzeit gesamt: %f", sumSoll)
	   	sumIst, _ := getFullTimeAccountings(user)
	   	log.Printf("Istzeit gesamt: %f", sumIst)

	   	dif := sumIst - sumSoll
	   	valDif := math.Round(dif*100) / 100
	   	log.Printf("Dif: %f", valDif)
	*/
	hoursPM, months_passed, err := getEmployeementMonths(user)
	if err != nil {
		return generateResponse("handleStartRecalcResponse", true, "Error on recalc - cannot get Epmployee data: "+err.Error())
	}

	sumSoll := hoursPM * months_passed
	sumSoll = math.Round(sumSoll*100) / 100

	sumIst, _ := getFullTimeAccountings(user)
	sumIst = math.Round(sumIst*100) / 100

	valDif := sumIst - sumSoll
	valDif = math.Round(valDif*100) / 100

	log.Printf("User Monate: %.2f, Stunden/Monat: %.2f,  User Soll: %.2f, User Ist: %.2f, User Dif: %.2f ", months_passed, hoursPM, sumSoll, sumIst, valDif)

	var summ SummCalculation
	summ.PassedMonths = fmt.Sprintf("%.2f", months_passed)
	summ.HourPerMonth = fmt.Sprintf("%.2f", hoursPM)
	summ.SollZeit = fmt.Sprintf("%.2f", sumSoll)
	summ.IstZeit = fmt.Sprintf("%.2f", sumIst)
	summ.DifZeit = fmt.Sprintf("%.2f", valDif)

	return generateResponse("handleStartRecalcResponse", false, summ)
}
