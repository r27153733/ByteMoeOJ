package unsafetool

import "unsafe"

type AnyPointer[T any] struct {
	typ unsafe.Pointer
	p   *T
}

func AnyToPointer[T any](a any) *T {
	return (*AnyPointer[T])(unsafe.Pointer(&a)).p
}

func PointerAnyToAny[T any](a *any) {
	ap := (*AnyPointer[T])(unsafe.Pointer(a))
	p := ap.p
	*a = *new(T)
	ap.p = p
}

func AnyToPointerAny[T any](a *any) {
	p := AnyToPointer[T](a)
	*a = p
}
