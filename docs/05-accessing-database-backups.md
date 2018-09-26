# Access to Database Backups

> All commands shown on this page assume you've configured the AWS CLI with
> a profile called `moodle` and that profile has access to your Moodle's S3
> bucket.

Backups of your databases are uploaded to your Moodle's S3 bucket and made
available to you under the /db-backups/ path.  Retention of database backups
is at least:

* one backup per day is available for two months
* one backup per month is available for seven years

The name of your Moodle's S3 bucket is provided by the Blackboard Open LMS
Enterprise team with the confirmation that your access has been provisioned.
You'll need to substitute that in place of `example-s3-bucket` in the
examples shown below.

To obtain a list of database backups available in your Moodle's S3 bucket:

```bash
$ aws --profile moodle s3 ls s3://example-s3-bucket/db-backups/
```

If there's a backup in that list named
`example_mdl_prod-daily_dump-20180919.directory_format.tar`, it can be
downloaded to the current working directory with the following command:

```bash
$ aws --profile moodle s3 cp s3://example-s3-bucket/db-backups/example_mdl_prod-daily_dump-20180919.directory_format.tar .
```

It's also possible to maintain a local copy of all database backups using
the s3 sync command.  Run periodically, this will only copy new database
backups to your local copy (in a directory named `db-backups`):

```bash
$ aws --profile moodle s3 sync s3://example-s3-bucket/db-backups/ ./db-backups/
```

Note that the access provided to your Moodle's S3 bucket is read-only and
this cannot be changed.

Next: [Plain Format Database Backups](06-plain-format-backups.md)
