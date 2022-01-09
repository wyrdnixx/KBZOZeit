package models

type UserDevice struct {
	Uuid   string `json:"Uuid"`
	Name   string `json:"Name"`
	Status string `json:"Status"`
}
type User struct {
	Name    string `json:"Name"`
	Enabled int    `json:"Enabled"`
}
type Users struct {
	User []User
}

type ReceivedMessage struct {
	MsgType string `json:"MsgType"`
	//MsgData []struct{} // `json:"MsgData"`
}

type TimeAccountingMessage struct {
	MsgType  string `json:MsgType`
	Name     string `json:Name`
	Typ      string `json:Typ`
	FromDate string `json:fromDate`
	ToDate   string `json:toDate`
}

type GetAccountingsMessage struct {
	MsgType  string `json:"MsgType"`
	User     string `json:"Name"`
	FromDate string `json:"FromDate"`
	ToDate   string `json:"ToDate"`
}
