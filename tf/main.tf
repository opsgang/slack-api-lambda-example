# terraform

provider "aws" {
    region = "eu-west-1"
}

variable "channel_id" {
    description = "see api channels.list method to get ids"
}

variable "icon_url" {
    default     = "http://blog.edtechie.net/wp-content/uploads/2015/07/kingofcomedy.jpg"
    description = "The icon/image to display when winning results are posted back to slack"
}

variable "results_posted_by" {
    default     = "Rupert Pupkin Speaks!"
    description = "The name to display when the winning results are posted back to slack"
}

resource "aws_iam_role" "iam_for_example_lambda" {
  name = "iam_for_example_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "pupkin" {
  filename         = "pupkin.zip"
  function_name    = "pupkin"
  role             = "${aws_iam_role.iam_for_example_lambda.arn}"
  handler          = "pupkin"
  source_code_hash = "${base64sha256(file("pupkin.zip"))}"
  runtime          = "go1.x"

  # API_KEY for slack must also be set, but this is a secret ...
  environment {
    variables = {
      CHANNEL_ID        = "${var.channel_id}"
      ICON_URL          = "${var.icon_url}"
      RESULTS_POSTED_BY = "${var.results_posted_by}"
    }
  }
}

resource "aws_cloudwatch_event_rule" "ten_fifteen_am_BST" {
    name                = "ten-fifteen-am"
    description         = "Fires at 10:15 daily"
    schedule_expression = "cron(15 9 * * ? *)" # ... should account for timezone
}

resource "aws_cloudwatch_event_target" "invoke_pupkin_lambda_on_time" {
    rule      = "${aws_cloudwatch_event_rule.ten_fifteen_am_BST.name}"
    target_id = "pupkin"
    arn       = "${aws_lambda_function.pupkin.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_lambda" {
    statement_id  = "AllowExecutionFromCloudWatch"
    action        = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.pupkin.function_name}"
    principal     = "events.amazonaws.com"
    source_arn    = "${aws_cloudwatch_event_rule.ten_fifteen_am_BST.arn}"
}
