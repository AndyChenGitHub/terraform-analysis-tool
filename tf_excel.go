package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func addExcel(datas [][]interface{}) {
	f, err := excelize.OpenFile("tf_template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	sheetName := "分析报告" //+ cast.ToString(time.Now().Minute())

	_ = f.DeleteSheet(sheetName)
	_, _ = f.NewSheet(sheetName)
	for i, data := range datas {
		e := f.SetSheetRow(sheetName, "A"+strconv.Itoa(i+1), &data)
		if e != nil {
			fmt.Println(e)
		}
	}
	mergeCell(sheetName, f)
	// 根据指定路径保存文件
	if err = f.SaveAs("tf_template.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func mergeCell(sheetName string, f *excelize.File) {
	newRows, err := f.GetRows(sheetName)
	for i, _ := range newRows {
		r := i + 1 //excel 是从1开始的
		if r < len(newRows) && newRows[i][0] == newRows[i+1][0] {
			e := f.MergeCell(sheetName, "A"+strconv.Itoa(r), "A"+strconv.Itoa(r+1))
			if e != nil {
				fmt.Println(e)
			}
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readExcel() {
	//f, err := excelize.OpenFile("Book1.xlsx")
	f, err := excelize.OpenFile("tf_template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
	fmt.Println("===========================================")

	res := arrayTwoStringGroupsOf(rows, 3)
	fmt.Println(res)

	//按行赋值
	e := f.SetSheetRow("Sheet1", "A3", &[]interface{}{"alicloud_eip", "Elastic IP Address (EIP)", "ID", "tencentcloud_eip", "Elastic IP (EIP) - 弹性公网 IP", "ID"})
	if e != nil {
		fmt.Println(e)
	}
	// 根据指定路径保存文件
	if err = f.SaveAs("tf_template.xlsx"); err != nil {
		fmt.Println(err)
	}
}
func arrayTwoStringGroupsOf(arr [][]string, num int64) [][][]string {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][][]string{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][][]string, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}
func writeExcel() {
	f := excelize.NewFile()
	// 创建一个工作表
	index, _ := f.NewSheet("Sheet1")
	// 设置单元格的值
	f.SetCellValue("Sheet1", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)

	//按行赋值
	err := f.SetSheetRow("Sheet1", "A1", &[]interface{}{"39 - 38 = ", "39 - 38 = ", "39 - 38 = ", "39 - 38 = ", "39 - 38 = "})
	if err != nil {
		fmt.Println(err)
	}

	//设置列宽度
	err = f.SetColWidth("Sheet1", "A", "H", 16)
	if err != nil {
		fmt.Println(err)
	}
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err = f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
