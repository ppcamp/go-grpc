package log

import (
	"github.com/sirupsen/logrus"
)

// LogrusLogger implements the V method of the grpclog for the *logrus.Entry
//
// Example
// 	log := logrus.StandardLogger()
// 	log.SetLevel(logrus.DebugLevel)
// 	log.SetFormatter(
// 		&logrus.TextFormatter{
// 			ForceColors:      true,
// 			DisableColors:    false,
// 			TimestampFormat:  "2006-01-02 15:04:05",
// 			FullTimestamp:    true,
// 			QuoteEmptyFields: true,
// 			PadLevelText:     true,
// 		},
// 	)
//	grpclog.SetLoggerV2(&log.LogrusLogger{l})
type LogrusLogger struct{ *logrus.Logger }

// V returns the result of current level == l
func (s *LogrusLogger) V(l int) bool {
	return s.Level == logrus.Level(l)
}
