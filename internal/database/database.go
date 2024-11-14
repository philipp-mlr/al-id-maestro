package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
)

func InitDB(databaseFileName string) (*sqlx.DB, error) {
	db, err := open(databaseFileName)
	if err != nil {
		return nil, err
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func open(databaseFileName string) (*sqlx.DB, error) {
	var db *sqlx.DB

	// exactly the same as the built-in
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("../data/%s.db", databaseFileName))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w %s", err, databaseFileName)
	}

	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sqlx.DB) error {
	// claimed definition
	claimedchema := `
		CREATE TABLE IF NOT EXISTS claimed (
				entry_no INTEGER PRIMARY KEY AUTOINCREMENT,
				id INTEGER NOT NULL,
				"type" INTEGER NOT NULL,
				in_git BOOLEAN NOT NULL,
				expired BOOLEAN NOT NULL,
				created_at TEXT NOT NULL
				);
			`
	_, err := db.Exec(claimedchema)
	if err != nil {
		return err
	}

	// found definition
	foundObjectSchema := `
		CREATE TABLE IF NOT EXISTS found (
				id INTEGER NOT NULL,
				"type" TEXT NOT NULL,
				name TEXT NOT NULL,
				app_id TEXT NOT NULL,
				app_name TEXT NOT NULL,
				branch TEXT NOT NULL,
				repository TEXT NOT NULL,
				file_path TEXT NOT NULL,
				commit_id TEXT NOT NULL,
				created_at TEXT NOT NULL,
				CONSTRAINT found_pk PRIMARY KEY (id, "type", name, app_id, branch, repository)
				);
	`
	_, err = db.Exec(foundObjectSchema)
	if err != nil {
		return err
	}

	// allowed definition
	allowedSchema := `
		CREATE TABLE IF NOT EXISTS allowed (
				id INTEGER NOT NULL,
				"type" TEXT NOT NULL,
				CONSTRAINT allowed_pk PRIMARY KEY (id, "type")
				);
	`
	_, err = db.Exec(allowedSchema)
	if err != nil {
		return err
	}

	return nil
}

func InsertFoundObject(db *sqlx.DB, foundObject model.DiscoveredObject) error {
	stmt := `
		INSERT INTO 
			found (id, type, name, app_id, app_name, branch, repository, file_path, commit_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(stmt,
		foundObject.ID,
		foundObject.ObjectType,
		foundObject.Name,
		foundObject.AppID,
		foundObject.AppName,
		foundObject.Branch,
		foundObject.Repository,
		foundObject.FilePath,
		foundObject.CommitID,
		foundObject.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDiscoveredObjectsByBranchAndRepo(db *sqlx.DB, branch string, repository string) error {
	stmt := `
		DELETE FROM found
			WHERE branch = ?
			AND repository = ?
	`

	_, err := db.Exec(stmt, branch, repository)
	if err != nil {
		return err
	}

	return nil
}

func GetLastCommitID(db *sqlx.DB, branch string, repository string) string {
	stmt := `
		SELECT commit_id
			FROM found
			WHERE branch = ?
			AND repository = ?
			ORDER BY created_at DESC
			LIMIT 1
	`

	var commitID string
	err := db.Get(&commitID, stmt, branch, repository)
	if err != nil {
		return ""
	}

	return commitID
}

func InsertClaimedObject(db *sqlx.DB, claimedObject model.ClaimedObject) error {
	stmt := `
		INSERT INTO claimed (id, type, in_git, expired, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(stmt,
		claimedObject.ID,
		claimedObject.ObjectType,
		claimedObject.InGit,
		claimedObject.Expired,
		claimedObject.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func SelectClaimedObjects(db *sqlx.DB, offset uint64) ([]model.ClaimedObject, error) {
	stmt := `
		SELECT *
			FROM claimed
			LIMIT 50 OFFSET 50 * ?
	`

	claimed := []model.ClaimedObject{}
	err := db.Select(&claimed, stmt, offset)
	if err != nil {
		return nil, err
	}

	return claimed, nil
}

func SelectDuplicates(db *sqlx.DB, offset uint64) ([]model.DiscoveredObject, error) {
	stmt := `
		SELECT DISTINCT      f1.id,
				f1.type,
				f1.name,
				f1.app_name,
				f1.app_id,
				f1.repository,
				f1.file_path
			FROM "found" f1
			JOIN "found" f2 ON f1.id = f2.id
			AND f1.type = f2.type
			AND f1.file_path != f2.file_path
			LIMIT 50
		OFFSET 50 * ?
	`

	duplicates := []model.DiscoveredObject{}
	err := db.Select(&duplicates, stmt, offset)
	if err != nil {
		return nil, err
	}

	return duplicates, nil
}

func SelectDiscoveredObjects(db *sqlx.DB, offset uint64) ([]model.DiscoveredObject, error) {
	stmt := `
		SELECT id, "type", name, app_id, app_name, branch, repository, file_path, commit_id, created_at
		FROM "found"
		ORDER BY id, type ASC
		LIMIT 50
		OFFSET 50 * ?;
	`

	found := []model.DiscoveredObject{}
	err := db.Select(&found, stmt, offset)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func SelectClaimedObjectsNotFoundInDiscoveredObjects(db *sqlx.DB) ([]model.ClaimedObject, error) {
	stmt := `
		SELECT  c.id,
				c."type",
				c.in_git,
				c.expired,
				c.created_at
		FROM      claimed c
		LEFT JOIN found f ON c.id = f.id
		AND       c.type = f.type
		WHERE     f.id IS NULL
		AND       f.type IS NULL
		AND	      c.expired = false
	`

	claimed := []model.ClaimedObject{}
	err := db.Select(&claimed, stmt)
	if err != nil {
		return nil, err
	}

	return claimed, nil
}

func UpdateClaimedObjectsSetExpired(db *sqlx.DB, id uint, objectType model.ObjectType, value bool) error {
	stmt := `
	UPDATE claimed
		SET expired = ?
		WHERE id = ?
		AND type = ?
	`

	_, err := db.Exec(stmt, value, id, objectType)
	if err != nil {
		return err
	}

	return nil
}

func UpdateClaimedObjectsNotFoundDiscoveredObjects(db *sqlx.DB) error {
	stmt := `
	UPDATE claimed
	SET in_git = false
	WHERE NOT EXISTS (
		SELECT 1
		FROM found f
		WHERE claimed.id = f.id
		AND claimed.type = f.type
	);
	`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func UpdateClaimedObjectsFoundInDiscoveredObjects(db *sqlx.DB) error {
	stmt := `
	UPDATE claimed
	SET in_git = true
	WHERE EXISTS (
		SELECT 1
		FROM found f
		WHERE claimed.id = f.id
		AND claimed.type = f.type
	);
	`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func SelectDistinctDiscoveredObjectsByType(db *sqlx.DB, objectType model.ObjectType) ([]model.DiscoveredObject, error) {
	stmt := `
	SELECT DISTINCT id, type
		FROM found
		WHERE type = ?
	`

	found := []model.DiscoveredObject{}
	err := db.Select(&found, stmt, objectType)
	if err != nil {
		return nil, err
	}

	return found, nil
}

func SelectDistinctClaimedObjectsByType(db *sqlx.DB, objectType model.ObjectType) ([]model.ClaimedObject, error) {
	stmt := `
	SELECT DISTINCT id, type
		FROM claimed
		WHERE type = ?
		AND expired = false
	`

	claimed := []model.ClaimedObject{}
	err := db.Select(&claimed, stmt, objectType)
	if err != nil {
		return nil, err
	}

	return claimed, nil
}

func DeleteDiscoveredObjectNotInBranches(db *sqlx.DB, branches []model.Branch, repository string) (int64, error) {
	stmt := `
	DELETE FROM found
		WHERE branch NOT IN (?)
	`

	branchNames := []string{}
	for _, b := range branches {
		branchNames = append(branchNames, b.Name)
	}

	query, args, err := sqlx.In(stmt, branchNames)
	if err != nil {
		return 0, err
	}

	query += "AND repository = ?"

	args = append(args, repository)

	r, err := db.Exec(query, args...)

	count, _ := r.RowsAffected()

	if err != nil {
		return count, err
	}

	return count, nil
}

func DeleteDiscoveredObjectsNotInRepositories(db *sqlx.DB, repositories []string) error {
	stmt := `
	DELETE FROM found
		WHERE repository NOT IN (?)
	`

	query, args, err := sqlx.In(stmt, repositories)
	if err != nil {
		return err
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
