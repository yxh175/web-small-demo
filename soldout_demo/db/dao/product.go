package dao

import (
	"context"
	"gin-mall/cache_demo/db/model"

	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{NewDBClient(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

// 获取商品
func (dao *ProductDao) GetProduct(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where("p_id=?", id).First(&product).Error
	return
}

// CreateProduct 创建商品
func (dao *ProductDao) CreateProduct(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).
		Create(&product).Error
}

// DeleteProduct 删除商品
func (dao *ProductDao) DeleteProduct(id uint) error {
	return dao.DB.Model(&model.Product{}).
		Where("p_id = ?", id).
		Delete(&model.Product{}).
		Error
}

// UpdateProduct 更新商品
func (dao *ProductDao) UpdateProduct(id uint, product *model.Product) error {
	return dao.DB.Model(&model.Product{}).
		Where("p_id=?", id).Updates(&product).Error
}
