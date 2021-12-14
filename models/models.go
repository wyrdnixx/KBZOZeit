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
