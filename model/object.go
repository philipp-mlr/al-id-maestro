package model

type Object struct {
	// gorm.Model
	ID                   uint       `gorm:"primarykey;autoIncrement:false"`
	Name                 string     `gorm:"primarykey"`
	ObjectTypeID         uint       `gorm:"primarykey;autoIncrement:false"`
	ObjectType           ObjectType `gorm:"foreignKey:ObjectTypeID;references:ID"`
	AppID                string
	BranchRepositoryName string
	BranchName           string
	App                  App `gorm:"foreignKey:AppID,BranchRepositoryName,BranchName;references:ID,BranchRepositoryName,BranchName"`
}

func NewAlObject(id uint, objectType ObjectType, objectName string, app App) *Object {
	return &Object{
		//Model:      gorm.Model{ID: id},
		ID:           id,
		ObjectTypeID: objectType.ID,
		ObjectType:   objectType,
		Name:         objectName,
		App:          app,
		AppID:        app.ID,
	}
}
