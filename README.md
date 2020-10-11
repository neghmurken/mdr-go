# MDR-GO #

A Golang implementation of the MDR chat protocol (https://github.com/neghmurken/mdr) as a toy-project for the 13th KNPLabs hackathon

## Requirements

 - Docker
 - inotify-tools (optional)

## How to dev

Build the binary with

```shell
$ make build
```

This will produce a `mdr` binary in the `/bin` folder

Use the project watcher to automatically rebuild the binary if a go file is modified

```shell
$ make watch
```

## How to run

```shell
$ make run

... or

$ /bin/mdr
```
