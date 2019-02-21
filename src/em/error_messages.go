package em

import (
	"fmt"

	"github.com/sxarp/c_compiler_go/src/tok"
)

type errMsg struct {
	tt *tok.TokenType
	t  *tok.Token
}

func (e *errMsg) Set(tt *tok.TokenType, t *tok.Token) { e.t, e.tt = t, tt }

func (e *errMsg) Message() string {
	return fmt.Sprintf("Something went wrong around (%d, %d). Parser got [%s], while expecting [%s].",
		e.t.Row, e.t.Col, e.t.Val(), e.tt.Str())
}

var EM = errMsg{}
