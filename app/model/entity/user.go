package entity

type User struct {
	Id       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}
