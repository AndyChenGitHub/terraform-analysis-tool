// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count       = var.vpc_id != "" ? 0 : var.create ? 1 : 0
  name        = var.vpc_name
  cidr_block  = var.vpc_cidr
  description = var.vpc_description
  tags = merge(
    {
      "Name" = format("%s", var.vpc_name)
    },
    var.vpc_tags,
  )
}

// According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count             = local.create_sub_resources ? length(var.vswitch_cidrs) : 0
  vpc_id            = var.vpc_id != "" ? var.vpc_id : concat(alicloud_vpc.vpc.*.id, [""])[0]
  cidr_block        = var.vswitch_cidrs[count.index]
  availability_zone = element(var.availability_zones, count.index)
  name              = length(var.vswitch_cidrs) > 1 || var.use_num_suffix ? format("%s%03d", var.vswitch_name, count.index + 1) : var.vswitch_name
  tags = merge(
    {
      Name = format(
        "%s%03d",
        var.vswitch_name,
        count.index + 1
      )
    }
  )
}

