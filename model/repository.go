package model

import (
	"time"
)

type Repository struct {
	Name            string   `gorm:"primaryKey" yaml:"name"`
	URL             string   `yaml:"url"`
	RemoteName      string   `yaml:"remoteName"`
	AuthToken       string   `gorm:"-" yaml:"authToken"`
	ExcludeBranches []string `gorm:"-" yaml:"excludeBranches"`
	LastScan        time.Time
	Branches        []Branch `gorm:"foreignKey:RepositoryName"`
}
