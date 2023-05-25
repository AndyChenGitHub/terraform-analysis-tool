terraform {
  backend "s3" {
    bucket         = "congenica-terraform-state"
    key            = "cn-prod/cn-prod.tfstate"
    region         = "eu-west-1"
    encrypt        = true
    profile        = "bastion"
    dynamodb_table = "tf-state-terraform" # each terraform landscape should be its own line in this table
    role_arn       = "arn:aws:iam::024764271014:role/Devops_STS_Role"
  }
}

provider "alicloud" {
  version = "= 1.113" # pinned to stop deprecation warnings
  region  = "cn-beijing"
  profile = "alicloud_terraform"
  assume_role {
    role_arn = "acs:ram::1244841646161695:role/terraform"
  }
}

module "vpc" {
  source             = "t12-vpc"
  vpc_name           = "cn-prod"
  vpc_cidr           = var.vpc-cidr
  availability_zones = var.availability-zones
  vswitch_cidrs      = var.subnet-cidrs
  vpc_tags           = var.tags
  nat_type           = "Enhanced"
  vswitch_tags = {}
}

# Create a new ECS instance for a VPC
resource "alicloud_security_group" "group" {
  name        = "tf_test_foo"
  description = "foo"
  vpc_id      = var.deployment-vpc-id
}

resource "alicloud_instance" "instance" {
  key_name = var.key-name
  security_groups = alicloud_security_group.group.*.id
  instance_type              = "ecs.xn4.small"
  image_id                   = "centos_7_7_x64_20G_alibase_20200426.vhd"
  instance_name              = "generic-1"
  internet_max_bandwidth_out = 10
  host_name = join("", [var.tags.Environment, "-generic-node"])
  tags        = var.tags
}

// Import an existing public key to build a alicloud key pair
resource "alicloud_key_pair" "key" {
  key_name   = join("", ["terraform-", var.tags.Environment])
  public_key = file(var.public-key-path)

  tags = var.tags
}
