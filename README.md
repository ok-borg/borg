<p align="center"><img src="https://github.com/fabric-8/borg/raw/master/borg_mascot.png" alt=""></p>

BORG - A terminal based search engine for bash commands
===
![](https://img.shields.io/badge/cruft-guaranteed-green.svg)

Borg was built out of the frustration of having to leave the terminal to search for bash commands.
Its succinct output also makes it easier to glance over sources

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

Just do a `go install`, or if you don't have Go installed there are builds in the build folder - pick one according to your environment.

It should be something like (this one is for linux):

```
wget https://github.com/crufter/borg/blob/master/builds/borg_linux_amd64\?raw\=true -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Now you are ready to rock! Query with `borg "my query"`, because there is a server listening to your questions and eager to help!
Keep querying, and let me know if you want something to be improved.

#### How it works

The client connects to a server at borg.crufter.com, but you can host your own if you want to (see daemon folder).
Self hosting will become less appealing once people start contributing their own content to the database though.

##### Features

- only querying works for now

#### Future plans

- add a way to add entries and rate solutions
