package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {

	http.HandleFunc("/genFile", genFile)     //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: go run main.go <path_to_terraform_module_file>")
	//	os.Exit(1)
	//}
	//filePath := os.Args[1]
	//filePath := "/Users/andy/go/src/github/devops-terraform-master/prod/cn-prod/main.tf"
	//filePath := "https://github.com/terraform-tencentcloud-modules/terraform-tencentcloud-vpc/blob/master/main.tf"
	//filePath := "/Users/andy/go/src/github/terraform-analysis-tool/tf_example/web-vpc/main.tf"
	//filePath := "/Users/andy/go/src/github/terraform-analysis-tool/tf_example/local/main.tf"
	//read, err := tfRead(filePath)

}

func genFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("开始接收入参")
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	filePath := strings.Join(r.Form["path"], "")
	tfRead(filePath)
	m := make(map[string]string)
	m["msg"] = "success"
	marshal, _ := json.Marshal(m)
	w.Write(marshal)
}
