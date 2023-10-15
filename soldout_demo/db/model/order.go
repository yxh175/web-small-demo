package model

import "time"

type Order struct {
	Id         uint
	Gid        int
	Name       string
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
}
