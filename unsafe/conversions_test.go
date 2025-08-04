package unsafeext

import "testing"

func TestBytesToString(t *testing.T) {
	b := []byte{'e', 'x', 't', 'e', 'n', 'd', 'e', 'r'}
	s := BytesToString(b)
	expected := string(b)
	if s != expected {
		t.Fatalf("expected '%s' got '%s'", expected, s)
	}
}

func TestStringToBytes(t *testing.T) {
	s := "extender"
	b := StringToBytes(s)
	if string(b) != s {
		t.Fatalf("expected '%s' got '%s'", s, string(b))
	}
}
