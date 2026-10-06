#ifndef TREE_SITTER_BINDINGS_H_
#define TREE_SITTER_BINDINGS_H_

#include "api.h"

TSLogger stderr_logger_new(bool include_lexing);

typedef struct
{
    int read_function_id;
    char *previous_content;
} ParsePayload;

// ParseLimits emulates the parser timeout and cancellation flag that were
// deprecated in tree-sitter 0.25 (and removed in 0.26) on top of the
// ts_parser_parse_with_options progress callback.
//
// cancel_flag points to C-allocated memory owned by the Go Parser; a non-zero
// value halts the parse. timeout_micros == 0 means "no limit".
typedef struct
{
    const size_t *cancel_flag;
    uint64_t timeout_micros;
} ParseLimits;

extern char *callReadFunc(int id, uint32_t byteIndex, TSPoint position, uint32_t *bytesRead);
TSTree *call_ts_parser_parse(TSParser *self, const TSTree *old_tree, int read_function_id, TSInputEncoding encoding, ParseLimits limits);
TSTree *call_ts_parser_parse_string(TSParser *self, const TSTree *old_tree, const char *string, uint32_t length, ParseLimits limits);

#endif
