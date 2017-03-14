<p align="center"><img height="180px" width="180px" src="https://github.com/fabric-8/borg/raw/master/assets/borg_mascot.png" alt=""></p>

BORG - Search and save shell snippets without leaving your terminal 
===
![cruft guaranteed](https://img.shields.io/badge/cruft-guaranteed-green.svg) [![Travis CI](https://api.travis-ci.org/ok-borg/borg.svg?branch=master)](https://travis-ci.org/ok-borg/borg) [![Go Report Card](https://goreportcard.com/badge/github.com/ok-borg/borg)](https://goreportcard.com/report/github.com/ok-borg/borg) [![Slack Status](http://ok-b.org:1492/badge.svg)](http://ok-b.org:1492)

Borg was built out of the frustration of having to leave the terminal to search and click around for bash snippets.
Borg's succint output also makes it easy to glance over multiple snippets quickly.

### Search

```
borg "list only files"
```

```shell
(1) Bash: How to list only files?
        [a] find . -maxdepth 1 -type f
        [b] ls -l | egrep -v '^d'
            ls -l | grep -v '^d'

(2) List only common parent directories for files
        [a] # read a line into the variable "prefix", split at slashes
            IFS=/ read -a prefix
            # while there are more lines, one after another read them into "next",
            # also split at slashes
            while IFS=/ read -a next; do
                new_prefix=()
                # for all indexes in prefix
                for ((i=0; i < "${#prefix[@]}"; ++i)); do
                    # if the word in the new line matches the old one
                    if [[ "${prefix[i]}" == "${next[i]}" ]]; then
        ...
```

Use `borg pipeto less` to pipe the results straight to `less`, or any other program of your choosing.

Can't find what you are looking for? Be a good hacker and contribute your wisdom to the hive mind - [add your own snippets or tweak the existing ones.](https://github.com/ok-borg/borg/tree/master/docs)

### Install

The following releases only let you search, to use add/edit install from source, releases are coming soon.

```
brew install borg
```

For linux, download a release manually [releases](https://github.com/ok-borg/borg/releases)

```
wget https://github.com/ok-borg/borg/releases/download/v0.0.3/borg_linux_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Or download a release manually for Mac:

```
wget https://github.com/ok-borg/borg/releases/download/v0.0.3/borg_darwin_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

### Rate results: `worked`

When you see a result that worked for you, you can use the `worked` command to give feedback:

```
borg worked 12
```

Once you do this the result will rank higher for similar queries - it is especially useful if you find a good result that you think are too down in the result list.

### Advanced usage

For more commands and their explanations, please see [advanced usage](https://github.com/ok-borg/borg/tree/master/docs)

### How does borg work?

The client connects to a server at ok-b.org, but you can host your own if you want to (see daemon folder).

Self hosting will become less appealing once people start contributing their own content to the database though.

### Explanation for ui

- `()` denotes hits for your query
- `[]` denotes snippets found for a given query
- `...` under a `[]` means more lines to display (use the `-f` flag for full display, see more about usage below)

### Credits

The borg mascot has been delivered to you by the amazing [Fabricio Rosa Marques](https://dribbble.com/fabric8).

### Community:

##### Running with docker

You can use the [dockerized borg client](https://github.com/juhofriman/borg-docker) if you don't want to install anything on your host!

### Outages

- We've had a multi day outage on Jan 18 - Jan 28. It ruined our uptime stats for the next couple of thousand years.
