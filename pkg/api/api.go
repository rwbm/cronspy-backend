package api

import (
	"cronspy/backend/pkg/api/job"
	jt "cronspy/backend/pkg/api/job/transport"
	"cronspy/backend/pkg/api/user"
	ut "cronspy/backend/pkg/api/user/transport"
	"cronspy/backend/pkg/util/config"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/server"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Start starts the API service
func Start(cfg *config.Configuration) (err error) {

	// create DB connection
	connString := createMySQLConnectionString(cfg.Database.Address, cfg.Database.Username, cfg.Database.Password, cfg.Database.DefaultDB)
	ds, errDB := gorm.Open(cfg.Database.Driver, connString)
	if errDB != nil {
		return errDB
	}
	ds.DB().SetMaxOpenConns(cfg.Database.MaxOpenConnections)
	ds.DB().SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	ds.DB().SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifeTime) * time.Second)

	// default logger
	logger := log.New()

	// http server
	e := server.New()

	// +++++++++++ SERVICES ++++++++++++
	//

	ut.NewHTTP(user.Initialize(ds, nil, logger, cfg.Server.TokenExpiration, cfg.Server.DryRun), jwtSigningKey, jwtSigningMethod, e)
	jt.NewHTTP(job.Initialize(ds, nil, logger), jwtSigningKey, jwtSigningMethod, e)

	//
	// +++++++++++++++++++++++++++++++++

	// start HTTP server
	server.Start(e,
		&server.Config{
			ServiceName:         cfg.Server.Name,
			Port:                cfg.Server.Port,
			ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
			WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		},
		logger)

	return
}

func createMySQLConnectionString(address, username, password, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True", username, password, address, dbName)
}
