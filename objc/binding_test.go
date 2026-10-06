package objc_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/objc"
	"github.com/stretchr/testify/assert"
)

const code = `@interface Foo : NSObject
- (void)bar;
@end`

const expected = `(translation_unit (class_interface (identifier) superclass: (identifier) (method_declaration (method_type (type_name (primitive_type))) (identifier))))`

func TestGrammar(t *testing.T) {
	assert := assert.New(t)

	n, err := sitter.ParseCtx(context.Background(), []byte(code), objc.GetLanguage())
	assert.NoError(err)
	assert.Equal(expected, n.String())
}
