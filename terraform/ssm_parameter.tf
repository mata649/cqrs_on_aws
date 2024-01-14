

resource "aws_ssm_parameter" "users_tb_name" {
  name        = "users_tb_name"
  description = "The name of the users table for DynamoDB"
  type        = "String"
  value       = aws_dynamodb_table.users_tb.name
}

resource "aws_ssm_parameter" "tasks_tb_name" {
  name        = "tasks_tb_name"
  description = "The name of the tasks table for DynamoDB"
  type        = "String"
  value       = aws_dynamodb_table.tasks_tb.name
}
resource "aws_ssm_parameter" "exchange" {
  name        = "exchange"
  description = "The name of the topic of sns"
  type        = "String"
  value       = aws_sns_topic.domain_sns_topic.arn
}

resource "aws_ssm_parameter" "task_created_queue" {
  name        = "task_created_queue"
  description = "The name of the  task created queue of sqs"
  type        = "String"
  value       = aws_sqs_queue.notify_subscribers_on_task_created.arn
}
