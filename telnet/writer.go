package telnet

import (
	"io"
	"sync"

	"github.com/lostsnow/cloudrain/telnet/internal"
)

type writerEntry struct {
	writer io.WriteCloser
	lock   sync.RWMutex
}

func NewWriterEntry(writer io.WriteCloser) *writerEntry {
	return &writerEntry{writer: writer}
}

func (w *writerEntry) Write(p []byte) (int, error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	cnt, err := w.writer.Write(p)
	if err != nil {
		if err := w.writer.Close(); err != nil {
			internal.Log.Println(err)
		}
		internal.Log.Println("detaching session.")
		return 0, err
	}

	return cnt, nil
}

func (w *writerEntry) Close() {
	if err := w.writer.Close(); err != nil {
		internal.Log.Println(err)
	}
}
