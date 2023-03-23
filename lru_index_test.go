package diskv

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLruIndexOrder(t *testing.T) {
	defer os.RemoveAll("index-test-db-1")
	d := New(Options{
		BasePath:     "index-test",
		CacheSizeMax: 1024,
		Index:        newLruIndex("index-test-db-1"),
	})
	defer d.EraseAll()

	v := []byte{'1', '2', '3'}
	d.Write("a", v)
	if !d.isIndexed("a") {
		t.Fatalf("'a' not indexed after write")
	}
	d.Write("b", v)
	d.Write("m", v)
	d.Write("d", v)
	d.Write("A", v)
	d.Write("d", v)

	expectedKeys := []string{"a", "b", "m", "A", "d"}
	keys := make([]string, 0)
	for _, key := range d.Index.Keys("", 100) {
		keys = append(keys, key)
	}

	if !cmpStrings(keys, expectedKeys) {
		t.Fatalf("got %s, expected %s", keys, expectedKeys)
	}
}

func TestLruIndexLoad(t *testing.T) {
	defer os.RemoveAll("index-test-db-2")
	d1 := New(Options{
		BasePath:     "index-test",
		CacheSizeMax: 1024,
		Index:        newLruIndex("index-test-db-2"),
	})

	val := []byte{'1', '2', '3'}
	keys := []string{"a", "b", "c", "d", "e", "f", "g"}
	for _, key := range keys {
		d1.Write(key, val)
	}

	d1.Index.(*lruIndex).db.Close()
	d1.Index = nil

	d2 := New(Options{
		BasePath:     "index-test",
		CacheSizeMax: 1024,
		Index:        newLruIndex("index-test-db-2"),
	})
	defer d2.EraseAll()

	// check d2 has properly loaded existing d1 data
	for _, key := range keys {
		if !d2.isIndexed(key) {
			t.Fatalf("key '%s' not indexed on secondary", key)
		}
	}

	// cache one
	if readValue, err := d2.Read(keys[0]); err != nil {
		t.Fatalf("%s", err)
	} else if bytes.Compare(val, readValue) != 0 {
		t.Fatalf("%s: got %s, expected %s", keys[0], readValue, val)
	}

	// make sure it got cached
	for i := 0; i < 10 && !d2.isCached(keys[0]); i++ {
		time.Sleep(10 * time.Millisecond)
	}
	if !d2.isCached(keys[0]) {
		t.Fatalf("key '%s' not cached", keys[0])
	}

	// kill the disk
	d1.EraseAll()

	// cached value should still be there in the second
	if readValue, err := d2.Read(keys[0]); err != nil {
		t.Fatalf("%s", err)
	} else if bytes.Compare(val, readValue) != 0 {
		t.Fatalf("%s: got %s, expected %s", keys[0], readValue, val)
	}

	// but not in the original
	if _, err := d1.Read(keys[0]); err == nil {
		t.Fatalf("expected error reading from flushed store")
	}
}

func TestLruIndexKeysEmptyFrom(t *testing.T) {
	defer os.RemoveAll("index-test-db-3")
	d := New(Options{
		BasePath:     "index-test",
		CacheSizeMax: 1024,
		Index:        newLruIndex("index-test-db-3"),
		IndexLess:    strLess,
	})
	defer d.EraseAll()

	keys := []string{"a", "c", "z", "b", "x", "b", "y"}
	for _, k := range keys {
		d.Write(k, []byte("1"))
	}

	want := []string{"a", "c", "z", "x", "b", "y"}
	have := d.Index.Keys("", 99)
	if !reflect.DeepEqual(want, have) {
		t.Errorf("want %v, have %v", keys, have)
	}
}
