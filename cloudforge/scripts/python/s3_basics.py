from datetime import datetime
from typing import Any

import boto3
from botocore.exceptions import ClientError

# --- Configuration ---
ENDPOINT_URL = "http://localhost:4566"
REGION_NAME = "us-east-1"
BUCKET_NAME = "cloudforge-artifacts"

# --- Boto3 Client for S3 ---
s3_client: Any = boto3.client(
    "s3",
    endpoint_url=ENDPOINT_URL,
    region_name=REGION_NAME,
    aws_access_key_id="test",
    aws_secret_access_key="test",
)


def create_bucket_with_versioning(s3_client: Any, bucket_name: str) -> None:
    """Creates an S3 bucket if missing and enables versioning"""
    try:
        s3_client.create_bucket(Bucket=bucket_name)
        print(f"✓ Bucket '{bucket_name}' created successfully.")
    except ClientError as exc:
        code = exc.response.get("Error", {}).get("Code")
        if code not in ("BucketAlreadyOwnedByYou", "BucketAlreadyExists"):
            raise
        print(f"i Bucket '{bucket_name}' already exists.")
    except Exception as e:
        print(f"✗ An unexpected error occurred: {e}")

    s3_client.put_bucket_versioning(
        Bucket=bucket_name, VersioningConfiguration={"Status": "Enabled"}
    )
    print(f"✓ Versioning enabled on bucket '{bucket_name}'.")


def upload_file(s3_client: Any, bucket: str, key: str, body: bytes) -> None:
    """
    Upload a file to S3. Body must be bytes or file-like object.
    """
    print(f"Uploading file {key} to bucket {bucket}")
    try:
        s3_client.put_object(Bucket=bucket, Key=key, Body=body)
        s3_path = f"s3://{bucket}/{key}"
        print(f"✓ Successfully uploaded file to {s3_path}")
    except Exception as e:
        print(f"✗ File upload failed: {e}")


if __name__ == "__main__":
    print("\n--- Setting up S3 Resources in LocalStack ---")
    create_bucket_with_versioning(s3_client=s3_client, bucket_name=BUCKET_NAME)

    timestamp = datetime.now().isoformat()
    file_content = f"This is a test artifact uploaded at {timestamp}."
    file_key = f"{timestamp}.txt"
    body = file_content.encode("utf-8")
    upload_file(
        s3_client=s3_client,
        bucket=BUCKET_NAME,
        key=file_key,
        body=body,
    )
    print("\nS3 setup complete.")
