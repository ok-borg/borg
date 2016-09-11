![](./borghead.png) BORG - A terminal based search engine for bash commands
===

Borg was built out of the frustration of having to leave the terminal to search for bash commands.
This is how it looks like:

```
$ borg "list all files in dir"
(1) create list of all files in every subdirectories in bash
        [11] find . -type f -exec md5 {} \;
        [12] #!/bin/sh
             DIR=${1:-`pwd`}
             SPACING=${2:-|}
             cd $DIR
             for x in * ; do
                 [ -d "$DIR/$x" ] &&  echo "$SPACING\`-{$x" && $0 "$DIR/$x" "$SPACING  " || \
                 echo "$SPACING $x : MD5=" && md5sum "$DIR/$x"
             done

(2) Bash: How to list only files?
        [21] find . -maxdepth 1 -type f
        [22] ls -l | egrep -v '^d'
             ls -l | grep -v '^d'
        [23]     find * -maxdepth 0 -type f  # find -L * ... includes symlinks to files
             -
             fls f        # list all files in current dir.
             fls d -tA ~  #  list dirs. in home dir., including hidden ones, most recent first
             fls f^l /usr/local/bin/c* # List matches that are files, but not (^) symlinks (l)
             -
             [sudo] npm install fls -g
        [24] find . -maxdepth 1 -type f|ls -lt|less
```

(Only displaying the first 2 hits here, but that's configurable, by default it's 5)

Some explanations for the UI:
- () denotes hits for your query
- [] denotes possible solutions
- ... under a [] means more lines to display
- a - in a solution means separate code examples extracted from the same source

#### Installation

Just do a `go install`, I will provide binaries to download later to make it easy to install for those who don't have go installed.

##### Feature

- only querying works for now

#### Future plans

- add a way to add entries and rate solutions
