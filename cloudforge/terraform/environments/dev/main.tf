terraform {
  required_version = ">= 1.5"
}

# --- Provider Configuration for LocalStack ---
# This block tells Terraform HOW and WHERE to build the infrastructure.
# We are pointing it directly at our LocalStack container.
provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  # Redirecting AWS API calls to localhost.
  endpoints {
    s3   = "http://localhost:4566"
    ec2  = "http://localhost:4566"
    iam  = "http://localhost:4566"
    sts  = "http://localhost:4566"
    kms  = "http://localhost:4566"
    logs = "http://localhost:4566"
    # Add endpoints for other services used in the module
  }
}

# --- Module Usage ---
# Here we call our VPC module, passing in values for the dev environment.
module "vpc" {
  source = "../../modules/vpc" # Path to the module code

  project_name = "cloudforge"
  environment  = "dev"
  vpc_cidr     = "10.10.0.0/16" # A specific CIDR for dev
}

# --- Root Outputs ---
# Expose the module's output at the environment level.
output "dev_vpc_id" {
  value = module.vpc.vpc_id
}
