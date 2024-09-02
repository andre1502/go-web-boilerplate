package database

import (
	"boilerplate/utils/config"
	"boilerplate/utils/logger"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"moul.io/zapgorm2"
)

type MySQL struct {
	Orm *gorm.DB
}

func NewMySQL(config *config.Config) *MySQL {
	my := &MySQL{}
	my.Connect(config)

	return my
}

func (my *MySQL) Connect(config *config.Config) {
	var err error

	log := zapgorm2.New(logger.Logger)
	log.SetAsDefault()
	log.LogLevel = logger.DBLevel
	log.SlowThreshold = time.Second * 5
	log.IgnoreRecordNotFoundError = true

	// get default db
	dsnWrite := my.getDsn(config.MySQL.Default, config.Timezone)

	my.Orm, err = gorm.Open(mysql.Open(dsnWrite), &gorm.Config{
		Logger:         log,
		TranslateError: true,
	})

	if err != nil {
		logger.Sugar.Fatal(err)
	}

	my.Orm.Use(
		my.initDbResolver(config.MySQL.Connections, config.Timezone).
			SetConnMaxIdleTime(time.Duration(config.MySQL.ConnMaxIdleTime) * time.Minute).
			SetConnMaxLifetime(time.Duration(config.MySQL.ConnMaxLifeTime) * time.Minute).
			SetMaxIdleConns(config.MySQL.MaxIdleConns).
			SetMaxOpenConns(config.MySQL.MaxOpenConns),
	)

	logger.Sugar.Debug("MySQL connected.")
}

func (my *MySQL) initDbResolver(configs []config.MySQLConnections, timezone string) *dbresolver.DBResolver {
	resolver := &dbresolver.DBResolver{}

	for _, config := range configs {
		resolver.Register(dbresolver.Config{
			Sources:           my.initDialetor(config.Writes, timezone),
			Replicas:          my.initDialetor(config.Reads, timezone),
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}, strings.Join(config.Datas, ", "))
	}

	return resolver
}

func (my *MySQL) initDialetor(configs []config.MySQLConnection, timezone string) []gorm.Dialector {
	dialector := []gorm.Dialector{}

	for _, config := range configs {
		dsn := my.getDsn(config, timezone)

		dialector = append(dialector, mysql.Open(dsn))
	}

	return dialector
}

func (my *MySQL) getDsn(config config.MySQLConnection, timezone string) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
		config.Username, config.Password, config.Network, config.Host, config.Port, config.Schema, config.Charset, timezone,
	)
}
