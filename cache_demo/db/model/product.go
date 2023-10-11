package model

// Product 结构体对应 "products" 表
type Product struct {
	ID    uint `gorm:"primaryKey"`
	PID   uint
	Name  string `gorm:"not null"`
	Price float64
	Count int `gorm:"not null"`
}
