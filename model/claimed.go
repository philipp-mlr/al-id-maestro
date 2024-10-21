package model

import "time"

type Claimed struct {
	EntryNo    uint       `db:"entry_no"`
	ID         uint       `db:"id"`
	ObjectType ObjectType `db:"type"`
	InGit      bool       `db:"in_git"`
	Expired    bool       `db:"expired"`
	CreatedAt  string     `db:"created_at"`
}

func NewClaimedObject(id uint, objectType ObjectType) *Claimed {
	return &Claimed{
		ID:         id,
		ObjectType: objectType,
		InGit:      false,
		Expired:    false,
		CreatedAt:  time.Now().Format(time.RFC1123),
	}
}
