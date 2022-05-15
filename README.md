# Git Auto Sync

GitAutoSync is a simple command line program to automatically commit changes
to your git repo, and always keep that repo up to date. This way you can use any editor with your text files, and never need to worry about
comitting and remembering to push and pull changes.

## Installation

* OSX - `brew install GitJournal/tap/git-auto-sync`
* Linux - Download the [latest release](https://github.com/GitJournal/git-auto-sync/releases/latest)
* Windows - Download the [latest release](https://github.com/GitJournal/git-auto-sync/releases/latest)

## How to use?

GitAutoSync comes with a manual and daemon mode. It's recommended to start with the manual
mode to ensure authentication is working correctly. It internally just calls the `git` executable
so, if that works, `git-auto-sync` should just work.

You can test it out by running `git-auto-sync sync` to commit, pull, rebase and push any changes.
If there are no changes, it will just attempt to pull, rebase and push.

Once you're satisfied that `git-auto-sync` is working for you. You can run `git-auto-sync daemon add <repoPath>` to start a background daemon which will continously monitor that repo for any changes
in the file system and accordingly sync the changes.

This daemon will be automatically started as a system process.

You can check if it is running by looking for a process called `git-auto-sync-daemon`

### Merge Conflicts

GitAutoSync current only supports rebases, and doesn't yet attempt to do a merge. In the case of a
rebase conflict, it will abort and stop syncing that repo. It will send a system notification
to inform you of the conflict.

### Ignored Files

It currently ignores all hidden files, files ignored by git, and additional temporary swap files
created by vim, emacs and similar editors.

## Similar Projects

- [Obsidian Git](https://github.com/denolehov/obsidian-git)
- [VS Code GitDoc](https://marketplace.visualstudio.com/items?itemName=vsls-contrib.gitdoc)
- [Git Annex](https://git-annex.branchable.com/)
- [Git Sync](https://github.com/simonthum/git-sync)
