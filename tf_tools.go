package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

func tfRead(filePath string) {
	blocks := getBlocks(filePath)
	paths := getSourcePath(blocks)
	var datas [][]interface{}
	//data := []interface{}{"alicloud_eip", "Elastic IP Address (EIP)", "ID", "tencentcloud_eip", "Elastic IP (EIP) - 弹性公网 IP", "ID"}
	for _, path := range paths {
		newPath := getPath(filePath, path)
		resourceBlocks := getBlocks(newPath)
		for _, block := range resourceBlocks {
			switch block.Type {
			case "resource":
				fmt.Printf("Resource: %s\n", block.Labels[0])
			case "data":
				fmt.Printf("Data Source: %s\n", block.Labels[0])
			case "variable":
				fmt.Printf("Variable: %s\n", block.Labels[0])
			case "module":
				fmt.Printf("module: %s\n", block.Labels[0])
			case "provider":
				fmt.Printf("provider: %s\n", block.Labels[0])
			}

			for _, attr := range block.Body.Attributes {
				if attr.Name == "count" {
					continue
				}
				typeName := reflect.TypeOf(attr.Expr).Elem().Name()
				if typeName == "ScopeTraversalExpr" {
					vb := attr.Expr.Variables()
					fmt.Printf("  %s = %s.%s \n", attr.Name, vb[0][0].(hcl.TraverseRoot).Name, vb[0][1].(hcl.TraverseAttr).Name)
				} else if typeName == "TemplateExpr" {
					exp := attr.Expr.(*hclsyntax.TemplateExpr)
					v, _ := exp.Parts[0].Value(nil)
					fmt.Printf("  %s = %s \n", attr.Name, v.AsString())
				} else if typeName == "FunctionCallExpr" {
					fmt.Printf("  %s = %s \n", attr.Name, "element")
				}
				datas = append(datas, []interface{}{block.Labels[0], "", attr.Name})
			}
		}
	}
	addExcel(datas)
}

func getPath(mainPath string, path string) string {
	newPath := ""
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
	return newPath + path[index*2:] + "/main.tf"
}

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

func getSourcePath(blocks hclsyntax.Blocks) []string {
	var paths []string
	for _, block := range blocks {
		switch block.Type {
		case "resource":
			fmt.Printf("Resource: %s\n", block.Labels[0])
		case "data":
			fmt.Printf("Data Source: %s\n", block.Labels[0])
		case "variable":
			fmt.Printf("Variable: %s\n", block.Labels[0])
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
			fmt.Printf("provider: %s\n", block.Labels[0])
		}
	}
	return paths
	//取参数和值
	//for _, attr := range block.Body.Attributes {
	//	//是否为变量值
	//	typeName := reflect.TypeOf(attr.Expr).Elem().Name()
	//	if typeName == "ScopeTraversalExpr" {
	//		vb := attr.Expr.Variables()
	//		fmt.Printf("  %s = %s.%s \n", attr.Name, vb[0][0].(hcl.TraverseRoot).Name, vb[0][1].(hcl.TraverseAttr).Name)
	//	} else if typeName == "TemplateExpr" {
	//		exp := attr.Expr.(*hclsyntax.TemplateExpr)
	//		v, _ := exp.Parts[0].Value(nil)
	//		fmt.Printf("  %s = %s \n", attr.Name, v.AsString())
	//	} else if typeName == "FunctionCallExpr" {
	//		fmt.Printf("  %s = %s \n", attr.Name, "element")
	//	}
	//
	//}
}
