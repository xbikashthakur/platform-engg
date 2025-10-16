import boto3
from datetime import datetime

# --- Configuration ---
ENDPOINT_URL = 'http://localhost:4566'
REGION_NAME = 'us-east-1'
BUCKET_NAME = 'cloudforge-artifacts'

# --- Boto3 Client for S3 ---
s3_client = boto3.client(
    's3',
    endpoint_url=ENDPOINT_URL,
    region_name=REGION_NAME,
    aws_access_key_id='test',
    aws_secret_access_key='test'
)

def create_bucket_with_versioning(bucket_name):
    """Creates an S3 bucket and enables versioning."""
    try:
        s3_client.create_bucket(Bucket=bucket_name)
        print(f"✓ Bucket '{bucket_name}' created successfully.")

        # Versioning is a best practice to protect against accidental deletions
        s3_client.put_bucket_versioning(
            Bucket=bucket_name,
            VersioningConfiguration={'Status': 'Enabled'}
        )
        print(f"✓ Versioning enabled on bucket '{bucket_name}'.")

    except s3_client.exceptions.BucketAlreadyOwnedByYou:
        print(f"i Bucket '{bucket_name}' already exists. Skipping creation.")
    except Exception as e:
        print(f"✗ An unexpected error occurred: {e}")

def upload_file(bucket_name):
    """Uploads a simple text file to the specified bucket."""
    timestamp = datetime.now().isoformat()
    file_content = f"This is a test artifact uploaded at {timestamp}."
    file_key = f"{timestamp}.txt"

    try:
        s3_client.put_object(
            Bucket=bucket_name,
            Key=file_key,
            Body=file_content.encode('utf-8') # Content must be bytes
        )
        print(f"✓ Successfully uploaded file to s3://{bucket_name}/{file_key}")
    except Exception as e:
        print(f"✗ File upload failed: {e}")


if __name__ == '__main__':
    print(f"\n--- Setting up S3 Resources in LocalStack ---")
    create_bucket_with_versioning(BUCKET_NAME)
    upload_file(BUCKET_NAME)
    print("\nS3 setup complete.")