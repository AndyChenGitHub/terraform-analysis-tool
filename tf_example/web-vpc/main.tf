
provider "alicloud" {
  version = "= 1.113" # pinned to stop deprecation warnings
  region  = "cn-beijing"
  profile = "alicloud_terraform"
  assume_role {
    role_arn = "acs:ram::1244841646161695:role/terraform"
  }
}
module "vpc" {
  source  = "alibaba/vpc/alicloud"

  create            = true
  vpc_name          = "my-env-vpc"
  vpc_cidr          = "10.10.0.0/16"

  availability_zones = ["cn-hangzhou-e", "cn-hangzhou-f", "cn-hangzhou-g"]
  vswitch_cidrs      = ["10.10.1.0/24", "10.10.2.0/24", "10.10.3.0/24"]

}
