package util

import (
	"bytes"
	"fmt"
)

type BufferWrapper struct {
	*bytes.Buffer
}

func NewBufferWrapper() *BufferWrapper {
	return &BufferWrapper{&bytes.Buffer{}}
}

func (b *BufferWrapper) Write(s string, a ...interface{}) {
	b.WriteString(fmt.Sprintf(s, a...))
}

func (b *BufferWrapper) Writeln(s string, a ...interface{}) {
	b.WriteString(fmt.Sprintf(s+"\n", a...))
}
