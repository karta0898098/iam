package logging

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Level defines log levels.
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

var (
	// Teal ...
	Teal = Color("\033[1;36m%s\033[0m")
	// Yellow ...
	Yellow = Color("\033[35m%s\033[0m")
	// Green
	Green = Color("\033[32m%s\033[0m")
)

var (
	DefaultLoggerConfig = &Config{
		Debug:  false,
		Level:  InfoLevel,
		writer: os.Stdout,
	}
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

// Color ...
func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type severityHook struct{}

// Run ...
func (h severityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	e.Float64("timestamp", float64(time.Now().UnixNano()/int64(time.Millisecond))/1000)
}

func SetupWithOption(opts ...Option) zerolog.Logger {
	logger := setup(DefaultLoggerConfig, opts...)
	log.Logger = logger

	return logger
}

func Setup(cfg Config) zerolog.Logger {
	logger := setup(&cfg)
	log.Logger = logger

	return logger
}

func setup(config *Config, opts ...Option) zerolog.Logger {
	var (
		logger zerolog.Logger
	)

	if config == nil {
		config = DefaultLoggerConfig
	}

	if config.writer == nil {
		config.writer = DefaultLoggerConfig.writer
	}

	for _, opt := range opts {
		opt.apply(config)
	}

	if config.Env == "" {
		config.Env = os.Getenv("LOG_ENV")
	}

	if lv, err := strconv.Atoi(os.Getenv("LOG_LEVEL")); err == nil {
		config.Level = Level(lv)
	}

	zerolog.DisableSampling(true)
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	level := config.Level
	if config.Debug {
		output := zerolog.ConsoleWriter{
			Out: config.writer,
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s=", Teal(i))
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatTimestamp = func(i interface{}) string {
			t := fmt.Sprintf("%v", i)
			millisecond, err := strconv.ParseInt(fmt.Sprintf("%s", i), 10, 64)
			if err == nil {
				t = time.Unix(millisecond/1000, 0).Local().Format("2006/01/02 15:04:05")
			}
			return colorize(t, colorCyan)
		}
		output.FormatCaller = func(i interface{}) string {
			var c string
			if cc, ok := i.(string); ok {
				c = cc
			}
			if len(c) > 0 {
				cwd, err := os.Getwd()
				if err == nil {
					c = strings.TrimPrefix(c, cwd)
					c = strings.TrimPrefix(c, "/")
				}
				c = colorize(c, colorGreen)

				if c != "" {
					c = fmt.Sprintf("%s %s", " >", c)
				}
			}
			return c
		}

		output.PartsOrder = []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.MessageFieldName,
			zerolog.CallerFieldName,
		}
		logger = zerolog.New(output)
	} else {
		logger = zerolog.New(config.writer)
	}

	if config.App != "" {
		logger = logger.With().Str("app", config.App).Logger()
	}

	if config.Env != "" {
		logger = logger.With().Str("env", config.Env).Logger()
	}

	logger = logger.
		Hook(severityHook{}).
		With().
		Timestamp().
		Logger().
		Level(zerolog.Level(level))

	return logger
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
