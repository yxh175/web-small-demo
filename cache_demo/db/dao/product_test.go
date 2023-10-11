package dao

import (
	"context"
	"fmt"
	"gin-mall/cache_demo/db/model"
	"testing"

	"gorm.io/gorm"
)

func TestProductDB(t *testing.T) {
	InitMySQL()
	productDao := NewProductDao(context.Background())

	// 增
	err := productDao.CreateProduct(&model.Product{
		PID:   1,
		Name:  "手机",
		Price: 123,
		Count: 22,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 查
	product, err := productDao.GetProduct(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(product)

	// 改
	err = productDao.UpdateProduct(1, &model.Product{
		Name:  "新手机",
		Price: 200,
		Count: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	// 查
	product, err = productDao.GetProduct(1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("----- 新查询 ------")
	fmt.Println(product)

	// 删
	err = productDao.DeleteProduct(1)
	if err != nil {
		t.Fatal(err)
	}

	// 查
	product, err = productDao.GetProduct(1)
	if err != nil && err != gorm.ErrRecordNotFound {
		t.Fatal(err)
	}
	fmt.Println("----- 新查询 ------")
	fmt.Println(product)
}
