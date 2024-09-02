package database

import (
	"boilerplate/utils/config"
)

type Database struct {
	config *config.Config
	MySQL  *MySQL
	Redis  *Redis
}

func NewDatabase(config *config.Config) *Database {
	return &Database{
		config: config,
		MySQL:  NewMySQL(config),
		Redis:  NewRedis(config),
	}
}
