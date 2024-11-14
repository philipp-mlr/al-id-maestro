package licensedList

import (
	"github.com/philipp-mlr/al-id-maestro/config"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
)

func NewLicensedObjectList(config *config.Config) (*model.LicensedObjectList, error) {
	var allowed model.LicensedObjectList

	for _, c := range config.ConfigIDRanges {
		for i := c.StartID; i <= c.EndID; i++ {
			allowed = append(allowed, model.LicensedObject{
				ID:         i,
				ObjectType: c.ObjectType,
			})
		}
	}

	return &allowed, nil
}
