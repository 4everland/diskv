package diskv

import (
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"os"
	"sort"
	"time"
)

type lruIndex struct {
	path string
	db   *leveldb.DB
}

func newLruIndex(path string) Index {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &lruIndex{db: db, path: path}
}

func (s *lruIndex) Initialize(less LessFunction, keys <-chan string) {
	var err error
	if s.db != nil {
		if err = s.db.Close(); err != nil {
			panic(err)
		}
	}
	_ = os.RemoveAll(s.path)
	if s.db, err = leveldb.OpenFile(s.path, nil); err != nil {
		panic(err)
	}
}

func (s *lruIndex) Insert(key string) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(time.Now().UnixMicro()))
	s.db.Put([]byte(key), buf, nil)
}

func (s *lruIndex) Delete(key string) {
	s.db.Delete([]byte(key), nil)
}

func (s *lruIndex) Keys(from string, n int) []string {
	var (
		iter = s.db.NewIterator(&util.Range{
			Start: []byte(from),
		}, nil)

		items = make(lruIndexItems, 0)
		i     = 0
	)

	for iter.Next() {
		items = append(items, lruIndexItem{
			Key:       string(iter.Key()),
			UpdatedAt: binary.BigEndian.Uint64(iter.Value()),
		})

		if i++; n > 0 && i >= n {
			break
		}
	}

	sort.Sort(items)
	iter.Release()

	return items.GetKeys()
}

type lruIndexItem struct {
	Key       string
	UpdatedAt uint64
}

type lruIndexItems []lruIndexItem

func (s lruIndexItems) Len() int           { return len(s) }
func (s lruIndexItems) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s lruIndexItems) Less(i, j int) bool { return s[i].UpdatedAt < s[j].UpdatedAt }

func (s lruIndexItems) GetKeys() []string {
	keys := make([]string, len(s))
	for i, item := range s {
		keys[i] = item.Key
	}

	return keys
}
