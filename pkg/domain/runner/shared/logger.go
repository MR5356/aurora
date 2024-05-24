package shared

import (
	"github.com/hashicorp/go-hclog"
	"os"
)

var Logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Trace,
	Output:     os.Stderr,
	JSONFormat: true,
})
