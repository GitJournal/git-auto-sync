package common

import (
	"github.com/ztrue/tracerr"
)

func push(repoConfig RepoConfig) error {
	bi, err := fetchBranchInfo(repoConfig.RepoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if bi.UpstreamBranch == "" || bi.UpstreamRemote == "" {
		return nil
	}

	_, err = GitCommand(repoConfig, []string{"push", bi.UpstreamRemote, bi.UpstreamBranch})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
