package model

import "time"

type Order struct {
	Id         uint
	Uid        int
	Pid        int
	Name       string
	Count      int
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
}
