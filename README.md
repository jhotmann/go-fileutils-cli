# FileUtils (FU)

A cross-platform collection of command line utilities for file interactions.

### Features
- Powered by a Django-like templating engine [pongo2](https://github.com/flosch/pongo2)
- Command history and undo abilities (TODO)
- Glob file matching
- Auto-indexing when multiple files result in the same output
- Favorite commands for quick access

### Opinionated Defaults
- Missing directories are automatically created (like `mkdir -p`)
- Using the `rename` alias keeps files in the same directory
- Outputting to an existing file will result in a prompt instead of overwriting (unless `--force` option used)

## Usage

`fu [command] [command options]`

| Command | Alias(es) | Description |
| ----- | ----- | ----- |
| [cp](#copy) | copy | copy one or more files/directories to a destination (with variable support) |
| [favorites](#favorites) | f, fav, favourites | run, view, and edit favorited commands |
| [hash](#hash) | md5, sha1, sha256, sha512 | get the hash of one or more files (use the appropriate alias for the algorithm you need) |
| help | | view help (works with individual commands as well) |
| [history](#history) | h | view, undo, re-run, copy, and favorite past commands |
| [ln](#link) | link, mklink | create soft or hard links to one or more files (with variable support) |
| [mv](#move) | move, rename | move/rename one or more files/directories (with variable support) |
| [undo](#undo) | u | undo the last undoable command that hasn't already been undone |

## Installation

## Move

## Copy

## Link

## Variables And Filters

## Hash

## History

## Undo

## Favorites