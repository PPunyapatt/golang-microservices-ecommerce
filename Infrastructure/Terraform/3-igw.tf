resource "aws_internet_gateway" "igw-eks-main" {
    vpc_id = aws_vpc.vpc-eks-main.id

    tags = {
        Name = "${local.env}-igw-eks-main"
    }
}