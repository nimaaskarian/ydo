<div align="center">
   
# ![ydo](https://raw.githubusercontent.com/nimaaskarian/ydo/refs/heads/master/cmd/webgui/static/imgs/icon-48x48.png)

![Github top language](https://img.shields.io/github/languages/top/nimaaskarian/ydo?style=flat&color=blue)

###### a frictionless, fast and featureful cli to-do application
seriously, y do?


[Getting started](#getting-started) •
[Installation](#installation) •
[Configuration](#configuration) 

</div>

---
ydo is a to-do app, with a command line interface heavily inspired by
[taskwarrior](https://taskwarrior.org/); but uses [yaml](yaml.org) for both
configurations and task files themselves. This makes ydo easily configured,
accessible, and fast.

> this app is under heavy development and will change frequently as for now.
> consult help pages of the build you use

# Getting started
If you haven't, [install](#installation) the application first.

## Adding tasks
You can add tasks by running `ydo add your task comes here`, for example for
adding a task to buy groceries you can just:
```console
nima@foo:~$: ydo add buy groceries
Task "t1" added
- [ ] t1: buy groceries
```

## Understanding keys, and editing a task
Each task has a unique string as its key. This can be automatically generated by
`ydo`, as you seen in the command above, or specified by you!

You have to use this key as an identifier of that task.
For example, if I wanna edit the task I've just created, I have to use `ydo edit` command.
It kinda works like add, but you have to also mention the key you want to edit.
```console
nima@foo:~$: ydo edit t1 buy groceries from the local shop
- [ ] t1: buy groceries from the local shop
```

## Adding dependencies
Using the `-D <key>` flag on `ydo add`, you can add a dependency to a task that
you've already created. You can also use it `n` times like so: `-D <key1> -D <key2> ... -D <keyn>`.
```console
nima@foo:~$: ydo add -D t1 buy milk
Task "t2" added
- [ ] t1: buy groceries from the local shop
   - [ ] t2: buy milk
nima@foo:~$: ydo add -D t1 buy bread
Task "t3" added
- [ ] t1: buy groceries from the local shop
   - [ ] t2: buy milk
   - [ ] t3: buy bread
```

## Setting a task as done
You just need to run `ydo do <keys>`. You can specify one or more keys.

If no task is specified all the tasks would be set as done (after asking you
to confirm angrily)

```console
nima@foo:~$: ydo do t2 t3
- [ ] t1: buy groceries from the local shop
   - [x] t2: buy milk
   - [x] t3: buy bread
nima@foo:~$: ydo do t1
- [x] t1: buy groceries from the local shop
   - [x] t2: buy milk
   - [x] t3: buy bread
```
`ydo rm`, `ydo undo`, `ydo yaml`, `ydo md` work the same way with keys.
Use `ydo help` for a complete help, and `ydo <command> --help` to see a complete
help for a subcommand (for example `ydo add --help`).

## Tf-idf key generation
You can use `--tfidf` flag for `ydo add` to use tf-idf for automatic key
generation. You can also use the [config file](https://github.com/nimaaskarian/ydo/blob/master/config.yaml) option to enable and set a criteria
for this. (copy config file to `~/.config/ydo/config.yaml` on *nix,
`%APPDATA%\ydo\config.yaml` on windows)



# Installation

## from source
in *nix, run the command below. in windows just go build and use the `ydo.exe` file
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

# License
GPL-3.0

also the cmd/webgui/tailwind.css' theme section is the modernized version of [wheatjs' gruvbox tailwind theme](https://github.com/wheatjs/gruvbox-tailwind-theme/) 
