# 2.0.0
- Rewrote in Go for faster syncing, especially under bursts of many simultaneous
  file changes (e.g. a git checkout or branch switch).
- Copies/deletes now run through a bounded, concurrent worker pool instead of
  firing unbounded per-event work.
- Rapid repeat events for the same path are debounced into a single copy.
- Directories matched by `.csignore` are no longer watched at all, instead of
  being watched and filtered per-event.
- Ships as a single static binary; no Node/npm runtime required.

# 1.1.1
- Updated Node packages.

# 1.1.0
- Added option to ignore files through `.csignore` file

# 1.0.0
- Initial release

