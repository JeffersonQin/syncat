package database

import (
	"database/sql"
	"github.com/JeffersonQin/syncat/pkg/config"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

// Global database instance
var db *sql.DB

// Get the database connection string
func getDbConnectionString(dbConfig config.SyncatDBConfig) string {
	return "file:" + dbConfig.Filename + "?cache=shared&mode=rwc"
}

// SQL statements for creating tables
var createTableSql = []string{ /* file entry table */ `
	CREATE TABLE IF NOT EXISTS "entries" (
		"id"		INTEGER PRIMARY KEY AUTOINCREMENT,
		"path" 		VARCHAR(512) NOT NULL,
		"hash_md5" 	VARCHAR(32) NOT NULL,
		"timestamp" DATETIME NOT NULL,
		"size" 		INTEGER NOT NULL,
	    "is_dir"    INTEGER NOT NULL,
		"deleted" 	INTEGER NOT NULL,
		"uuid"		VARCHAR(36) NOT NULL
	)
	`,
	/*
	 * client table. server will use this table to store clients,
	 * and clients will use this table to store their own uuid allocated by the server.
	 * for clients, only the first row with id = 1 will be used.
	 */`
	CREATE TABLE IF NOT EXISTS "clients" (
		"id"		INTEGER PRIMARY KEY AUTOINCREMENT,
		"uuid"		VARCHAR(36) NOT NULL
	)
	`,
	/*
	 * last sync table. server will use this table to store the last syncing status of each client.
	 * for clients, using cid = 1 will store the last syncing status of the server.
	 * this table will be helpful when deciding which files to sync.
	 */`
	CREATE TABLE IF NOT EXISTS "last_sync" (
		"fid"		INTEGER NOT NULL,
		"cid"		INTEGER NOT NULL,
		"path" 		VARCHAR(512) NOT NULL,
		"hash_md5" 	VARCHAR(32) NOT NULL,
		"timestamp" DATETIME NOT NULL,
		"size" 		INTEGER NOT NULL,
	    "is_dir"    INTEGER NOT NULL,
		"deleted" 	INTEGER NOT NULL,
		"uuid"		VARCHAR(36) NOT NULL,
	    PRIMARY KEY (fid, cid)
	)
	`,
}

// LoadDatabase Load database
func LoadDatabase() error {
	dbConfig := config.GetConfig().Db
	// Create db file if not exists
	err := os.MkdirAll(filepath.Dir(dbConfig.Filename), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(dbConfig.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		_ = f.Close()
		return err
	}
	// Close db file
	err = f.Close()
	if err != nil {
		return err
	}
	// Connect to the database
	db, err = sql.Open("sqlite3", getDbConnectionString(dbConfig))
	if err != nil {
		_ = db.Close()
		return err
	}
	// Ping
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return err
	}
	// Create the tables if not exist
	for _, s := range createTableSql {
		_, err = db.Exec(s)
		if err != nil {
			_ = db.Close()
			return err
		}
	}
	return nil
}

// CloseDatabase Close database
func CloseDatabase() error {
	return db.Close()
}

// QueryExistsClientUuid Query if the client uuid exists on server
func QueryExistsClientUuid(uuid string) (bool, error) {
	row, err := db.Query("SELECT 1 FROM `clients` WHERE `uuid` = ? LIMIT 1", uuid)
	if err != nil {
		return false, err
	}
	defer func(row *sql.Rows) {
		_ = row.Close()
	}(row)
	return row.Next(), nil
}

// AllocateNewClientUuid Allocate a new uuid for the client on the server
func AllocateNewClientUuid() (string, error) {
	var uuidStr string
	for {
		uuidStr = uuid.NewString()
		exists, err := QueryExistsClientUuid(uuidStr)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
	}
	_, err := db.Exec("INSERT INTO `clients` (`uuid`) VALUES (?)", uuidStr)
	if err != nil {
		return "", err
	}
	return uuidStr, nil
}

// QueryClientUuid Query the client's own uuid
func QueryClientUuid() (string, error) {
	row, err := db.Query("SELECT `uuid` FROM `clients` WHERE `id` = 1 LIMIT 1")
	if err != nil {
		return "", err
	}
	defer func(row *sql.Rows) {
		_ = row.Close()
	}(row)
	if !row.Next() {
		return "", nil
	}
	var uuidStr string
	err = row.Scan(&uuidStr)
	if err != nil {
		return "", err
	}
	return uuidStr, nil
}

// UpdateClientUuid Update the client's own uuid
func UpdateClientUuid(uuid string) error {
	_, err := db.Exec("UPDATE `clients` SET `uuid` = ? WHERE `id` = 1", uuid)
	return err
}
