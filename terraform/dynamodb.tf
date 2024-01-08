resource "aws_dynamodb_table" "users_tb" {
  name           = var.users_tb_name
  billing_mode   = "PROVISIONED"
  hash_key       = "id"
  read_capacity  = 1
  write_capacity = 1
  attribute {
    name = "id"
    type = "S"
  }
  attribute {
    name = "email"
    type = "S"
  }
  global_secondary_index {
    name = "email_index"
    hash_key = "email"
    write_capacity = 1
    read_capacity = 1
    projection_type = "ALL"
  }

}
