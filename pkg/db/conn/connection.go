package conn

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/karta0898098/iam/pkg/errors"
)

type DatabaseType string

const (
	// MySQL ...
	MySQL DatabaseType = "mysql"
	// Postgres ...
	Postgres DatabaseType = "postgres"
)

type Database struct {
	Debug          bool         `mapstructure:"debug"`
	Host           string       `mapstructure:"host"`
	User           string       `mapstructure:"user"`
	Port           int          `mapstructure:"port"`
	Password       string       `mapstructure:"password"`
	Name           string       `mapstructure:"name"`
	Type           DatabaseType `mapstructure:"type"`
	MaxIdleConns   int          `mapstructure:"max_idle_conns"`
	MaxOpenConns   int          `mapstructure:"max_open_conns"`
	MaxLifetimeSec int          `mapstructure:"max_lifetime"`
	ReadTimeout    string       `mapstructure:"read_timeout"`
	WriteTimeout   string       `mapstructure:"write_timeout"`
}

func SetupDatabase(database *Database) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	if database.WriteTimeout == "" {
		database.WriteTimeout = "10s"
	}

	if database.ReadTimeout == "" {
		database.ReadTimeout = "10s"
	}

	var dialector gorm.Dialector

	switch database.Type {
	case MySQL:
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true&readTimeout=%s&writeTimeout=%s", database.User, database.Password, database.Host+":"+strconv.Itoa(database.Port), database.Name, database.ReadTimeout, database.WriteTimeout)
		dialector = mysql.Open(dsn)
	case Postgres:
		dsn := fmt.Sprintf(`user=%s password=%s host=%s port=%d dbname=%s sslmode=disable `, database.User, database.Password, database.Host, database.Port, database.Name)
		dialector = postgres.Open(dsn)
	default:
		return nil, errors.New("Not support driver")
	}

	colorful := false
	logLevel := logger.Silent
	if database.Debug {
		colorful = true
		logLevel = logger.Info
	}

	newLogger := NewLogger(logger.Config{
		SlowThreshold: time.Second, // Slow SQL threshold
		LogLevel:      logLevel,    // Log level
		Colorful:      colorful,    // Disable color
	})

	var conn *gorm.DB

	// retry connect to database when db not response
	err := backoff.Retry(func() error {
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return err
		}
		conn = db
		sqlDB, err := conn.DB()
		if err != nil {
			return err
		}

		err = sqlDB.Ping()
		if err != nil {
			return err
		}
		return nil
	}, bo)

	if err != nil {
		return nil, err
	}

	// set default idle conn
	if database.MaxIdleConns == 0 {
		database.MaxIdleConns = 10
	}

	if database.MaxOpenConns == 0 {
		database.MaxOpenConns = 20
	}

	if database.MaxLifetimeSec == 0 {
		database.MaxLifetimeSec = 14400
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(database.MaxLifetimeSec))
	return conn, nil
}
