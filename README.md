<span align=center>
  
  # Gigit

  ###### gigit user/repo

</span>

---

Gigit, a tool for downloading repositories at a reasonable speed.
It is written in Go, making it very efficient in starting any project.

## Installation

```
go install github.com/nazhard/gigit/cmd/gigit@latest
```

## Example

By default, gigit will download repos from GitHub, our favorite git as a service.

```sh
gigit user/repo
```

#### Spesific branch, commit hash, tag

You can use specific branches, commits, or tags with a `#`

```sh
gigit user/repo#dev

gigit user/repo#691c0bf

# on spesific tag, "v" is required
gigit user/repo#v1.0.0
```

#### Subdir

Get sub directory of a repository.

```sh
gigit user/repo/dir

gigit nazhard/gigit/cmd/gigit
```

#### Commands

Clone instead of download.

With cloning, you will get a .git folder (because it's just git clone 😂)

```sh
gigit clone user/repo
```

Clone with `--depth=1` if you just want to fix typo

```sh
gigit c1 user/repo
# or
gigit 1 user/repo
```

More documentation at [pkg.go.dev/github.com/nazhard/gigit](https://pkg.go.dev/github.com/nazhard/gigit)

## Why not use degit instead?

I don't know.
I was originally using degit with pnpm, and I felt this way:

- Slow
- Buggy
- Good

Honestly, it's a good project.
But, it doesn't seem to be maintained.

From there I thought about making something similar, with some improvements (maybe).

---

Contributors are welcome! 🤗

Inspired by [degit](https://github.com/Rich-Harris/degit)
