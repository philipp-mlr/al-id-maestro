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

		err = scanAndInsertRepositoryFiles(db, config, branch)
		if err != nil {
			log.Printf("Error traversing repository %s %s", config.RepositoryURL, err)
			continue
		}

		log.Print("Done\n\n")
	}

	return nil
}

func scanAndInsertRepositoryFiles(db *sqlx.DB, config *model.RemoteConfiguration, branch model.Branch) error {
	var found []model.Found
	var apps []model.AppJsonFile

	countALFiles := 0

	err := Walk(config, func(f *object.File) error {
		if !f.Mode.IsFile() {
			return nil
		}

		if strings.Contains(f.Name, "app.json") {
			content, err := f.Contents()
			if err != nil {
				return err
			}

			app := model.AppJsonFile{}

			err = json.Unmarshal([]byte(content), &app)
			if err != nil {
				return err
			}

			p := strings.Replace(f.Name, "app.json", "", 1)
			app.BasePath = p

			apps = append(apps, app)

			return nil
		}

		if filepath.Ext(f.Name) != ".al" {
			return nil
		}

		countALFiles++

		lines, err := f.Lines()
		if err != nil {
			return err
		}

		objectType, objectName, objectId, err := findMatches(&lines, f.Name)
		if err != nil {
			log.Printf("Error finding matches: %s", err)
			return nil
		}

		newFoundObject := model.NewFoundObject(
			uint(objectId),
			objectType,
			objectName,
			model.AppJsonFile{},
			branch,
			config.RepositoryName,
			f.Name)

		found = append(found, *newFoundObject)

		return nil
	})

	if err != nil {
		return err
	}

	for i, f := range found {
		for _, app := range apps {
			if strings.Contains(f.FilePath, app.BasePath) {
				found[i].AppID = app.ID
				found[i].AppName = app.Name
			}
		}
	}

	if err = deleteFoundObjectsByBranchAndRepo(db, branch.Name, config.RepositoryName); err != nil {
		return err
	}

	for _, f := range found {
		err = insertFoundObject(db, f)
		if err != nil {
			return fmt.Errorf("error inserting found object: %s %v", err, f)
		}
	}

	log.Printf("Found %s AL files in branch %s\n", strconv.Itoa(countALFiles), branch.Name)
	log.Printf("Found %s objects in branch %s\n", strconv.Itoa(len(found)), branch.Name)
	log.Printf("Found %s apps in branch %s\n", strconv.Itoa(len(apps)), branch.Name)

	return nil
}

func findMatches(lines *[]string, file string) (model.ObjectType, string, int, error) {
	pattern := regexp.MustCompile(`^(\w+) (\d{1,6}) "?"?([^"]*)?"?`)

	for _, line := range *lines {
		matches := pattern.FindStringSubmatch(line)

		if len(matches) == 4 {
			// matches[0] is the full match, matches[1], matches[2], matches[3] are the capture groups

			objectType := model.MapObjectType(matches[1])
			if objectType == model.Unknown {
				return model.Unknown, "", 0, fmt.Errorf("unknown object type %s for file %s", matches[1], file)
			}

			id, err := strconv.Atoi(matches[2])
			if err != nil {
				return model.Unknown, "", 0, fmt.Errorf("failed converting the object Id to type int for file %s", file)
			}

			objectName := matches[3]
			if objectName == "" {
				return model.Unknown, "", 0, fmt.Errorf("object name is empty for file %s", file)
			}

			return objectType, objectName, id, nil
		}
	}

	return model.Unknown, "", 0, fmt.Errorf("no object definition found in file %s", file)
}
