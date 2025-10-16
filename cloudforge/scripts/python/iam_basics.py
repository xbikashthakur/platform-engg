import boto3
import json


# --- Configuration ---
# In a real app, this would come from a config file
ENDPOINT_URL = 'http://localhost:4566'
REGION_NAME = 'us-east-1'
POLICY_NAME = 'CloudForgeS3AccessPolicyV2'
USER_NAME = 'cloudforge-user-5'


# --- Boto3 Client ---
iam_client = boto3.client(
    'iam',
    endpoint_url=ENDPOINT_URL,
    region_name=REGION_NAME,
    aws_access_key_id='test', # LocalStack credentials
    aws_secret_access_key='test' # LocalStack credentials
)

def create_user(username):
    """Creates an IAM user if they don't already exist."""
    try:
        iam_client.create_user(UserName=username)
        print(f"User '{username}' created successfuly.")
    except iam_client.exceptions.EntityAlreadyExistException:
        print(f"i User '{username}' already exists. Skipping user creation.")

def create_policy(policy_name):
    """Creates a policy to allow read/write access to specific S3 buckets."""
    # Policy document defines the permissions in JSON format
    policy_document = {
        "Version": "2012-10-17",
        "Statement": [{
            "Effect": "Allow",
            "Action": ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
            "Resource": "arn:aws:s3:::cloudforge-*/*"
        }]
    }
    try:
        iam_client.create_policy(
            PolicyName=policy_name,
            PolicyDocument=json.dumps(policy_document)
        )
        print(f"Policy '{policy_name}' created successfuly.")
    except iam_client.exceptions.EntityAlreadyExistException:
        print(f"i Policy '{policy_document}' already exists. Skipping policy creation.")



if __name__=='__main__':
    print("--- Setting up IAM Resources in LocalStack ---")
    create_user(USER_NAME)
    create_policy(POLICY_NAME)
    print("\n IAM setup completed.")

