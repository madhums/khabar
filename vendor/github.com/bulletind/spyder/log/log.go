package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_sentry"
)

var logger = logrus.New()

func Setup(debug bool, sentryDSN string) {
	logger.Formatter = &OurFormatter{}
	logger.Level = logrus.InfoLevel

	if debug {
		logger.Level = logrus.DebugLevel
	}

	if sentryDSN != "" {
		levels := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}

		hook, err := logrus_sentry.NewSentryHook(sentryDSN, levels)
		if err != nil {
			logger.Error("Unable to connect to sentry")
		} else {
			logger.Info("Adding Sentry Hook")
			hook.StacktraceConfiguration.Enable = true
			logger.Hooks.Add(hook)
		}
	}

	if debug {
		logrus.Info("Logging in DEBUG mode")
	}
}

func SetupWithExisting(existingLogger *logrus.Logger) {
	logger = existingLogger
}


// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		logger.Debug(args)
	}
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		logger.Info(args)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		logger.Warn(args)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.Error(args)
	}
}

func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		logger.Fatal(args)
	}
}
