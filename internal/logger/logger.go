package logger

import "go.uber.org/zap"

var Log = zap.NewNop()

func Init(level string) error {
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = atomicLevel
	Log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}
