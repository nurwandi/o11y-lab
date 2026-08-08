variable "region" {
  description = "AWS region"
  type        = string
  default     = "ap-southeast-3" # Jakarta
}

variable "github_repo" {
  description = "GitHub repo allowed to assume the CI role, as owner/name"
  type        = string
  default     = "nurwandi/o11y-lab"
}

variable "state_bucket_name" {
  description = "Globally-unique S3 bucket name that will hold the main stack's Terraform state"
  type        = string
}

variable "create_github_oidc_provider" {
  description = "Create the GitHub OIDC provider. Set to false if the account already has one (only one per account is allowed) and pass nothing else — the existing provider is looked up automatically."
  type        = bool
  default     = true
}

variable "tags" {
  description = "Tags applied to created resources"
  type        = map(string)
  default = {
    Project   = "o11y-lab"
    ManagedBy = "terraform"
  }
}
