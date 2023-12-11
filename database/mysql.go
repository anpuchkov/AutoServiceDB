package database

import (
	"database/sql"
	"fmt"
)

type ConnectionInfo struct {
	User   string // Username
	Passwd string // Password (requires User)
	Net    string // Network type
	Addr   string // Network address (requires Net)
	DBName string // Database name
}

func MySqlConnect(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s(%s)/%s", info.User, info.Passwd, info.Net, info.Addr, info.DBName))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %s", err)
	}
	return db, nil
}
