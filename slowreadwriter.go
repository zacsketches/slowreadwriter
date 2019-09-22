package slowreadwriter

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// SlowReadWriter is a tool to test slow reads on serial connections
// It implements the ReadWriter interface, but calls to Read are delayed
// by a random selection from an integer slice of millisecond delays passed
// to the SlowReadWriter at construction with the NewSlowReadWriter
// function.
type SlowReadWriter struct {
	p      []byte
	Delays []int
}

// NewSlowReadWriter is the preferred way to constructe the type.
func NewSlowReadWriter(delays []int) *SlowReadWriter {
	srw := &SlowReadWriter{
		Delays: delays,
	}

	return srw
}

func (sr *SlowReadWriter) Write(data []byte) (n int, err error) {
	l := len(sr.p)
	if l+len(data) > cap(sr.p) {
		// reallocate to double what's needed for future growth
		newSlice := make([]byte, (l+len(data))*2)
		copy(newSlice, sr.p)
		sr.p = newSlice
	}
	sr.p = sr.p[0 : l+len(data)]
	copy(sr.p[l:], data)

	return len(data), nil
}

func (sr SlowReadWriter) Read(p []byte) (n int, err error) {
	//find an index into the sr delays slice
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(sr.Delays))
	delay := sr.Delays[idx]
	t := time.NewTimer(time.Duration(delay) * time.Millisecond)
	<-t.C //block until the timer returns

	n = 0
	//copy some timer info into the read
	msg := fmt.Sprintf("%d-", delay)
	for i, b := range []byte(msg) {
		p[i] = b
		n++
	}
	//copy the sr bytes buffer into the read
	for i, b := range sr.p {
		p[i+len(msg)] = b
		n++
	}
	return n, io.EOF
}

// PrintBuffer is a direct copy of the SlowReadWriter bytes buffer
// to stdout.  It does not include the delay.
func (sr SlowReadWriter) PrintBuffer() {
	io.Copy(os.Stdout, bytes.NewReader(sr.p))
}

// PrintBufferln is a direct copy of the SlowReadWriter bytes buffer
// to stdout followed by a new line.  It does not include the delay.
func (sr SlowReadWriter) PrintBufferln() {
	sr.PrintBuffer()
	io.Copy(os.Stdout, bytes.NewReader([]byte{'\n'}))
}
