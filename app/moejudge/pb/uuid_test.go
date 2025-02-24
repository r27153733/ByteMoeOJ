package pb

import (
	"fmt"
	"github.com/r27153733/ByteMoeOJ/lib/uuid"
	"testing"
)

func TestName(t *testing.T) {
	a := &UUID{
		Hi: 5437209496599500033,
		Lo: 2251311124549082293,
	}
	id := ToUUID(a)
	pbUUID := ToPbUUID(id)
	if pbUUID.Hi != a.Hi || pbUUID.Lo != a.Lo {
		t.Fatal()
	}
	fmt.Println(id.String())
}

func TestName1(t *testing.T) {

	id := uuid.ParseOrZero("0194ea8e-27ae-7fa0-ae09-a28a77cd984c")
	pbUUID := ToPbUUID(id)
	fmt.Println(pbUUID)
	if id != ToUUID(pbUUID) {
		t.Fatal()
	}
}
