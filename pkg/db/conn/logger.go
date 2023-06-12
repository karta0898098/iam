package conn

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	LogLevel                            logger.LogLevel
	Config                              logger.Config
	SlowThreshold                       time.Duration
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewLogger(config logger.Config) logger.Interface {
	var (
		infoStr      = "%s "
		warnStr      = "%s "
		errStr       = "%s "
		traceStr     = "%s "
		traceWarnStr = "%s "
		traceErrStr  = "%s "
	)

	l := &gormLogger{
		LogLevel:      config.LogLevel,
		SlowThreshold: config.SlowThreshold,
		Config:        config,
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}

	return l
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return g
}

func (g *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	var (
		currentLogger zerolog.Logger
	)

	currentLogger = log.With().Logger()

	if g.LogLevel >= logger.Info {
		currentLogger.
			Info().
			Msgf(msg, data...)
	}
}

func (g *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	var (
		currentLogger zerolog.Logger
	)

	currentLogger = log.With().Logger()

	if g.LogLevel >= logger.Warn {
		currentLogger.Warn().Msgf(msg, data...)
	}
}

func (g *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	var (
		currentLogger zerolog.Logger
		errors        []error
	)

	currentLogger = log.With().Logger()

	if g.LogLevel >= logger.Error {
		for i := 0; i < len(data); i++ {
			if err, ok := data[i].(error); ok && err != nil {
				errors = append(errors, err)
			}
		}
		currentLogger.Error().Errs("errors", errors).Msgf(msg, data...)
	}
}

func (g *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	var (
		currentLogger zerolog.Logger
	)

	currentLogger = log.Ctx(ctx).With().Logger()

	if g.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && g.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				currentLogger.
					Error().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", -1).
					Err(err).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Error().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Err(err).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			}
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
			if rows == -1 {
				currentLogger.
					Warn().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", -1).
					Str("sql_slow_log", slowLog).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Warn().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Str("sql_slow_log", slowLog).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			}
		case g.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				currentLogger.
					Info().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", -1).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			} else {
				currentLogger.
					Info().
					Float64("sql_elapsed", float64(elapsed.Nanoseconds())/1e6).
					Int64("sql_row", rows).
					Str("caller", utils.FileWithLineNum()).
					Msg(sql)
			}
		}
	}
}
