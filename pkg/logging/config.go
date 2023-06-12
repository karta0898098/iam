package logging

import (
	"io"
)

// Config for setup log want output setting
type Config struct {
	Env   string `mapstructure:"env"`
	App   string `mapstructure:"app"`
	Debug bool   `mapstructure:"debug"`
	Level Level  `mapstructure:"level"`

	writer io.Writer
}

// An Option is passed to Config
type Option interface {
	apply(*Config)
}

type setEnv struct{ env string }

func (opt *setEnv) apply(c *Config) { c.Env = opt.env }

type setApp struct{ app string }

func (opt *setApp) apply(c *Config) { c.App = opt.app }

type setDebug struct{ debug bool }

func (opt *setDebug) apply(c *Config) { c.Debug = opt.debug }

type setLevel struct{ level Level }

func (opt *setLevel) apply(c *Config) { c.Level = opt.level }

type setWriter struct{ writer io.Writer }

func (opt *setWriter) apply(c *Config) { c.writer = opt.writer }

func WithEnv(env string) Option {
	return &setEnv{
		env: env,
	}
}

func WithApp(app string) Option {
	return &setApp{
		app: app,
	}
}

func WithDebug(debug bool) Option {
	return &setDebug{
		debug: debug,
	}
}

func WithLevel(level Level) Option {
	return &setLevel{
		level: level,
	}
}

func WithOutput(w io.Writer) Option {
	return &setWriter{writer: w}
}
