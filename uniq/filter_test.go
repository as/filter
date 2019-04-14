package uniq

import "testing"

func TestFilter(t *testing.T) {
	f := NewFilter()
	if !f.PutAll("a", "b", "c") {
		t.Fatal("failed to put fresh value list")
	}
	if f.PutAll("a", "b", "c") {
		t.Fatal("duplicate insertion of a,b,c")
	}
	if f.PutAll("a") {
		t.Fatal("duplicate partial insertion of a")
	}
	if !f.PutAll("d") {
		t.Fatal("failed to put fresh values after successful insertion of unrelated values")
	}
	f.DelAll("d")
	if !f.PutAll("d") {
		t.Fatal("failed to delete 'd'")
	}
	f.DelAll("a", "b", "c")
	if !f.PutAll("a", "b", "c") {
		t.Fatal("failed to delete a,b,c")
	}
}

func testPutDel(f *Filter, put, del []string, done chan bool) {
	for {
		select {
		case <-done:
		default:
		}
		f.PutAll(put...)
		f.DelAll(del...)
	}
}

func TestFilterConcurrent(t *testing.T) {
	const N = 1000000
	f := NewFilter()
	c := make(chan bool)
	defer close(c)
	go testPutDel(f, []string{"a", "b"}, []string{"c", "d", "e"}, c)
	go testPutDel(f, []string{"c", "d", "e"}, []string{"a", "b"}, c)
	for i := 0; i < N; i++ {
		if !f.PutAll("z") {
			t.Fatal("failed to put z")
		}
		f.DelAll("z")
	}
}

func BenchmarkFilter(b *testing.B) {
	f := NewFilter()
	for n := 0; n < b.N; n++ {
		f.PutAll("z")
		f.DelAll("z")
	}
}

func BenchmarkFilterConcurrentLoadx2(b *testing.B) {
	f := NewFilter()
	c := make(chan bool)
	defer close(c)
	go testPutDel(f, []string{"a", "b"}, []string{"c", "d", "e"}, c)
	go testPutDel(f, []string{"c", "d", "e"}, []string{"a", "b"}, c)
	for n := 0; n < b.N; n++ {
		f.PutAll("z")
		f.DelAll("z")
	}
}
