resource "aws_subnet" "private-1-zone1" {
    vpc_id            = aws_vpc.vpc-eks-main.id
    cidr_block        = "10.0.0.0/24"
    availability_zone = local.zone1

    tags = {
      "Name" = "${local.env}-private-1-zone1"
      "kubernetes.io/role/internal-elb" = "1"
      "kubernetes.io/cluster/${local.eks_name}" = "shared"
    }
}

resource "aws_subnet" "private-1-zone2" {
    vpc_id            = aws_vpc.vpc-eks-main.id
    cidr_block        = "10.0.1.0/24"
    availability_zone = local.zone2

    tags = {
      "Name" = "${local.env}-private-1-zone2"
      "kubernetes.io/role/internal-elb" = "1"
      "kubernetes.io/cluster/${local.eks_name}" = "shared"
    }
}

resource "aws_subnet" "public-1-zone1" {
    vpc_id            = aws_vpc.vpc-eks-main.id
    cidr_block        = "10.0.2.0/24"
    availability_zone = local.zone1
    map_public_ip_on_launch = false

    tags = {
      "Name" = "${local.env}-public-1-zone1"
      "kubernetes.io/role/elb" = "1"
      "kubernetes.io/cluster/${local.eks_name}" = "shared"
    }
  
}

resource "aws_subnet" "public-1-zone2" {
    vpc_id            = aws_vpc.vpc-eks-main.id
    cidr_block        = "10.0.3.0/24"
    availability_zone = local.zone2
    map_public_ip_on_launch = false

    tags = {
      "Name" = "${local.env}-public-1-zone2"
      "kubernetes.io/role/elb" = "1"
      "kubernetes.io/cluster/${local.eks_name}" = "shared"
    }
  
}
