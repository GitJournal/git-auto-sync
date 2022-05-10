- Add the logs of operation to .git/logs/auto-sync

- Make this a daemon that just auto-commits for now
  - Get the config from /etc/git-auto-sync
  - Ensure it is valid?
  - https://github.com/kardianos/service

- Add a merge option

- Use
  - https://goreleaser.com/intro/

- https://stackoverflow.com/questions/21705950/running-external-commands-through-os-exec-under-another-user
  https://groups.google.com/g/golang-nuts/c/bcjk9ncP5ac

Commands -
* watch <repo>
  - Allow configuration to be changed
  - Passed via command line args
* sync
* test <repo>
  - Tests for write-access to the repo

* daemon - Run it in the background
  - add/ls/rm

## First Release

* Daemon + enable / disable
* Handle rebase / merge / merge-conflict (commit everything, pick newest, notify)
* Daemon - proper logs

* CLI auto-complete
* Allow running the 'sync' command from a non top level directory


* Add a merge strategy where it doesn't take the YAML header changes into account
