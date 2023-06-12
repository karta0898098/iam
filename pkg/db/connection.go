package db

import (
	"gorm.io/gorm"

	"github.com/karta0898098/iam/pkg/db/conn"
)

type RWConfig struct {
	Read  conn.Database `mapstructure:"read"`
	Write conn.Database `mapstructure:"write"`
}

type Connection interface {
	ReadDB() *gorm.DB
	WriteDB() *gorm.DB
}

type RWConnection struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

func (conn *RWConnection) ReadDB() *gorm.DB {
	return conn.readDB
}

func (conn *RWConnection) WriteDB() *gorm.DB {
	return conn.writeDB

}

func NewRWConnection(config RWConfig) (Connection, error) {
	var (
		c Connection
	)
	readDB, err := conn.SetupDatabase(&config.Read)
	if err != nil {
		return c, err
	}
	writeDB, err := conn.SetupDatabase(&config.Write)
	if err != nil {
		return c, err
	}

	return &RWConnection{
		readDB:  readDB,
		writeDB: writeDB,
	}, nil
}

type Config struct {
	Conn conn.Database `mapstructure:"conn"`
}

type connection struct {
	db *gorm.DB
}

func (c *connection) ReadDB() *gorm.DB {
	return c.db
}

func (c *connection) WriteDB() *gorm.DB {
	return c.db
}

func NewConnection(config Config) (Connection, error) {
	var (
		c Connection
	)
	db, err := conn.SetupDatabase(&config.Conn)
	if err != nil {
		return c, err
	}

	return &connection{
		db: db,
	}, nil
}
