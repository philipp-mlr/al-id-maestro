package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/philipp-mlr/al-id-maestro/internal/model"
	"github.com/philipp-mlr/al-id-maestro/internal/objectType"
)

func InitDB(databaseFileName string) (*sqlx.DB, error) {
	db, err := open(databaseFileName)
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		log.Println("Migration failed: ", err)
	}

	err = createSchema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func open(databaseFileName string) (*sqlx.DB, error) {
	var db *sqlx.DB

	file := fmt.Sprintf("./data/%s.db", databaseFileName)
	log.Println("Opening database file: ", file)

	// exactly the same as the built-in
	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(db *sqlx.DB) error {
	// claimed object definition
	claimedObjectSchema := `
		CREATE TABLE IF NOT EXISTS claimed_object (
				entry_no INTEGER PRIMARY KEY AUTOINCREMENT,
				id INTEGER NOT NULL,
				"type" INTEGER NOT NULL,
				in_git BOOLEAN NOT NULL,
				expired BOOLEAN NOT NULL,
				source TEXT NOT NULL,
				created_at TEXT NOT NULL
				);
			`
	_, err := db.Exec(claimedObjectSchema)
	if err != nil {
		return err
	}

	// discovered object definition
	discoveredObjectSchema := `
		CREATE TABLE IF NOT EXISTS discovered_object (
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
	_, err = db.Exec(discoveredObjectSchema)
	if err != nil {
		return err
	}

	return nil
}

func migrate(db *sqlx.DB) error {
	_, err := db.Exec(`ALTER TABLE claimed RENAME TO claimed_object`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`ALTER TABLE found RENAME TO discovered_object`)
	if err != nil {
		return err
	}

	return nil
}

func InsertDiscoveredObject(db *sqlx.DB, discoveredObject model.DiscoveredObject) error {
	stmt := `
		INSERT INTO 
			discovered_object (id, type, name, app_id, app_name, branch, repository, file_path, commit_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(stmt,
		discoveredObject.ID,
		discoveredObject.ObjectType,
		discoveredObject.Name,
		discoveredObject.AppID,
		discoveredObject.AppName,
		discoveredObject.Branch,
		discoveredObject.Repository,
		discoveredObject.FilePath,
		discoveredObject.CommitID,
		discoveredObject.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDiscoveredObjectsByBranchAndRepo(db *sqlx.DB, branch string, repository string) error {
	stmt := `
		DELETE FROM discovered_object
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
			FROM discovered_object
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
		INSERT INTO claimed_object (id, type, in_git, expired, source, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(stmt,
		claimedObject.ID,
		claimedObject.ObjectType,
		claimedObject.InGit,
		claimedObject.Expired,
		claimedObject.Source,
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
			FROM claimed_object
			LIMIT 50 OFFSET 50 * ?
	`

	claimedObject := []model.ClaimedObject{}
	err := db.Select(&claimedObject, stmt, offset)
	if err != nil {
		return nil, err
	}

	return claimedObject, nil
}

func SelectDuplicates(db *sqlx.DB, offset uint64) ([]model.DiscoveredObject, error) {
	stmt := `
		WITH ranked_dupes AS (
			SELECT 
				dupes.*,
				do.branch,
				ROW_NUMBER() OVER (PARTITION BY dupes.id, dupes.type, dupes.name, dupes.app_id, dupes.repository, dupes.file_path ORDER BY do.branch) AS rn
			FROM (
				SELECT DISTINCT 
					do1.id,
					do1.type,
					do1.name,
					do1.app_id,
					do1.app_name,
					do1.repository,
					do1.file_path
				FROM "discovered_object" do1
				JOIN "discovered_object" do2 
					ON do1.id = do2.id
				AND do1.type = do2.type
				AND do1.file_path != do2.file_path
			) AS dupes
			JOIN "discovered_object" do
				ON do.id = dupes.id
			AND do.type = dupes.type
			AND do.name = dupes.name
			AND do.app_id = dupes.app_id
			AND do.repository = dupes.repository
			AND do.file_path = dupes.file_path
		)
		SELECT 
			ranked_dupes.id, 
			ranked_dupes.type,
			ranked_dupes.name,
			ranked_dupes.app_id,
			ranked_dupes.app_name,
			ranked_dupes.repository,
			ranked_dupes.branch,
			ranked_dupes.file_path
		FROM ranked_dupes
		WHERE rn = 1
		ORDER BY id ASC
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
		FROM "discovered_object"
		ORDER BY id, type ASC
		LIMIT 50
		OFFSET 50 * ?;
	`

	discoveredObject := []model.DiscoveredObject{}
	err := db.Select(&discoveredObject, stmt, offset)
	if err != nil {
		return nil, err
	}

	return discoveredObject, nil
}

func SelectClaimedObjectsNotFoundInDiscoveredObjects(db *sqlx.DB) ([]model.ClaimedObject, error) {
	stmt := `
		SELECT  c.id,
				c."type",
				c.in_git,
				c.expired,
				c.created_at
		FROM      claimed_object c
		LEFT JOIN discovered_object d ON c.id = d.id
		AND       c.type = d.type
		WHERE     d.id IS NULL
		AND       d.type IS NULL
		AND	      c.expired = false
	`

	claimedObject := []model.ClaimedObject{}
	err := db.Select(&claimedObject, stmt)
	if err != nil {
		return nil, err
	}

	return claimedObject, nil
}

func UpdateClaimedObjectsSetExpired(db *sqlx.DB, id uint, objectType objectType.Type, value bool) error {
	stmt := `
	UPDATE claimed_object
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
	UPDATE claimed_object
	SET in_git = false
	WHERE NOT EXISTS (
		SELECT 1
		FROM discovered_object d
		WHERE claimed_object.id = d.id
		AND claimed_object.type = d.type
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
	UPDATE claimed_object
	SET in_git = true
	WHERE EXISTS (
		SELECT 1
		FROM discovered_object d
		WHERE claimed_object.id = d.id
		AND claimed_object.type = d.type
	);
	`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func SelectDistinctDiscoveredObjectsByType(db *sqlx.DB, t objectType.Type) ([]model.DiscoveredObject, error) {
	stmt := `
	SELECT DISTINCT id, type
		FROM discovered_object
		WHERE type = ?
	`

	discoveredObject := []model.DiscoveredObject{}
	err := db.Select(&discoveredObject, stmt, t)
	if err != nil {
		return nil, err
	}

	return discoveredObject, nil
}

func SelectDistinctClaimedObjectsByType(db *sqlx.DB, t objectType.Type) ([]model.ClaimedObject, error) {
	stmt := `
	SELECT DISTINCT id, type
		FROM claimed_object
		WHERE type = ?
		AND expired = false
	`

	claimedObject := []model.ClaimedObject{}
	err := db.Select(&claimedObject, stmt, t)
	if err != nil {
		return nil, err
	}

	return claimedObject, nil
}

func DeleteDiscoveredObjectNotInBranches(db *sqlx.DB, branches []model.Branch, repository string) (int64, error) {
	stmt := `
	DELETE FROM discovered_object
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
	DELETE FROM discovered_object
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

func GetClaimCountByDate(db *sqlx.DB, date string) (int, error) {
	stmt := `
		SELECT COUNT(*)
			FROM claimed_object
			WHERE created_at LIKE ?
	`

	date += "%"

	var count int
	err := db.Get(&count, stmt, date)
	if err != nil {
		return 0, err
	}

	return count, nil
}
