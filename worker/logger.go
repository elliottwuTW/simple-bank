package worker

import "context"

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

// func (logger *Logger) Print(level zerolog.Level, args ...interface{}) {
// 	log.WithLevel(level).Msg(fmt.Sprint(args...))
// }

// 為了符合給 Redis.SetLogger 時用的方法
func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	// log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *Logger) Debug(args ...interface{}) {
	// logger.Print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	// logger.Print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	// logger.Print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	// logger.Print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	// logger.Print(zerolog.FatalLevel, args...)
}
