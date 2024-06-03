package service

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/philipp-mlr/al-id-maestro/model"
	"gorm.io/gorm"
)

func ScanRepository(repository *model.Repository, db *model.DB, objects *model.Objects) error {
	log.Println("Processing repository", repository.Name)
	log.Println("Getting remote branches")

	branches, err := getRemoteBranches(repository)
	if err != nil {
		return err
	}

	db.Database.Find(&repository.Branches, model.Branch{RepositoryName: repository.Name})

	log.Println("Cloning or opening repo...")

	path := buildRepoPath(repository.URL)

	repo, err := cloneOrOpenRepo(repository.AuthToken, repository.URL, path)
	if err != nil {
		return err
	}

	for _, branch := range branches {
		log.Println("Processing branch", branch.Name)

		update := false
		for _, b := range repository.Branches {
			if b.RepositoryName == branch.RepositoryName && b.Name == branch.Name && b.LastCommit != branch.LastCommit {
				update = true
				break
			}
		}

		if !update {
			log.Printf("Branch %v is up to date", branch.Name)
		} else {
			err := checkoutBranch(repo, repository.AuthToken, branch.Name)
			if err != nil {
				return err
			}

			err = traverseRepo(db, path, objects, branch)
			if err != nil {
				return err
			}
		}

		branch.LastScan = time.Now()
		db.Save(branch)
	}

	repository.LastScan = time.Now()
	db.Save(repository)
	db.Save(branches)

	return nil
}

func buildRepoPath(url string) string {
	path := strings.Replace(url, "https://", "", 1)
	path = strings.Replace(path, "/", "_", -1)
	path = "./repo/" + path

	return path
}

func traverseRepo(db *model.DB, dir string, objects *model.Objects, branch model.Branch) error {
	total := new(int)
	*total = 0

	// TODO: Delete all objects from the current branch
	// TODO: Read from db
	queryObject := new(model.Object)
	tx := db.Database.First(queryObject, model.Object{App: model.App{Branch: branch}})

	fillInitial := tx.Error == gorm.ErrRecordNotFound

	var appData []model.App
	appData, err := readApp(dir, branch)
	if err != nil {
		return err
	}
	db.Save(appData)

	log.Printf("Found %v apps in branch %v ", len(appData), branch.Name)

	for _, app := range appData {
		log.Println("Checking app: ", app.Name)
		err = filepath.Walk(app.BasePath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			// Check if the file has a ".al" extension
			if !info.IsDir() && filepath.Ext(info.Name()) == ".al" {
				if err := addMatchesAndSort(db, path, total, objects, fillInitial, app); err != nil {
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

	if fillInitial {
		sort.Sort(objects)
	}

	return nil
}

func readApp(dir string, branch model.Branch) ([]model.App, error) {
	appData := []model.App{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a ".al" extension
		if !info.IsDir() && info.Name() == "app.json" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var app model.App
			err = json.Unmarshal(content, &app)
			if err != nil {
				return err
			}

			app.Branch = branch

			//appInfo.BasePath = path
			p := strings.Replace(path, "app.json", "", 1)
			p = strings.ReplaceAll(p, "\\", "/")
			p = "./" + p
			app.BasePath = p

			appData = append(appData, app)

		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return appData, nil
}

func addMatchesAndSort(db *model.DB, filePath string, total *int, objects *model.Objects, fillInitial bool, app model.App) error {
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

			objectType := strings.ToLower(matches[1])
			objectName := strings.ToLower(matches[3])
			id, err := strconv.Atoi(matches[2])
			if err != nil {
				log.Println("Error converting ID to int for file", filePath)
				continue
			}

			newObjType := model.NewObjectType(objectType)
			newObj := *model.NewAlObject(uint(id), *newObjType, objectName, app)

			if fillInitial {
				objects.Objects = append(objects.Objects, newObj)
				db.Save(&newObj)
			} else {
				index := objects.BinarySearch(newObj)
				if index == -1 {
					objects.Objects = append(objects.Objects, newObj)
					db.Save(&newObj)
					sort.Sort(objects)
					log.Printf("Found new object: ID: %d, Type: %s, Name: %s\n", id, objectType, objectName)
				}
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
