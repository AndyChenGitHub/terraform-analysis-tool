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
		if r < len(newRows) && newRows[i][0] == newRows[i+1][0] && newRows[i][0] != "" {
			f.MergeCell(sheetName, "A"+strconv.Itoa(r), "A"+strconv.Itoa(r+1))
		}
		if r < len(newRows) && newRows[i][1] == newRows[i+1][1] && newRows[i][1] != "" {
			f.MergeCell(sheetName, "B"+strconv.Itoa(r), "B"+strconv.Itoa(r+1))
		}
		if r < len(newRows) && len(newRows[i+1]) > 3 && newRows[i][3] != "" && newRows[i][3] == newRows[i+1][3] {
			f.MergeCell(sheetName, "D"+strconv.Itoa(r), "D"+strconv.Itoa(r+1))
		}
		if r < len(newRows) && len(newRows[i+1]) > 4 && len(newRows[i]) > 4 && newRows[i][4] != "" && newRows[i][4] == newRows[i+1][4] {
			f.MergeCell(sheetName, "E"+strconv.Itoa(r), "E"+strconv.Itoa(r+1))
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
