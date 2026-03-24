terraform {
  required_version = ">= 1.0.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    key     = "lambda-user-authorizer.tfstate"
    region  = "us-east-1"
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_ecr_repository" "this" {
  name = "tech-challenge-user-authorizer-repo"
}

data "aws_iam_role" "lab_role" {
  name = "LabRole"
}

resource "aws_lambda_function" "this" {
  function_name = "tech-challenge-user-authorizer"
  role          = data.aws_iam_role.lab_role.arn
  package_type  = "Image"
  image_uri     = "${data.aws_ecr_repository.this.repository_url}:${var.image_tag}"

  reserved_concurrent_executions = 3

  timeout     = 30
  memory_size = 128

  environment {
    variables = {
      JWT_SECRET         = var.jwt_secret
      JWT_ISSUER         = var.jwt_issuer
      DYNAMODB_TABLE_NAME = var.dynamodb_table_name
      DYNAMODB_ENDPOINT  = var.dynamodb_endpoint
    }
  }

  image_config {
    command = ["bootstrap"]
  }
}

resource "aws_lambda_function_url" "this" {
  function_name      = aws_lambda_function.this.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "allow_public_url" {
  statement_id           = "AllowPublicAccess"
  action                 = "lambda:InvokeFunctionUrl"
  function_name          = aws_lambda_function.this.function_name
  principal              = "*"
  function_url_auth_type = "NONE"
}
