package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestVpcModuleIntegration(t *testing.T) {
	// This tells Go to run this test in parallel with others
	t.Parallel()

	// Define the Terraform options.
	// We are testing our 'dev' ENVIRONMENT, not the module directly,
	// because the environment has the LocalStack provider configured.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../terraform/environments/dev",

		// We don't need to pass variables since our 'dev' environment has them defined.
		// But if we were testing the module directly, we would do it here:
		// Vars: map[string]interface{}{
		//  "project_name": "terratest",
		//  "environment":  "test",
		// },
	}

	// --- Test Lifecycle ---
	// At the end of the test, run 'terraform destroy' to clean up any resources.
	// 'defer' ensures this runs even if the test fails.
	defer terraform.Destroy(t, terraformOptions)

	// Run 'terraform init' and 'terraform apply'.
	// Terratest will fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// --- Assertions ---
	// Run 'terraform output' to get the value of an output variable.
	vpcId := terraform.Output(t, terraformOptions, "dev_vpc_id")

	// Assert that the output is not empty.
	// The 'assert' library provides helpful functions for writing test checks.
	assert.NotEmpty(t, vpcId, "VPC ID should not be empty")
	assert.Contains(t, vpcId, "vpc-", "VPC ID should have the 'vpc-' prefix")
}