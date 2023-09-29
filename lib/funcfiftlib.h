#pragma once

#include <stdint.h>
#include <stdbool.h>



/// Callback used to retrieve additional source files or data.
///
/// @param _kind The kind of callback (a string).
/// @param _data The data for the callback (a string).
/// @param o_contents A pointer to the contents of the file, if found. Allocated via malloc().
/// @param o_error A pointer to an error message, if there is one. Allocated via malloc().
///
/// The callback implementor must use malloc() to allocate storage for
/// contents or error. The callback implementor must use free() to free
/// said storage after func_compile returns.
///
/// If the callback is not supported, *o_contents and *o_error must be set to NULL.
typedef void (*CStyleReadFileCallback)(char* _kind, char* _data, char** o_contents, char** o_error);

const char *func_compile(char *config_json, CStyleReadFileCallback callback);


extern void fileReader(char* _kind, char* _data, char** o_contents, char** o_error);
static inline const char * compile(char *config_json) {
    return func_compile(config_json, fileReader);
}
