# Advanced usage

### Add

First obtain an oauth token by loggin in with github at [ok-b.org](http://ok-b.org).

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

### Edit

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

### Who can add/edit what?

Any logged in user can edit any content. We trust you with not being a vandal.

### Flags

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
-h (= "ok-b.org")
    Server to connect to
-l  (= 5)
    Result list limit. Defaults to 5
-p  (= false)
    Private search. Your search won't leave a trace. Pinky promise. Don't use this all the time if you want to see the search result relevancy improved
```

