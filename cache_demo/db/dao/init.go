package dao

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var _db *gorm.DB

func InitMySQL() {
	// 假装两个mysql服务器，进行读写分离, 电商读多写少
	read1_dsn := "root:1234@tcp(localhost:3306)/cache_demo?parseTime=true"
	read2_dsn := "root:1234@tcp(localhost:3306)/cache_demo?parseTime=true"
	write1_dsn := "root:1234@tcp(localhost:3306)/cache_demo?parseTime=true"
	write2_dsn := "root:1234@tcp(localhost:3306)/cache_demo?parseTime=true"

	db, err := gorm.Open(mysql.Open(write1_dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_db = db
	_db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(write2_dsn)},
			Replicas: []gorm.Dialector{mysql.Open(read1_dsn), mysql.Open(read2_dsn)},
			Policy:   dbresolver.RandomPolicy{},
			// print sources/replicas mode in logger
			TraceResolverMode: true,
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	_db = _db.Set("gorm:table_options", "charset=utf8mb4")
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
