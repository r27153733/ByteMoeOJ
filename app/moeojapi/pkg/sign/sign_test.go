package sign

import "testing"

func TestSign(t *testing.T) {
	s := NewSign("test")
	res := s.Sign([]byte("666 test"))
	verify := s.Verify(res[:], []byte("666 test"))
	if verify != true {
		t.Fatal("verify failed")
	}
}
