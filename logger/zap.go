package logger

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	Log *zap.SugaredLogger
}

func newZapLogger(v *viper.Viper) (Logger, error) {
	if v == nil {
		return nil, ErrEmptyLoggerConfig
	}

	cores := []zapcore.Core{}

	for _, k := range keys(v) {
		c := newZapCore(v.Sub(k))
		if c == nil {
			continue
		}
		cores = append(cores, c)
	}

	if len(cores) == 0 {
		return nil, ErrNoValidLoggerConfig
	}

	return &zapLogger{
		Log: zap.New(zapcore.NewTee(cores...)).Sugar(),
	}, nil
}

func newZapCore(v *viper.Viper) zapcore.Core {
	if v == nil {
		return nil
	}
	encoder := getEncoder(v.Sub("encoder"))
	level := getZapLevel(v.GetString("level"))

	var writer zapcore.WriteSyncer
	if v.GetBool("rotate") {
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   v.GetString("out"),
			MaxSize:    v.GetInt("max-size"),
			MaxAge:     v.GetInt("max-days"),
			MaxBackups: v.GetInt("max-backups"),
			Compress:   true,
		})
	} else {
		w, closeOut, err := zap.Open(v.GetString("out"))
		if err != nil {
			closeOut()
			return nil
		}
		writer = w
	}

	return zapcore.NewCore(encoder, writer, level)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoder(v *viper.Viper) zapcore.Encoder {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = getTimeEncoder(v.GetString("time-format"))

	switch v.GetString("name") {
	case EncoderJson:
		return zapcore.NewJSONEncoder(cfg)
	default:
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(cfg)
	}
}

func getTimeEncoder(f string) zapcore.TimeEncoder {
	switch f {
	case "rfc3339nano", "RFC3339Nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "rfc3339", "RFC3339":
		return zapcore.RFC3339TimeEncoder
	case "iso8601", "ISO8601":
		return zapcore.ISO8601TimeEncoder
	case "millis":
		return zapcore.EpochMillisTimeEncoder
	case "nanos":
		return zapcore.EpochNanosTimeEncoder
	case "epoch":
		return zapcore.EpochTimeEncoder
	default:
		timeFormat = f
		return TimeEncoder
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(timeFormat))
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.Log.Debugf(format, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.Log.Debug(args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.Log.Infof(format, args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.Log.Info(args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.Log.Warnf(format, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.Log.Warn(args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.Log.Errorf(format, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.Log.Error(args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.Log.Fatalf(format, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.Log.Fatal(args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.Log.Fatalf(format, args...)
}

func (l *zapLogger) Panic(args ...interface{}) {
	l.Log.Panic(args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.Log.With(f...)
	return &zapLogger{newLogger}
}
