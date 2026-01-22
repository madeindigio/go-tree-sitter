package twig_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/twig"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	code := []byte(`{{ variable }}`)
	n, err := sitter.ParseCtx(context.Background(), code, twig.GetLanguage())
	assert.NoError(err)
	assert.NotNil(n)
}
