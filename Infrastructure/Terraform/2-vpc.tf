resource "aws_vpc" "vpc-eks-main" {
    cidr_block = "10.0.0.0/16"

    enable_dns_support = true
    enable_dns_hostnames = true

    tags = {
        Name = "${local.env}-vpc-eks-main"
    }
}