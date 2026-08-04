# CodeSync

This app is used to watch a folder for changes and update the target folder.  I built this app so that I could do development on my local environment and sync that code to a server over a network share.

## Prerequisites

- [Go](https://go.dev/) 1.21+
- [git](https://git-scm.com/)

## Getting Started

```
git clone https://github.com/clarmond/codesync
cd codesync
go build -o codesync .
```

## How to Use
`./codesync source_folder target_folder`

For example:

`./codesync C:\Users\dillinger\mcp\src \\encom\mcp\src`

`./codesync /Users/abradley/tron /Volumes/tron`

An optional `-concurrency N` flag controls how many files are copied/deleted in
parallel (default `16`). Raise it if you're syncing to a share that handles many
concurrent writes well; lower it if a large burst of changes (e.g. a git checkout)
seems to overwhelm the destination.

## Ignore Files
You can create a list of files to be ignored by CodeSync by creating a `.csignore` file.  (Dot files are excluded automatically.) This file follows the [.gitignore spec](https://git-scm.com/docs/gitignore).  For example:

```
node_modules/
README.md
```

## Notes
This is a _one-way_ sync.  Changes made on the target folder will _not_ be reflected back on the source folder.

It's recommended that you create a bat file or shell script (depending on your environment) to run the command for you.  For example:

```
#!/bin/bash
cd ~/codesync
./codesync ~/my-code /Volumes/dev-area
```

## Building for other platforms
Go cross-compiles without extra tooling. For example, to build a Windows binary from macOS/Linux:

```
GOOS=windows GOARCH=amd64 go build -o codesync.exe .
```
