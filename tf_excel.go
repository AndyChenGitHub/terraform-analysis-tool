package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

func addExcel(datas [][]string) {
	f, err := excelize.OpenFile("tf_template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	sheetName := "分析报告" //+ cast.ToString(time.Now().Minute())

	_ = f.DeleteSheet(sheetName)
	_, _ = f.NewSheet(sheetName)
	for i, data := range datas {
		rowData := data[1:8] // 移除第一个数组
		e := f.SetSheetRow(sheetName, "A"+strconv.Itoa(i+1), &rowData)

		borderStyle, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Underline: "single", Color: "0000FF"},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"}})
		if len(data) > 8 && data[8] != "" {
			err = f.SetCellStyle(sheetName, "B"+strconv.Itoa(i+1), "B"+strconv.Itoa(i+1), borderStyle)
			f.SetCellHyperLink(sheetName, "B"+strconv.Itoa(i+1), data[8], "External")
		}
		if len(data) > 9 && data[9] != "" {
			err = f.SetCellStyle(sheetName, "E"+strconv.Itoa(i+1), "E"+strconv.Itoa(i+1), borderStyle)
			f.SetCellHyperLink(sheetName, "E"+strconv.Itoa(i+1), data[9], "External")
		}

		//备注设置自动换行
		style, _ := f.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{WrapText: true}})
		f.SetCellStyle(sheetName, "G"+strconv.Itoa(i+1), "G"+strconv.Itoa(i+1), style)

		if e != nil {
			fmt.Println(e)
		}
	}

	mergeCell(sheetName, f)
	//表头样式
	borderStyle, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}, Fill: excelize.Fill{Type: "pattern", Color: []string{"808080"}, Pattern: 1}})
	err = f.SetCellStyle(sheetName, "A1", "G1", borderStyle)
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
		if r < len(newRows) && len(newRows[i+1]) > 6 && len(newRows[i]) > 6 && newRows[i][6] != "" && newRows[i][6] == newRows[i+1][6] {
			f.MergeCell(sheetName, "G"+strconv.Itoa(r), "G"+strconv.Itoa(r+1))
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
