resource "aws_sqs_queue" "queue-1" {
  name                      = "queue-1"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  sqs_managed_sse_enabled   = true

  tags = {
    Name        = "Queue 1"
    Environment = "DEV"
  }
}

resource "aws_sqs_queue" "queue-2" {
  name                      = "queue-2"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
  sqs_managed_sse_enabled   = true

  tags = {
    Name        = "Queue 2"
    Environment = "DEV"
  }
}
