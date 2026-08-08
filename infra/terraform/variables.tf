variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-3" # Jakarta
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidr" {
  description = "CIDR block for the public subnet"
  type        = string
  default     = "10.0.1.0/24"
}

variable "instance_type" {
  description = "EC2 instance type (k3s + LGTM + app fit comfortably in 8GB)"
  type        = string
  default     = "t3.large"
}

variable "root_volume_gb" {
  description = "Root EBS volume size (GiB)"
  type        = number
  default     = 30
}

variable "operator_cidr" {
  description = "Your public IP as a CIDR (e.g. 203.0.113.4/32). Restricts SSH and the k3s API to you."
  type        = string
}

variable "key_name" {
  description = "Name of an existing EC2 key pair for SSH access"
  type        = string
}

variable "gitops_repo_url" {
  description = "Git repo ArgoCD pulls manifests from"
  type        = string
  default     = "https://github.com/nurwandi/o11y-lab.git"
}

variable "tags" {
  description = "Tags applied to all resources"
  type        = map(string)
  default = {
    Project   = "o11y-lab"
    ManagedBy = "terraform"
  }
}
