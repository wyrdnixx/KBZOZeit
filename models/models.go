package models

type UserDevice struct {
	Uuid     int    `json:"Uuid"`
	UserName string `json:"UserName"`
}
type User struct {
	Name    string `json:"Name"`
	Enabled int    `json:"Enabled"`
}
type Users struct {
	User []User
}
