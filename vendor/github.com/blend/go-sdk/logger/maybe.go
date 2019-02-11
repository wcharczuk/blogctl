package logger

// MaybeTrigger triggers an event if the logger is set.
func MaybeTrigger(log Triggerable, e Event) {
	if log == nil {
		return
	}
	log.Trigger(e)
}

// MaybeSyncTrigger triggers an event if the logger is set.
func MaybeSyncTrigger(log SyncTriggerable, e Event) {
	if log == nil {
		return
	}
	log.SyncTrigger(e)
}

// MaybeInfof triggers Infof if the logger is set.
func MaybeInfof(log OutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Infof(format, args...)
}

// MaybeSyncInfof triggers SyncInfof if the logger is set.
func MaybeSyncInfof(log SyncOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncInfof(format, args...)
}

// MaybeDebugf triggers Debugf if the logger is set.
func MaybeDebugf(log OutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Debugf(format, args...)
}

// MaybeSyncDebugf triggers SyncDebugf if the logger is set.
func MaybeSyncDebugf(log SyncOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncDebugf(format, args...)
}

// MaybeWarningf triggers Warningf if the logger is set.
func MaybeWarningf(log ErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Warningf(format, args...)
}

// MaybeSyncWarningf triggers SyncWarningf if the logger is set.
func MaybeSyncWarningf(log SyncErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncWarningf(format, args...)
}

// MaybeWarning triggers Warning if the logger is set.
func MaybeWarning(log ErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Warning(err)
}

// MaybeSyncWarning triggers SyncWarning if the logger is set.
func MaybeSyncWarning(log SyncErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.SyncWarning(err)
}

// MaybeErrorf triggers Errorf if the logger is set.
func MaybeErrorf(log ErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Errorf(format, args...)
}

// MaybeSyncErrorf triggers SyncErrorf if the logger is set.
func MaybeSyncErrorf(log SyncErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncErrorf(format, args...)
}

// MaybeError triggers Error if the logger is set.
func MaybeError(log ErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Error(err)
}

// MaybeSyncError triggers SyncError if the logger is set.
func MaybeSyncError(log SyncErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.SyncError(err)
}

// MaybeFatalf triggers Fatalf if the logger is set.
func MaybeFatalf(log ErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.Fatalf(format, args...)
}

// MaybeSyncFatalf triggers SyncFatalf if the logger is set.
func MaybeSyncFatalf(log SyncErrorOutputReceiver, format string, args ...interface{}) {
	if log == nil {
		return
	}
	log.SyncFatalf(format, args...)
}

// MaybeFatal triggers Fatal if the logger is set.
func MaybeFatal(log ErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.Fatal(err)
}

// MaybeSyncFatal triggers SyncFatal if the logger is set.
func MaybeSyncFatal(log SyncErrorReceiver, err error) {
	if log == nil || err == nil {
		return
	}
	log.SyncFatal(err)
}
