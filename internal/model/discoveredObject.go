package model

import (
	"time"

	"github.com/philipp-mlr/al-id-maestro/internal/objectType"
)

type DiscoveredObject struct {
	ID         uint            `db:"id"`
	ObjectType objectType.Type `db:"type"`
	Name       string          `db:"name"`
	AppID      string          `db:"app_id"`
	AppName    string          `db:"app_name"`
	Branch     string          `db:"branch"`
	Repository string          `db:"repository"`
	FilePath   string          `db:"file_path"`
	CommitID   string          `db:"commit_id"`
	CreatedAt  string          `db:"created_at"`
}

func NewFoundObject(id uint, objectType objectType.Type, objectName string, app AppJsonFile, branch Branch, repository string, filePath string) *DiscoveredObject {
	return &DiscoveredObject{
		ID:         id,
		ObjectType: objectType,
		Name:       objectName,
		AppID:      app.ID,
		AppName:    app.Name,
		Branch:     branch.Name,
		Repository: repository,
		FilePath:   filePath,
		CommitID:   branch.CommitID,
		CreatedAt:  time.Now().Format(time.RFC1123),
	}
}
