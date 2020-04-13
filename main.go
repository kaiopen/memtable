package memtable

import (
	"log"
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

	var i *Item
	hasExpiry := duration > 0
	if hasExpiry {
		i = &Item{
			data:      data,
			hasExpiry: true,
			delChan:   make(chan bool, 1),
		}
	} else {
		i = &Item{
			data:      data,
			hasExpiry: false,
			delChan:   nil,
		}
	}

	m.items.Store(key, i)
	log.Printf("MemTable: %s update.\n", key)

	if hasExpiry {
		go func() {
			defer log.Println("time func end.")
			tricker := time.NewTicker(time.Second * time.Duration(duration))
			select {
			case <-tricker.C:
				close(i.delChan)
				m.items.Delete(key)
				log.Println("MemTable: time out and delete.")
			case <-i.delChan:
				close(i.delChan)
				log.Println("MemTable: del chan got.")
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
		}
		m.items.Delete(key)
		log.Println("MemTable: delete manually.")
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
