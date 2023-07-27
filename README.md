打包

mac:
GOOS=darwin GOARCH=amd64 go build -o terraform_analysis_tool_macos  

windows:
GOOS=windows GOARCH=amd64 go build -o terraform_analysis_tool_win.exe