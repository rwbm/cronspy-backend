package api

import (
	"cronspy/backend/pkg/util/config"
	"cronspy/backend/pkg/util/log"
	"cronspy/backend/pkg/util/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) (err error) {

	// create DB connection
	// ds, errDB := mysql.New(
	// 	cfg.Database.Address,
	// 	cfg.Database.DefaultDB,
	// 	cfg.Database.Username,
	// 	cfg.Database.Password,
	// 	cfg.Database.MaxOpenConnections,
	// 	cfg.Database.MaxIdleConnecrtions,
	// 	cfg.Database.MaxLifeTime)
	// if errDB != nil {
	// 	return errDB
	// }

	logger := log.New()

	// http server
	e := server.New(cfg.Server.Debug)

	// start HTTP server
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	}, logger)

	return
}
