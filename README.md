<p align="center"><img height="180px" width="180px" src="https://github.com/fabric-8/borg/raw/master/borg_mascot.png" alt=""></p>

BORG - A terminal based search engine for bash commands
===
![cruft guaranteed](https://img.shields.io/badge/cruft-guaranteed-green.svg) [![Travis CI](https://api.travis-ci.org/crufter/borg.svg?branch=master)](https://travis-ci.org/crufter/borg)

Borg was built out of the frustration of having to leave the terminal to search for bash commands.
Its succinct output also makes it easier to glance over multiple snippets of code.

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
- `()` denotes hits for your query
- `[]` denotes possible solutions
- `...` under a `[]` means more lines to display (use the `-f` flag for full display, see more about usage below)
- a `-` in a solution means separate code snippets extracted from the same source

#### State

Please keep in mind that this is in a really-really early phase.
Glitches are expected and there are a lot of low hanging fruits I can go after.
The relevancy of the results, the interface, and the available features in general will be greatly improved in the coming weeks.

#### Installation

Just do a `go install`, or if you don't have Go installed check out the [releases](https://github.com/crufter/borg/releases) and download the appropriate binary for your system. 

For example, for Linux:

```
wget https://github.com/crufter/borg/releases/download/v0.0.1/borg_linux_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Or Mac:

```
wget https://github.com/crufter/borg/releases/download/v0.0.1/borg_darwin_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Now you are ready to rock! Query with `borg "my query"`, because there is a server listening to your questions and eager to help!
Keep querying, and let me know if you want something to be improved.

##### Running with docker

You can use the [dockerized borg client](https://github.com/juhofriman/borg-docker) if you don't want to install anything on your host!

#### How it works

The client connects to a server at borg.crufter.com, but you can host your own if you want to (see daemon folder).
Self hosting will become less appealing once people start contributing their own content to the database though.

#### Features

- only querying works for now

#### Future plans

- add a way to add public, private and organisation private entries
- enabling users to rate solutions
- way more, but first let's tackle the ones above =)

#### Credits

The borg mascot has been delivered to you by the amazing [Fabricio Rosa Marques](https://dribbble.com/fabric8).

#### Usage

Borg supports gnu flags, so flags are supported both before and after the arguments, so all of the followings are valid:

```
borg -l 1 -f "md5 Mac"
borg "md5 Mac" -l 1 -f
borg -f "md5 Mac" -l 1
```

But what do they do?

```
-f  (= false)
    Print full results, ie. no more '...'
-h (= "borg.crufter.com")
    Server to connect to
-l  (= 5)
    Result list limit. Defaults to 5
-p  (= false)
    Private search. Your search won't leave a trace. Pinky promise. Don't use this all the time if you want to see the search result relevancy improved
```
