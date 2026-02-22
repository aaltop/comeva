package validate

import (
	"comeva/internal/logging"
	"os"
)

var debugLogger = logging.NewDebugLogger(os.Stdout)
