import boto3
from moto import mock_aws  # Changed from mock_iam for newer moto versions


# The @mock_aws decorator intercepts any boto3 calls
# made within this function and redirects them to the
# in-memory mock AWS environment instead of the real one.
@mock_aws
def test_iam_user_creation() -> None:
    """
    Tests that an IAM user can be created successfully.
    """
    # 1. Setup: Create a boto3 client. moto ensures this client is a mock.
    iam_client = boto3.client("iam", region_name="us-east-1")
    test_username = "test-cloudforge-user"

    # 2. Action: Call the function that creates the user.
    iam_client.create_user(UserName=test_username)

    # 3. Assertion: Check if the result is what we expect.
    response = iam_client.list_users()
    users = response["Users"]

    # We expect exactly one user to have been created.
    assert len(users) == 1
    # We expect the created user's name to match our test username.
    assert users[0]["UserName"] == test_username


@mock_aws
def test_iam_policy_creation() -> None:
    """
    Tests that a simple IAM policy can be created.
    """
    # 1. Setup
    iam_client = boto3.client("iam", region_name="us-east-1")
    test_policy_name = "TestPolicy"
    policy_document = (
        '{"Version": "2012-10-17", '
        '"Statement": [{"Effect": "Allow", "Action": "*", "Resource": "*"}]}'
    )

    # 2. Action
    iam_client.create_policy(
        PolicyName=test_policy_name, PolicyDocument=policy_document
    )

    # 3. Assertion
    response = iam_client.list_policies(Scope="Local")
    policies = response["Policies"]
    assert len(policies) == 1
    assert policies[0]["PolicyName"] == test_policy_name
