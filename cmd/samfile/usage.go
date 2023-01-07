package main

func usage(versionName string) string {
	return versionName + `

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
`
}
