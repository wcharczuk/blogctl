package logger

import (
	"context"
)

// MaybeSet returns if the logger instance is set.
func MaybeSet(log interface{}) bool {
	if log == nil {
		return false
	}
	if typed, ok := log.(*Logger); ok {
		return typed != nil
	}
	return true
}

// MaybeTrigger triggers an event if the logger is set.
func MaybeTrigger(ctx context.Context, log Triggerable, e Event) {
	if !MaybeSet(log) {
		return
	}
	log.Trigger(ctx, e)
}

// MaybeInfo triggers Info if the logger is set.
func MaybeInfo(log InfoReceiver, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Info(args...)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log InfofReceiver, format string, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Infof(format, args...)
}

// MaybeDebug triggers Debug if the logger is set.
func MaybeDebug(log DebugReceiver, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Debug(args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log DebugfReceiver, format string, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Debugf(format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log WarningfReceiver, format string, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Warningf(format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log WarningReceiver, err error) {
	if !MaybeSet(log) || err == nil {
		return
	}
	log.Warning(err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log ErrorfReceiver, format string, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Errorf(format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log ErrorReceiver, err error) {
	if !MaybeSet(log) || err == nil {
		return
	}
	log.Error(err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log FatalfReceiver, format string, args ...interface{}) {
	if !MaybeSet(log) {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log FatalReceiver, err error) {
	if !MaybeSet(log) || err == nil {
		return
	}
	log.Fatal(err)
}
