package main

func main() {

	//if len(os.Args) < 2 {
	//	fmt.Println("Usage: go run main.go <path_to_terraform_module_file>")
	//	os.Exit(1)
	//}
	//filePath := os.Args[1]
	filePath := "/Users/andy/go/src/github/devops-terraform-master/prod/cn-prod/main.tf"
	tfRead(filePath)
}
