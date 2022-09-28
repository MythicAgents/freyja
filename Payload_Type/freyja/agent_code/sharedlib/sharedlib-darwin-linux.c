#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <wchar.h>
#include <assert.h>
#include <pthread.h>
#include "freyja-darwin-10.06-amd64.h" //Change the header file if something different was used
// To build :
// 1. Build a c-archive in golang: go build -buildmode=c-archive -o freyja-darwin-10.06-amd64.a -tags=[profile] freyja.go
// 2. Execute: ranlib freyja-darwin-10.06-amd64.a
// 3. Build a shared lib (darwin): clang -shared -framework Foundation -framework CoreGraphics -framework Security -framework ApplicationServices -framework OSAKit -fpic sharedlib-darwin-linux.c freyja-darwin-10.06-amd64.a -o freyja.dylib

// Test Dylib execution with python3
// python3
// import ctypes
// ctypes.CDLL("./freyja.dylib")

__attribute__ ((constructor)) void initializer()
{
	pthread_attr_t  attr;
    pthread_t       posixThreadID;
    int             returnVal;

    returnVal = pthread_attr_init(&attr);
    assert(!returnVal);
    returnVal = pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);
    assert(!returnVal);
    pthread_create(&posixThreadID, &attr, &RunMain, NULL);
}
