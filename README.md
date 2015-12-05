# The Platinum Searcher [![Build Status](https://travis-ci.org/monochromegane/the_platinum_searcher.svg?branch=master)](https://travis-ci.org/monochromegane/the_platinum_searcher) [![wercker status](https://app.wercker.com/status/59ef90ac217537abc0994546958037f3/s/master "wercker status")](https://app.wercker.com/project/bykey/59ef90ac217537abc0994546958037f3)

A code search tool similar to `ack` and `the_silver_searcher(ag)`. It supports multi platforms and multi encodings.

## Features

- It searches code about 3–5× faster than `ack`.
- It searches code as fast as `the_silver_searcher(ag)`.
- It ignores file patterns from your `.gitignore` and `.hgignore`.
- It searches `UTF-8`, `EUC-JP` and `Shift_JIS` files.
- It provides binaries for multi platform (Mac OS X, Windows, Linux).

### Benchmarks

```sh
ack go  6.24s user 1.06s system 99%  cpu 7.304 total # ack:
ag go   0.88s user 1.39s system 221% cpu 1.027 total # ag:  It's faster than ack
pt go   1.09s user 1.01s system 235% cpu 0.892 total # pt:  It's faster than ag!!
```

## Installation

### Developer

```sh
$ go get -u github.com/monochromegane/the_platinum_searcher/...
```

### User

Download from the following url.

- [https://github.com/monochromegane/the_platinum_searcher/releases](https://github.com/monochromegane/the_platinum_searcher/releases)

Or, you can use Homebrew (Only MacOSX).

```sh
$ brew install pt
```

`pt` is an alias for `the_platinum_searcher` in Homebrew.

## Usage

```sh
$ # Recursively searchs for PATTERN in current directory.
$ pt PATTERN

$ # You can specified PATH and some OPTIONS.
$ pt OPTIONS PATTERN PATH
```

## Configuration

### .ptignore

If you put .ptignore on $HOME or current directory, pt ignore the files.
The file is .gitignore format.

### .ptconfig.toml

If you put .ptconfig.toml on $HOME or current directory, pt use option in the file.
The file is TOML format like the following.

```toml
color = true
context = 3
ignore = ["dir1", "dir2"]
```

## Editor Integration

### Vim + Unite.vim

You can use pt with [Unite.vim](https://github.com/Shougo/unite.vim).

```vim
nnoremap <silent> ,g :<C-u>Unite grep:. -buffer-name=search-buffer<CR>
if executable('pt')
  let g:unite_source_grep_command = 'pt'
  let g:unite_source_grep_default_opts = '--nogroup --nocolor'
  let g:unite_source_grep_recursive_opt = ''
  let g:unite_source_grep_encoding = 'utf-8'
endif
```

### Emacs + pt.el

You can use pt with [pt.el](https://github.com/bling/pt.el), which can be installed from [MELPA](http://melpa.milkbox.net/).

## Contribution

1. Fork it
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request
