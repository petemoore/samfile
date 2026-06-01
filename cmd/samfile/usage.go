package main

func usage(versionName string) string {
	return versionName + `

Manipulate files in SAM Coupé floppy disk images.

  Usage:
    samfile add -i IMAGE -f FILE -c -l LOAD_ADDRESS [-e EXECUTION_ADDRESS]
    samfile add -i IMAGE -f FILE -s MODE
    samfile basic-to-text [--lossy]
    samfile text-to-basic
    samfile screen-to-png [--mode MODE] [--format FORMAT]
    samfile cat -i IMAGE -f FILE [-c] [-a FORMAT] [--loops N]
    samfile extract -i IMAGE [-t TARGET] [-c] [-a FORMAT] [--loops N]
    samfile ls -i IMAGE
    samfile --help
    samfile --version

  Targets:
    add                   Adds a file from the host file system to the SAM Disk
                          image file.
    basic-to-text         Read a SAM Basic encoded file from stdin and output
                          plain text listing to stdout.
    text-to-basic         Read plain-text SAM BASIC source from stdin and
                          output the tokenised program body (suitable for
                          piping into 'samfile basic-to-text' to verify
                          the round-trip).
    screen-to-png         Read a raw SAM SCREEN$ body from stdin and write a
                          PNG (default) or GIF image to stdout. Only MODE 4
                          (256x192, 16 colours) is decoded; a raw body has no
                          directory entry so the mode defaults to MODE 4 (use
                          --mode to override). Flashing screens (two differing
                          palettes) animate when --format gif is given;
                          otherwise palette A is rendered.
    cat                   Output a single file from a SAM Disk image file to
                          stdout. With -c, BASIC files are detokenised to text,
                          SCREEN$ files are emitted as PNG bytes, and E-Tracker
                          modules are rendered to audio.
    extract               Extracts all files from a SAM Disk image file to a
                          local directory. With -c, BASIC files are written as
                          '<name>.txt' (detokenised), SCREEN$ files as
                          '<name>.png', and E-Tracker modules as audio files
                          ('<name>.mp3' etc.); other types extract raw.
    ls                    Lists files on SAM Disk image file.

  Options:
    -i IMAGE              The raw floppy disk image (.mgt format / 819200 bytes)
                          On linux a floppy disk image can be created by running
                            dd if=/dev/fd0u800 of=image.mgt conv=noerror,sync
                          If /dev/fd0u800 does not exist it can be created with
                            sudo mknod /dev/fd0u800 b 2 120
    -t TARGET             An existing directory to write all files to. Defaults
                          to current directory.
    -f FILE               A single file inside the disk image (add/cat). For
                          'screen-to-png' the SCREEN$ body is read from stdin
                          instead.
    -c                    For 'add': the input file is a code file. For 'cat'
                          and 'extract': convert known types on output — SAM
                          BASIC files are detokenised to text, SCREEN$ files
                          are rendered to PNG, and E-Tracker modules are
                          rendered to audio; other types stay raw.
    -a FORMAT, --audio-format FORMAT
                          With -c, output format for E-Tracker tunes:
                          wav, flac, or mp3 [default: mp3]. Provenance tags are
                          written to mp3 (ID3v2) and wav (INFO); flac is
                          currently untagged.
    --loops N             With -c, times to repeat an E-Tracker tune's loop
                          body (after its one-shot intro) [default: 4].
    -s MODE, --screen MODE
                          Add the input file as a SCREEN$ (display dump) file
                          with the given SAM MODE (1..4). The input file (-f)
                          must be the raw SCREEN$ body.
    --mode MODE           (screen-to-png) SAM MODE (1..4) of the SCREEN$ body
                          on stdin. Defaults to 4. Only MODE 4 decodes.
    --format FORMAT       (screen-to-png) Output image format: png (default)
                          or gif. gif animates flashing two-palette screens.
    -l LOAD_ADDRESS       Load address of code file on the SAM Disk image.
    -e EXECUTION_ADDRESS  Execution address of code file on the SAM Disk image.
    --help                Display this help text.
    --version             Display the release version of samfile.
    --lossy               (basic-to-text) Emit the byte-for-byte
                          equivalent of the SAM ROM's LLIST routine
                          (matches stream 3 / printer output). Filters
                          attribute control sequences (INK, PAPER,
                          BRIGHT, etc.), wraps at column 80 with a
                          6-space continuation indent, terminates lines
                          with CR LF, and emits '>' after the first
                          line's number. Use for SAM ROM oracle
                          comparisons; without --lossy the output is
                          round-trip-faithful through text-to-basic.

  Examples:

    Extract SAM Basic file 'SCREENS' from disk image 'fred27.mgt' and write to
    file 'SCREENS.basic' as plain text listing:

    $ samfile cat -i fred27.mgt -f SCREENS | samfile basic-to-text > SCREENS.basic

    Render a SCREEN$ body to a PNG:

    $ samfile cat -i fred27.mgt -f TITLE | samfile screen-to-png > TITLE.png

    Add a MODE 4 SCREEN$ body 'title.scr' to a disk image:

    $ samfile add -i fred27.mgt -f title.scr -s 4

  SAMFile source code:
    https://github.com/petemoore/samfile
`
}
