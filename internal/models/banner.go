package models

type Banner struct {
	ID   uint64 `gorm:"serial;primary_key"`
	Name string `gorm:"name"`
}
