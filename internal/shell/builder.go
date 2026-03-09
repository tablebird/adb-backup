package shell

import (
	"strings"
)

type Builder interface {
	Build() string
}

type cmdBuilder struct {
	cmd []string
}

func (b *cmdBuilder) append(cmd string) Builder {
	b.cmd = append(b.cmd, cmd)
	return b
}

func (b *cmdBuilder) multiAppend(cmd string) Builder {
	if len(b.cmd) > 0 {
		b.cmd = append(b.cmd, "&&", cmd)
	} else {
		b.cmd = append(b.cmd, cmd)
	}
	return b
}

func (b *cmdBuilder) Build() string {
	return strings.Join(b.cmd, " ")
}

func RunBuilder(s Shell, builder Builder) (string, error) {
	cmd := builder.Build()
	// bu, ok := builder.(cmdBuilder)

	return s.RunCommand(cmd)
}
