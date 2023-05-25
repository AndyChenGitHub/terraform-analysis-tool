#############################################################################################
###  instance vars  ###

variable "key-name" {
  description = "name of the keypair to use for this instance"
  type        = string
}

# variable "key-name" {
#   description = "name of the ssh keypair to use"
#   type        = string
# }

# variable "instance-type" {
#   description = "size (type) of machine to use. Check Alicloud mapping of machine types for region"
#   type        = string
# }

# variable "ami" {
#   description = "ami to use for this machine"
#   type        = string
# }

variable "deployment-vpc-id" {
  description = "Alicloud ID of the VPC we want to deploy this resource into"
  type        = string
}

variable "deployment-vswitch-id" {
  description = "Alicloud ID of the vswitch we want to deploy this resource into"
  type        = string
}

# variable "external-ssh-security-group-id" {
#   description = "Alicloud ID of the security group that provides external access to SSH ports"
#   type        = string
# }

# variable "internal-ssh-security-group-id" {
#   description = "Alicloud ID of the security group that provides internal access to SSH ports"
#   type        = string
# }

# variable "ansible-user-groups" {
#   description = "ansible user groups to deploy to the node"
#   type        = string
# }
# variable "route53-hostedzone-zoneid" {
#   description = "zone-id of the hosted zone created for the deployment"
#   type        = string
# }
#############################################################################################

### TAG VALUES ###

variable "tags" {
  type = map(string)
  default = {
    Environment      = "cn-prod"
    Environment-FQDN = "cn-prod.k8s.congenica.com.cn"
    EnvironmentType  = "prod"
    Owner            = "devops"
    Managed_by       = "terraform"
  }
  description = "A map of tag values to use throughout the deployment."
}


variable "public-key-path" {
  description = "Path to the public key for the keypair"
  type        = string
}

variable "vpc-cidr" {
  type        = string
  description = "cidr ip address range assigned to VPC"
  default     = "10.68.0.0/16"
}

variable "availability-zones" {
  type        = list(string)
  default     = ["cn-beijing-h", "cn-beijing-g", "cn-beijing-f"]
  description = "A list of availability zones in which to create subnets"
}

variable "subnet-cidrs" {
  type        = list(string)
  default     = ["10.68.1.0/24", "10.68.2.0/24", "10.68.3.0/24"]
  description = "A list of subnet cidrs which to apply to the vswitchs, in order"
}