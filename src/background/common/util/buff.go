package util

import (
	"io"
)

type BuffReadWriter struct {
	Buff []byte
}

func (w *BuffReadWriter) Write(p []byte) (n int, err error) {
	w.Buff = make([]byte, len(p))
	l := copy(w.Buff, p[:])
	return l, nil
}

func (r *BuffReadWriter) Read(p []byte) (n int, err error) {
	l := 0
	buf := len(p)
	tal := len(r.Buff)
	if buf > tal {
		l = copy(p, r.Buff[:tal])
		return l, io.EOF
	} else {
		l = copy(p, r.Buff[:buf])
		r.Buff = r.Buff[buf:]
	}
	return l, nil
}

func (r *BuffReadWriter) Close() error {
	return nil
}
