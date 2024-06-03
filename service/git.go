package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/philipp-mlr/al-id-maestro/model"
)

func getRemoteBranches(repository *model.Repository) ([]model.Branch, error) {
	branches := []model.Branch{}

	remote := git.NewRemote(
		memory.NewStorage(),
		&config.RemoteConfig{
			Name: repository.RemoteName,
			URLs: []string{repository.URL},
		})

	rfs, err := remote.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: repository.Name,
			Password: repository.AuthToken,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, rf := range rfs {
		if rf.Name().IsRemote() || rf.Name().IsBranch() {
			if !IsExcludeBranch(rf.Name().Short(), repository.ExcludeBranches) {
				branches = append(branches, model.Branch{
					Name:           rf.Name().Short(),
					RepositoryName: repository.Name,
					LastCommit:     rf.Hash().String(),
				})
			}
		}
	}

	return branches, nil
}

func IsExcludeBranch(branch string, excludeBranches []string) bool {
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

func checkoutBranch(repo *git.Repository, authToken string, branchName string) error {
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

	if err != nil {
		log.Println(err)

		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)
		err = fetchOrigin(repo, mirrorRemoteBranchRefSpec, &http.BasicAuth{
			Username: "token",
			Password: authToken,
		})
		if err != nil {
			return err
		}

		err = wt.Checkout(&branchCoOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchOrigin(repo *git.Repository, refSpecStr string, authContext *http.BasicAuth) error {
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
		Auth:     authContext,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}
