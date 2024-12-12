package models

import "time"

type Stats struct {
	ID        uint64    `gorm:"serial;primary_key"`
	Count     uint64    `gorm:"count"`
	Timestamp time.Time `gorm:"timestamp"`
}
