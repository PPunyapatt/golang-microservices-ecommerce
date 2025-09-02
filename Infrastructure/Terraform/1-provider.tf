provider "aws" {
    region = local.region
}

terraform {
  required_version = ">=1.13.1"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.10.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.16.0" 
    }
  }
}