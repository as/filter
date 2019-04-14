package uniq

import "testing"

func TestCtr(t *testing.T) {
	f := NewCtr()
	if f.Inc("a", "b", "c") != 1 {
		t.Fatal("max isn't 1")
	}
	f.Dec("a", "b", "c")
	if f.Inc("a", "b", "c") != 1 {
		t.Fatal("max isn't 1")
	}
	if f.Inc("a", "b", "c") != 2 {
		t.Fatal("max isn't 2")
	}
	if f.Inc("c") != 3 {
		t.Fatal("max isn't 3")
	}
	f.Dec("c")
	f.Dec("a", "b", "c")
	f.Dec("a", "b", "c")
	if f.Inc("a", "b", "c") != 1 {
		t.Fatal("max isn't 1")
	}
}
