package twig

//#cgo CFLAGS: -I. -Itree_sitter
//#include "parser.h"
//TSLanguage *tree_sitter_twig();
import "C"
import (
	"unsafe"

	sitter "github.com/madeindigio/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_twig())
	return sitter.NewLanguage(ptr)
}
