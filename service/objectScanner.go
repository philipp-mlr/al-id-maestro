package service

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/model"
)

func Sync(db *sqlx.DB, config *model.Config) error {
	for _, config := range config.RemoteConfiguration {
		log.Printf("Processing repository %s", config.RepositoryURL)

		err := scanRepository(&config, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanRepository(config *model.RemoteConfiguration, db *sqlx.DB) error {
	branches, err := getRemoteBranches(config)
	if err != nil {
		return err
	}

	deleteRemovedBranches(db, branches, config.RepositoryName)

	for i, branch := range branches {
		log.Printf("%d/%d | syncing branch %s", i+1, len(branches), branch.Name)

		commitID := getLastCommitID(db, branch.Name, config.RepositoryName)

		update := commitID == "" || commitID != branch.CommitID

		if !update {
			log.Print("Nothing to update\n\n")
			continue
		}

		err = checkout(config.Git, config.AuthContext, branch.Name, config.RemoteName)
		if err != nil {
			log.Printf("Error durching checkout on branch %s %s", branch.Name, err)
			continue
		}

		err = pull(config.Git, config.AuthContext, config.RemoteName, branch.Name)
		if err != nil {
			log.Printf("Error pulling branch %s %s", branch.Name, err)
			continue
		}

		err = traverseRepo(db, config, branch)
		if err != nil {
			log.Printf("Error traversing repository %s %s", config.RepositoryURL, err)
			continue
		}

		log.Print("Done\n\n")
	}

	return nil
}

func traverseRepo(db *sqlx.DB, config *model.RemoteConfiguration, branch model.Branch) error {
	appData, err := getAppJsonFiles(config)
	if err != nil {
		return err
	}

	for _, app := range appData {
		deleteFoundObjects(db, app.ID, branch.Name, config.RepositoryName)

		err := Walk(config, func(f *object.File) error {
			if !f.Mode.IsFile() {
				return nil
			}

			if !strings.Contains(f.Name, app.BasePath) {
				return nil
			}

			if filepath.Ext(f.Name) != ".al" {
				return nil
			}

			lines, err := f.Lines()
			if err != nil {
				return err
			}

			if err := findAndInsertMatches(db, &lines, f.Name, app, branch, config.RepositoryName); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func getAppJsonFiles(config *model.RemoteConfiguration) ([]model.AppJsonFile, error) {
	apps := []model.AppJsonFile{}

	err := Walk(config, func(f *object.File) error {
		if !f.Mode.IsFile() {
			return nil
		}

		if !strings.Contains(f.Name, "app.json") {
			return nil
		}

		app := model.AppJsonFile{}

		content, err := f.Contents()
		if err != nil {
			return err
		}

		err = json.Unmarshal([]byte(content), &app)
		if err != nil {
			return err
		}

		p := strings.Replace(f.Name, "app.json", "", 1)
		app.BasePath = p

		apps = append(apps, app)

		return nil
	})

	return apps, err
}

func findAndInsertMatches(db *sqlx.DB, lines *[]string, fileName string, app model.AppJsonFile, branch model.Branch, repository string) error {
	pattern := regexp.MustCompile(`^(\w+) (\d{1,6}) "?([^"]*)"?$`)

	for _, line := range *lines {

		matches := pattern.FindStringSubmatch(line)

		if len(matches) == 4 {
			// matches[0] is the full match, matches[1], matches[2], matches[3] are the capture groups

			objectType := matches[1]
			objectName := matches[3]
			id, err := strconv.Atoi(matches[2])
			if err != nil {
				return fmt.Errorf("failed converting the object Id to type int for file %s", fileName)
			}

			foundObject := *model.NewFoundObject(uint(id), model.MapObjectType(objectType), objectName, app, branch, repository, fileName)

			err = insertFoundObject(db, foundObject)

			return err
		}
	}

	return nil
}
