package main

func main() {

	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: go run main.go <path_to_terraform_module_file>")
	//	os.Exit(1)
	//}
	//filePath := os.Args[1]
	filePath := "/Users/andy/go/src/github/devops-terraform-master/prod/cn-prod/main.tf"
	//filePath := "https://github.com/terraform-tencentcloud-modules/terraform-tencentcloud-vpc/blob/master/main.tf"
	//filePath := "/Users/andy/go/src/github/terraform-analysis-tool/tf_example/web-vpc/main.tf"
	//filePath := "/Users/andy/go/src/github/terraform-analysis-tool/tf_example/local/main.tf"
	//filePath := "/Users/andy/go/src/github/terraform-analysis-tool/tf_example/web-vpc-aws/main.tf"

	tfRead(filePath)
}
