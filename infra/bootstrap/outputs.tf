output "ci_role_arn" {
  description = "ARN of the CI role — set this as the AWS_CI_ROLE_ARN repo variable for GitHub Actions"
  value       = aws_iam_role.ci.arn
}

output "state_bucket" {
  description = "S3 bucket name for the main stack's backend — pass to `terraform init -backend-config`"
  value       = aws_s3_bucket.state.id
}

output "oidc_provider_arn" {
  description = "GitHub OIDC provider ARN in use"
  value       = local.oidc_provider_arn
}
