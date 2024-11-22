package model

import (
	"time"

	"github.com/philipp-mlr/al-id-maestro/internal/objectType"
)

type Source string

const (
	Empty Source = ""
	GUI   Source = "gui"
	API   Source = "api"
)

type ClaimedObject struct {
	EntryNo    uint            `db:"entry_no"`
	ID         uint            `db:"id"`
	ObjectType objectType.Type `db:"type"`
	InGit      bool            `db:"in_git"`
	Expired    bool            `db:"expired"`
	Source     Source          `db:"source"`
	CreatedAt  string          `db:"created_at"`
}

func NewClaimedObject(id uint, objectType objectType.Type, source Source) *ClaimedObject {
	return &ClaimedObject{
		ID:         id,
		ObjectType: objectType,
		InGit:      false,
		Expired:    false,
		Source:     source,
		CreatedAt:  time.Now().Format(time.RFC1123),
	}
}
