package model

import (
	"time"
)

type Repository struct {
	Name            string   `yaml:"name"`
	URL             string   `yaml:"url"`
	RemoteName      string   `yaml:"remoteName"`
	AuthToken       string   `yaml:"authToken"`
	ExcludeBranches []string `yaml:"excludeBranches"`
	LastScan        time.Time
	Branches        []Branch
}
