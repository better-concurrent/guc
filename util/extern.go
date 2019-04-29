package util

import (
	"unsafe"
)

// Active spinning runtime support.
// runtime_canSpin returns true is spinning makes sense at the moment.
//go:linkname SyncRuntimeCanSpin sync.runtime_canSpin
func SyncRuntimeCanSpin(i int) bool

// runtime_doSpin does active spinning.
//go:linkname SyncRuntimeDoSpin sync.runtime_doSpin
func SyncRuntimeDoSpin()

// nanotime
//go:linkname SyncRuntimeNanoTime sync.runtime_nanotime
func SyncRuntimeNanoTime() int64

// Semacquire waits until *s > 0 and then atomically decrements it.
// It is intended as a simple sleep primitive for use by the synchronization
// library and should not be used directly.
//go:linkname SyncRuntimeSemacquire sync.runtime_Semacquire
func SyncRuntimeSemacquire(s *uint32)

// SemacquireMutex is like Semacquire, but for profiling contended Mutexes.
// If lifo is true, queue waiter at the head of wait queue.
//go:linkname SyncRuntimeSemacquireMutex sync.runtime_SemacquireMutex
func SyncRuntimeSemacquireMutex(s *uint32, lifo bool)

// Semrelease atomically increments *s and notifies a waiting goroutine
// if one is blocked in Semacquire.
// It is intended as a simple wakeup primitive for use by the synchronization
// library and should not be used directly.
// If handoff is true, pass count directly to the first waiter.
//go:linkname SyncRuntimeSemrelease sync.runtime_Semrelease
func SyncRuntimeSemrelease(s *uint32, handoff bool)

// ====== hash func ========
// var algarray = [alg_max]typeAlg{
//	alg_NOEQ:     {nil, nil},
//	alg_MEM0:     {memhash0, memequal0},
//	alg_MEM8:     {memhash8, memequal8},
//	alg_MEM16:    {memhash16, memequal16},
//	alg_MEM32:    {memhash32, memequal32},
//	alg_MEM64:    {memhash64, memequal64},
//	alg_MEM128:   {memhash128, memequal128},
//	alg_STRING:   {strhash, strequal},
//	alg_INTER:    {interhash, interequal},
//	alg_NILINTER: {nilinterhash, nilinterequal},
//	alg_FLOAT32:  {f32hash, f32equal},
//	alg_FLOAT64:  {f64hash, f64equal},
//	alg_CPLX64:   {c64hash, c64equal},
//	alg_CPLX128:  {c128hash, c128equal},
//}

//go:linkname Memhash0 runtime.memhash0
func Memhash0(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Memhash8 runtime.memhash8
func Memhash8(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Memhash16 runtime.memhash16
func Memhash16(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Memhash32 runtime.memhash32
func Memhash32(p unsafe.Pointer, seed uintptr) uintptr

//go:linkname Memhash64 runtime.memhash64
func Memhash64(p unsafe.Pointer, seed uintptr) uintptr

//go:linkname Memhash128 runtime.memhash128
func Memhash128(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Strhash runtime.strhash
func Strhash(a unsafe.Pointer, h uintptr) uintptr

//go:linkname Interhash runtime.interhash
func Interhash(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Nilinterhash runtime.nilinterhash
func Nilinterhash(p unsafe.Pointer, h uintptr)

//go:linkname F32hash runtime.f32hash
func F32hash(p unsafe.Pointer, h uintptr)

//go:linkname F64hash runtime.f64hash
func F64hash(p unsafe.Pointer, h uintptr)

//go:linkname C64hash runtime.c64hash
func C64hash(p unsafe.Pointer, h uintptr) uintptr

//go:linkname C128hash runtime.c128hash
func C128hash(p unsafe.Pointer, h uintptr) uintptr

//go:linkname Memequal0 runtime.memequal0
func Memequal0(p, q unsafe.Pointer) bool

//go:linkname Memequal8 runtime.memequal8
func Memequal8(p, q unsafe.Pointer) bool

//go:linkname Memequal16 runtime.memequal16
func Memequal16(p, q unsafe.Pointer) bool

//go:linkname Memequal32 runtime.memequal32
func Memequal32(p, q unsafe.Pointer) bool

//go:linkname Memequal64 runtime.memequal64
func Memequal64(p, q unsafe.Pointer) bool

//go:linkname Memequal128 runtime.memequal128
func Memequal128(p, q unsafe.Pointer) bool

//go:linkname Strequal runtime.strequal
func Strequal(p, q unsafe.Pointer) bool

//go:linkname Interequal runtime.interequal
func Interequal(p, q unsafe.Pointer) bool

//go:linkname Nilinterequal runtime.nilinterequal
func Nilinterequal(p, q unsafe.Pointer) bool

//go:linkname F32equal runtime.f32equal
func F32equal(p, q unsafe.Pointer) bool

//go:linkname F64equal runtime.f64equal
func F64equal(p, q unsafe.Pointer) bool

//go:linkname C64equal runtime.c64equal
func C64equal(p, q unsafe.Pointer) bool

//go:linkname C128equal runtime.c128equal
func C128equal(p, q unsafe.Pointer) bool
