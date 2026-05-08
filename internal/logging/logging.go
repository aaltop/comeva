// Package logging contains more controlled methods of printing to outputs.
package logging

import (
	"io"
	"log"
	"os"
)

// LevelLogger contains loggers with levels of criticality.
type LevelLogger struct {
	level int

	// Logger corresponding to a level of criticality.
	Debug, Info, Warning, Error, Critical *log.Logger
}

func NewDefaultLevelLogger() (logger *LevelLogger) {
	var e error
	logger, e = NewLevelLogger(0, 0, 0, 0, 0, 0, log.Writer(), func(out io.Writer) *log.Logger {
		return log.New(out, "", 0)
	})

	if e != nil {
		panic(e)
	}
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

	logger.Debug = loggerFactory(newLevelWriter(debug, &logger.level, out))
	logger.Info = loggerFactory(newLevelWriter(info, &logger.level, out))
	logger.Warning = loggerFactory(newLevelWriter(warning, &logger.level, out))
	logger.Error = loggerFactory(newLevelWriter(er, &logger.level, out))
	logger.Critical = loggerFactory(newLevelWriter(critical, &logger.level, out))
	return
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

func (writer *levelWriter) Write(p []byte) (n int, err error) {
	if writer.level >= *writer.externalLevel {
		return writer.out.Write(p)
	}
	return
}

func (writer *levelWriter) SetExternalLevel(level *int) (e error) {
	writer.externalLevel = level
	return
}

// NewDebugLogger creates a logger useful for debugging purposes.
func NewDebugLogger(out io.Writer) (logger *log.Logger) {
	return log.New(out, "", log.Llongfile)
}
