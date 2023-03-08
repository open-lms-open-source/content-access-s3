#
# IAM user for Open LMS Enterprise S3 Content Access
#

provider "aws" {
  region = "ap-southeast-2"
}

resource "aws_iam_user" "content-access" {
  name = "content-access"
}

resource "aws_iam_access_key" "content-access" {
  user = "${aws_iam_user.content-access.name}"
}

resource "aws_iam_user_policy_attachment" "content-access" {
  user       = "${aws_iam_user.content-access.name}"
  policy_arn = "${aws_iam_policy.content-access.arn}"
}

resource "aws_iam_policy" "content-access" {
  name        = "content-access"
  description = "Policy allowing read-only access to Moodle content"
  policy      = "${data.aws_iam_policy_document.content-access.json}"
}

data "aws_iam_policy_document" "content-access" {
  statement {
    sid = "OpenLMSEnterpriseReadOnlyAccess"

    actions = [
      "sts:AssumeRole",
    ]

    resources = [
      # the full resource ARN to use here will be provided by the Open LMS
      # Enterprise team with the confirmation that your access has been
      # provisioned
      "arn:aws:iam::XXXXXXXXXXXX:role/XXXReadOnly",
    ]
  }
}

output "aws_access_key_id" {
  value = "${aws_iam_access_key.content-access.id}"
}

output "aws_secret_access_key" {
  value = "${aws_iam_access_key.content-access.secret}"
}
