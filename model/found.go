package model

import (
	"time"
)

type Found struct {
	ID         uint       `db:"id"`
	ObjectType ObjectType `db:"type"`
	Name       string     `db:"name"`
	AppID      string     `db:"app_id"`
	AppName    string     `db:"app_name"`
	Branch     string     `db:"branch"`
	Repository string     `db:"repository"`
	FilePath   string     `db:"file_path"`
	CommitID   string     `db:"commit_id"`
	CreatedAt  string     `db:"created_at"`
}

func NewFoundObject(id uint, objectType ObjectType, objectName string, appId string, appName string, branch string, repository string, filePath string, commitID string) *Found {
	return &Found{
		ID:         id,
		ObjectType: objectType,
		Name:       objectName,
		AppID:      appId,
		AppName:    appName,
		Branch:     branch,
		Repository: repository,
		FilePath:   filePath,
		CommitID:   commitID,
		CreatedAt:  time.Now().Format(time.RFC1123),
	}
}
