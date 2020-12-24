resource "aws_iam_role" "oops_ec2_iam_role" {
  name               = "prod-oops-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
  tags = {
    Name         = "prod-oops-role"
    Environment  = "prod"
    ManagingTeam = "SysEng"
    CostCenter   = "SysEng"
  }
}

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/api-permissions-reference.html
resource "aws_iam_policy" "oops_iam_policy" {
  name   = "prod-oops-policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
          "dynamodb:DeleteItem",
          "dynamodb:GetItem",
          "dynamodb:PutItem"
      ],
      "Resource": "${aws_dynamodb_table.oops_table.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "oops_attach" {
  role   = aws_iam_role.oops_ec2_iam_role.arn
  policy = aws_iam_policy.oops_iam_policy.arn
}