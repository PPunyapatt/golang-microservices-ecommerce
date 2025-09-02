resource "aws_route_table" "private-rt" {
  vpc_id = aws_vpc.vpc-eks-main.id

  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.nat-gw.id
  }

  tags = {
    Name = "${local.env}-private-rt"
  }
}

resource "aws_route_table" "public-rt" {
  vpc_id = aws_vpc.vpc-eks-main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw-eks-main.id
  }

  tags = {
    Name = "${local.env}-public-rt"
  }
}

resource "aws_route_table_association" "private-zone1" {
  subnet_id      = aws_subnet.private-1-zone1.id
  route_table_id = aws_route_table.private-rt.id
}

resource "aws_route_table_association" "private-zone2" {
  subnet_id      = aws_subnet.private-1-zone2.id
  route_table_id = aws_route_table.private-rt.id
}

resource "aws_route_table_association" "public-zone1" {
  subnet_id      = aws_subnet.public-1-zone1.id
  route_table_id = aws_route_table.public-rt.id
}

resource "aws_route_table_association" "public-zone2" {
  subnet_id      = aws_subnet.public-1-zone2.id
  route_table_id = aws_route_table.public-rt.id
  
}
