package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *Once) recover() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}

var once Once
func test(){
	once.Do(doSomething)
}
func main() {
	var wg sync.WaitGroup
	wg.Add(100)

	for i:=0;i<100;i++ {
		go test()
	}

	wg.Wait()
	fmt.Println("finish")
}

func doSomething()  {
	fmt.Println("i'm doing sth")
}




