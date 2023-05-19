package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

// 读取TF文件
func tfRead(filePath string) {
	blocks := getBlocks(filePath)
	paths, providerName, datas := getSourcePath(blocks)
	//为了避免重复查询数据库，数据前置
	products := getProducts("")
	productsRule := getProductResourceRules("")
	productsDataRule := getProductDataSourceRules("")
	resources := getResources(providerName)
	dataSources := getDataSources(providerName)
	for _, path := range paths {
		newPaths := getPath(filePath, path)
		for _, newPath := range newPaths {
			resourceBlocks := getBlocks(newPath)
			for _, block := range resourceBlocks {

				ty := "d"
				var productName, tencentCloudProductName, tencentCloudResources, remark, otherTfUrl, tfUrl = "", "", "", "", "", ""

				var tencentCloudStackByte []byte
				if block.Type == "resource" {
					resRule := getResourceRules(block.Labels[0], providerName, productsRule)
					pd := getTencentCloudProduct(resRule.ProductName, providerName, products)
					res := getResourcesByOtherResource(block.Labels[0], providerName, resources)
					productName = resRule.ProductName
					tencentCloudProductName = pd.TencentcloudProductName
					tencentCloudResources = res.TencentcloudResource
					remark = res.Remark
					otherTfUrl = res.OtherTfUrl
					tfUrl = res.TencentcloudTfUrl
					ty = "r"
				} else if block.Type == "data" {
					resRule := getDataSourceRules(block.Labels[0], providerName, productsDataRule)
					pd := getTencentCloudProduct(resRule.ProductName, providerName, products)
					dataRes := getDataSourcesByOtherResource(block.Labels[0], providerName, dataSources)
					productName = resRule.ProductName
					tencentCloudProductName = pd.TencentcloudProductName
					tencentCloudResources = dataRes.TencentcloudDataSource
					remark = "data资源，" + dataRes.Remark
					otherTfUrl = dataRes.OtherTfUrl
					tfUrl = dataRes.TencentcloudTfUrl
				} else {
					continue
				}
				if tencentCloudResources != "" {
					tencentCloudStackByte = getTencentCloudStackMarkdown(tencentCloudResources, ty)
					tencentCloudResources = "tencentcloud_" + tencentCloudResources
				}
				//获取参数内容
				for _, attr := range block.Body.Attributes {
					if attr.Name == "count" {
						continue
					}
					if filterRepeat(datas, productName, block.Labels[0], attr.Name) {
						continue
					}
					//typeName := reflect.TypeOf(attr.Expr).Elem().Name()
					//if typeName == "ScopeTraversalExpr" {
					//	vb := attr.Expr.Variables()
					//	fmt.Printf("  %s = %s.%s \n", attr.Name, vb[0][0].(hcl.TraverseRoot).Name, vb[0][1].(hcl.TraverseAttr).Name)
					//} else if typeName == "TemplateExpr" {
					//	exp := attr.Expr.(*hclsyntax.TemplateExpr)
					//	v, _ := exp.Parts[0].Value(nil)
					//	fmt.Printf("  %s = %s \n", attr.Name, v.AsString())
					//} else if typeName == "FunctionCallExpr" {
					//	fmt.Printf("  %s = %s \n", attr.Name, "element")
					//}
					arg := getTencentCloudStackArg(tencentCloudStackByte, attr.Name)
					data := []string{
						productName + block.Labels[0] + attr.Name, //用于索引排序
						productName,
						block.Labels[0],
						attr.Name,
						tencentCloudProductName,
						tencentCloudResources,
						arg,
						remark,
						otherTfUrl,
						tfUrl,
					}
					datas = append(datas, data)
				}
			}
		}

	}

	title := datas[:1]    // 取title
	rowDatas := datas[1:] // 取后面的数据，进行排序
	sort.Slice(rowDatas, func(i, j int) bool {
		//前两列一样的情况下比较第二列
		if rowDatas[i][0] == rowDatas[j][0] {
			return rowDatas[i][1] < rowDatas[j][1]
		}
		return rowDatas[i][0] < rowDatas[j][0]
	})
	title = append(title, rowDatas...)
	addExcel(title)
}

// 过滤重复的资源,重复不需要计算
func filterRepeat(datas [][]string, productName string, sourcesName string, arg string) bool {
	for _, ds := range datas {
		if ds[1] == productName && ds[2] == sourcesName && ds[3] == arg {
			return true
		}
	}

	return false
}

func getDataSourcesByOtherResource(otherResource string, providerName string, res []DataSources) DataSources {
	var r DataSources
	for _, d := range res {
		if providerName+"_"+d.OtherDataSource == otherResource {
			return d
		}
	}
	return r
}

func getResourcesByOtherResource(otherResource string, providerName string, res []Resources) Resources {
	var r Resources
	for _, d := range res {
		if providerName+"_"+d.OtherResource == otherResource {
			r = d
			break
		}
	}
	return r
}

func getResourceRules(resource string, providerName string, ruls []ProductResourceRule) ProductResourceRule {
	var r ProductResourceRule
	for _, d := range ruls {
		if providerName+"_"+d.ResourceName == resource && d.Company == providerName {
			r = d
			break
		}
	}
	return r
}

func getDataSourceRules(resource string, providerName string, ruls []ProductDataSourceRule) ProductDataSourceRule {
	var r ProductDataSourceRule
	for _, d := range ruls {
		if providerName+"_"+d.DataSourceName == resource && d.Company == providerName {
			r = d
			break
		}
	}
	return r
}

func getTencentCloudProduct(otherProduct string, providerName string, products []Products) Products {
	var r Products
	for _, d := range products {
		if d.OtherProductName == otherProduct && d.Company == providerName {
			r = d
			break
		}
	}
	return r
}

// 获取module下面的resource路径，主要是main入口
func getPath(mainPath string, path string) []string {
	newPath := ""
	var retrunPaths []string
	sp := strings.Split(path, "/")
	sp1 := strings.Split(mainPath, "/")
	index := 1
	for _, s := range sp {
		if s == ".." {
			index++
		}
	}

	for i, s := range sp1 {
		max := len(sp1) - index
		if i < max {
			newPath += s + "/"
		}
	}

	root := newPath + path[index*2:]
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".tf" {
			retrunPaths = append(retrunPaths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return retrunPaths
}

// 获取tf里面的block
func getBlocks(filePath string) hclsyntax.Blocks {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, filePath)
	if diags.HasErrors() {
		fmt.Printf("Error parsing HCL: %s\n", diags)
		os.Exit(1)
	}

	body, _ := file.Body.(*hclsyntax.Body)
	return body.Blocks
}

// 读取URL文件
func getURLBlocks(url string) hclsyntax.Blocks {
	filePath := "main.tf"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	f, err := os.Create(filePath)
	io.Copy(f, resp.Body)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, filePath)
	if diags.HasErrors() {
		fmt.Printf("Error parsing HCL: %s\n", diags)
		os.Exit(1)
	}

	body, _ := file.Body.(*hclsyntax.Body)
	return body.Blocks
}

// 获取程序入口
func getSourcePath(blocks hclsyntax.Blocks) ([]string, string, [][]string) {
	var paths []string
	var providerName string
	var datas [][]string
	for _, block := range blocks {
		switch block.Type {
		case "module":
			for _, attr := range block.Body.Attributes {
				//是否为变量值
				typeName := reflect.TypeOf(attr.Expr).Elem().Name()
				if typeName == "TemplateExpr" {
					exp := attr.Expr.(*hclsyntax.TemplateExpr)
					v, _ := exp.Parts[0].Value(nil)
					if attr.Name == "source" {
						paths = append(paths, v.AsString())
					}
				}
			}
		case "provider":
			providerName = block.Labels[0]
		case "locals":
		default:
			if providerName == "" {
				providerName = "ali" //默认阿里
			}
			for _, attr := range block.Body.Attributes {
				if attr.Name == "count" {
					continue
				}
				datas = append(datas, []string{"ID", block.Labels[0], "", attr.Name})
			}
		}

	}
	newRow := []string{"ID", providerName + "_product", providerName, providerName + "_arg", "tencentcloud_product", "tencentcloud", "tencentcloud_arg", "remark"}
	datas = append([][]string{newRow}, datas...)
	return paths, providerName, datas
}
