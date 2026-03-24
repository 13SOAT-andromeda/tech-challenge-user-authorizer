variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "image_tag" {
  description = "ECR image tag"
  type        = string
}

variable "dynamodb_table_name" {
  description = "The name of the DynamoDB table for session management"
  type        = string
  default     = "user-authentication-token"
}

variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
}

variable "jwt_issuer" {
  description = "JWT issuer claim"
  type        = string
  sensitive   = true
}

variable "session_table_name" {
  description = "DynamoDB session table name"
  type        = string
}

variable "dynamodb_endpoint" {
  description = "DynamoDB endpoint override (e.g. for LocalStack)"
  type        = string
  default     = ""
}
