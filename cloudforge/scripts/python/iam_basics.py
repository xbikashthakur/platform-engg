import json
from typing import Any, Mapping

import boto3
from botocore.exceptions import ClientError

# --- Configuration ---
# In a real app, this would come from a config file
ENDPOINT_URL = "http://localhost:4566"
REGION_NAME = "us-east-1"
POLICY_NAME = "CloudForgeS3AccessPolicyV4"
USER_NAME = "cloudforge-user-7"

# Policy document defines the permissions in JSON format
policy_document = {
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
            "Resource": "arn:aws:s3:::cloudforge-*/*",
        }
    ],
}

# --- Boto3 Client ---
iam_client: Any = boto3.client(
    "iam",
    endpoint_url=ENDPOINT_URL,
    region_name=REGION_NAME,
    aws_access_key_id="test",  # LocalStack credentials
    aws_secret_access_key="test",  # LocalStack credentials
)


def create_user(iam_client: Any, username: str) -> str:
    """
    Create an IAM user if missing. Return the user's ARN.
    Idempotent: returns existing ARN when the user exists.
    """
    try:
        resp: Mapping[str, Any] = iam_client.create_user(UserName=username)
        print(f"✓ User '{username}' created successfully.")
        return str(resp["User"]["Arn"])
    except ClientError as exc:
        code = exc.response.get("Error", {}).get("Code")
        if code in ("EntityAlreadyExists", "EntityAlreadyExistsException"):
            resp = iam_client.get_user(UserName=username)
            print(f"i User '{username}' already exists.")
            return str(resp["User"]["Arn"])

        print(f"i Error while creating user: '{username}'")
        raise


def create_policy(
    iam_client: Any, policy_name: str, policy_document: dict[str, Any]
) -> str:
    """
    Create a managed policy or return existing ARN.
    """
    policy_arn: str = f"arn:aws:iam::000000000000:policy/{policy_name}"
    try:
        iam_client.get_policy(PolicyArn=policy_arn)
        print(f"i Policy '{policy_name}' already exists.")
        return policy_arn
    except ClientError as exc:
        code = exc.response.get("Error", {}).get("Code")
        if code in ("NoSuchEntity", "NoSuchEntityException"):
            policy_doc_json = json.dumps(policy_document)
            resp = iam_client.create_policy(
                PolicyName=policy_name, PolicyDocument=policy_doc_json
            )
            print(f"✓ Policy '{policy_name}' created successfully.")
            return str(resp["Policy"]["Arn"])

        print(f"i Error while creating policy: '{policy_name}'")
        raise


if __name__ == "__main__":
    print("--- IAM Resources setup in LocalStack ---")
    create_user(iam_client=iam_client, username=USER_NAME)
    create_policy(
        iam_client=iam_client,
        policy_name=POLICY_NAME,
        policy_document=policy_document,
    )
    print("\n IAM setup completed.")
