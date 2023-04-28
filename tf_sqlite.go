package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type TencentcloudAlibabacloudProduct struct {
	gorm.Model
	tencentcloudProductName string
	alibabacloudProductName string
	UpdateTime              time.Time
}
type TencentcloudAlibabacloudResource struct {
	gorm.Model
	alibabacloudResource string
	tencentcloudResource string
	sourcesType          string
	tencentcloudTfUrl    string
	UpdateTime           time.Time
}

func getTencentcloudAlibabacloudProduct() []TencentcloudAlibabacloudProduct {
	var products []TencentcloudAlibabacloudProduct
	getDB().Select(&products, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	return products
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
