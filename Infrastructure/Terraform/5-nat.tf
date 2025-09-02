resource "aws_eip" "nat-eip" {
    domain = "vpc"

    tags = {
      Name = "${local.env}-nat-eip"
    }
}

resource "aws_nat_gateway" "nat-gw" {
    allocation_id = aws_eip.nat-eip.id
    subnet_id     = aws_subnet.public-1-zone1.id

    tags = {
      Name = "${local.env}-nat-gw"
    }

    depends_on = [aws_internet_gateway.igw-eks-main]
  
}
