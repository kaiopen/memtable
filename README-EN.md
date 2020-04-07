# memtable
A memory table for Go: A thread-safe map with automatic deletion.

[中文版](https://github.com/kaiopen/memtable/blob/master/README.md)

### Example
```go
// Create a memory table.
// Actually, it can be treated as a thread-safe map.
memtable := MemTable{}

// Insert an item with key, value and duration time (in second).
// If duration is no more than 0, the item will not be automatically deleted.
memtable.Update("1", "123", 10)

// Get an item by key. If the item does not exist, the `ok` is `false`.
if data, ok := memtable.Get("1"); ok {
    fmt.Printf("data: %s\n", data)  // data: 123
}

// Update an item.
// The item will be deleted firstly if exists. Then insert a new one.
memtable.Update("1", "456", 0)

// Delete an item according to key. Do nothing if the item does not exist.
memtable.Delete("1")
```