
// TestNativeLibrary API
#pragma once

// Portable export/import helpers for shared library builds.
#if defined(_WIN32) || defined(__CYGWIN__)
	#ifdef TNL_BUILD_SHARED
		#define TNL_API __declspec(dllexport)
	#elif defined(TNL_USE_SHARED)
		#define TNL_API __declspec(dllimport)
	#else
		#define TNL_API
	#endif
#else
	#if defined(__GNUC__) && (__GNUC__ >= 4)
		#define TNL_API __attribute__((visibility("default")))
	#else
		#define TNL_API
	#endif
#endif

#ifdef __cplusplus
extern "C" {
#endif

typedef void (*LogCallback)(int logLevel, const char* message, int identifiableInformation);

TNL_API void set_log_callback(LogCallback callback);
TNL_API int factorial(int n);

#ifdef __cplusplus
}
#endif
