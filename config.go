package wlog

import (
	"os"
	"time"
	"io"
)

type MultiConfig struct {
	Configs []Config `json:"config" yaml:"configs"`
}

func (c MultiConfig) Create(opts ...LoggerOpt) ([]*Logger, error) {
	var loggers []*Logger
	for _, cfg := range c.Configs {
		logger, err := cfg.Create(opts...)
		if err != nil {
			return nil, err
		}
		loggers = append(loggers, logger)
	}
	return loggers, nil
}

// EasyMultiConfig is a easy way to create multiple Loggers if
// all writers use the same Config.
//
// Note: if the length of Config's Paths is not 0, the Config also creates the first Logger.
type EasyMultiConfig struct {
	Config     Config     `json:"config" yaml:"config"`
	MultiPaths [][]string `json:"multi_paths" yaml:"multi_paths"`
}

func (c EasyMultiConfig) Create(opts ...LoggerOpt) ([]*Logger, error) {
	var loggers []*Logger
	if len(c.Config.Paths) != 0 {
		logger, err := c.Config.Create(opts...)
		if err != nil {
			return nil, err
		}
		loggers = append(loggers, logger)
	}
	cfg:=c.Config
	for _, paths := range c.MultiPaths {
		cfg.Paths=paths
		logger, err := cfg.Create(opts...)
		if err != nil {
			return nil, err
		}
		loggers = append(loggers, logger)
	}
	return loggers, nil
}

type Config struct {
	// MinLevel is the string representation of minimum logging level.
	// it's default value is "debug", and supported values are as follow:
	// "debug", "info", "warn", "error", "fatal", "panic".
	MinLevel string `json:"min_level" yaml:"min_level"`
	// Encoder is the type of chosen encoder, it's default value is "text",
	// and temporarily supported values are as follow: "text", "json".
	Encoder       string        `json:"encoder" yaml:"encoder"`
	EncoderConfig EncoderConfig `json:"encoder_config" yaml:"encoder_config"`
	// Paths is the descriptor of standard output, standard error or file paths
	// to write logs to.
	//
	// The temporarily supported values are as follow: "stdout", "stderr" and all existed file paths.
	// If the length of Paths is 0, then the logs are output to standard output by default.
	Paths        []string     `json:"paths" yaml:"paths"`
	FileConfig   FileConfig   `json:"file_config" yaml:"file_config"`
	WriterConfig WriterConfig `json:"writer_config" yaml:"writer_config"`
	// ErrWriter is the Path of writer to write internal errors to.
	// A standard error is the default writer.
	ErrWriter string `json:"err_writer" yaml:"err_writer"`
}

type EncoderConfig struct {
	ColorEnabled bool   `json:"color_enabled" yaml:"color_enabled"`
	LevelLower   bool   `json:"level_lower" yaml:"level_lower"`
	TimeDisabled bool   `json:"time_disable" yaml:"time_disable"`
	LineEnding   string `json:"line_ending" yaml:"line_ending"`
}

type FileConfig struct {
	MaxSize        int64 `json:"max_size" yaml:"max_size"`
	MaxRotatedSize int64 `json:"max_rotated_size" yaml:"max_rotated_size"`
	MaxRotatedDays int   `json:"max_rotated_days" yaml:"max_rotated_days"`
	DisableDaily   bool  `json:"disable_daily" yaml:"disable_daily"`
}

type WriterConfig struct {
	// FlushInterval is the flush interval in seconds of a TimingFlushWriter.
	FlushInterval int `json:"flush_interval" yaml:"flush_interval"`
	// MinBufSize is the minimum size of a BufWriter.
	MinBufSize int `json:"min_buf_size" yaml:"min_buf_size"`
	// MinBufSize is the maximum size of a BufWriter.
	MaxBufSize int `json:"max_buf_size" yaml:"max_buf_size"`
}

// NewConsoleConfig return a Config to create a Logger.
//
// It's minimum logging level is "InfoLvl".
//
// It uses a TextEncoder to encode the message,
// writes the data to the BufWriter with a buffer size of 4096 first,
// and flushes the buffered data to the standard output every 3 seconds.
func NewConsoleConfig() Config {
	return Config{
		MinLevel: "info",
	}
}

// NewFileConfig return a Config to create a Logger.
//
// It's minimum logging level is "InfoLvl".
//
// It uses a TextEncoder to encode the message,
// writes the data to the BufWriter with a buffer size of 4096 first,
// and flushes the buffered data to a file in the specified path every 3 seconds.
//
// The file is daily, it supports rotating by size and day, it's maximum size is "100 * 1 << 20".
func NewFileConfig(path string) Config {
	c := Config{
		MinLevel: "info",
		Paths:    []string{path},
	}
	return c
}

// NewMultiWriterConfig return a Config to create a Logger.
//
// It's minimum logging level is "InfoLvl".
//
// It uses a TextEncoder to encode the message, and writes the data to
// the standard output if consoleEnabled is true and multiple files in the specified paths.
//
// For the standard output or every file, the data is written to the BufWriter with a buffer size of 4096 first,
// and a TimingFlushWriter will flush the buffered data to it every 3 seconds.
//
// For every file, it is daily, it supports rotating by size and day, it's maximum size is "100 * 1 << 20".
func NewMultiWriterConfig(paths ...string) Config {
	c := Config{
		MinLevel: "info",
		Paths:    paths,
	}
	return c
}

// Create returns a Logger form the config and the opts.
func (c Config) Create(opts ...LoggerOpt) (*Logger, error) {
	errW, err := c.CreateErrWriter()
	if err != nil {
		return nil, err
	}
	encoder := c.CreateEncoder()
	writer := c.CreateWriter(SetFileErrW(errW))
	writer = c.WrapWriter(writer, SetBufErrW(errW))
	h := NewBaseHandler(writer, encoder)
	logger := NewLogger(h, []LoggerOpt{SetLogMinLvl(c.CreateLevel()), SetLogErrW(errW)}...)
	if len(opts) == 0 {
		return logger, nil
	}
	return logger.WithOpts(opts...), nil
}

// CreateErrWriter returns a io.Writer form the config to output the log's internal error.
func (c Config) CreateErrWriter() (io.Writer, error) {
	if c.ErrWriter == "" {
		return os.Stderr, nil
	}
	w, err := os.OpenFile(c.ErrWriter, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	return w, err
}

// CreateEncoder returns a Encoder form the config and the opts.
func (c Config) CreateEncoder(opts ...EncoderOpt) Encoder {
	ec := c.EncoderConfig
	cfgOpts := []EncoderOpt{SetLineEnding(ec.LineEnding)}
	if ec.ColorEnabled {
		cfgOpts = append(cfgOpts, EnableColor())
	}
	if ec.LevelLower {
		cfgOpts = append(cfgOpts, SetLevelLower())
	}
	if ec.TimeDisabled {
		cfgOpts = append(cfgOpts, DisableTime())
	}
	opts = append(cfgOpts, opts...)
	var encoder Encoder
	switch c.Encoder {
	case "json":
		encoder = NewJsonEncoder(opts...)
	default:
		encoder = NewTextEncoder(opts...)
	}
	return encoder
}

// CreateWriter returns a Writer form the config and the opts.
func (c Config) CreateWriter(opts ...FileWriterOpt) Writer {
	if len(c.Paths) == 0 {
		c.Paths = []string{"stdout"}
	}
	var writers []Writer
	fc := c.FileConfig
	cfgOpts := []FileWriterOpt{SetFileMaxSize(fc.MaxSize),
		SetFileMaxRotatedSize(fc.MaxRotatedSize),
		SetFileMaxRotatedDays(fc.MaxRotatedDays)}
	if fc.DisableDaily {
		cfgOpts = append(cfgOpts, DisableFileDaily())
	}
	opts = append(cfgOpts, opts...)
	for _, path := range c.Paths {
		switch path {
		case "stdout":
			writers = append(writers, NewIOWriter(os.Stdout))
		case "stderr":
			writers = append(writers, NewIOWriter(os.Stderr))
		default:
			writers = append(writers, NewFileWriter(path, opts...))
		}
	}
	if len(writers) == 1 {
		return writers[0]
	}
	return NewMultiWriter(writers)
}

// WrapWriter wraps the given Writer form the config and the opts.
// It wraps the Writer into a BufWriter first and finally returns a NewTimingFlushWriter.
func (c Config) WrapWriter(inner Writer, opts ...BufWriterOpt) *TimingFlushWriter {
	wc := c.WriterConfig
	cfgOpts := []BufWriterOpt{SetBufMinSize(wc.MinBufSize), SetBufMaxSize(wc.MaxBufSize)}
	opts = append(cfgOpts, opts...)
	bw := NewBufWriter(inner, opts...)
	return NewTimingFlushWriter(bw, time.Duration(wc.FlushInterval)*time.Second)
}

// CreateWriter returns a Level form the config.
func (c Config) CreateLevel() Level {
	var lvl Level
	switch c.MinLevel {
	case "info":
		lvl = DebugLvl
	case "warn":
		lvl = DebugLvl
	case "error":
		lvl = DebugLvl
	case "fatal":
		lvl = DebugLvl
	case "panic":
		lvl = DebugLvl
	default:
		lvl = DebugLvl
	}
	return lvl
}
