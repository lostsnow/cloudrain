package telnet

import (
	"fmt"

	"github.com/litsea/logger"
	"go.uber.org/atomic"
)

type SessionTracer interface {
	Created(s *Session)
	Closed(s *Session)
}

type SessionTrace struct {
	seqNo  *atomic.Int64
	active *atomic.Int64
}

func NewSessionTrace() *SessionTrace {
	return &SessionTrace{
		seqNo:  atomic.NewInt64(0),
		active: atomic.NewInt64(0),
	}
}

func (t *SessionTrace) Created(s *Session) {
	t.seqNo.Inc()
	t.active.Inc()
	logger.Infof("session %s started %s", s.RemoteIp, plural(t.active.Load(), t.seqNo.Load()))
}

func (t *SessionTrace) Closed(s *Session) {
	t.active.Dec()
	logger.Infof("session %s ended %s", s.RemoteIp, plural(t.active.Load(), t.seqNo.Load()))
}

func plural(value, total int64) string {
	if value == 0 {
		return fmt.Sprintf("(no active sessions, total %d)", total)
	} else {
		return fmt.Sprintf("(%d active sessions, total %d)", value, total)
	}
}
