package logging

import (
	"comeva/internal/io"
	"comeva/internal/io/ansi"
	baseIo "io"
)

// CreateColoredLogger creates a coloured logger which has suitable coloured
// output for each logging level.
func CreateColoredLogger() (logger *LevelLogger) {

	logger = NewLevelLoggerWithDefaults()
	var createModWriter = func(colorScheme *ansi.ColorScheme, modif ...func(in string) (out string)) baseIo.Writer {
		var modifier = io.NewDefaultModifier[string]()

		// doesn't allow splatting of it for some reason (needs to be
		// typed exactly the same, not just equivalent?)
		for _, fun := range modif {
			modifier.Add(fun)
		}
		modifier.Add(colorScheme.ApplyFore)
		return io.NewStringModifierWriter(
			modifier,
			logger.baseOut,
		)
	}

	logger.Configure(LevelConfigs{
		Debug: &LevelConfig{
			Out: createModWriter(ansi.BasicColorSchemes.Debug),
		},
		Info: &LevelConfig{
			Out: createModWriter(ansi.BasicColorSchemes.Info),
		},
		Warning: &LevelConfig{
			Out: createModWriter(ansi.BasicColorSchemes.Warning),
		},
		Error: &LevelConfig{
			Out: createModWriter(ansi.BasicColorSchemes.Error),
		},
		Critical: &LevelConfig{
			Out: createModWriter(ansi.BasicColorSchemes.Critical),
		},
	})
	return
}
