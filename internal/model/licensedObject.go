package model

import "sort"

type LicensedObject struct {
	ID         uint
	ObjectType ObjectType
	Used       bool
}

type LicensedObjectList []LicensedObject

func (list *LicensedObjectList) Len() int {
	return len(*list)
}

func (list *LicensedObjectList) Less(i, j int) bool {
	if (*list)[i].ID == (*list)[j].ID {
		return (*list)[i].ObjectType < (*list)[j].ObjectType // Compare by ObjectType if IDs are the same
	}
	return (*list)[i].ID < (*list)[j].ID // Compare by ID first
}

func (list *LicensedObjectList) Swap(i, j int) {
	(*list)[i], (*list)[j] = (*list)[j], (*list)[i]
}

func (list *LicensedObjectList) Sort() {
	sort.Sort(list)
}

func (list *LicensedObjectList) BinarySearch(targetID uint, targetType ObjectType) int {
	low, high := 0, len(*list)-1

	for low <= high {
		mid := low + (high-low)/2

		if (*list)[mid].ID == targetID {
			if (*list)[mid].ObjectType == targetType {
				return mid // Both ID and ObjectType match
			} else if (*list)[mid].ObjectType < targetType {
				low = mid + 1 // Search right half if ObjectType is smaller
			} else {
				high = mid - 1 // Search left half if ObjectType is larger
			}
		} else if (*list)[mid].ID < targetID {
			low = mid + 1 // Search right half if ID is smaller
		} else {
			high = mid - 1 // Search left half if ID is larger
		}
	}

	return -1 // Not found
}

func (list *LicensedObjectList) Filter(objectType ObjectType) LicensedObjectList {
	filteredList := LicensedObjectList{}

	for _, a := range *list {
		if a.ObjectType == objectType {
			filteredList = append(filteredList, LicensedObject{
				ID:         a.ID,
				ObjectType: a.ObjectType,
			})
		}
	}

	return filteredList
}
