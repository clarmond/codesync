# CodeSync

This app is used to watch a folder for changes and update the target folder.  I built this app so that I could do development on my local environment and sync that code to a server over a network share.

## Prerequisites

- [Node.js](https://nodejs.org/)
- [git](https://git-scm.com/)

## Getting Started

```
git clone https://github.com/clarmond/codesync
cd codesync
npm install
```

## How to Use
`npm start source_folder target_folder`

For example:

`npm start C:\Users\dillinger\mcp\src \\encom\mcp\src`

`npm start /Users/abradley/tron /Volumes/tron`

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
npm start ~/my-code /Volumes/dev-area
```
