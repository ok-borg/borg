<p align="center"><img height="180px" width="180px" src="https://github.com/fabric-8/borg/raw/master/assets/borg_mascot.png" alt=""></p>

BORG - Terminal based search for bash snippets
===
![cruft guaranteed](https://img.shields.io/badge/cruft-guaranteed-green.svg) [![Travis CI](https://api.travis-ci.org/ok-borg/borg.svg?branch=master)](https://travis-ci.org/ok-borg/borg) [![Go Report Card](https://goreportcard.com/badge/github.com/ok-borg/borg)](https://goreportcard.com/report/github.com/ok-borg/borg) [![Slack Status](http://ok-b.org:1492/badge.svg)](http://ok-b.org:1492)

Borg was built out of the frustration of having to leave the terminal to search and click around for bash snippets.
Borg's succint output also makes it easy to glance over multiple snippets quickly.

### Search

```
borg "find all txt"
```

```
(1) Find and delete .txt files in bash
        [a] find . -name "*.txt" | xargs rm
        [b] find . -name "*.txt" -exec rm {} \;
        [c] $ find  . -name "*.txt" -type f -delete

(2) bash loop through all find recursively in sub-directories
        [a] FILES=$(find public_html -type f -name '*.php')
        [b] FILES=`find public_html -type d`
```

Can't find what you are looking for? Be a good hacker and contribute your wisdom to the hive mind - [add your own snippets or tweak the existing ones.](https://github.com/ok-borg/borg/tree/master/docs)

### Install

The following releases only let you search, to use add/edit install from source, releases are coming soon.

```
brew install borg
```

For linux, download a release manually [releases](https://github.com/ok-borg/borg/releases)

```
wget https://github.com/ok-borg/borg/releases/download/v0.0.1/borg_linux_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Or download a release manually for Mac:

```
wget https://github.com/ok-borg/borg/releases/download/v0.0.1/borg_darwin_amd64 -O /usr/local/bin/borg
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
