package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func Log(level int, function string, text string) {
	t := time.Now()
	Timetext := t.Format(time.RFC3339)
	lvl := ""
	switch level {
	case 1:
		lvl = "Debug"
	case 2:
		lvl = "Warning"
	case 3:
		lvl = "Error"
	case 4:
		lvl = "Error-Fatal"
	default:
	}

	LogText := Timetext + " : " + lvl + " : " + function + " : " + text + "\n"

	file, err := os.OpenFile("Logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString(LogText); err != nil {
		log.Fatal(err)
	}

	if level == 4 {
		log.Fatal(LogText)
	} else {
		fmt.Println(LogText)
	}

	//fmt.Printf("%v : msg-lvl: %v msg-text: %v \n", t.Format(time.RFC3339), level, text)
}
