package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"runtime"
	"time"
)

func init() {
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&nested.Formatter{
		HideKeys:        false,
		FieldsOrder:     []string{"level"},
		TimestampFormat: time.DateTime,
		TrimMessages:    true,
		CallerFirst:     false,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return ""
			//return fmt.Sprintf(" %s:%d", frame.Function, frame.Line)
		},
	})
}
