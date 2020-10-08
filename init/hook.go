package init

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// SetupLogger function
func SetupLogger(logpath string, mode uint32) *logrus.Entry {

	/* CREATE A FORMATTER */
	formatter := new(logrus.TextFormatter)
	formatter.DisableColors = true
	formatter.ForceQuote = true
	formatter.FullTimestamp = true
	formatter.PadLevelText = true
	formatter.QuoteEmptyFields = true
	formatter.TimestampFormat = time.RFC3339

	newLogger := logrus.New()

	newLogger.SetFormatter(formatter)
	newLogger.SetLevel(logrus.Level(mode))

	if len(logpath) > 0 {
		if _, err := os.Stat(logpath); os.IsNotExist(err) {
			os.MkdirAll(logpath, os.ModePerm)
		}
		newLogger.AddHook(&WriterHook{
			Writer: newLoggerRotate(fmt.Sprintf("%s/errors.log", logpath)),
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
			},
		})
		newLogger.AddHook(&WriterHook{
			Writer: newLoggerRotate(fmt.Sprintf("%s/processing_tasks.log", logpath)),
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
			},
		})
		newLogger.AddHook(&WriterHook{
			Writer: newLoggerRotate(fmt.Sprintf("%s/warnings.log", logpath)),
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		})
		newLogger.AddHook(&WriterHook{
			Writer: newLoggerRotate(fmt.Sprintf("%s/debugs.log", logpath)),
			LogLevels: []logrus.Level{
				logrus.DebugLevel,
			},
		})
	} else {
		newLogger.SetOutput(os.Stdout)
	}
	newLogger.WithField("service", "[ Block - CHECK-CONNECTION ]")
	return logrus.NewEntry(newLogger)
}

func newLoggerRotate(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1024, // MB
		MaxBackups: 7,
		MaxAge:     1, //days
		LocalTime:  true,
		//Compress: true, // Disabled by default
	}
}
