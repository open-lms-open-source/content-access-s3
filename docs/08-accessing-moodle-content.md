# Access to Moodle Content

Moodle content is stored in your Moodle's S3 bucket in the same directory
structure that Moodle uses when storing content on local disk.  That is,
filenames are the SHA1 checksum of the file and the first four characters
are used as a two level directory structure.  For example, if the file has a
SHA1 checksum of `da39a3ee5e6b4b0d3255bfef95601890afd80709`, the file will
be stored in S3 at the following path:

```
da/39/da39a3ee5e6b4b0d3255bfef95601890afd80709
```

To map this to the original filename, you'll need to find the SHA1 hash in
the `mdl_files` table of the Moodle database.

Files may exist in the S3 bucket that don't exist in the `mdl_files` table,
as the files are not deleted from S3 (to facilitate restoring courses from
fileless backups).

The expected use cases for access to Moodle content are:

* obtaining individual files based on information found in database
  backups
* obtaining a copy of _all_ Moodle content for the purpose of keeping your
  own backup, setting up a copy of your Moodle on your own infrastructure,
  or migrating to another Moodle hosting provider

Given the SHA1 checksum, an individual file can be downloaded similarly to a
database backup:

```bash
$ aws --profile moodle s3 cp s3://example-s3-bucket/da/39/da39a3ee5e6b4b0d3255bfef95601890afd80709 .
```

To maintain a copy of all Moodle content that exists in your Moodle's S3
bucket, use the s3 sync command.  Run periodically, this will only copy new
Moodle content to your local copy:

```bash
$ aws --profile moodle s3 sync s3://example-s3-bucket/ . --exclude 'db-backups/*'
```

Note that this includes files that have been deleted from Moodle.

Next: [Fileless Course Backups](09-fileless-course-backups.md)
