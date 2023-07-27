package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var selectedFilePath string

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("terraform analysis tool")

	filePathLabel := widget.NewLabel("No file selected")

	chooseFileButton := widget.NewButton("Choose main.tf File", func() {
		fd := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				selectedFilePath = file.URI().Path()
				filePathLabel.SetText(selectedFilePath)
			}
		}, myWindow)
		fd.Show()
	})

	openFolderButton := widget.NewButton("Open File Folder", func() {
		if selectedFilePath != "" {
			if !strings.Contains(selectedFilePath, "main.tf") {
				dialog.ShowInformation("Info", "Please choose a main.tf File first.", myWindow)
				return
			}

			err := tfRead(selectedFilePath)
			if err != nil {
				dialog.ShowInformation("Info", err.Error(), myWindow)
			}
			path := GetAppPath() + "/tf_report.xlsx"
			openFileWithDefaultProgram(path)
		} else {
			dialog.ShowInformation("Info", "Please choose a main.tf File first.", myWindow)
		}
	})

	content := container.NewVBox(
		chooseFileButton,
		filePathLabel,
		openFolderButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 500))
	myWindow.ShowAndRun()
}

func openFileFolder(filePath string) {
	fileDir := filepath.Dir(filePath)

	switch runtime.GOOS {
	case "windows":
		exec.Command("explorer", "/select,"+filePath).Start()
	case "darwin":
		exec.Command("open", fileDir).Start()
	case "linux":
		exec.Command("xdg-open", fileDir).Start()
	default:
		panic("unsupported platform")
	}
}

func openFileWithDefaultProgram(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Run()
}
