package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type TencentCloudAlibabacloudProduct struct {
	gorm.Model
	tencentCloudProductName string
	alibabaCloudProductName string
	UpdateTime              time.Time
}
type TencentCloudAlibabacloudResource struct {
	gorm.Model
	alibabaCloudResource string
	tencentCloudResource string
	sourcesType          string
	tencentCloudTfUrl    string
	UpdateTime           time.Time
}

func getTencentCloudAlibabacloudProduct() []TencentCloudAlibabacloudProduct {
	var products []TencentCloudAlibabacloudProduct
	getDB().Select(&products)
	return products
}

func getTencentCloudAlibabacloudResource(sourcesType string) []TencentCloudAlibabacloudResource {
	var products []TencentCloudAlibabacloudResource
	getDB().Select(&products, "sources_type = ?", sourcesType) // 查找 code 字段值为 D42 的记录
	return products
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
