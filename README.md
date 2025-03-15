<div align="center">

# ydo
![Github top language](https://img.shields.io/github/languages/top/nimaaskarian/ydo?style=flat&color=blue)

###### a frictionless, fast and featureful cli to-do application


[Getting started](#getting-started) •
[Installation](#installation) •
[Configuration](#configuration) 

</div>

---
ydo is a to-do app, with a command line interface heavily inspired by
[taskwarrior](https://taskwarrior.org/); but uses [yaml](yaml.org) for both
configurations and task files themselves. this makes ydo easily configured,
accessible, and fast.

> this app is under heavy development and will change frequently as for now.
> consult help pages of the build you use

# Getting started
if you haven't, [install](#installation) the application first

then just run the command below
```shell
ydo tutorial
```

# Installation

## from source
*nix, run the command below. in windows just go build and use the ydo.exe file
```
git clone https://github.com/nimaaskarian/ydo && \
cd ydo && \
go build && \
sudo cp ydo /usr/bin/ydo
```

## from releases
head to [releases](https://github.com/nimaaskarian/ydo/releases), and download
the binary according to your operating system (ydo for *nix, ydo.exe for windows)

then copy the binary to one of the directories under your `PATH` variable.

# compared to other cli to-do applications
## ydo compared to taskwarrior has
- a lot more speed
- human readable file format (yaml) for both config and todos
- no database, easy to sync (using syncthing and what not)
- graph dependencies

## ydo compared to [c3](https://github.com/nimaaskarian/c3) has
- less speed :<
- no built in tui :<
- human readable dependency tree (and a more human readable file format)
- no tie to calcurse
- yaml
- graph instead of tree
- frictionless, and complete cli interface
