# samfile

[![CI](https://github.com/petemoore/samfile/actions/workflows/ci.yml/badge.svg)](https://github.com/petemoore/samfile/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/petemoore/samfile/v3.svg)](https://pkg.go.dev/github.com/petemoore/samfile/v3)
[![License](https://img.shields.io/badge/license-MIT-orange.svg)](https://opensource.org/licenses/MIT)

`samfile` is a tool for manipulating _individual files_ inside SAM Coupé floppy
disk images.

For reading, writing, creating and converting entire SAM disk images (and other
disk formats), see the excellent [samdisk](https://simonowen.com/samdisk)
utility.

```
$ samfile --help
samfile 3.0.0 [ revision: https://github.com/petemoore/samfile/commits/0a7a919d5a60b7881752c4f6c3cd5e9fb252e662 ]

Manipulate files in SAM Coupé floppy disk images.

  Usage:
    samfile add -i IMAGE -f FILE -c -l LOAD_ADDRESS [-e EXECUTION_ADDRESS]
    samfile basic-to-text
    samfile cat -i IMAGE -f FILE
    samfile extract -i IMAGE [-t TARGET]
    samfile ls -i IMAGE
    samfile --help
    samfile --version

  Targets:
    add                   Adds a file from the host file system to the SAM Disk
                          image file.
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
    -c                    File is a code file.
    -l LOAD_ADDRESS       Load address of code file on the SAM Disk image.
    -e EXECUTION_ADDRESS  Execution address of code file on the SAM Disk image.
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

To install the latest published version:

```
go install github.com/petemoore/samfile/v3/cmd/samfile@v3.0.0
```

To build and test from a local clone:

```
git clone https://github.com/petemoore/samfile.git
cd samfile
go install ./cmd/samfile     # installs samfile to $(go env GOPATH)/bin
go test ./...                # run unit and integration tests
go test -race ./...          # also run with the race detector (matches CI)
go vet ./...                 # static checks (matches CI)
```

## Releasing

See [`RELEASING.md`](RELEASING.md) for the cut-a-release procedure.
