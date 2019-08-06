# samfile

[![Build Status](https://img.shields.io/travis/petemoore/samfile.svg?style=flat-square&label=build+status)](https://travis-ci.org/petemoore/samfile)
[![GoDoc](https://godoc.org/github.com/petemoore/samfile?status.svg)](https://godoc.org/github.com/petemoore/samfile)
[![Coverage Status](https://coveralls.io/repos/petemoore/samfile/badge.svg?branch=master&service=github)](https://coveralls.io/github/petemoore/samfile?branch=master)
[![License](https://img.shields.io/badge/license-MIT-orange.svg)](https://opensource.org/licenses/MIT)

`samfile` is a tool for manipulating _individual files_ inside SAM Coupé floppy
disk images.

For reading, writing, creating and converting entire SAM disk images (and other
disk formats), see the excellent [samdisk](https://simonowen.com/samdisk)
utility.

```
$ samfile --help
samfile [ revision: https://github.com/petemoore/samfile/commits/6f52b243095c2f91ad5b843e7892c495ddf01ea3 ]

Manipulate files in SAM Coupé floppy disk images.

  Usage:
    samfile extract -i IMAGE [-t TARGET]
    samfile cat -i IMAGE -f FILE
    samfile ls -i IMAGE [-v|-n]
    samfile basic-to-text
    samfile --help
    samfile --version

  Targets:
    basic-to-text         Read a SAM Basic encoded file from stdin and output
                          plain text listing to stdout.
    cat                   Output a single file from a SAM Disk image file to
                          stdout.
    extract               Extracts all files from a SAM Disk image file to a
                          local directory.
    ls                    Lists files on SAM Disk image file.

  Options:
    -i IMAGE              The raw floppy disk image (.mgt format / 819200 bytes)
                          On linux a floppy disk image can be created by running
                            dd if=/dev/fd0u800 of=image.mgt conv=noerror,sync
                          If /dev/fd0u800 does not exist it can be created with
                            sudo mknod /dev/fd0u800 b 2 120
    -t TARGET             An existing directory to write all files to. Defaults
                          to current directory.
    -f FILE               A single file inside the disk image.
    -v                    Verbose output.
    -n                    Output names of files only.
    --help                Display this help text.
    --version             Display the release version of samfile.

  Examples:

    Extract SAM Basic file 'SCREENS' from disk image 'fred27.mgt' and write to
    file 'SCREENS.basic' as plain text listing:

    $ samfile cat -i fred27.mgt -f SCREENS | samfile basic-to-text > SCREENS.basic

  SAMFile source code:
    https://github.com/petemoore/samfile
```

## Installing

Download from https://github.com/petemoore/samfile/releases

## Building from source

```
go get github.com/petemoore/samfile/...
```
