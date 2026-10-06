#include "api.h"
#include "bindings.h"
#include "atomic.h"
#include "clock.h"
#include <string.h>
#include <stdio.h>

static void stderr_log(void *payload, TSLogType type, const char *msg)
{
    bool include_lexing = (bool)payload;
    switch (type)
    {
    case TSLogTypeParse:
        fprintf(stderr, "* %s\n", msg);
        break;
    case TSLogTypeLex:
        if (include_lexing)
            fprintf(stderr, "  %s\n", msg);
        break;
    }
}

TSLogger stderr_logger_new(bool include_lexing)
{
    TSLogger result;
    result.payload = (void *)include_lexing;
    result.log = stderr_log;
    return result;
}

/* Progress state shared with the parse progress callback. */
typedef struct
{
    const size_t *cancel_flag;
    TSClock end_clock;
} ProgressState;

static bool progress_callback(TSParseState *state)
{
    ProgressState *p = state->payload;
    if (p->cancel_flag && atomic_load(p->cancel_flag))
    {
        return true;
    }
    if (!clock_is_null(p->end_clock) && clock_is_gt(clock_now(), p->end_clock))
    {
        return true;
    }
    return false;
}

static TSTree *parse_with_limits(TSParser *self, const TSTree *old_tree, TSInput input, ParseLimits limits)
{
    ProgressState state = {limits.cancel_flag, clock_null()};
    if (limits.timeout_micros > 0)
    {
        state.end_clock = clock_after(clock_now(), duration_from_micros(limits.timeout_micros));
    }
    TSParseOptions options = {&state, progress_callback};
    return ts_parser_parse_with_options(self, old_tree, input, options);
}

const char *call_callReadFunc(void *payload, uint32_t byte_index, TSPoint position, uint32_t *bytes_read)
{
    ParsePayload *p = payload;
    if (p->previous_content != NULL)
    {
        free(p->previous_content);
    }
    p->previous_content = callReadFunc(p->read_function_id, byte_index, position, bytes_read);
    return p->previous_content;
}

TSTree *call_ts_parser_parse(TSParser *self, const TSTree *old_tree, int read_function_id, TSInputEncoding encoding, ParseLimits limits)
{
    ParsePayload payload = {read_function_id, NULL};
    TSInput input = {&payload, call_callReadFunc, encoding, NULL};
    TSTree *tree = parse_with_limits(self, old_tree, input, limits);
    if (payload.previous_content != NULL)
    {
        free(payload.previous_content);
    }
    return tree;
}

typedef struct
{
    const char *string;
    uint32_t length;
} StringInput;

static const char *string_input_read(void *payload, uint32_t byte, TSPoint point, uint32_t *length)
{
    (void)point;
    StringInput *self = payload;
    if (byte >= self->length)
    {
        *length = 0;
        return "";
    }
    *length = self->length - byte;
    return self->string + byte;
}

TSTree *call_ts_parser_parse_string(TSParser *self, const TSTree *old_tree, const char *string, uint32_t length, ParseLimits limits)
{
    StringInput payload = {string, length};
    TSInput input = {&payload, string_input_read, TSInputEncodingUTF8, NULL};
    return parse_with_limits(self, old_tree, input, limits);
}
