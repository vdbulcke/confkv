package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(debugEnabled bool) *zap.Logger {

	// l := zap.Must(zap.NewProduction())
	l := zap.Must(getCustomProdLogger())
	if debugEnabled {
		l = zap.Must(zap.NewDevelopment())
	}

	return l
}

func getCustomProdLogger() (*zap.Logger, error) {

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true, // enable stacktrace logger.WithOptions(zap.AddStacktrace())
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			// "pid": os.Getpid(),
		},
	}

	return config.Build()
}
