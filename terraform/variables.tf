variable "users_tb_name" {
  type = string
  default = "users"
}

variable "domain_sns_topic" {
  type = string
  default = "domain_sns_topic"
}

variable "notify_subscribers_on_task_created" {
  type = string
  default = "notify_subscribers_on_user_created"
  
}