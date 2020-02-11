package log

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Log represents zerolog logger
type Log struct {
	logger *zerolog.Logger
}

// New instantiates new zero logger
func New() *Log {
	z := zerolog.New(os.Stdout)
	return &Log{
		logger: &z,
	}
}

// Log with HTTP context logs using zerolog
func (z *Log) Log(ctx echo.Context, source, msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	params["source"] = source

	if id, ok := ctx.Get("client_id").(int); ok {
		params["client_id"] = id
	}

	if err != nil {
		params["error"] = err
		z.logger.Error().Fields(params).Msg(msg)
		return
	}

	z.logger.Info().Fields(params).Msg(msg)
}

// Debug logs with level DEBUG
func (z *Log) Debug(msg string, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	z.logger.Debug().Fields(params).Msg(msg)
}

// Info logs with level INFO
func (z *Log) Info(msg string, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	z.logger.Info().Fields(params).Msg(msg)
}

// Warn logs with level WARN
func (z *Log) Warn(msg string, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	z.logger.Warn().Fields(params).Msg(msg)
}

// Error logs with level ERROR
func (z *Log) Error(msg string, err error, params map[string]interface{}) {

	if params == nil {
		params = make(map[string]interface{})
	}

	if err != nil {
		params["error"] = err
	}

	z.logger.Error().Fields(params).Msg(msg)
}
