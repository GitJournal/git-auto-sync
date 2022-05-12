package common

import (
	"github.com/ztrue/tracerr"
)

func push(repoPath string) error {
	bi, err := fetchBranchInfo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if bi.UpstreamBranch == "" || bi.UpstreamRemote == "" {
		return nil
	}

	_, err = GitCommand(repoPath, []string{"push", bi.UpstreamRemote, bi.UpstreamBranch})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
