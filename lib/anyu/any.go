package anyu

import "unsafe"

type AnyPointer[T any] struct {
	typ unsafe.Pointer
	p   *T
}

func AnyToPointer[T any](a any) *T {
	return (*AnyPointer[T])(unsafe.Pointer(&a)).p
}

// PointerToValueAny 使用指向堆上的值的指针构建 any，避免堆分配。不支持 map、func、chan。
func PointerToValueAny[T any](p *T) (res any) {
	res = *new(T)
	resP := (*AnyPointer[T])(unsafe.Pointer(&res))
	resP.p = p
	return res
}
