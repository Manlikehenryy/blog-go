package models

type Blog struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Image  string `json:"image"`
	UserId string `json:"userId" gorm:"column:user_id"`
	User   User `json:"user" gorm:"foreignKey:UserId"`
}
