
# --- Provider Requirements ---
# This block specifies the Terraform version and the required providers.
terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# --- Input Variables ---
# These are the parameters for modules, making it reusable.
variable "project_name" {
  description = "The name of the project for tagging."
  type        = string
}

variable "environment" {
  description = "The deployment environment (eg. dev, staging, prod)."
  type        = string
}

variable "vpc_cidr" {
  description = "The CIDR block for the VPC."
  type        = string
  default     = "10.0.0.0/16"
}

# Data sources to get the current account ID and region for the KMS policy
data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# --- Resources ---
# These blocks define the actual infrastructure components to create.
resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr

  # Tags are crucial for identifying and managing resources
  tags = {
    Name        = "${var.project_name}-${var.environment}-vpc"
    Project     = var.project_name
    Environment = var.environment
  }
}

# Explicitly manage the default security group to restrict all traffic.
# By default, it allows all outbound traffic. This removes that rule.
resource "aws_default_security_group" "default" {
  vpc_id = aws_vpc.main.id

  # Empty ingress and egress blocks to remove the default rules.
  # This enforces a "deny-all" policy by default.
  ingress = []
  egress  = []

  tags = {
    Name = "${var.project_name}-${var.environment}-default-sg"
  }
}

# Create a KMS key to encrypt the CloudWatch Log Group
resource "aws_kms_key" "flow_logs_key" {
  description             = "KMS key for VPC flow logs"
  deletion_window_in_days = 7
  enable_key_rotation     = true

  # Define a policy based on the principle of least privilege
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      # Statement 1: Allows IAM administrators to manage the key
      {
        Sid    = "Enable IAM User Permissions",
        Effect = "Allow",
        Principal = {
          AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
        },
        Action   = "kms:*",
        Resource = "*"
      },
      # Statement 2: Allows CloudWatch Logs to use the key for encryption
      {
        Sid    = "Allow Cloudwatch Logs",
        Effect = "Allow",
        Principal = {
          Service = "logs.${data.aws_region.current.name}.amazonaws.com"
        },
        Action = [
          "kms:Encrypt*",
          "kms:Decrypt*",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:Describe*"
        ],
        Resource = "*"
      }
    ]
  })

  tags = {
    Name = "${var.project_name}-${var.environment}-vpc-flow-logs-key"
  }
}

# --- VPC Flow Logging ---
# 1. Create a CloudWatch Log Group to store flow logs
# Create a CloudWatch Log Group to store flow logs
resource "aws_cloudwatch_log_group" "flow_logs" {
  name              = "/aws/vpc-flow-logs/${var.project_name}-${var.environment}"
  retention_in_days = 365
  kms_key_id        = aws_kms_key.flow_logs_key.arn

  tags = {
    Name = "${var.project_name}-${var.environment}-vpc-flow-logs"
  }
}
# 2. Create an IAM Role that the Flow Logs service can assume
resource "aws_iam_role" "vpc_flow_log_role" {
  name = "${var.project_name}-${var.environment}-vpc-flow-log-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "vpc-flow-logs.amazonaws.com"
        }
      },
    ]
  })
}

# 3. Create a policy for the IAM role to allow writing to CloudWatch
resource "aws_iam_role_policy" "vpc_flow_log_policy" {
  name = "${var.project_name}-${var.environment}-vpc-flow-log-policy"
  role = aws_iam_role.vpc_flow_log_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogGroups",
          "logs:DescribeLogStreams"
        ]
        Effect   = "Allow"
        Resource = aws_cloudwatch_log_group.flow_logs.arn
      },
    ]
  })
}

# 4. Enable VPC Flow Logging for the VPC
resource "aws_flow_log" "main" {
  iam_role_arn    = aws_iam_role.vpc_flow_log_role.arn
  log_destination = aws_cloudwatch_log_group.flow_logs.arn
  traffic_type    = "ALL"
  vpc_id          = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-${var.environment}-vpc-flow-log"
  }
}

# Create a public subnet within the VPC
resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id
  # This function calculates a valid /24 subnet from the VPC's /16 range.
  # cidrsubnet(prefix, newbits, netnum)
  # It takes our VPC CIDR, adds 8 bits to the netmask (16+8=24), and creates the first subnet (netnum=0).
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, 0)
  availability_zone = "us-east-1a"
  # Set to false to avoid automatically assigning public IPs to instances.
  # This is a security best practice.
  map_public_ip_on_launch = false

  tags = {
    Name = "${var.project_name}-${var.environment}-public-subnet"
  }
}

# --- Outputs ---
# These are the return values of our module.
output "vpc_id" {
  description = "The ID of the created VPC."
  value       = aws_vpc.main.id
}

output "public_subnet_id" {
  description = "The ID of the public subnet."
  value       = aws_subnet.public.id
}
