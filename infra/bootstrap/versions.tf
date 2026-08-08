terraform {
  required_version = ">= 1.10"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.60"
    }
  }

  # Bootstrap uses LOCAL state on purpose: it CREATES the S3 bucket that the main
  # stack (../terraform) will use as its remote backend. Chicken-and-egg — so this
  # runs once, locally. State stays on your machine (and is gitignored).
}

provider "aws" {
  region = var.region
}
