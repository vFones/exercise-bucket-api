package logger

import (
	"context"
	"encoding/json"
	"os"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TraceIdKey string

const TraceId TraceIdKey = "traceIdKey"

// InitLogger
//
// initializes the logger with the given mode.
// If devMode is true, the logger will use zapcore.DebugLevel level, otherwise zapcore.InfoLevel level.
func InitLogger(devMode bool) {
	loggerLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if devMode {
		loggerLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	stdoutCore := func() zapcore.Core {
		stdoutSyncer := zapcore.Lock(os.Stdout)
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig()),
			stdoutSyncer,
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.ErrorLevel // Just Debug, Info, Warn
			}),
		)
	}

	stderrCore := func() zapcore.Core {
		stderrSyncer := zapcore.Lock(os.Stderr)
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig()),
			stderrSyncer,
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel // Just Error, DPanic, Panic, Fatal
			}),
		)
	}

	logger, err := config(loggerLevel).Build(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel), zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(stdoutCore(), stderrCore())
	}))
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func config(level zap.AtomicLevel) *zap.Config {
	return &zap.Config{
		Level:             level,
		Development:       level.Level() == zap.DebugLevel,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      "json",
		EncoderConfig: encoderConfig(),
		OutputPaths:   []string{"stdout"},
	}
}

// createZapFields
//
// internally used to transform fields to zap.Field
func createZapFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0)
	for _, field := range fields {
		switch f := field.(type) {
		case *LogValue:
			zapFields = append(zapFields, f.zapIt())
		case LogValue:
			zapFields = append(zapFields, f.zapIt())
		case zap.Field:
			zapFields = append(zapFields, f)
		case error:
			zapFields = append(zapFields, zap.Any("error", f.Error()))
		case string:
			var jsonObj map[string]interface{}
			if json.Unmarshal([]byte(f), &jsonObj) == nil {
				zapFields = append(zapFields, zap.Any("json", jsonObj))
			} else {
				zapFields = append(zapFields, zap.String("string", f))
			}
		default:
			typ := reflect.TypeOf(f)
			if typ == nil {
				continue
			}
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			zapFields = append(zapFields, zap.Any(typ.Name(), f))
		}
	}
	return zapFields
}

// LogValue
// Struct used to LogValue in a Json way:
//
//	{"key": value}
//
// see logger.Info for more information
type LogValue struct {
	Value interface{}
	Key   string
}

// NewLogValue
// create a logger.LogValue
func NewLogValue(key string, value interface{}) *LogValue {
	return &LogValue{Key: key, Value: value}
}

// zapIt
//
// internal function to transform logger.LogValue to zap.Field
func (lv *LogValue) zapIt() zap.Field {
	return zap.Any(lv.Key, lv.Value)
}

// Info logs an informational message.
//
// ctx is mandatory; use context.TODO() if unavailable.
//
// msg is mandatory; fill with a significant message string.
//
// fields are optional; when using them, DO NOT use primitive values but wrap them in logger.NewLogValue().
// Use any of logger.LogValue, error, or struct.
//
// Example:
//
//	ctx := context.WithValue(context.Background(), logger.TraceId, "12345")
//	logger.Info(ctx, "User logged in", logger.NewLogValue("userId", 42))
//	logger.Info(ctx, "test 1")
//
//	type testStruct struct {
//	    ValueOne   string `json:"value_one"`
//	    ValueTwo   int    `json:"value_two"`
//	    ValueThree bool   `json:"value_three"`
//	}
//
//	logger.Info(ctx, "test 2", testStruct{ValueOne: "value 1", ValueTwo: 123, ValueThree: true}) // struct are automatically tagged with struct name
//	logger.Info(ctx, "test 3", errors.New("New error")) // error are automatically tagged with "error" key
//
// Output:
//
//	{"level":"INFO","timestamp":"2025-03-10T10:36:37.958+0100","caller":"plain/main.go:43","message":"User logged in","userId":42,"logId":"23a3f92e-f36d-480f-b169-aa2dbe6465e7"}
//	{"level":"INFO","timestamp":"2025-03-10T10:36:37.958+0100","caller":"plain/main.go:44","message":"test 1","logId":"23a3f92e-f36d-480f-b169-aa2dbe6465e7"}
//	{"level":"INFO","timestamp":"2025-03-10T10:36:37.958+0100","caller":"plain/main.go:46","message":"test 2","testStruct":{"value_one":"value 1","value_two":123,"value_three":true},"logId":"23a3f92e-f36d-480f-b169-aa2dbe6465e7"}
//	{"level":"INFO","timestamp":"2025-03-10T10:36:37.958+0100","caller":"plain/main.go:47","message":"test 3","error":"new error","logId":"23a3f92e-f36d-480f-b169-aa2dbe6465e7"}
func Info(ctx context.Context, message string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().Info(message, appendRequestId(ctx, zapFields...)...)
}

// Infof logs a formatted informational message.
//
// msg is mandatory; fill with a significant message string.
//
// fields are optional;
// Example:
//
//	world := "radical"
//	logger.Infof("hello %s", world)
//
// Output:
//
//	{"level":"INFO","timestamp":"2025-03-10T10:56:10.381+0100","caller":"plain/main.go:50","message":"hello radical"}
func Infof(msg string, fields ...interface{}) {
	zap.S().Infof(msg, fields...)
}

// InfoNoCaller
// Just avoid appending "caller" in log
// see logger.Info for more information
//
// Example:
//
//	logger.InfoNoCaller(ctx, "hello radical")
//
// Output:
//
//	{"level":"INFO","timestamp":"2025-03-10T10:58:17.017+0100","message":"hello radical","logId":"1f7f1abe-721a-4d93-a57c-52bca6f10b44"}
func InfoNoCaller(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().WithOptions(zap.WithCaller(false)).Info(msg, appendRequestId(ctx, zapFields...)...)
}

// Debug
// output logs to debug level
//
// see logger.Info for more information
func Debug(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().Debug(msg, appendRequestId(ctx, zapFields...)...)
}

// Error
// output logs to error level
//
// see logger.Error for more information
func Error(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().Error(msg, appendRequestId(ctx, zapFields...)...)
}

// ErrorStackSkipNoCaller
// output logs to error level
//
// # No caller or Stack Trace are appended to log
//
// see logger.Info for more information
func ErrorStackSkipNoCaller(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().WithOptions(zap.AddStacktrace(zap.FatalLevel), zap.WithCaller(false)).Error(msg, appendRequestId(ctx, zapFields...)...)
}

// Fatal
// exit program and output logs to fatal level
//
// see logger.Info for more information
func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	zapFields := createZapFields(fields...)
	zap.L().Fatal(msg, appendRequestId(ctx, zapFields...)...)
}

func Sync() error {
	return zap.S().Sync()
}

// appendRequestId
// used internally to append "logId" from context if previously set
func appendRequestId(ctx context.Context, fields ...zap.Field) []zap.Field {
	if ctx == nil {
		return fields
	}
	if fields == nil {
		fields = make([]zap.Field, 0)
	}
	traceId := ctx.Value(TraceId)
	if traceId != nil {
		fields = append(fields, zap.Any("logId", traceId))
	}
	return fields
}
