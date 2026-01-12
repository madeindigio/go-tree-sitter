package vue2_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/vue2"
	"github.com/stretchr/testify/assert"
)

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	code := `<template>
  <div>{{ message }}</div>
</template>

<script>
export default {
  data() {
    return {
      message: 'Hello Vue!'
    }
  }
}
</script>

<style>
div {
  color: red;
}
</style>`

	n, err := sitter.ParseCtx(context.Background(), []byte(code), vue2.GetLanguage())
	assert.NoError(err)
	assert.NotNil(n)
	assert.Equal("document", n.Type())
}
