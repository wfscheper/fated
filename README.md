# FateD

FateD is a dice roller for [Fate Core](http://www.evilhat.com/home/fate-core/)
that operates in two modes: dice or deck.

## Getting started

Pre-built binaries for MacOS, Linux, and Windows are available on
[github](https://github.com/wfscheper/fated/releases).

Building this project requires Go to be installed. On MacOS with Homebrew you
can just run `brew install go`. For linux either check your distrubutions
package manager or follow the instructions at https://golang.org/doc/install.

Running it then should be as simple as:

```console
$ make
$ bin/fated
```

## Running fated

To draw a single Fate card::

```console
$ bin/fated draw
```

To roll dice instead::

```console
$ fated roll
```

If you expect to draw or roll multiple times and want a cleaner interface, use
the `-f, --foregroun` option::

```console
$ fated draw --foreground
$ fated roll --foreground
```
