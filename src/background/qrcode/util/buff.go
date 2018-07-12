package util

import (
	"io"
)

type BuffReadWriter struct {
	Buff []byte
}

func (w *BuffReadWriter) Write(p []byte) (n int, err error) {
	if w.Buff == nil {
		w.Buff = []byte{}
	}
	w.Buff = append(w.Buff, p...)
	return len(w.Buff), nil
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

func (w *BuffReadWriter) WriteByte(c byte) error {
	if w.Buff == nil {
		w.Buff = []byte{}
	}
	w.Buff = append(w.Buff, c)
	return nil
}

func (f *BuffReadWriter) Flush() error {
	return nil
}
