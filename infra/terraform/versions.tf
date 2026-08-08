terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.60"
    }
  }

  # Remote state in S3 with native lockfile (no DynamoDB). The bucket name is
  # provided at init time so it isn't hard-coded:
  #   terraform init -backend-config="bucket=<from bootstrap output>"
  backend "s3" {
    key          = "env/terraform.tfstate"
    region       = "ap-southeast-3"
    encrypt      = true
    use_lockfile = true
  }
}

provider "aws" {
  region = var.region

  default_tags {
    tags = var.tags
  }
}
