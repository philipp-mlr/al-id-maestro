package service

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/model"
	"gopkg.in/yaml.v3"
)

func NewConfig(db *sqlx.DB) (*model.Config, error) {
	c := model.Config{RemoteConfiguration: []model.RemoteConfiguration{}}

	file, err := os.ReadFile("./data/config.yml")
	if err != nil {
		return nil, fmt.Errorf("error opening config.yml: %v", err)
	}

	err = yaml.Unmarshal([]byte(file), &c)
	if err != nil {
		return nil, fmt.Errorf("error reading config.yml: %v", err)
	}

	err = validateIDRanges(&c.ConfigIDRanges)
	if err != nil {
		return nil, err
	}

	var repositories []string
	for _, r := range c.RemoteConfiguration {
		repositories = append(repositories, r.RepositoryName)
	}

	if len(repositories) == 0 {
		return nil, fmt.Errorf("error: no repositories defined in config.yml")
	}

	err = deleteMissingRepositories(db, repositories)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func validateIDRanges(idRangeConfig *[]model.ConfigIDRange) error {
	if len(*idRangeConfig) != len(model.GetObjectTypes()) {
		return fmt.Errorf("error: object type %s is not configured", model.GetObjectTypes())
	}

	for i, c := range *idRangeConfig {
		(*idRangeConfig)[i].ObjectType = model.MapObjectType(string(c.ObjectType))

		if c.ObjectType == "" {
			return fmt.Errorf("error: configuration invalid object type %v", c.ObjectType)
		}

		if c.StartID <= 0 || c.EndID <= 0 {
			return fmt.Errorf("error: configuration id ranges may not be negative or zero")
		}

		if c.StartID > c.EndID {
			return fmt.Errorf("error: configuration invalid ID range for object type %v: from %v to %v", c.ObjectType, c.StartID, c.EndID)
		}

		// check overlapping ranges
		for j, c2 := range *idRangeConfig {
			if i == j {
				continue
			}

			if c.ObjectType != c2.ObjectType {
				continue
			}

			if c.StartID >= c2.StartID && c.StartID <= c2.EndID {
				return fmt.Errorf("error: configuration overlapping ID ranges for object type %v: %v-%v and %v-%v", c.ObjectType, c.StartID, c.EndID, c2.StartID, c2.EndID)
			}

			if c.EndID >= c2.StartID && c.EndID <= c2.EndID {
				return fmt.Errorf("error: configuration overlapping ID ranges for object type %v: %v-%v and %v-%v", c.ObjectType, c.StartID, c.EndID, c2.StartID, c2.EndID)
			}
		}
	}

	return nil
}
