resource "aws_iam_role" "eks" {
  name               = "${local.env}-${local.eks_name}-cluster"
  assume_role_policy = <<POLICY
    {
        "Version": "2012-10-17",
        "Statement": [
            {
            "Effect": "Allow",
            "Principal": {
                "Service": "eks.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
            }
        ]
    }
    POLICY
  
}

resource "aws_iam_role_policy_attachment" "eks" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role = aws_iam_role.eks.name
}

resource "aws_eks_cluster" "eks" {
  name     = local.eks_name
  role_arn = aws_iam_role.eks.arn
  version  = local.eks_version

  vpc_config {
    endpoint_private_access = false
    endpoint_public_access = true

    subnet_ids = [
        aws_subnet.private-1-zone1.id, 
        aws_subnet.private-1-zone2.id
    ]
  }

  access_config {
    authentication_mode =  "API"
    bootstrap_cluster_creator_admin_permissions = true
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks
  ]
}
