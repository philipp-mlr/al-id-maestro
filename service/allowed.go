package service

import (
	"github.com/philipp-mlr/al-id-maestro/model"
)

func NewAllowList(config *model.Config) (*model.AllowedList, error) {
	var allowed model.AllowedList

	for _, c := range config.ConfigIDRanges {
		for i := c.StartID; i <= c.EndID; i++ {
			allowed = append(allowed, model.Allowed{
				ID:         i,
				ObjectType: c.ObjectType,
			})
		}
	}

	return &allowed, nil
}
