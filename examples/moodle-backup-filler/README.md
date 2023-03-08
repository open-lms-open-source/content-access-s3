# moodle-backup-filler

This directory contains source code for `moodle-backup-filler`, which can be
used to convert Open LMS Enterprise "fileless" backups to standard Moodle
course backups that can be restored into Moodle outside of the Open LMS
Enterprise environment.

## Fileless Backups

Standard Moodle course backups are made up of three parts:

* metadata, or information about the backup (mostly XML files)
* data exported from the Moodle database (XML files)
* files referenced by the course (the "files/" directory within the backup)

All files stored by Moodle are immutable since the filenames are a checksum of
the content.  Open LMS Enterprise exploits this immutability by omitting files
from the course backup.  Instead, files are kept in Moodle's file store (even
after being deleted from Moodle).  When a course is restored from a fileless
backup, the files are still in the file store and are simply associated with
the restored course.

This approach reduces disk capacity compared with standard Moodle course
backups by deduplicating files across multiple backups, while also providing
a speed boost for backups and restores.

The downside to all this, however, is _**if fileless course backups are
restored to a Moodle that doesn't already have the files it references, the
files will not be accessible**_.

## Adding Files to Fileless Backups

To mitigate the downside, fileless backups can be converted to standard
Moodle backups by adding the files.  The metadata includes a flag indicating
that it's a fileless backup, so that needs to be modified too.
`moodle-backup-filler` performs these tasks, converting fileless backups to
standard Moodle course backups using files retrieved from S3, HTTP, or local
disk.

### Building moodle-backup-filler

You'll need a [Golang](https://golang.org/) build environment plus GNU Make.
You'll also need a local copy of this repository.  We've successfully built
the software on MacOS and Linux platforms using Golang 1.11, though it may
also be possible to build on other platforms and other versions of Golang.

To build, ensure you're in the directory containing this README.md, and run
`make` from a terminal window.  You should see something like this:

```
$ make
▶ running gofmt…
▶ setting GOPATH…
▶ building github.com/golang/dep/cmd/dep…
▶ retrieving dependencies…
▶ building github.com/golang/lint/golint…
▶ running golint…
▶ building executable…
```

Once complete, all going well you'll find `moodle-backup-filler` built for
your platform in the `bin/` directory:

```
$ ls bin
moodle-backup-filler
```

### Using moodle-backup-filler

You'll require one or more Open LMS Enterprise fileless backups on local disk,
and access to the files referenced by them (either by having them on local
disk, or by having access to the Open LMS Enterprise S3 bucket that contains
them or a copy of it).

To convert a single fileless backup in the file `in.mbz` using files on
local disk stored in a directory called `files` and writing the resulting
Moodle course backup to `out.mbz`:

```bash
$ moodle-backup-filler --source in.mbz --dest out.mbz --contentbase files
```

To convert all fileless backups from a directory `in` using files on local
disk stored in a directory called `files` and writing the resulting Moodle
course backups to a directory `out`:

```bash
$ moodle-backup-filler --sourcedir in --destdir out --contentbase files
```

A TOML format configuration file can be used in place of command line
options.  An example configuration file can be found
[here](moodle-backup-filler.toml).  To use a configuration file, specify it
with the --config command line option:

```bash
$ moodle-backup-filler --config moodle-backup-filler.toml
```

If an option is specified through both a configuration file and the command
line, the command line takes precedence.

## Contributing

Please read [CONTRIBUTING.md](../../CONTRIBUTING.md) for details on the
process for submitting pull requests for this project.

## License

This project is licensed under the GNU General Public License - see the
[LICENSE.md](../../LICENSE) file for details.
