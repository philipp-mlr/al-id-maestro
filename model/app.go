package model

type App struct {
	ID                   string `json:"id" gorm:"primaryKey"`
	Name                 string `json:"name"`
	BasePath             string
	BranchRepositoryName string `gorm:"primaryKey"`
	BranchName           string `gorm:"primaryKey"`
	Branch               Branch `gorm:"ForeignKey:BranchRepositoryName,BranchName;References:RepositoryName,Name"`
}

func NewApp(branch Branch, appName string, id string, basePath string) *App {
	return &App{
		Branch:   branch,
		BasePath: basePath,
		Name:     appName,
		ID:       id,
	}
}
