variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "dynamodb_table_name" {
  description = "The name of the DynamoDB table for session management"
  type        = string
  default     = "user-authentication-token"
}

variable "image_tag" {
  description = "ECR image tag"
  type        = string
}
