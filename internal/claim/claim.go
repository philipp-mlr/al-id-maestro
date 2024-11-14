package claim

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/philipp-mlr/al-id-maestro/internal/database"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
)

var claimLock sync.Mutex

func ClaimObjectID(db *sqlx.DB, allowedList *model.LicensedObjectList, objectType model.ObjectType) (*model.ClaimedObject, error) {
	claimLock.Lock()
	defer claimLock.Unlock()

	objectTypeAllowedList := allowedList.Filter(objectType)

	if len(objectTypeAllowedList) == 0 {
		return nil, fmt.Errorf("%s is not configured", objectType)
	}

	objectTypeAllowedList.Sort()

	found, err := database.SelectDistinctDiscoveredObjectsByType(db, objectType)
	if err != nil {
		return nil, err
	}

	claimed, err := database.SelectDistinctClaimedObjectsByType(db, objectType)
	if err != nil {
		return nil, err
	}

	if len(found) >= len(objectTypeAllowedList) || len(claimed) >= len(objectTypeAllowedList) {
		return nil, fmt.Errorf("no free %s IDs", objectType)
	}

	for _, foundObject := range found {
		i := objectTypeAllowedList.BinarySearch(foundObject.ID, foundObject.ObjectType)

		if i != -1 {
			objectTypeAllowedList[i].Used = true
		} else {
			return nil, fmt.Errorf("found object %d of type %s not in allowed list", foundObject.ID, foundObject.ObjectType)
		}
	}

	for _, claimedObject := range claimed {
		i := objectTypeAllowedList.BinarySearch(claimedObject.ID, claimedObject.ObjectType)

		if i != -1 {
			objectTypeAllowedList[i].Used = true
		} else {
			return nil, fmt.Errorf("found object %d of type %s not in allowed list", claimedObject.ID, claimedObject.ObjectType)
		}
	}

	for _, allowedObject := range objectTypeAllowedList {
		if !allowedObject.Used {
			c := model.NewClaimedObject(allowedObject.ID, allowedObject.ObjectType)

			err = database.InsertClaimedObject(db, *c)
			if err != nil {
				return nil, err
			}

			return c, nil
		}
	}

	return nil, nil
}

func UpdateClaimed(db *sqlx.DB) error {
	claimLock.Lock()
	defer claimLock.Unlock()

	err := database.UpdateClaimedObjectsNotFoundDiscoveredObjects(db)
	if err != nil {
		return err
	}

	err = database.UpdateClaimedObjectsFoundInDiscoveredObjects(db)
	if err != nil {
		return err
	}

	c, err := database.SelectClaimedObjectsNotFoundInDiscoveredObjects(db)
	if err != nil {
		return err
	}

	for _, claimed := range c {
		isOlder, err := isOlderThan3Days(claimed.CreatedAt)
		if err != nil {
			return err
		}

		if isOlder {
			err = database.UpdateClaimedObjectsSetExpired(db, claimed.ID, claimed.ObjectType, true)
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
