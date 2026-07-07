// Package logging contains more controlled methods of printing to outputs.
package logging

import (
	"io"
	"log"
	"os"

	errorHelpers "github.com/aaltop/comeva/internal/errors"
)

// LevelLogger contains loggers with levels of criticality.
type LevelLogger struct {
	// Logger is the underlying logger. This can be manipulated directly;
	// however, the final output of [LevelLogger] should be set using
	// [LevelLogger.SetOutput]. In particular, setting the output of
	// [LevelLogger.Logger] will not work to change the output of [LevelLogger].
	Logger *log.Logger

	level int

	// The final Writer in a possible chain of writers; the actual output.
	baseOut io.Writer
	// Used as the base prefix for the [LevelLogger.Logger], set using
	// [LevelLogger.SetPrefix].
	basePrefix string

	// Output corresponding to a level of criticality.
	debugWriter, infoWriter, warningWriter, errorWriter, criticalWriter *levelWriter
}

func NewDefaultLevelLogger() (logger *LevelLogger) {
	logger = errorHelpers.Panic2(NewLevelLogger(0, 0, 0, 0, 0, 0, log.Writer(), func(out io.Writer) *log.Logger {
		return log.New(out, "", 0)
	}))
	return logger
}

// NewLevelLogger creates a new [LevelLogger].
//
//   - debug, info, warning er, and critical are the levels for each logging level.
//
//   - external is the logging level currently in effect.
//
//   - out is the writer used for output.
//
//   - loggerFactory creates new loggers. its argument `out` is expected to be
//     set as the out argument of the returned logger.
//
// The expected form of loggerFactory is
//
//	loggerFactory = func(out io.Writer) *log.Logger {
//	  return log.New(out, prefix, flag)
//	}
//
// with prefix and flag chosen by the user, thus simply allowing for
// better control of the underlying logger.
func NewLevelLogger(debug, info, warning, er, critical, external int, out io.Writer, loggerFactory func(out io.Writer) *log.Logger) (logger *LevelLogger, e error) {
	logger = &LevelLogger{
		level: external,
	}

	logger.errorWriter = newLevelWriter(er, &logger.level, logger.baseOut)
	// default to error level
	logger.Logger = loggerFactory(logger.errorWriter)

	logger.baseOut = out
	logger.basePrefix = logger.Logger.Prefix()

	logger.debugWriter = newLevelWriter(debug, &logger.level, logger.baseOut)
	logger.infoWriter = newLevelWriter(info, &logger.level, logger.baseOut)
	logger.warningWriter = newLevelWriter(warning, &logger.level, logger.baseOut)
	logger.criticalWriter = newLevelWriter(critical, &logger.level, logger.baseOut)
	return
}

// NewLevelLoggerWithDefaults returns a [LevelLogger] with some suitable
// default values set.
func NewLevelLoggerWithDefaults() (logger *LevelLogger) {
	logger = errorHelpers.Panic2(NewLevelLogger(
		DEBUG, INFO, WARNING, ERROR, CRITICAL, ERROR,
		log.Writer(),
		func(out io.Writer) *log.Logger { return log.New(out, "", 0) },
	))
	return
}

// SetLevel sets the current logging level of the logger.
func (logger *LevelLogger) SetLevel(level *int) (e error) {

	defer func() {
		e = errorHelpers.HandleReturn(recover())
	}()

	errorHelpers.Return(logger.debugWriter.SetExternalLevel(level))
	errorHelpers.Return(logger.infoWriter.SetExternalLevel(level))
	errorHelpers.Return(logger.warningWriter.SetExternalLevel(level))
	errorHelpers.Return(logger.errorWriter.SetExternalLevel(level))
	errorHelpers.Return(logger.criticalWriter.SetExternalLevel(level))
	return
}

// LevelConfig is configuration for a particular level of a logger.
type LevelConfig struct {
	// Out is the output used for a given level.
	Out io.Writer

	// Level is used as the logging level. If left nil, the level
	// is not changed.
	Level *int
}

// LevelConfigs is the set of configs used to configure the levels
// of a [LevelLogger].
type LevelConfigs struct {
	Debug, Info, Warning, Error, Critical *LevelConfig
}

// Configure configures each of the logging levels of the logger.
func (logger *LevelLogger) Configure(config LevelConfigs) (e error) {

	defer func() {
		e = errorHelpers.HandleReturn(recover())
	}()

	if config.Debug != nil {
		errorHelpers.Return(logger.debugWriter.Configure(*config.Debug))
	}
	if config.Info != nil {
		errorHelpers.Return(logger.infoWriter.Configure(*config.Info))
	}
	if config.Warning != nil {
		errorHelpers.Return(logger.warningWriter.Configure(*config.Warning))
	}
	if config.Error != nil {
		errorHelpers.Return(logger.errorWriter.Configure(*config.Error))
	}
	if config.Critical != nil {
		errorHelpers.Return(logger.criticalWriter.Configure(*config.Critical))
	}
	return
}

func (logger *LevelLogger) SetOutput(out io.Writer) {
	logger.baseOut = out
}

// SetPrefix sets the base prefix for the logger. The actual prefix
// may contain more.
func (logger *LevelLogger) SetPrefix(prefix string) {
	logger.basePrefix = prefix
}

// Debug returns a logger with logging level at debug.
func (logger *LevelLogger) Debug() (internalLogger *log.Logger) {
	logger.Logger.SetOutput(logger.debugWriter)
	logger.Logger.SetPrefix("DEBUG " + logger.basePrefix)
	return logger.Logger
}

// Info returns a logger with logging level at info.
func (logger *LevelLogger) Info() (internalLogger *log.Logger) {
	logger.Logger.SetOutput(logger.infoWriter)
	logger.Logger.SetPrefix("INFO " + logger.basePrefix)
	return logger.Logger
}

// Warning returns a logger with logging level at warning.
func (logger *LevelLogger) Warning() (internalLogger *log.Logger) {
	logger.Logger.SetOutput(logger.warningWriter)
	logger.Logger.SetPrefix("WARNING " + logger.basePrefix)
	return logger.Logger
}

// Error returns a logger with logging level at error.
func (logger *LevelLogger) Error() (internalLogger *log.Logger) {
	logger.Logger.SetOutput(logger.errorWriter)
	logger.Logger.SetPrefix("ERROR " + logger.basePrefix)
	return logger.Logger
}

// Critical returns a logger with logging level at critical.
func (logger *LevelLogger) Critical() (internalLogger *log.Logger) {
	logger.Logger.SetOutput(logger.criticalWriter)
	logger.Logger.SetPrefix("CRITICAL " + logger.basePrefix)
	return logger.Logger
}

// levelWriter implements the [io.Writer] interface, wrapping another writer
// and including a level that controls whether the writer should write.
// The writing depends on a externally set level: if the set level is greater
// than or equal to the external level, writing is performed.
type levelWriter struct {
	// level is the level of the writer.
	level int

	// externalLevel is an externally set level.
	externalLevel *int

	out io.Writer
}

func newDefaultLeverWriter() (writer *levelWriter) {
	var dum = 0
	return &levelWriter{
		level:         0,
		externalLevel: &dum,
		out:           os.Stdout,
	}
}

func newLevelWriter(level int, externalLevel *int, out io.Writer) (writer *levelWriter) {
	writer = newDefaultLeverWriter()
	writer.level = level
	writer.externalLevel = externalLevel
	writer.out = out
	return
}

func (writer *levelWriter) SetExternalLevel(level *int) (e error) {
	writer.externalLevel = level
	return
}

// Configure configures the writer.
func (writer *levelWriter) Configure(config LevelConfig) (e error) {
	if config.Level != nil {
		writer.level = *config.Level
	}
	if config.Out != nil {
		writer.out = config.Out
	}
	return
}

func (writer *levelWriter) Write(p []byte) (n int, err error) {
	if writer.level >= *writer.externalLevel {
		return writer.out.Write(p)
	}
	return
}

// NewDebugLogger creates a logger useful for debugging purposes.
func NewDebugLogger(out io.Writer) (logger *log.Logger) {
	return log.New(out, "", log.Llongfile)
}
