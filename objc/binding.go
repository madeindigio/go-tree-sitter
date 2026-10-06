package objc

//#include "parser.h"
//const TSLanguage *tree_sitter_objc(void);
import "C"
import (
	"unsafe"

	sitter "github.com/madeindigio/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_objc())
	return sitter.NewLanguage(ptr)
}
