githelpy [![Build Status](https://travis-ci.org/petems/githelpy.svg?branch=master)](https://travis-ci.org/petems/githelpy) [![codecov](https://codecov.io/gh/petems/githelpy/branch/master/graph/badge.svg)](https://codecov.io/gh/petems/githelpy) [![codebeat badge](https://codebeat.co/badges/cc515300-053e-4b62-8184-645be6e6aa2f)](https://codebeat.co/projects/github-com-antham-githelpy)
======

githelpy analyze commits messages to ensure they follow defined pattern.

[![asciicast](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2.png)](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2)

## Summary

* [Setup](#setup)
* [Usage](#usage)
* [Practical Usage](#practical-usage)
* [Third Part Libraries](#third-part-libraries)

## Setup

Download from release page according to your architecture githelpy binary : https://github.com/petems/githelpy/releases

### Define a file .githelpy.toml

Create a file ```.githelpy.toml``` at the root of your project, for instance :

```toml
[config]
exclude-merge-commits=true
check-summary-length=true

[matchers]
all="(?:ref|feat|test|fix|style)\\(.*?\\) : .*?\n(?:\n?(?:\\* |  ).*?\n)*"

[examples]
a_simple_commit="""
[feat|test|ref|fix|style](module) : A commit message
"""
an_extended_commit="""
[feat|test|ref|fix|style](module) : A commit message

* first line
* second line
* and so on...
"""
```

#### Config

* ```exclude-merge-commits``` : if set to true, will not check commit mesage for merge commit
* ```check-summary-length``` : if set to true, check commit summary length is 50 characters

#### Matchers

You can define as many matchers you want, naming is up to you, they will all be compared against a commit message till one match.

#### Examples

Provided to help user to understand where is the problem, like matchers you can define as many examples as you want, they all will be displayed to the user if an error occured.

If you defined for instance  :

```
a_simple_commit="""
[feat|test|ref|fix|style](module) : A commit message
"""
```

this example will be displayed to the user like that :

```
A simple commit :

[feat|test|ref|fix|style](module) : A commit message
```

key is used as a title, underscore are replaced with withespaces.

## Usage

```bash
Ensure your commit messages are consistent

Usage:
  githelpy [command]

Available Commands:
  check       Check ensure a message follows defined patterns
  version     App version

Flags:
      --config string    (default ".githelpy.toml")
  -h, --help            help for githelpy

Use "githelpy [command] --help" for more information about a command.
```

### check

```bash
Check ensure a message follows defined patterns

Usage:
  githelpy check [flags]
  githelpy check [command]

Available Commands:
  commit      Check commit message
  message     Check message
  range       Check messages in commit range

Flags:
  -h, --help   help for check

Global Flags:
      --config string    (default ".githelpy.toml")

Use "githelpy check [command] --help" for more information about a command.


You need to provide two commit references to run matching for instance :
```

#### check commit

Check one comit from its commit ID, doesn't support short ID currently :

```githelpy check commit aeb603ba83614fae682337bdce9ee1bad1da6d6e```

#### check message

Check a message, useful for script for instance when you want to use it with git hooks :

```githelpy check message "Hello"```

#### check range

Check a commit range, useful if you want to use it with a CI to ensure all commits in branch are following your conventions :

* with relative references                             : ```githelpy check range master~2^ master```
* with absolute references                             : ```githelpy check range dev test```
* with commit ids (doesn't support short ID currently) : ```githelpy check range 7bbb37ade3ff36e362d7e20bf34a1325a15b 09f25db7971c100a8c0cfc2b22ab7f872ff0c18d```

## Practical usage

If your system isn't described here and you find a way to have githelpy working on it, please improve this documentation by doing a PR for the next who would like to do the same.

### Git hook

It's possible to use githelpy to validate each commit when you are creating them. To do so, you need to use the ```commit-msg``` hook, you can replace default script with this one :

```
#!/bin/sh

githelpy check message "$(cat "$1")";
```

### Travis

In travis, all history isn't cloned, default depth is 50 commits, you can change it : https://docs.travis-ci.com/user/customizing-the-build#Git-Clone-Depth.

First, we download the binary from the release page according to the version we want and we add in ```.travis.yml``` :

```yaml
before_install:
  - wget -O /tmp/githelpy https://github.com/petems/githelpy/releases/download/v2.0.0/githelpy_linux_386 && chmod 777 /tmp/githelpy
```

We can add a perl script in our repository to analyze the commit range against master for instance (master reference needs to be part of cloned history):

```perl
#!/bin/perl

`git ls-remote origin master` =~ /([a-f0-9]{40})/;

my $refHead = `git rev-parse HEAD`;
my $refTail = $1;

chomp($refHead);
chomp($refTail);

if ($refHead eq $refTail) {
    exit 0;
}

system "githelpy check range $refTail $refHead";

if ($? > 0) {
    exit 1;
}
```

And finally in ```.travis.yml```, make it crashs when an error occured :

```yaml
script: perl test-branch-commit-messages-in-travis.pl
```

### CircleCI

In CircleCI, there is an environment variable that describe current branch : ```CIRCLE_BRANCH``` (https://circleci.com/docs/environment-variables/).

First, we download the binary from the release page according to the version we want and we add in ```circle.yml``` :

```yaml
dependencies:
  pre:
    - wget -O /home/ubuntu/bin/githelpy https://github.com/petems/githelpy/releases/download/v2.0.0/githelpy_linux_386 && chmod 777 /home/ubuntu/bin/githelpy
```

And in ```test``` we can run githelpy against master for instance :

```
test:
  override:
    - githelpy check range master $CIRCLE_BRANCH
```

## Third Part Libraries

### Nodejs

* [githelpyjs](https://github.com/dschnare/githelpyjs) : A Nodejs wrapper for githelpy
