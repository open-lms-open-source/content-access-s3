# Directory Format Database Backups

Open LMS Enterprise databases have grown significantly over the years and the
amount of time required to restore a plain format backup has become
problematic.  The Open LMS Enterprise team regularly restores production
backups to the staging environment, and several clients download their database
backups and restore them every day.  Both of these can benefit from faster
restore times than plain format backups can provide.

To address this, "directory format" backups were introduced.  These are a
directory containing one file for each table and blob in the database plus a
table of contents describing the objects in a machine-readable format that
`pg_restore` can read.  To make the backup easy to manage, that directory is
archived using `tar` and that archive is given a filename ending in
`.directory_format.tar`.

The full filename of a directory format database backup has the following
parts, separated by hyphens:

* the database name, such as `example_mdl_prod`
* the string `daily_dump`
* the date the backup run was started in the format `YYYYMMDD`
* an optional time the backup run was started in the format `HHMM`
* an optional sequence number for the backup (an alternative to the time
  for the second and subsequent backups of a given database in one day)
* the `.directory_format.tar` suffix

Directory format database backups can be restored to an empty PostgreSQL
database by extracting the archive and then passing that to `pg_restore`.
For example:

```bash
$ tar xf example_mdl_prod-daily_dump-20180919.directory_format.tar
$ pg_restore \
    --username=postgres \
    --schema=public \
    --format=directory \
    --jobs=2 \
    --clean \
    --if-exists \
    --no-owner \
    --no-acl \
    --dbname=example_mdl_prod \
    --role=example_mdl_prod \
    example_mdl_prod-daily_dump-20180919
```

Depending on the configuration and location of the PostgreSQL server, you
may need to pass additional or different options to `pg_restore`.

Note that you'll need sufficient disk capacity available in the working
directory to extract the archive.

Restores from a directory format backup have several advantages over
restores from plain format backups:

* parallel jobs can be enabled with `-j `_number-of-jobs_ to run the most
  time-consuming parts of a restore using multiple concurrent jobs.  This
  can dramatically reduce the time to restore the database to a server
  running on a multiprocessor machine
* restore only specific tables with `-t `_table_ to save time restoring
  tables that aren't required

Next: [Access to Log Files](08-accessing-logs.md)
