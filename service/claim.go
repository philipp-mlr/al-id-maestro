package service

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/model"
)

func ClaimObjectID(db *sqlx.DB, allowed *model.AllowedList, objectType model.ObjectType) (*model.Claimed, error) {
	objectTypeAllowedList := model.FilterAllowedList(allowed, objectType)

	if len(objectTypeAllowedList) == 0 {
		return nil, fmt.Errorf("%s is not configured", objectType)
	}

	model.SortAllowedList(&objectTypeAllowedList)

	// get found objects by type
	// get claimed objects by type
	// eliminate objects in ojbectTypeAllowedList that are in the results of the above two queries

	found, err := selectDistinctFoundByType(db, objectType)
	if err != nil {
		return nil, err
	}

	claimed, err := selectDistinctClaimedByType(db, objectType)
	if err != nil {
		return nil, err
	}

	if len(found) >= len(objectTypeAllowedList) || len(claimed) >= len(objectTypeAllowedList) {
		return nil, fmt.Errorf("No free %s IDs.", objectType)
	}

	for _, foundObject := range found {
		i := model.BinarySearchAllowed(&objectTypeAllowedList, foundObject.ID, foundObject.ObjectType)

		if i != -1 {
			objectTypeAllowedList[i].Used = true
		} else {
			return nil, fmt.Errorf("found object %d of type %s not in allowed list", foundObject.ID, foundObject.ObjectType)
		}
	}

	for _, claimedObject := range claimed {
		i := model.BinarySearchAllowed(&objectTypeAllowedList, claimedObject.ID, claimedObject.ObjectType)

		if i != -1 {
			objectTypeAllowedList[i].Used = true
		} else {
			return nil, fmt.Errorf("found object %d of type %s not in allowed list", claimedObject.ID, claimedObject.ObjectType)
		}
	}

	for _, allowedObject := range objectTypeAllowedList {
		if !allowedObject.Used {
			c := model.NewClaimedObject(allowedObject.ID, allowedObject.ObjectType)

			err = insertClaimedObject(db, *c)
			if err != nil {
				return nil, err
			}

			return c, nil
		}
	}

	return nil, nil
}

func UpdateClaimed(db *sqlx.DB) error {
	err := updateClaimedNotInGit(db)
	if err != nil {
		return err
	}

	err = updateClaimedInGit(db)
	if err != nil {
		return err
	}

	c, err := selectClaimedNotInFound(db)
	if err != nil {
		return err
	}

	for _, claimed := range c {
		isOlder, err := isOlderThan3Days(claimed.CreatedAt)
		if err != nil {
			return err
		}

		if isOlder {
			err = updateClaimedExpired(db, claimed.ID, claimed.ObjectType, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isOlderThan3Days(createdDate string) (bool, error) {
	t, err := time.Parse(time.RFC1123, createdDate)

	if err != nil {
		return false, err
	}

	return time.Since(t).Hours() > 72, nil
}
