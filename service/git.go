package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/philipp-mlr/al-id-maestro/model"
)

func getRemoteBranches(config *model.RemoteConfiguration) ([]model.Branch, error) {
	branches := []model.Branch{}

	remote := git.NewRemote(
		memory.NewStorage(),
		&gitConfig.RemoteConfig{
			Name: config.RemoteName,
			URLs: []string{config.RepositoryURL},
		})

	rfs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: config.RemoteName,
			Password: config.GithubAuthToken,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, rf := range rfs {
		if rf.Name().IsRemote() || rf.Name().IsBranch() {
			if !isExcludeBranch(rf.Name().Short(), config.ExcludeBranches) {
				branches = append(branches,
					model.Branch{
						Name:           rf.Name().Short(),
						RepositoryName: config.RepositoryName,
						CommitID:       rf.Hash().String(),
					})
			}
		}
	}

	return branches, nil
}

func isExcludeBranch(branch string, excludeBranches []string) bool {
	for _, excludeBranch := range excludeBranches {
		if strings.Contains(branch, excludeBranch) {
			return true
		}
	}

	return false
}

func cloneOrOpenRepo(authToken string, url string, path string) (*git.Repository, error) {
	var repo *git.Repository

	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Println("Cloning repo...")
		repo, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
			Auth: &http.BasicAuth{
				Username: "token",
				Password: authToken,
			},
			Progress: os.Stdout,
		})
	} else {
		log.Println("Opening existing repo...")
		repo, err = git.PlainOpen(path)
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func checkoutBranch(repo *git.Repository, authContext http.BasicAuth, branchName string, remoteName string) error {
	branchRefName := plumbing.NewBranchReferenceName(branchName)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Checkout(&branchCoOpts)
	if err == nil {
		return nil
	}

	log.Println("Branch not found locally, trying to fetch...")

	mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
	err = fetch(repo, mirrorRemoteBranchRefSpec, authContext, remoteName)
	if err != nil {
		return err
	}

	err = wt.Checkout(&branchCoOpts)
	if err != nil {
		return err
	}

	return nil
}

func fetch(repo *git.Repository, refSpecStr string, authContext http.BasicAuth, remoteName string) error {
	remote, err := repo.Remote(remoteName)
	if err != nil {
		return err
	}

	var refSpecs []gitConfig.RefSpec
	if refSpecStr != "" {
		refSpecs = []gitConfig.RefSpec{gitConfig.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{RefSpecs: refSpecs, Auth: &authContext}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch failed: %v", err)
		}
	}

	return nil
}

func pull(repo *git.Repository, authContext http.BasicAuth, remoteName string, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    remoteName,
		Force:         true,
		Auth:          &authContext,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Print("Branch already up to date")
		} else {
			return err
		}
	}

	return nil
}
