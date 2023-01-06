package repository

import "database/sql"

type DatabaseManager interface {
	SetClient(client *sql.DB)
	GetClient() *sql.DB
}

type DatabaseManagerImpl struct {
	Client *sql.DB
}

func NewDatabaseManager() DatabaseManager {
	return &DatabaseManagerImpl{}
}

func (d *DatabaseManagerImpl) SetClient(client *sql.DB) {
	d.Client = client
}

func (d *DatabaseManagerImpl) GetClient() *sql.DB {
	return d.Client
}
