package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Products struct {
	TencentcloudProductName string    `db:"tencentcloud_product_name"`
	OtherProductName        string    `db:"other_product_name"`
	Company                 string    `db:"company"`
	UpdateTime              time.Time `db:"update_time"`
}
type Resources struct {
	OtherResource        string    `db:"other_resource"`
	TencentcloudResource string    `db:"tencentcloud_resource"`
	TencentcloudTfUrl    string    `db:"tencentcloud_tf_url"`
	OtherTfUrl           string    `db:"other_tf_url"`
	Company              string    `db:"company"`
	UpdateTime           time.Time `db:"update_time"`
	Remark               string    `db:"remark"`
}

type DataSources struct {
	OtherDataSource        string    `db:"other_data_source"`
	TencentcloudDataSource string    `db:"tencentcloud_data_source"`
	TencentcloudTfUrl      string    `db:"tencentcloud_tf_url"`
	OtherTfUrl             string    `db:"other_tf_url"`
	Company                string    `db:"company"`
	UpdateTime             time.Time `db:"update_time"`
	Remark                 string    `db:"remark"`
}

type ProductResourceRule struct {
	ProductName  string `db:"product_name"`
	ResourceName string `db:"resource_name"`
	Company      string `db:"company"`
}

type ProductDataSourceRule struct {
	ProductName    string `db:"product_name"`
	DataSourceName string `db:"data_source_name"`
	Company        string `db:"company"`
}

func getProducts(company string) []Products {
	var products []Products
	clauses := make([]clause.Expression, 0)
	if company != "" {
		clauses = append(clauses, clause.Eq{Column: "company", Value: company})
	}
	getDB().Debug().Clauses(clauses...).Find(&products) // and tencentcloud_product_name!=""
	return products
}

func getResources(company string) []Resources {
	var resources []Resources
	clauses := make([]clause.Expression, 0)
	if company != "" {
		clauses = append(clauses, clause.Eq{Column: "company", Value: company})
	}
	getDB().Clauses(clauses...).Find(&resources)
	return resources
}
func getDataSources(company string) []DataSources {
	var resources []DataSources
	clauses := make([]clause.Expression, 0)
	if company != "" {
		clauses = append(clauses, clause.Eq{Column: "company", Value: company})
	}
	getDB().Clauses(clauses...).Find(&resources)
	return resources
}

func getProductResourceRules(company string) []ProductResourceRule {
	var productResourceRule []ProductResourceRule
	clauses := make([]clause.Expression, 0)
	if company != "" {
		clauses = append(clauses, clause.Eq{Column: "company", Value: company})
	}
	getDB().Debug().Clauses(clauses...).Find(&productResourceRule)
	return productResourceRule
}

func getProductDataSourceRules(company string) []ProductDataSourceRule {
	var productResourceRule []ProductDataSourceRule
	clauses := make([]clause.Expression, 0)
	if company != "" {
		clauses = append(clauses, clause.Eq{Column: "company", Value: company})
	}
	getDB().Debug().Clauses(clauses...).Find(&productResourceRule)
	return productResourceRule
}

func getDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
