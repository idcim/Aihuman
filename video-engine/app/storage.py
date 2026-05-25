import os
import json
import time
from pathlib import Path

import boto3
from botocore.client import Config
from botocore.exceptions import ClientError


class ObjectStorage:
    def __init__(self) -> None:
        self.endpoint = os.getenv("OSS_ENDPOINT", "http://minio:9000")
        self.public_url = os.getenv("OSS_PUBLIC_URL", "http://localhost:19000").rstrip("/")
        self.bucket = os.getenv("OSS_BUCKET", "aihuman-assets")
        self.client = boto3.client(
            "s3",
            endpoint_url=self.endpoint,
            aws_access_key_id=os.getenv("OSS_ACCESS_KEY", "aihuman"),
            aws_secret_access_key=os.getenv("OSS_SECRET_KEY", "aihuman_dev_password"),
            config=Config(signature_version="s3v4"),
            region_name="us-east-1",
        )

    def ensure_bucket(self) -> None:
        for _ in range(30):
            try:
                self.client.head_bucket(Bucket=self.bucket)
                self.ensure_public_read()
                return
            except ClientError as error:
                code = error.response.get("Error", {}).get("Code")
                if code in {"404", "NoSuchBucket", "NotFound"}:
                    self.client.create_bucket(Bucket=self.bucket)
                    self.ensure_public_read()
                    return
            except Exception:
                time.sleep(1)
        self.client.create_bucket(Bucket=self.bucket)
        self.ensure_public_read()

    def ensure_public_read(self) -> None:
        policy = {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {"AWS": ["*"]},
                    "Action": ["s3:GetObject"],
                    "Resource": [f"arn:aws:s3:::{self.bucket}/*"],
                }
            ],
        }
        self.client.put_bucket_policy(Bucket=self.bucket, Policy=json.dumps(policy))

    def upload_file(self, path: Path, object_key: str) -> str:
        self.ensure_bucket()
        self.client.upload_file(
            str(path),
            self.bucket,
            object_key,
            ExtraArgs={"ContentType": "video/mp4"},
        )
        return f"{self.public_url}/{self.bucket}/{object_key}"
