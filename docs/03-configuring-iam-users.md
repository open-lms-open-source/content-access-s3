# Configuring IAM Users

To provide you with access, the Open LMS Enterprise team will configure
read-only access to the account number you provided.  The administrator of your
account can then delegate that access to IAM roles or users within your
account.

The following IAM policy should be used to provide access:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "OpenLMSEnterpriseReadOnlyAccess",
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "arn:aws:iam::XXXXXXXXXXXX:role/XXXReadOnly"
        }
    ]
}
```

The full resource ARN will be provided by the Open LMS Enterprise team with the
confirmation that your access has been provisioned.

An example [Terraform](https://terraform.io/) configuration for creating an
IAM user with this policy attached can be found in
[examples/iam-user](../examples/iam-user/main.tf).

Next: [Configuring AWSCLI](04-configuring-awscli.md)
