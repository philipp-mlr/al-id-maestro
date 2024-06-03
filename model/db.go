package model

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Connection string
	Config     *gorm.Config
	Database   *gorm.DB
}

func NewDB(connection string, config *gorm.Config) *DB {
	return &DB{
		Connection: connection,
		Config:     config,
	}
}

func (db *DB) Migrate() error {
	err := db.Database.AutoMigrate(&Object{}, &Repository{}, &App{}, &Branch{}, &ClaimHistory{}, &ObjectType{})
	return err
}

func (db *DB) Connect() error {
	database, err := gorm.Open(sqlite.Open(db.Connection), db.Config)

	if err != nil {
		return err
	}

	db.Database = database

	return nil
}

func (db *DB) Create(value interface{}) error {
	tx := db.Database.Create(value)
	return tx.Error
}

func (db *DB) Save(value interface{}) error {
	tx := db.Database.Save(value)
	fmt.Println(tx.Statement.SQL.String())
	return tx.Error
}

func (db *DB) Find(value interface{}, query interface{}) error {
	tx := db.Database.Find(value, query)
	return tx.Error
}
