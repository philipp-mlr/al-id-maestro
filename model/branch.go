package model

import (
	"time"
)

type Branch struct {
	RepositoryName string     `gorm:"primaryKey"`
	Repository     Repository `gorm:"foreignKey:RepositoryName;references:Name"`
	Name           string     `gorm:"primaryKey"`
	LastCommit     string
	LastScan       time.Time
	Changed        bool `gorm:"-"`
}
