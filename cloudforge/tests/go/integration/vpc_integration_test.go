//go:build integration
// +build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/xbikashthakur/platform-engg/cloudforge/tests/go/helpers"
)

func TestVpcModuleIntegration(t *testing.T) {
	// This tells Go to run this test in parallel with others
	t.Parallel()

	// 1. Find the absolute path to the repository root. This is our reliable anchor.
	repoRoot, err := helpers.FindRepoRoot()
	if err != nil {
		t.Fatalf("Failed to find repository root: %v", err)
	}
	t.Logf("Repository root found at: %s", repoRoot)

	// 2. Get the Terraform directory path from the environment variable.
	// the path to be relative to the repository root.
	terraformDirRelative := os.Getenv("TERRAFORM_DIR")
	t.Logf("TERRAFORM_DIR (relative to repo): %s", terraformDirRelative)

	// 3. If the env var is not set, provide a sensible default.
	if terraformDirRelative == "" {
		terraformDirRelative = "cloudforge/terraform/environments/dev"
		t.Logf("TERRAFORM_DIR not set, using default: %s", terraformDirRelative)
	}

	// 4. Join the repo root with the relative path to get the final, correct absolute path.
	terraformDirAbsolute := filepath.Join(repoRoot, terraformDirRelative)
	t.Logf("Final absolute TerraformDir for Terratest: %s", terraformDirAbsolute)

	terraformOptions := &terraform.Options{
        TerraformDir:     terraformDirAbsolute,  // Path to the Terraform configuration
        TerraformBinary:  "terraform",   // Default is 'tofu' from OpenTofu, explicitly set to use Terraform CLI
        // Additional fields as needed, e.g., Vars for variables
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
