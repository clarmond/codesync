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

`npm start ~/apps/my-code /Volumes/mounts/dev-area`

## Notes
This is a _one-way_ sync.  Changes made on the target folder will _not_ be reflected back on the source folder.

It's recommended that you create a bat file or shell script (depending on your environment) to run the command for you.  For example:

```
#!/bin/bash
cd ~/codesync
npm start ~/my-code /Volumes/mounts/dev-area
```

## License
Copyright 2023 Chad Armond

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
