package log

type Loggable struct {
	log *Logger
}

func (l *Loggable) Log() *Logger {
	if l.log == nil {
		l.log = newLogger(2)
	}
	return l.log
}
