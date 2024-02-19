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

# resource "aws_sqs_queue" "fila_de_saida_de_erros" {
#   name                      = "fila-de-saida-de-erros"
#   delay_seconds             = 0
#   max_message_size          = 262144
#   message_retention_seconds = 86400
#   sqs_managed_sse_enabled   = true
#   receive_wait_time_seconds = 10
#   redrive_policy = jsonencode({
#     deadLetterTargetArn = aws_sqs_queue.fila_de_saida_de_erros_deadletter.arn
#     maxReceiveCount     = 4
#   })

#   tags = {
#     Name        = "Fila de Saída de Erros"
#     Environment = "DEV"
#   }
# }

# resource "aws_sqs_queue" "fila_de_saida_de_sucesso_deadletter" {
#   name                    = "fila-de-saida-de-sucesso-deadletter"
#   sqs_managed_sse_enabled = true

# }

# resource "aws_sqs_queue" "fila_de_saida_de_sucesso" {
#   name                      = "fila-de-saida-de-sucesso"
#   delay_seconds             = 0
#   max_message_size          = 262144
#   message_retention_seconds = 86400
#   receive_wait_time_seconds = 10
#   sqs_managed_sse_enabled   = true
#   redrive_policy = jsonencode({
#     deadLetterTargetArn = aws_sqs_queue.fila_de_saida_de_sucesso_deadletter.arn
#     maxReceiveCount     = 4
#   })

#   tags = {
#     Name        = "Fila de Saída de Sucesso"
#     Environment = "DEV"
#   }
# }
