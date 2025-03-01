package uuid

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewUUID(t *testing.T) {
	u1 := NewUUIDPtr()
	u2 := NewUUIDPtr()

	if u1 == nil || u2 == nil {
		t.Fatalf("UUID should not be nil")
	}

	// UUIDs should be unique
	if Equal(u1, u2) {
		t.Fatalf("Generated UUIDs should be unique: %s == %s", u1.String(), u2.String())
	}

	// Check version and variant bits
	if (u1[6] >> 4) != 4 {
		t.Errorf("Expected version 4 UUID, but got version %d", u1[6]>>4)
	}
	if (u1[8] >> 6) != 2 {
		t.Errorf("Expected variant 10, but got %b", u1[8]>>6)
	}
}

func TestUUIDString(t *testing.T) {
	u := NewUUID()
	s := u.String()

	if len(s) != 36 {
		t.Fatalf("UUID string should be 36 characters long, got %d", len(s))
	}

	// Validate hyphen positions
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		t.Errorf("Invalid UUID format: %s", s)
	}
}

func TestParseBytes(t *testing.T) {
	u := NewUUID()
	s := u.String()
	b := []byte(s)

	parsed, err := ParseBytes(b)
	if err != nil {
		t.Fatalf("Failed to parse valid UUID string: %v", err)
	}

	if u != parsed {
		t.Errorf("Parsed UUID does not match original: %s != %s", u.String(), parsed.String())
	}

	// Test invalid formats
	invalidCases := [][]byte{
		[]byte(""),             // Empty string
		[]byte("invalid-uuid"), // Incorrect format
		[]byte("12345678123456781234567812345678"),     // Missing hyphens
		[]byte("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"), // Invalid characters
	}

	for _, invalid := range invalidCases {
		_, err = ParseBytes(invalid)
		if err == nil {
			t.Errorf("Expected error for invalid UUID: %s", invalid)
		}
	}
}

func TestEqual(t *testing.T) {
	u1 := NewUUIDPtr()
	u2 := NewUUIDPtr()
	marshal, err := json.Marshal(&u1)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
	if Equal(u1, u1) == false || Equal(nil, nil) == false {
		t.Errorf("Equal should return true for the same UUID")
	}

	if Equal(u1, u2) == true || Equal(u1, nil) == true || Equal(nil, u2) == true {
		t.Errorf("Equal should return false for different UUIDs")
	}
}

func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u := NewUUIDPtr()
		ReleaseUUIDBuf(u)
	}
}

func TestZeroAlloc(t *testing.T) {
	avg := testing.AllocsPerRun(100, func() {
		u := NewUUIDPtr()
		ReleaseUUIDBuf(u)
	})

	if avg > 0 {
		t.Fatal()
	}
}

func TestZeroAllocSQL(t *testing.T) {
	avg := testing.AllocsPerRun(100, func() {
		u := NewUUID()
		value, err := u.Value()
		if err != nil {
			t.Fatal(err)
		}
		_ = value
		err = u.Scan(value)
		if err != nil {
			t.Fatal(err)
		}
	})

	if avg > 0 {
		fmt.Println(avg)
		t.Fatal()
	}
}
