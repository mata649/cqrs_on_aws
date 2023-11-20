

resource "aws_ssm_parameter" "users_tb_name" {
  name        = "users_tb_name"
  description = "The name of the users table for DynamoDB"
  type        = "String"
  value       = aws_dynamodb_table.users_tb.name
}
