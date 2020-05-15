# Configuring the AWS CLI

AWS provides excellent documentation on configuring the AWS CLI, which can
be found
[here](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

To configure the AWS CLI for access to your Moodle's S3 bucket, you'll
require the following:

* the access key ID and secret access key for your IAM user; these can be
  generated from the AWS IAM Console as detailed in the documentation linked
  above
* the ARN string for your assigned MFA device, if using multi-factor
  authentication; this can be found in the AWS IAM Console on the same page
  that access key IDs are managed
* the ARN string for the role you'll use to access your Moodle's S3 bucket;
  this is provided to you by the Open LMS Enterprise team with the confirmation
  that your access has been provisioned
* the AWS CLI tool must be installed on the computer that will be
  downloading files from the S3 bucket; documentation can be found
  [here](https://docs.aws.amazon.com/cli/latest/userguide/installing.html)

Once you have those, a profile can be created for accessing your Moodle's S3
bucket with the following commands (where `$` represents the command prompt
and `X`s represent required information as detailed above):

```bash
$ aws configure set aws_access_key_id XXXXXXXXXXXXXXXXXXXX --profile moodle
$ aws configure set aws_secret_access_key XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --profile moodle
$ aws configure set region ap-southeast-2 --profile moodle
$ aws configure set role_arn arn:aws:iam::XXXXXXXXXXXX:role/XXXXXXXXXXX --profile moodle
$ aws configure set source_profile moodle --profile moodle
$ aws configure set mfa_serial arn:aws:iam::XXXXXXXXXXXX:mfa/XXXXXXXX --profile moodle
```

The final command shown above is optional and should only be executed if
you're using multi-factor authentication.

Next: [Access to Database Backups](05-accessing-database-backups.md)
