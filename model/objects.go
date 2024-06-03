package model

type Objects struct {
	Objects []Object
}

func (a Objects) Len() int {
	return len(a.Objects)
}

func (a Objects) Less(i, j int) bool {
	if a.Objects[i].ID != a.Objects[j].ID {
		return a.Objects[i].ID < a.Objects[j].ID
	}
	if a.Objects[i].ObjectType.ID != a.Objects[j].ObjectType.ID {
		return a.Objects[i].ObjectType.ID < a.Objects[j].ObjectType.ID
	}
	return a.Objects[i].Name < a.Objects[j].Name
}

func (a Objects) Swap(i, j int) {
	a.Objects[i], a.Objects[j] = a.Objects[j], a.Objects[i]
}

func (a Objects) BinarySearch(target Object) int {
	return a.binarySearch(target, false)
}

func (a Objects) BinarySearchOmitName(target Object) int {
	return a.binarySearch(target, true)
}

func (a Objects) binarySearch(target Object, omitName bool) int {
	left, right := 0, a.Len()-1
	for left <= right {
		mid := left + (right-left)/2
		comparison := a.compareAlObjects(a.Objects[mid], target, omitName)
		if comparison == 0 {
			return mid
		} else if comparison < 0 {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1 // Not found
}

func (a Objects) compareAlObjects(obj1, obj2 Object, omitName bool) int {
	if obj1.ID < obj2.ID {
		return -1
	} else if obj1.ID > obj2.ID {
		return 1
	}

	if obj1.ObjectType.ID < obj2.ObjectType.ID {
		return -1
	} else if obj1.ObjectType.ID > obj2.ObjectType.ID {
		return 1
	}

	if omitName {
		return 0
	}

	if obj1.Name < obj2.Name {
		return -1
	} else if obj1.Name > obj2.Name {
		return 1
	}

	return 0
}
