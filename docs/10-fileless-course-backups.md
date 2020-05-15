# Fileless Course Backups

Every night, Open LMS Enterprise performs automated backups of every course
that's changed since the last backup.  To allow backups to be efficiently
created and stored at scale, a customisation has been made to Moodle that
excludes files from the backup.

When that course backup is restored into Open LMS Enterprise, the files already
exist in S3 and the course restore is successful.  If the course backups are
required outside of Open LMS Enterprise -- if you wish to import them to
another Moodle hosting provider for example -- they either need to be
regenerated with files included or the fileless backups need to have files
added to them.

Example code can be found [here](../examples/moodle-backup-filler/) to add
files to fileless backups.  See the
[README.md](../examples/moodle-backup-filler/README.md) included with that
software for further information on using it.

Next: [Using S3 Event Notifications](11-s3-event-notifications.md)
