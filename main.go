package memtable

import (
	"sync"
	"time"
)

// Item is table item.
type Item struct {
	data      interface{}
	hasExpiry bool
	delChan   chan bool
}

// MemTable is a table in member.
type MemTable struct {
	items sync.Map
}

// Update updates or creates one item in table. The item will be deleted
// firstly it exists. If duration is larger than 0, the item will be deleted automatically after `duration` seconds.
func (m *MemTable) Update(key interface{}, data interface{}, duration int64) {
	m.Delete(key)

	i := &Item{
		data:      data,
		hasExpiry: duration > 0,
		delChan:   make(chan bool),
	}
	m.items.Store(key, i)

	if i.hasExpiry {
		go func() {
			tricker := time.NewTicker(time.Second * time.Duration(duration))
			select {
			case <-tricker.C:
				close(i.delChan)
				m.items.Delete(key)
			case <-i.delChan:
				close(i.delChan)
				m.items.Delete(key)
			}
		}()
	}
}

// Delete deletes a item in table if exists or do nothing.
func (m *MemTable) Delete(key interface{}) {
	if i, ok := m.items.Load(key); ok {
		item := i.(*Item)
		if item.hasExpiry {
			item.delChan <- true
		} else {
			close(item.delChan)
			m.items.Delete(key)
		}
	}
}

// Get gets an item.
func (m *MemTable) Get(key interface{}) (data interface{}, ok bool) {
	i, ok := m.items.Load(key)
	if ok {
		data = i.(*Item).data
	}
	return
}
