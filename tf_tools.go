package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

type ResourcesData struct {
	ProviderName     string
	ProductsDataRule []ProductDataSourceRule
	Products         []Products
	Resources        []Resources
	DataSources      []DataSources
	ProductsRule     []ProductResourceRule
}

// tfRead 读取TF文件
func tfRead(filePath string) error {
	blocks := getLocalBlocks(filePath)
	//为了避免重复查询数据库，数据前置
	paths, datas, resourcesData := getSourceMain(blocks, filePath)
	if resourcesData.ProviderName == "" {
		println("provider not find！ ")

		return fmt.Errorf("provider not find！ ")
	}

	for _, path := range paths {
		newPaths := getSourcePath(filePath, path)
		for _, newPath := range newPaths {
			var resourceBlocks hclsyntax.Blocks
			if strings.Contains(newPath, "https://raw.githubusercontent.com") {
				resourceBlocks = getWebBlocks(newPath)
			} else {
				resourceBlocks = getLocalBlocks(newPath)
			}

			for _, block := range resourceBlocks {
				data := getResourceData(*block, resourcesData)
				datas = append(datas, data...)
			}
		}
	}

	newDatas := DeduplicateStringSlice(datas)
	title := newDatas[:1]    // 取title
	rowDatas := newDatas[1:] // 取后面的数据，进行排序
	sort.Slice(rowDatas, func(i, j int) bool {
		return rowDatas[i][0] < rowDatas[j][0]
	})

	title = append(title, rowDatas...)
	err := addExcel(title)
	if err != nil {
		return err
	}

	return nil
}

// getResourceData 获取资源数据，并返回报表数据
func getResourceData(block hclsyntax.Block, resourcesData ResourcesData) [][]string {
	ty := "d"
	var datas [][]string
	var productName, tencentCloudProductName, tencentCloudResources, describe, remark, otherTfUrl, tfUrl = "", "", "", "", "", "", ""

	var tencentCloudStackByte, otherMarkByte []byte
	if block.Type == "resource" {
		resRule := getResourceRules(block.Labels[0], resourcesData.ProviderName, resourcesData.ProductsRule)
		pd := getTencentCloudProduct(resRule.ProductName, resourcesData.ProviderName, resourcesData.Products)
		res := getResourcesByOtherResource(block.Labels[0], resourcesData.ProviderName, resourcesData.Resources)
		productName = resRule.ProductName
		tencentCloudProductName = pd.TencentcloudProductName
		tencentCloudResources = res.TencentcloudResource
		remark = res.Remark
		otherTfUrl = res.OtherTfUrl
		tfUrl = res.TencentcloudTfUrl
		ty = "r"
	} else if block.Type == "data" {
		resRule := getDataSourceRules(block.Labels[0], resourcesData.ProviderName, resourcesData.ProductsDataRule)
		pd := getTencentCloudProduct(resRule.ProductName, resourcesData.ProviderName, resourcesData.Products)
		dataRes := getDataSourcesByOtherResource(block.Labels[0], resourcesData.ProviderName, resourcesData.DataSources)
		productName = resRule.ProductName
		tencentCloudProductName = pd.TencentcloudProductName
		tencentCloudResources = dataRes.TencentcloudDataSource
		remark = "data资源，" + dataRes.Remark
		otherTfUrl = dataRes.OtherTfUrl
		tfUrl = dataRes.TencentcloudTfUrl
	} else {
		return datas
	}

	otherResources := strings.Replace(block.Labels[0], resourcesData.ProviderName+"_", "", 1)
	if resourcesData.ProviderName == "alicloud" {
		otherMarkByte = getAliyunMarkdown(otherResources, ty)
	} else if resourcesData.ProviderName == "aws" {
		otherMarkByte = getAwsMarkdown(otherResources, ty)
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
		if resourcesData.ProviderName == "alicloud" {
			describe = getAliyunArgDesc(otherMarkByte, attr.Name)
		} else if resourcesData.ProviderName == "aws" {
			describe = getAwsArgDesc(otherMarkByte, attr.Name)
		}

		arg, describeTen := getTencentCloudStackArgDesc(tencentCloudStackByte, attr.Name)
		describe = describe + describeTen

		data := []string{
			productName + "_" + block.Labels[0] + "_" + attr.Name, //用于索引排序
			productName,
			block.Labels[0],
			attr.Name,
			tencentCloudProductName,
			tencentCloudResources,
			arg,
			describe,
			remark,
			otherTfUrl,
			tfUrl,
		}
		datas = append(datas, data)
	}
	return datas
}

// DeduplicateStringSlice 过滤重复的资源,重复不需要计算
func DeduplicateStringSlice(slice [][]string) [][]string {
	seen := make(map[string]bool)
	deduplicated := make([][]string, 0, len(slice))
	for _, s := range slice {
		if !seen[s[0]] {
			seen[s[0]] = true
			deduplicated = append(deduplicated, s)
		}
	}
	return deduplicated
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

// getSourcePath 获取module下面的resource路径，主要是main入口
func getSourcePath(mainPath string, path string) []string {
	newPath := ""
	var retrunPaths []string
	sp := strings.Split(path, "/")
	// 当最后一个为 alicloud 且第一个为alibaba 或 terraform-alicloud-modules 认为是官网的
	// terraform-alicloud-modules/rds/alicloud,alibaba/security-group/alicloud//modules/http-80
	if len(sp) > 2 && strings.Contains("alibaba,terraform-alicloud-modules", sp[0]) && sp[2] == "alicloud" {
		retrunPaths = append(retrunPaths, "https://raw.githubusercontent.com/terraform-alicloud-modules/terraform-alicloud-"+sp[1]+"/master/main.tf")
		return retrunPaths
	} else if len(sp) > 2 && strings.Contains("aws,terraform-aws-modules", sp[0]) && sp[2] == "aws" {
		retrunPaths = append(retrunPaths, "https://raw.githubusercontent.com/terraform-aws-modules/terraform-aws-"+sp[1]+"/master/main.tf")
		return retrunPaths
	}

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
	if len(sp) == 1 {
		root = newPath + path
	}
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

// getLocalBlocks 获取本地tf里面的block
func getLocalBlocks(filePath string) hclsyntax.Blocks {
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

// getWebBlocks 读取官网的TF文件
func getWebBlocks(url string) hclsyntax.Blocks {
	content := getCloudMarkdown(url)
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(content, "main.tf")
	if diags.HasErrors() {
		fmt.Printf("Error parsing HCL: %s\n", diags)
		os.Exit(1)
	}

	body, _ := file.Body.(*hclsyntax.Body)
	return body.Blocks
}

// getSourceMain 获取程序入口
func getSourceMain(blocks hclsyntax.Blocks, filePath string) ([]string, [][]string, ResourcesData) {
	var paths []string
	var datas [][]string
	resourcesData := ResourcesData{}

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
			getProvider(*block, &resourcesData)
		case "locals":
		default:
			//一个tf语法中，一定得有 provider 信息
			if resourcesData.ProviderName != "" {
				data := getResourceData(*block, resourcesData)
				datas = append(datas, data...)
			}

			//for _, attr := range block.Body.Attributes {
			//	if attr.Name == "count" {
			//		continue
			//	}
			//	datas = append(datas, []string{"ID", block.Labels[0], "", attr.Name})
			//}
		}

	}

	// 没有找到providerName的时候，重新根据路径查找
	if resourcesData.ProviderName == "" {
		blocksProvider := getLocalBlocks(strings.Replace(filePath, "main.tf", "providers.tf", 1))
		for _, block := range blocksProvider {
			if block.Type == "provider" {
				getProvider(*block, &resourcesData)
				break
			}
		}
	}

	newRow := []string{"ID", resourcesData.ProviderName + "_product", resourcesData.ProviderName, resourcesData.ProviderName + "_arg", "tencentcloud_product", "tencentcloud", "tencentcloud_arg", "arg_describe", "remark"}
	datas = append([][]string{newRow}, datas...)
	return paths, datas, resourcesData
}

func getProvider(block hclsyntax.Block, resourcesData *ResourcesData) {
	resourcesData.ProviderName = block.Labels[0]

	resourcesData.Products = getProducts("")
	resourcesData.ProductsRule = getProductResourceRules("")
	resourcesData.ProductsDataRule = getProductDataSourceRules("")

	resourcesData.Resources = getResources(resourcesData.ProviderName)
	resourcesData.DataSources = getDataSources(resourcesData.ProviderName)
}
