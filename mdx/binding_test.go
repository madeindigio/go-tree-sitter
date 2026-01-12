package mdx_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/mdx"
	"github.com/stretchr/testify/assert"
)

func TestCanLoadGrammar(t *testing.T) {
	language := mdx.GetLanguage()
	if language == nil {
		t.Errorf("Error loading MDX grammar")
	}
}

func TestMDX(t *testing.T) {
	assert := assert.New(t)

	parser := sitter.NewParser()
	parser.SetLanguage(mdx.GetLanguage())

	content := []byte(`# Hello MDX

import Component from './Component'

<Component prop="value">
  This is **MDX** content
</Component>
`)

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	assert.NoError(err)
	assert.NotNil(tree)
	assert.NotNil(tree.RootNode())

	root := tree.RootNode()
	assert.Equal("document", root.Type())
	assert.True(root.ChildCount() > 0)
}

func TestMDXWithJSX(t *testing.T) {
	assert := assert.New(t)

	parser := sitter.NewParser()
	parser.SetLanguage(mdx.GetLanguage())

	content := []byte(`export const Thing = () => <div>Hello</div>

<Thing />
`)

	tree, err := parser.ParseCtx(context.Background(), nil, content)
	assert.NoError(err)
	assert.NotNil(tree)

	root := tree.RootNode()
	assert.Equal("document", root.Type())
}
