#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status.

echo "=== 🚀 Phase 0 Verification Script 🚀 ==="

# 1. Check for core tools
echo -e "\n[1/7] Verifying core tool installations..."
go version
python3 --version
terraform --version
docker --version
aws --version
pre-commit --version
echo "✓ All tools are installed."

# 2. Check LocalStack container
echo -e "\n[2/7] Verifying LocalStack container..."
if [ ! "$(docker ps -q -f name=cloudforge-localstack)" ]; then
    echo "✗ LocalStack container is not running. Please run 'docker-compose up -d'."
    exit 1
fi
echo "✓ LocalStack container is running."

# 3. Check CloudForge CLI build
echo -e "\n[3/7] Building and testing CloudForge CLI..."
go build -o /tmp/cloudforge cmd/cloudforge/main.go
/tmp/cloudforge config.yaml
echo "✓ CloudForge CLI built and ran successfully."

# 4. Run pre-commit hooks
echo -e "\n[4/7] Running pre-commit checks..."
pre-commit run --all-files
echo "✓ Pre-commit hooks passed."

# 5. Validate Terraform
echo -e "\n[5/7] Validating Terraform configuration..."
terraform -chdir=terraform/environments/dev init -backend=false -reconfigure > /dev/null
terraform -chdir=terraform/environments/dev validate
echo "✓ Terraform configuration is valid."

# 6. Run all tests (Go and Python)
echo -e "\n[6/7] Running Go and Python tests..."
go test -v -timeout 30s ./... # Short timeout for unit tests
pytest tests/python/
echo "✓ All unit tests passed."

# 7. Run Terratest (Optional, as it's slow)
read -p $'\n[7/7] Run slow Terratest integration test? (y/N): ' choice
case "$choice" in
  y|Y )
    echo "Running Terratest... (this may take a few minutes)"
    go test -v -timeout 15m ./tests/go/...
    echo "✓ Terratest integration test passed."
    ;;
  * )
    echo "i Skipped Terratest."
    ;;
esac

echo -e "\n=== 🎉 Congratulations! Phase 0 Verification Complete! 🎉 ==="