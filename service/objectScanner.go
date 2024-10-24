package service

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/philipp-mlr/al-id-maestro/model"
)

func Scan(db *sqlx.DB, config *model.Config) error {
	for _, config := range config.RemoteConfiguration {
		log.Println("Scanning repository ", config.RepositoryURL)
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

	path := buildRepoPath(config.RepositoryURL)

	repo, err := cloneOrOpenRepo(config.GithubAuthToken, config.RepositoryURL, path)
	if err != nil {
		log.Printf("Error cloning or opening repository %s %s", config.RepositoryURL, err)
		return err
	}

	for i, branch := range branches {
		log.Printf("Processing branch %s | %d/%d", branch.Name, i+1, len(branches))

		commitID := getLastCommitID(db, branch.Name, config.RepositoryName)

		update := commitID == "" || commitID != branch.CommitID

		if !update {
			log.Printf("Branch %v is up to date", branch.Name)
			continue
		}

		authContext := http.BasicAuth{
			Username: "token",
			Password: config.GithubAuthToken,
		}

		err = checkout(repo, authContext, branch.Name, config.RemoteName)
		if err != nil {
			log.Printf("Error checking out branch %s %s", branch.Name, err)
			continue
		}

		err = pull(repo, authContext, config.RemoteName, branch.Name)
		if err != nil {
			log.Printf("Error pulling branch %s %s", branch.Name, err)
			continue
		}

		err = traverseRepo(db, path, branch.Name, config.RepositoryName, branch.CommitID)
		if err != nil {
			log.Printf("Error traversing repository %s %s", config.RepositoryURL, err)
			continue
		}
	}

	return nil
}

func buildRepoPath(url string) string {
	path := strings.Replace(url, "https://", "", 1)
	path = strings.Replace(path, "/", "_", -1)
	path = "./data/repo/" + path

	return path
}

func traverseRepo(db *sqlx.DB, dir string, branch string, repository string, commitID string) error {
	total := new(int)
	*total = 0

	appData, err := getAppJsonFiles(dir)
	if err != nil {
		return err
	}

	log.Printf("Found %v apps in branch %v ", len(appData), branch)

	for _, app := range appData {
		deleteFoundObjects(db, app.ID, branch, repository)
		log.Println("Checking app: ", app.Name)
		err = filepath.Walk(app.BasePath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			// Check if the file has a ".al" extension
			if !info.IsDir() && filepath.Ext(info.Name()) == ".al" {
				if err := findAndInsertMatches(db, path, total, app.ID, app.Name, branch, repository, commitID); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	log.Println("Total objects found: ", *total)

	return nil
}

func getAppJsonFiles(dir string) ([]model.AppJsonFile, error) {
	apps := []model.AppJsonFile{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "app.json" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			app := model.AppJsonFile{}

			err = json.Unmarshal(content, &app)
			if err != nil {
				return err
			}

			p := strings.Replace(path, "app.json", "", 1)
			p = strings.ReplaceAll(p, "\\", "/")
			p = "./" + p
			app.BasePath = p

			apps = append(apps, app)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func findAndInsertMatches(db *sqlx.DB, filePath string, total *int, appId string, appName string, branch string, repository string, commitID string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	pattern := regexp.MustCompile(`^(\w+) (\d{1,6}) "?([^"]*)"?$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := pattern.FindStringSubmatch(line)
		if len(matches) == 4 {
			// matches[0] is the full match, matches[1], matches[2], matches[3] are the capture groups

			objectType := matches[1]
			objectName := matches[3]
			id, err := strconv.Atoi(matches[2])
			if err != nil {
				log.Println("Error converting ID to int for file", filePath)
				continue
			}

			foundObject := *model.NewFoundObject(uint(id), model.MapObjectType(objectType), objectName, appId, appName, branch, repository, filePath, commitID)

			err = insertFoundObject(db, foundObject)

			if err == sqlite3.ErrConstraintUnique {
				log.Printf("Object %v already exists in the database", foundObject.Name)
			} else if err != nil {
				return err
			}

			*total = *total + 1
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
