<p align="center"><img height="180px" width="180px" src="https://github.com/fabric-8/borg/raw/master/borg_mascot.png" alt=""></p>

BORG - A terminal based search engine for bash snippets 
===
![cruft guaranteed](https://img.shields.io/badge/cruft-guaranteed-green.svg) [![Travis CI](https://api.travis-ci.org/crufter/borg.svg?branch=master)](https://travis-ci.org/crufter/borg) [ok-b.org](http://ok-b.org)

Borg was built out of the frustration of having to leave the terminal to search for bash snippets.

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

##### Add

First obtain an oauth token by loggin in with github at [ok-b.org](http://ok-b.org). Copy the token you can find on [ok-b.org/me](http://ok-b.org/me) and log in with borg:

```
borg login my3XamPleT0k3n
```

You are ready to save your own content

```
borg new
```

A vim window opens and lets you save your snippet. For example:

```
How to grep for a file in current directory

ls | grep mySearchTerm
```

Save and exit vim.

##### Edit

Using our search example, typing `borg edit 1` will present you with an editor window containing:

```
Find and delete .txt files in bash
 
[a]
find . -name "*.txt" | xargs rm

[b]
find . -name "*.txt" -exec rm {} \;

[c]
10 $ find  . -name "*.txt" -type f -delete
```

Let's say you want to remove the second snippet because your don't like it. Modify it so it becomes:

```
Find and delete .txt files in bash
 
[a]
find . -name "*.txt" | xargs rm

[c]
10 $ find  . -name "*.txt" -type f -delete
```

Save and exit.

(Do not care about the incorrect alphabetical order, it's ok)

#### Installation

```
brew intall borg
```

For linux, download a release manually [releases](https://github.com/crufter/borg/releases)

```
wget https://github.com/crufter/borg/releases/download/v0.0.1/borg_linux_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

Or download a release manually for Mac:

```
wget https://github.com/crufter/borg/releases/download/v0.0.1/borg_darwin_amd64 -O /usr/local/bin/borg
chmod 755 /usr/local/bin/borg
```

##### Who can add/edit what?

Any logged in user can edit any content. We trust you with not being a vandal.

#### How it works

The client connects to a server at borg.crufter.com, but you can host your own if you want to (see daemon folder).
Self hosting will become less appealing once people start contributing their own content to the database though.

#### Features

Command line:
- search, add, edit content

Web:
- login, search, add, edit content

#### Future plans

- add a way to save private and organisation private entries
- enable users to rate results
- after a lot of lot of things make borg your own notebook/private search engine for anything

### Explanation for ui

- `()` denotes hits for your query
- `[]` denotes snippets found for a given query
- `...` under a `[]` means more lines to display (use the `-f` flag for full display, see more about usage below)

#### Usage

Borg supports gnu flags, so flags are supported both before and after the arguments, so all of the followings are valid:

```
borg -l 30 -f "md5 Mac"
borg "md5 Mac" -l30 -f
borg -f "md5 Mac" -l30
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

#### Credits

The borg mascot has been delivered to you by the amazing [Fabricio Rosa Marques](https://dribbble.com/fabric8).

### Community:

##### Running with docker

You can use the [dockerized borg client](https://github.com/juhofriman/borg-docker) if you don't want to install anything on your host!
<p align="center"><img height="180px" width="180px" src="https://github.com/fabric-8/borg/raw/master/borg_mascot.png" alt=""></p>

BORG - A terminal based search engine for bash snippets 
===
![cruft guaranteed](https://img.shields.io/badge/cruft-guaranteed-green.svg) [![Travis CI](https://api.travis-ci.org/crufter/borg.svg?branch=master)](https://travis-ci.org/crufter/borg) [ok-b.org](http://ok-b.org)

Borg was built out of the frustration of having to leave the terminal to search for bash snippets.
Its succinct output also makes it easier to glance over multiple snippets of code.

This is how it looks like:

```
$ borg "find all txt"
(1) Find and delete .txt files in bash
        [a] find . -name "*.txt" | xargs rm
        [b] find . -name "*.txt" -exec rm {} \;
        [c] $ find  . -name "*.txt" -type f -delete

(2) bash loop through all find recursively in sub-directories
        [a] FILES=$(find public_html -type f -name '*.php')
        [b] FILES=`find public_html -type d`
```

(Only displaying the first 2 hits here, but that's configurable, by default it's 5)

Some explanations for the UI:
- `()` denotes hits for your query
- `[]` denotes snippets found for a given query
- `...` under a `[]` means more lines to display (use the `-f` flag for full display, see more about usage below)

#### State

Please keep in mind that this is in a really-really early phase. Glitches are expected and there are a lot of low hanging fruits we can go after.

#### The web client

There is a web client available under [ok-b.org](http://ok-b.org) primarily for obtaining an oauth token, but it also provides more or less the same functionality what the command line tool does.

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

On Mac you can also use HomeBrew (however, probably brew releases will lag behind github releases, if you want the latest, see above):

```
brew intall borg
```

Now you are ready to rock! Query with `borg "my query"`, because there is a server listening to your questions and eager to help!
Keep querying, and let me know if you want something to be improved.

##### Search

```
borg "my search query"
```

##### Save your own snippet

First obtain an oauth token by loggin in with github at [ok-b.org](http://ok-b.org). Copy the token you can find on [ok-b.org/me](http://ok-b.org/me) and log in with borg:

```
borg login my3XamPleT0k3n
```

You are ready to save your own content

```
borg new
```

A vim window opens and lets you save your snippet. For example:

```
How to grep for a file in current directory

ls | grep mySearchTerm
```

Type `:wq`/`:x` and your snippet is already up in the cloud (publicly, private entries are coming)

##### Edit

Editing is really easy, to save you from having to type a long id to edit, you can use the index from the search result (remember `(1)` `(2)` etc?)

Remember the search example which gave back a list, the first element being this:

```
(1) Find and delete .txt files in bash
        [a] find . -name "*.txt" | xargs rm
        [b] find . -name "*.txt" -exec rm {} \;
        [c] $ find  . -name "*.txt" -type f -delete
```

Type `borg edit 1` and you will be presented with an editor window with content like this:

```
Find and delete .txt files in bash
 
[a]
find . -name "*.txt" | xargs rm

[b]
find . -name "*.txt" -exec rm {} \;

[c]
10 $ find  . -name "*.txt" -type f -delete
```

Edit it as you like, let's say you want to remove the second snippet because your don't like it. Modify it so it becomes:

```
Find and delete .txt files in bash
 
[a]
find . -name "*.txt" | xargs rm

[c]
10 $ find  . -name "*.txt" -type f -delete
```

Do not mind that `c` does not follow `a` in the alphabet. Just save, borg is smart enough to know what to do!

##### Who can add/edit what?

Any logged in user can edit any content. We trust you with not being a vandal.

#### How does borg work? Can it be mine and only mine?

The client by default connects to a server at borg.crufter.com, but you can host your own if you want to (see daemon folder).
Self hosting will become less appealing once people start contributing their own content to the database though.

##### Running with docker

You can use the [dockerized borg client](https://github.com/juhofriman/borg-docker) if you don't want to install anything on your host!

#### Features

Command line:
- search, add, edit content

Web:
- login, search, add, edit content

#### Future plans

- add a way to save private and organisation private entries
- enable users to rate results
- after a lot of lot of things make borg your own notebook/private search engine for anything

#### Credits

The borg mascot has been delivered to you by the amazing [Fabricio Rosa Marques](https://dribbble.com/fabric8).

#### Usage

Borg supports gnu flags, so flags are supported both before and after the arguments, so all of the followings are valid:

```
borg -l 30 -f "md5 Mac"
borg "md5 Mac" -l30 -f
borg -f "md5 Mac" -l30
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
