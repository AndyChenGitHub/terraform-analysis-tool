terraform {
  backend "s3" {
    bucket         = "congenica-terraform-state"
    key            = "us-prod/us-prod.tfstate"
    region         = "eu-west-1"
    encrypt        = true
    profile        = "sso-services-admin"
    dynamodb_table = "tf-state-terraform"
  }
  required_version = ">= 1.0"
}



provider "aws" {
  region  = "us-east-1"
  profile = "sso-prod-admin"
}
