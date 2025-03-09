Test Project - 60


# start the main app
go run .

# run tests
go test

# run vuejs developmend
docker compose up


# example JSON Messages

{
	"type":"timebooking",
	"content": {
		"from": "01.01.2020 00:00:00",
		"to": "01.01.2020 00:01:00"
	}
}


{
	"type":"clocking",
	"content":"clockIn"
}