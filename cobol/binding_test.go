package cobol_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/cobol"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	code := `IDENTIFICATION DIVISION.
PROGRAM-ID. HELLO-WORLD.`

	n, err := sitter.ParseCtx(context.Background(), []byte(code), cobol.GetLanguage())
	assert.NoError(err)
	assert.NotNil(n)
	assert.Equal("start", n.Type())
}
