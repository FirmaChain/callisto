package v5

import (
	"fmt"
	"log"
	"os"

	utils "github.com/forbole/callisto/v4/modules/utils"
	"github.com/forbole/callisto/v4/types"
)

// Migrate implements database.Migrator
func (db *Migrator) Migrate() error {
	msgTypes, err := db.getMsgTypesFromMessageTable()
	if err != nil {
		return fmt.Errorf("error while getting message types rows: %s", err)
	}

	for _, msgType := range msgTypes {
		// migrate message types
		err = db.migrateMsgTypes(types.NewMessageType(
			msgType.Type,
			utils.GetModuleNameFromTypeURL(msgType.Type),
			utils.GetMsgFromTypeURL(msgType.Type),
			msgType.Height))

		if err != nil {
			return err
		}
	}
	// read the sql migration file
	sqlBytes, err := os.ReadFile("database/migrate/v5/migrationv5.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	// execute scripts
	err = db.migrateDbSchema(string(sqlBytes))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	return nil
}

// getMsgTypesFromMessageTable retrieves messages types stored in database inside message table
func (db *Migrator) getMsgTypesFromMessageTable() ([]MessageRow, error) {
	smt := "SELECT DISTINCT ON (type) type, transaction_hash, height FROM message ORDER BY type DESC"
	var rows []MessageRow
	err := db.SQL.Select(&rows, smt)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// migrateMsgTypes stores the given message type inside the database
func (db *Migrator) migrateMsgTypes(msg *types.MessageType) error {
	// Step 1: Create the table and indexes (no parameters, can be run as plain exec)
	schemaStmt := `
	CREATE TABLE IF NOT EXISTS message_type
	(
		type      TEXT   NOT NULL UNIQUE,
		module    TEXT   NOT NULL,
		label     TEXT   NOT NULL,
		height    BIGINT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS message_type_module_index ON message_type (module);
	CREATE INDEX IF NOT EXISTS message_type_type_index ON message_type (type);
	`

	_, err := db.SQL.Exec(schemaStmt)
	if err != nil {
		return err
	}

	// Step 2: Insert the message type (with parameters)
	insertStmt := `
	INSERT INTO message_type(type, module, label, height) 
	VALUES ($1, $2, $3, $4) 
	ON CONFLICT (type) DO NOTHING;`

	_, err = db.SQL.Exec(insertStmt, msg.Type, msg.Module, msg.Label, msg.Height)
	return err
}

func (db *Migrator) migrateDbSchema(migrateSQL string) error {
	_, err := db.SQL.Exec(migrateSQL)
	return err
}
