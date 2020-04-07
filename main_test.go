package memtable

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type Test struct {
	key      string
	data     string
	duration int64
}

func TestMemtable(t *testing.T) {
	tests := [...]Test{
		{key: "1", data: "123", duration: 10},
		{key: "2", data: "234", duration: 0},
		{key: "3", data: "345", duration: -1},
		{key: "4", data: "456", duration: 5},
	}
	memtable := MemTable{}
	for _, test := range tests {
		memtable.Update(test.key, test.data, test.duration)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		// Print table per second.
		ticker := time.NewTicker(time.Second * 1)
		d := 0
		for {
			select {
			case <-ticker.C:
				d++
				fmt.Printf("%d seconds\n", d)
				for _, test := range tests {
					if data, ok := memtable.Get(test.key); ok {
						fmt.Printf(
							"key: %s, data: %s, duration: %d\n",
							test.key, data, test.duration,
						)
					} else {
						fmt.Printf("%s does not exist.\n", test.key)
					}
				}
				if d == 12 {
					wg.Done()
					return
				}
			}
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(time.Second * 2)
		memtable.Delete("3")
		fmt.Println("The item with key 3 is deleted.")
		wg.Done()
	}(wg)

	wg.Wait()
	fmt.Print("End")
}
