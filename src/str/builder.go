package str

import (
	"bytes"
	"fmt"
)

type Builder struct {
	b bytes.Buffer
}

func (b *Builder) Put(s string) *Builder {
	b.b.WriteString(s)

	return b
}

func (b *Builder) Write(s string) {
	b.Put(fmt.Sprintf("%s\n", s))
}

func (b *Builder) Str() string {
	return b.b.String()
}

func (b *Builder) Nr() {
	b.Write("")
}
