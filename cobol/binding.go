package cobol

//#include "parser.h"
//TSLanguage *tree_sitter_COBOL();
import "C"
import (
	"unsafe"

	sitter "github.com/madeindigio/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_COBOL())
	return sitter.NewLanguage(ptr)
}
