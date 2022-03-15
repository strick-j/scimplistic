package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/strick-j/scimplistic/utils"
)

////////////// DATABASE SETUP FUNCTIONS //////////////////////////////////////////////////////////////////////

var database Database
var err error

//Database encapsulates database
type Database struct {
	db *sql.DB
}

//Begins a transaction
func (db Database) begin() (tx *sql.Tx) {
	tx, err := db.db.Begin()
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "begin"}).Error(err)
		return nil
	}
	return tx
}

func (db Database) prepare(q string) (*sql.Stmt, error) {
	stmt, err := db.db.Prepare(q)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "prepare"}).Error(err)
		return nil, err
	}
	return stmt, nil
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(q, args...)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "query"}).Error(err)
		return nil
	}
	return rows
}

func init() {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DatabaseConnect"}).Trace("Starting Database Connection")
	values, err := utils.ReadConfig("settings.json")
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "DatabaseConnect"}).Error(err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", values.DatabaseIP, values.DatabasePort, values.DatabaseUser, values.DatabasePass, values.DatabaseName)

	// Open database
	database.db, err = sql.Open("postgres", psqlconn)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "init"}).Error(err)
	}
}

// Close function closes this database connection
func Close() {
	database.db.Close()
}

////////////// MANAGEMENT FUNCTIONS //////////////////////////////////////////////////////////////////////

func AddAction(action string, resourceType string, resource string, actionResult string) error {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "AddAction"}).Trace("Starting AddAction process")
	err = genQuery(`insert into "actions"("action", "resourcetype", "resource", "result", "performed") values($1,$2,$3,$4,$5)`, action, resourceType, resource, actionResult, time.Now())
	if err != nil {
		return err
	}
	return nil
}

////////////// GENERIC DB FUNCTIONS //////////////////////////////////////////////////////////////////////

func genQuery(sql string, args ...interface{}) error {
	log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GenericQuery"}).Trace("Inside Generic Query")
	SQL, err := database.prepare(sql)
	if err != nil {
		return err
	}
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GenericQuery"}).Error(err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GenericQuery"}).Error(err)
			return err
		}
		log.WithFields(log.Fields{"Category": "Server Processes", "Function": "GenericQuery"}).Info("Generic Query commit successful")
	}
	return err
}
