resource "aws_sqs_queue" "users_queue" {
  name = var.users_created_queue
  receive_wait_time_seconds = 20
  message_retention_seconds =18400 
}


