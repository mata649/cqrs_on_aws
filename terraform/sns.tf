resource "aws_sns_topic" "domain_sns_topic" {
  name = var.domain_sns_topic
}

resource "aws_sns_topic_subscription" "notify_subscribers_on_task_created_subscription" {
  topic_arn = aws_sns_topic.domain_sns_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.notify_subscribers_on_task_created.arn
  filter_policy = jsonencode({
    "event_type" = ["events.task.created"]
  })

}
