package dao

import (
	"context"
	"gin-mall/soldout_demo/db/model"

	"gorm.io/gorm"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{NewDBClient(ctx)}
}

func (od *OrderDao) AddOrder(order *model.Order) error {
	return od.DB.Model(&model.Order{}).Create(order).Error
}
