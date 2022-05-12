package common

import (
	"errors"
	"os/exec"

	"github.com/ztrue/tracerr"
)

var errNoGitAuthorEmail = errors.New("Missing git author email")
var errNoGitAuthorName = errors.New("Missing git author name")

func ensureGitAuthor(repoPath string) error {
	_, err := GitCommand(repoPath, []string{"config", "user.email"})
	if err != nil {
		var exerr *exec.ExitError
		if errors.As(err, &exerr) && exerr.ExitCode() == 1 {
			return errNoGitAuthorEmail
		}
		return tracerr.Wrap(err)
	}

	_, err = GitCommand(repoPath, []string{"config", "user.name"})
	if err != nil {
		var exerr *exec.ExitError
		if errors.As(err, &exerr) && exerr.ExitCode() == 1 {
			return errNoGitAuthorName
		}
		return tracerr.Wrap(err)
	}

	return nil
}
