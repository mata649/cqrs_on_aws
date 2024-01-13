resource "aws_sqs_queue" "notify_subscribers_on_task_created" {
  name = var.notify_subscribers_on_task_created
  receive_wait_time_seconds = 20
  message_retention_seconds =18400 
}

resource "aws_sqs_queue_policy" "notify_subscribers_on_task_created_subscription" {
  queue_url = aws_sqs_queue.notify_subscribers_on_task_created.id
  policy    = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "sns.amazonaws.com"
      },
      "Action": [
        "sqs:SendMessage"
      ],
      "Resource": [
        "${aws_sqs_queue.notify_subscribers_on_task_created.arn}"
      ],
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sns_topic.domain_sns_topic.arn}"
        }
      }
    }
  ]
}
EOF
}
