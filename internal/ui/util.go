package ui

import (
	"bytes"
	"fmt"
)

type bufferWrapper struct {
	*bytes.Buffer
}

func newBufferWrapper() *bufferWrapper {
	return &bufferWrapper{&bytes.Buffer{}}
}

func (b *bufferWrapper) writeln(s string, a ...interface{}) {
	b.WriteString(fmt.Sprintf(s+"\n", a...))
}
