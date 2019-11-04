package telnet

import (
	"io"
	"log"
	"sync"

	"github.com/lostsnow/cloudrain/telnet/internal"
)

type MultiWriterEntry struct {
	writer io.WriteCloser
}

type multiWriter struct {
	entries []MultiWriterEntry
	lock    sync.RWMutex
}

func NewMultiWriterEntry(writer io.WriteCloser) *MultiWriterEntry {
	return &MultiWriterEntry{writer}
}

func NewMultiWriter(writer io.WriteCloser) *multiWriter {
	return &multiWriter{entries: []MultiWriterEntry{{writer}}}
}

func (w *multiWriter) attach(writer io.WriteCloser) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if len(w.entries) >= MaxSessionConnections {
		log.Println("max connections for session reached.")
		return ErrMaxConnection
	}

	w.entries = append(w.entries, MultiWriterEntry{writer})
	return nil
}

func (w *multiWriter) Write(p []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	entries := w.entries
	var err error
	var lstCnt, cnt int

	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		cnt, err = e.writer.Write(p)
		if err != nil {
			if err := e.writer.Close(); err != nil {
				internal.Log.Println(err)
			}
			entries[i] = entries[len(entries)-1]
			entries = entries[0 : len(entries)-1]
			internal.Log.Println("detaching session.")
		} else {
			lstCnt = cnt
		}
	}

	w.entries = entries

	if len(entries) == 0 {
		return cnt, err
	}

	return lstCnt, nil
}

func (w *multiWriter) Close() {
	for _, e := range w.entries {
		if err := e.writer.Close(); err != nil {
			internal.Log.Println(err)
		}
	}
	w.entries = nil
}
