package diskv

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestLruAddRemove(t *testing.T) {
	c := newLru(1, func(key string) error { return nil })
	k := "a"
	if err := c.Add(k); err != nil {
		t.Fatalf("write: %s", err)
	}
	if _, ok := c.items[k]; !ok {
		t.Fatalf("add key failed: %s", k)
	}
	c.Remove(k)
	if _, ok := c.items[k]; ok {
		t.Fatalf("remove key failed: %s", k)
	}
}

func TestLruEvict(t *testing.T) {
	var (
		d = New(Options{
			BasePath: "test-data",
			LruSize:  2,
		})

		keys     = []string{"a", "b", "c", "d"}
		v        = []byte{'1'}
		expected = map[string]bool{"a": true, "b": true, "c": false, "d": false}
	)
	defer d.EraseAll()
	for _, k := range keys {
		if err := d.Write(k, v); err != nil {
			t.Fatalf("write: %s: %s", k, err)
		}
	}
	for key, evict := range expected {
		b, err := d.Read(key)
		if evict {
			if errors.Is(err, os.ErrNotExist) {
				t.Logf("%s evict got: %v", key, evict)
			} else {
				t.Fatalf("%s expected error: %s got: %v", key, os.ErrNotExist, err)
			}
		} else {
			if err != nil {
				t.Fatalf("%s expected error: nil got: %s", key, err)
			}

			if bytes.Equal(v, b) {
				t.Logf("%s value got: %s", key, b)
			} else {
				t.Logf("%s expected value: %s got: %v", key, v, b)
			}
		}
	}
}
