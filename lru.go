package diskv

import "container/list"

type evictFunction func(key string) error

type lru struct {
	size      int
	evictList *list.List
	items     map[string]*list.Element
	fn        evictFunction
}

func newLru(size int, fn evictFunction) *lru {
	return &lru{
		size:      size,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
		fn:        fn,
	}
}

func (c *lru) Add(key string) error {
	if e, ok := c.items[key]; ok {
		c.evictList.MoveToFront(e)
	} else {
		c.items[key] = c.evictList.PushFront(key)
		if c.evictList.Len() > c.size {
			if e = c.evictList.Back(); e != nil {
				if v, ok := e.Value.(string); ok {
					return c.fn(v)
				}
				c.evictList.Remove(e)
			}
		}
	}

	return nil
}

func (c *lru) Read(key string) {
	if e, ok := c.items[key]; ok {
		c.evictList.MoveToFront(e)
	}
}

func (c *lru) Remove(key string) {
	if e, ok := c.items[key]; ok {
		c.evictList.Remove(e)
		delete(c.items, key)
	}
}

func (c *lru) Purge() {
	c.evictList.Init()
	c.items = make(map[string]*list.Element)
}
