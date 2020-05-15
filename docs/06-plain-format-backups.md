# Plain Format Database Backups

The format originally used for Open LMS Enterprise database backups was a gzip
compressed text file containing SQL statements to recreate the database and its
contents.  These are referred to as "plain format" backups and have filenames
ending in `.sql.gz`.

The full filename of a plain format database backup has the following parts,
separated by hyphens:

* the database name, such as `example_mdl_prod`
* the string `daily_dump`
* the date the backup run was started in the format `YYYYMMDD`
* an optional time the backup run was started in the format `HHMM`
* an optional sequence number for the backup (an alternative to the time
  for the second and subsequent backups of a given database in one day)
* the `.sql.gz` suffix

Plain format database backups can be restored to an empty PostgreSQL
database by uncompressing them and using psql to run the SQL contained
inside them.  For example:

```bash
$ zcat example_mdl_prod-daily_dump-20180919.sql.gz | psql example_mdl_prod
```

Depending on the configuration and location of the PostgreSQL server, you
may need to pass additional options to `psql`.

Next: [Directory Format Database Backups](07-directory-format-backups.md)
