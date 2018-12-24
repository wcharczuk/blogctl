package logger

// MaybeTrigger triggers an event if the logger is set.
func MaybeTrigger(log *Logger, e Event) {
	if log == nil {
		return
	}
	log.Trigger(e)
}

// MaybeSyncTrigger triggers an event if the logger is set.
func MaybeSyncTrigger(log *Logger, e Event) {
	if log == nil {
		return
	}
	log.Trigger(e)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Infof(format, args...)
}

// MaybeSyncInfof triggers SyncInfof if the logger is set.
func MaybeSyncInfof(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncInfof(format, args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debugf(format, args...)
}

// MaybeSyncDebugf triggers SyncDebugf if the logger is set.
func MaybeSyncDebugf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncDebugf(format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Warningf(format, args...)
}

// MaybeSyncWarningf triggers SyncWarningf if the logger is set.
func MaybeSyncWarningf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncWarningf(format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log *Logger, err error) {
	if log == nil {
		return
	}
	log.Warning(err)
}

// MaybeSyncWarning triggers SyncWarning if the logger is set.
func MaybeSyncWarning(log *Logger, err error) {
	if log == nil {
		return
	}
	log.SyncWarning(err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Errorf(format, args...)
}

// MaybeSyncErrorf triggers SyncErrorf if the logger is set.
func MaybeSyncErrorf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncErrorf(format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log *Logger, err error) {
	if log == nil {
		return
	}
	log.Error(err)
}

// MaybeSyncError triggers SyncError if the logger is set.
func MaybeSyncError(log *Logger, err error) {
	if log == nil {
		return
	}
	log.SyncError(err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeSyncFatalf triggers SyncFatalf if the logger is set.
func MaybeSyncFatalf(log *Logger, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncFatalf(format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log *Logger, err error) {
	if log == nil {
		return
	}
	log.Fatal(err)
}

// MaybeSyncFatal triggers SyncFatal if the logger is set.
func MaybeSyncFatal(log *Logger, err error) {
	if log == nil {
		return
	}
	log.SyncFatal(err)
}
