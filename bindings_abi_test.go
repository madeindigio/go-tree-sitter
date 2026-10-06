package sitter_test

import (
	"context"
	"testing"

	sitter "github.com/madeindigio/go-tree-sitter"
	"github.com/madeindigio/go-tree-sitter/cpp"
	"github.com/madeindigio/go-tree-sitter/csharp"
	"github.com/madeindigio/go-tree-sitter/golang"
	"github.com/madeindigio/go-tree-sitter/java"
	"github.com/madeindigio/go-tree-sitter/kotlin"
	"github.com/madeindigio/go-tree-sitter/objc"
	"github.com/madeindigio/go-tree-sitter/ruby"
	"github.com/madeindigio/go-tree-sitter/rust"
	"github.com/madeindigio/go-tree-sitter/swift"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRuntimeABIRange pins the ABI window of the vendored runtime
// (tree-sitter 0.25.x: ABI 15, still loading ABI 13/14 grammars).
func TestRuntimeABIRange(t *testing.T) {
	assert.Equal(t, 15, sitter.LanguageVersion)
	assert.Equal(t, 13, sitter.MinCompatibleLanguageVersion)
}

// TestGrammarsLoadAndParse loads several grammars on the vendored runtime and
// parses a snippet in each one, checking the ABI accessor and that the tree
// has no syntax errors.
func TestGrammarsLoadAndParse(t *testing.T) {
	cases := []struct {
		name     string
		lang     *sitter.Language
		source   string
		rootType string
	}{
		{"golang", golang.GetLanguage(), "package main\n\nfunc main() { println(\"hi\") }\n", "source_file"},
		{"swift", swift.GetLanguage(), "func greet(name: String) -> String {\n  return \"Hello \\(name)\"\n}\n", "source_file"},
		{"kotlin", kotlin.GetLanguage(), "fun main() {\n    val x = 1\n    println(x)\n}\n", "source_file"},
		{"rust", rust.GetLanguage(), "fn main() {\n    let x: i32 = 1;\n    println!(\"{}\", x);\n}\n", "source_file"},
		{"java", java.GetLanguage(), "class A {\n  int f(int x) { return x + 1; }\n}\n", "program"},
		{"csharp", csharp.GetLanguage(), "class A { int F(int x) => x + 1; }\n", "compilation_unit"},
		{"ruby", ruby.GetLanguage(), "def f(x)\n  x + 1\nend\n", "program"},
		{"cpp", cpp.GetLanguage(), "int f(int x) { return x + 1; }\n", "translation_unit"},
		{"objc", objc.GetLanguage(), "@interface A : NSObject\n- (void)f;\n@end\n", "translation_unit"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			abi := tc.lang.ABIVersion()
			assert.GreaterOrEqual(t, abi, uint32(sitter.MinCompatibleLanguageVersion))
			assert.LessOrEqual(t, abi, uint32(sitter.LanguageVersion))
			if abi < 15 {
				assert.Equal(t, "", tc.lang.Name(), "ABI < 15 languages carry no name")
			} else {
				assert.NotEmpty(t, tc.lang.Name())
			}

			parser := sitter.NewParser()
			defer parser.Close()
			require.NoError(t, parser.TrySetLanguage(tc.lang))

			tree, err := parser.ParseCtx(context.Background(), nil, []byte(tc.source))
			require.NoError(t, err)
			defer tree.Close()

			root := tree.RootNode()
			assert.Equal(t, tc.rootType, root.Type())
			assert.False(t, root.HasError(), "unexpected syntax error: %s", root.String())
			assert.Greater(t, root.NamedChildCount(), uint32(0))
		})
	}
}

// TestQueryOnABI14Grammar runs a query against an ABI 14 grammar to make sure
// the 0.25 query engine still works with older languages.
func TestQueryOnABI14Grammar(t *testing.T) {
	lang := golang.GetLanguage()
	src := []byte("package main\n\nfunc a() {}\nfunc b() {}\n")

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(lang)
	tree, err := parser.ParseCtx(context.Background(), nil, src)
	require.NoError(t, err)
	defer tree.Close()

	q, err := sitter.NewQuery([]byte("(function_declaration name: (identifier) @name)"), lang)
	require.NoError(t, err)
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	qc.Exec(q, tree.RootNode())

	var names []string
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			names = append(names, c.Node.Content(src))
		}
	}
	assert.Equal(t, []string{"a", "b"}, names)
}

// TestOperationLimitParseInput checks that the emulated timeout (progress
// callback) also applies to callback-based parsing and that the parser can be
// reused afterwards.
func TestOperationLimitParseInput(t *testing.T) {
	src := []byte("package main\n")
	for i := 0; i < 2000; i++ {
		src = append(src, []byte("func f() { a := 1; b := a + 2; _ = b }\n")...)
	}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(golang.GetLanguage())
	parser.SetOperationLimit(1)

	read := func(offset uint32, _ sitter.Point) []byte {
		if int(offset) >= len(src) {
			return nil
		}
		end := int(offset) + 1024
		if end > len(src) {
			end = len(src)
		}
		return src[offset:end]
	}
	tree, err := parser.ParseInputCtx(context.Background(), nil, sitter.Input{Read: read, Encoding: sitter.InputEncodingUTF8})
	assert.Nil(t, tree)
	assert.ErrorIs(t, err, sitter.ErrOperationLimit)

	parser.SetOperationLimit(0)
	parser.Reset()
	tree, err = parser.ParseCtx(context.Background(), nil, src)
	require.NoError(t, err)
	defer tree.Close()
	assert.False(t, tree.RootNode().HasError())
}
