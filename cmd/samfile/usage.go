package main

func usage(versionName string) string {
	return versionName + `

Manipulate files in SAM Coupé floppy disk images.

  Usage:
    samfile extract [--dest TARGET] IMAGE [FILE]
    samfile ls IMAGE
    samfile --help
    samfile --version

  Targets:
    extract               Extracts one or more files from a SAM Disk image file.
    ls                    Lists files on SAM Disk image file.

  Options:
    --dest TARGET         When specifying a FILE to extract from a Disk image,
                          TARGET can be either the file to save to, or an
                          existing directory to write the file to with its
                          original name from the disk image. When extracting all
                          files (no FILE is specfied) TARGET should be an
                          existing directory to write all files to. Defaults to
                          current directory.
    --help                Display this help text.
    --version             Display the release version of samfile.
    IMAGE                 The raw floppy disk image (.mgt format / 819200 bytes)
                          On linux a floppy disk image can be created by running
                            dd if=/dev/fd0u800 of=image.mgt conv=noerror,sync
                          If /dev/fd0u800 does not exist it can be created with
                            sudo mknod /dev/fd0u800 b 2 120
    FILE                  A single file to extract from the disk image. To
                          extract ALL files from the disk image, omit FILE.
`
}
