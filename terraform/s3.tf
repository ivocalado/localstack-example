resource "aws_s3_bucket" "bucket_1" {
  bucket = "bucket-1"
  tags = {
    Name        = "Bucket 1"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket" "bucket_2" {
  bucket = "bucket-2"
  tags = {
    Name        = "Bucket 2"
    Environment = "Dev"
  }
}
