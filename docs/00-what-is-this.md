# About the Open LMS Enterprise S3 Content Access Service

The S3 Content Access Service provides Open LMS Enterprise customers with
secure access to their Moodle instances' database backups, various log files,
and Moodle file store by downloading them directly from the AWS S3 bucket in
which they're stored.

This service is intended to replace existing methods of access to database
backups and Moodle content.  They will, however, run in parallel until
Open LMS Enterprise instances are migrated to AWS.

### Features

* Authentication through Amazon Identity and Access Management (IAM)
* Data is stored in a highly scalable, highly available, fast data storage
  service, Amazon S3
* Data is encrypted at rest
* Option to integrate with other AWS services to take action when new files
  are available
* New database backups available daily, usually by 8am (Australian central time)
* Log files available within 1 - 2 days
* Content from the Moodle file store available in real time

Next: [Prerequisites](01-prerequisites.md)
