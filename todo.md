* Handle merge / merge-conflict (commit everything, pick newest, notify)
  - Add the logs of operation to .git/logs/auto-sync
* Daemon - proper logs
* Daemon - status

* CLI auto-complete

* Add a merge strategy where it doesn't take the YAML header changes into account
* Add a way to only sync some-subfolder in a repo
* Ensure that the 'daemon' process can be found
  - It seems it wasn't being installed by the homebrew package

* RPM / DEB packages - https://netdevops.me/2021/building-and-publishing-deb/rpm-packages-with-goreleaser-and-gemfury/

* Check the state of the git repo before
  - Make sure it isn't in the middle of a rebase / merge
  - Make sure that head commit is pointing to a branch
  - Make sure that branch has a proper remote tracking branch

* Write about how to handle the case with an ssh-agent running

* When stopping the daemon, do not do anything if it is not running

* Avoid these thousands of commits each time one saves a file
  - Try to batch them together?
  - Or optionally, it should only be saved when the editor is closed
    - For vim this means figuring out the backupfile directory
