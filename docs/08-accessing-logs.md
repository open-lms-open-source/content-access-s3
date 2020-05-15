# Access to Log Files

> All commands shown on this page assume you've configured the AWS CLI with a
> profile called `moodle` and that profile has access to your Moodle's S3
> bucket.

Your Moodle's access logs, collected from the reverse proxy cluster and
aggregated into a single file, are uploaded to your Moodle's S3 bucket and made
available to you under the /logs/rp-access/ path.  Similarly, logs from
Moodle's "admin cron" are uploaded to the same S3 bucket and made available to
you under the /logs/cron/ path.  All logs published to the S3 bucket are
retained for at least 12 months.

The name of your Moodle's S3 bucket is provided by the Open LMS Enterprise team
with the configuration that your access has been provisioned.  You'll need to
substitute that in place of `example-s3-bucket` in the examples shown below.

To obtain a list of Moodle instance names for which access log files are
available in your Moodle's S3 bucket:

```bash
$ aws --profile moodle s3 ls s3://example-s3-bucket/logs/rp-access/
```

To obtain a list of access log files for a Moodle instance named `example-mdl-prod`:

```bash
$ aws --profile moodle s3 ls s3://example-s3-bucket/logs/rp-access/example-mdl-prod/
```

If there's an access log file in that list named `rp-access-2020-05-12.log.gz`,
it can be downloaded to the current working directory with the following
command:

```bash
$ aws --profile moodle s3 cp s3://example-s3-bucket/logs/rp-access/example-mdl-prod/rp-access-2020-50-12.log.gz .
```

It's also possible to maintain a local copy of all access log files using the
s3 sync command.  Run periodically, this will only copy new access log files to
your local copy (in a directory named `access-logs`):

```bash
$ aws --profile moodle s3 sync s3://example-s3-bucket/logs/rp-access/ ./access-logs/
```

Access to cron logs is achieved in the same way, except rather than
`/logs/rp-access/` you would use `/logs/cron/`.

Note that the access provided to your Moodle's S3 bucket is read-only and this
cannot be changed.

Next: [Access to Moodle Content](09-accessing-moodle-content.md)
