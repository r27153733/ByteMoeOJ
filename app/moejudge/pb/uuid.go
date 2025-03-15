package pb

import (
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"unsafe"
)

func ToUUIDPointer(id *UUID) *uuid.UUID {
	if id == nil {
		return nil
	}

	return (*uuid.UUID)(unsafe.Pointer(&(id.Hi)))
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
