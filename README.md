# Git Auto Sync

## Existing Work
- https://marketplace.visualstudio.com/items?itemName=vsls-contrib.gitdoc
- auto sync script
- Obsidian's git plugin

Lets write this in GoLang as that's the easiest to do. It'll also work across
all platforms, and it will be simple to deploy. Very low memory footprint.

Later we can add a GUI based on electron if required or hell, even Qt.

- List of paths where the repo exists
  - It will sync all branches that are being tracked?
  - Just the current branch?

  - Polling Frequency
  - Watch for FS changes

  - Optional pausing of this checking when certain apps are open
  - Optional pausing for certain files until program is closed
    - .fileName.swp (for vim)
      .fileName.swap (for ?)
      ~fileName (for emacs?)
      or fileName~ for kate
    - Ignore all hidden files
    - Stops 'FS' changes, but polling should still happen for that file
      in case we don't close the editor for a long time or if the editor crashes

- Provide a SystemD service file
- Provide one for MacOS

## Daemon

There can be a git-auto-sync-daemon or `gasd`. I'm a bit scared about name collision. But it seems to be fine.

* https://github.com/takama/daemon
  - Has Windows support

* https://github.com/kardianos/service
* https://ieftimov.com/posts/four-steps-daemonize-your-golang-programs/

See if I can figure out 'onSuspend' and 'onResume' events. Also onShutdown and onBoot.

- https://github.com/prashantgupta24/mac-sleep-notifier
- How to do so on Linux?
  - https://github.com/coreos/go-systemd
  - logind signal - PrepareForSleep & PrepareForShutdown
  - Possibly even inhibit shutdown / suspend
