# Resumen de Implementación: Soporte para Twig y Blade

## Análisis de la Implementación de Rust

La implementación de soporte para lenguajes en go-tree-sitter sigue un patrón estándar:

### Estructura de Directorios
```
language_name/
├── binding.go          # Interfaz Go que expone GetLanguage()
├── binding_test.go     # Tests del parser
├── parser.c           # Parser generado por tree-sitter
├── parser.h           # Header del parser (opcional, algunos lo necesitan)
├── scanner.c          # Scanner léxico (si el lenguaje lo requiere)
└── tree_sitter/       # Headers adicionales (para algunos lenguajes)
    ├── parser.h
    ├── array.h
    └── alloc.h
```

### Patrón de binding.go
```go
package language_name

//#include "parser.h"
//TSLanguage *tree_sitter_language_name();
import "C"
import (
	"unsafe"
	sitter "github.com/madeindigio/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_language_name())
	return sitter.NewLanguage(ptr)
}
```

Para lenguajes con directorios tree_sitter adicionales:
```go
//#cgo CFLAGS: -I. -Itree_sitter
//#include "parser.h"
```

### Registro en grammars.json
Cada lenguaje debe registrarse en `_automation/grammars.json`:
```json
{
  "language": "language_name",
  "url": "https://github.com/user/tree-sitter-language",
  "files": ["parser.c", "scanner.c"],
  "reference": "v1.0.0",
  "revision": "commit_hash",
  "updateBasedOn": "tag"
}
```

## Implementación de Twig ✅

### Estado: COMPLETADO Y FUNCIONAL

### Archivos Añadidos:
- `twig/binding.go`
- `twig/binding_test.go`
- `twig/parser.c` (1.4 MB)
- `twig/parser.h`
- `twig/scanner.c`
- `twig/tree_sitter/parser.h`

### Configuración:
- **Repositorio**: https://github.com/kaermorchen/tree-sitter-twig
- **Versión**: v0.7.0
- **Commit**: 40d17f0eb990215e12531abe29ee7691d7ca99a5
- **Language Version**: Compatible (versión 14)

### Prueba de Funcionamiento:
```bash
cd /www/MCP/Remembrances/go-tree-sitter
go test ./twig -v
# PASS
```

### Ejemplo de Uso:
```go
import (
    sitter "github.com/madeindigio/go-tree-sitter"
    "github.com/madeindigio/go-tree-sitter/twig"
)

parser := sitter.NewParser()
parser.SetLanguage(twig.GetLanguage())

code := []byte(`{{ variable }}`)
tree, _ := parser.ParseCtx(context.Background(), nil, code)
root := tree.RootNode()
// root.String() = "(template (output (variable)))"
```

## Implementación de Blade ❌

### Estado: NO COMPATIBLE

### Problema:
Laravel Blade tiene incompatibilidades fundamentales con la versión de tree-sitter usada por go-tree-sitter:

**Versión v0.12.3 (más reciente):**
- Usa tree-sitter v0.25.6 con **Language Version 15**
- go-tree-sitter soporta hasta **Language Version 14**

**Versiones v0.9.x - v0.11.x (más antiguas):**
- Usan Language Version 14 ✓
- **PERO** usan formato de macro `REDUCE(symbol, children)` con 2 argumentos
- go-tree-sitter requiere `REDUCE(symbol, children, precedence, prod_id)` con 4 argumentos

### Análisis de Versiones Probadas:

| Versión | Language ABI | Formato REDUCE | Compatible |
|---------|-------------|----------------|------------|
| v0.12.3 | 15          | 2 args         | ❌ No      |
| v0.11.0 | 14          | 2 args         | ❌ No      |
| v0.9.2  | 14          | 2 args         | ❌ No      |

### Errores Encontrados:

**v0.12.3:**
```
TREE_SITTER_LANGUAGE_VERSION 14 (máximo soportado)
Blade LANGUAGE_VERSION 15 (requerido)
```

**v0.11.0 y v0.9.2:**
```c
parser.c:327158:71: error: macro "REDUCE" requires 4 arguments, but only 2 given
327158 |   [3] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_blade, 0),
       |                                                                       ^
```

### Comparación de Formatos:

**Blade (todas las versiones):**
```c
REDUCE(sym_blade, 0)
REDUCE(aux_sym_blade_repeat1, 2)
```

**go-tree-sitter (Rust, Python, etc.):**
```c
REDUCE(sym_source_file, 0, 0, 0)
REDUCE(aux_sym_module_repeat1, 2, 0, 0)
```

### Registro en grammars.json:
```json
{
  "language": "blade",
  "url": "https://github.com/EmranMR/tree-sitter-blade",
  "files": ["parser.c", "scanner.c", "tag.h"],
  "reference": "v0.12.3",
  "revision": "cc764dadcbbceb3f259396fef66f970c72e94f96",
  "updateBasedOn": "tag",
  "note": "INCOMPATIBLE: All versions use incompatible REDUCE macro format. v0.12+ requires ABI 15, older versions (v0.9-v0.11) use 2-arg REDUCE vs required 4-arg format."
}
```

### Solución:
No hay solución directa sin modificar significativamente el parser.h o sin que tree-sitter-blade regenere sus parsers con una versión compatible de tree-sitter. Blade requeriría:
1. Regeneración del parser con tree-sitter compatible con go-tree-sitter
2. O actualización completa de go-tree-sitter a tree-sitter 0.25+

## Cambios Realizados

### 1. Directorio twig/ (nuevo)
- Soporte completo para plantillas Twig
- Tests funcionando correctamente

### 2. _automation/grammars.json
- Añadida entrada para Twig (funcional)
- Añadida entrada para Blade (marcada como incompatible)

## Comandos de Verificación

### Verificar Twig:
```bash
go test ./twig -v
```

### Actualizar gramáticas (cuando esté disponible):
```bash
go run _automation/main.go update twig
# Para Blade, esperar actualización de go-tree-sitter primero
```

## Notas Técnicas

### Diferencias entre Versiones de Tree-sitter:
- **v0.20.x**: Language Version 13-14 (go-tree-sitter actual)
- **v0.25.x**: Language Version 15 (Blade requiere esto)

### Archivos Críticos del Proyecto:
- `api.h`: Define TREE_SITTER_LANGUAGE_VERSION (actualmente 14)
- `parser.h`: Header estándar para parsers
- `array.h`: Utilidades para arrays dinámicos (usado por HTML, Svelte, Blade)

### CGo Flags Comunes:
- Sin flags adicionales: La mayoría de lenguajes
- `-I. -Itree_sitter`: Lenguajes con headers en subdirectorio tree_sitter/

## Conclusión

✅ **Twig**: Implementado exitosamente y funcionando.
❌ **Blade**: Preparado pero esperando actualización de go-tree-sitter para soportar Language Version 15.

Los archivos están organizados siguiendo el patrón establecido del proyecto, facilitando futuras actualizaciones.
