- Add the logs of operation to .git/auto-sync

- Make this a daemon that just auto-commits for now
  - Get the config from /etc/git-auto-sync
  - Ensure it is valid?
  - https://github.com/kardianos/service

- Add a merge option
- Push the changes

- Use
  - https://goreleaser.com/intro/

- Add a 'watch' command

- Make sure the author / email is set!
  - Check via git config user.email & user.name

- https://stackoverflow.com/questions/21705950/running-external-commands-through-os-exec-under-another-user
  https://groups.google.com/g/golang-nuts/c/bcjk9ncP5ac


1. Ignore files which .gitignore ignores

Commands -
* watch <repo>
  - Allow configuration to be changed
  - Passed via command line args
* sync

* enable
* disable
* status
* daemon - Run it in the background

* watch add/ls/rm
  -

## First Release
* Daemon + enable / disable
* Handle rebase / merge / merge-conflict (commit everything, pick newest, notify)
