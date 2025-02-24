package pb

import (
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"reflect"
	"unsafe"
)

func init() {
	t := reflect.TypeOf(UUID{})
	field, _ := t.FieldByName("Hi")
	Offset = field.Offset
}

var Offset uintptr

func ToUUIDPointer(id *UUID) *uuid.UUID {
	if id == nil {
		return nil
	}
	return (*uuid.UUID)((unsafe.Pointer)((uintptr)(unsafe.Pointer(id)) + Offset))
}

func ToUUID(id *UUID) uuid.UUID {
	if id == nil {
		return uuid.UUID{}
	}
	return *ToUUIDPointer(id)
}

func ToPbUUID(id uuid.UUID) *UUID {
	p := new(UUID)
	SetUUID(p, id)
	return p
}

func SetUUID(dst *UUID, src uuid.UUID) {
	p := ToUUIDPointer(dst)
	*p = src
}
