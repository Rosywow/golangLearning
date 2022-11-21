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

func (o *Once) Do(from string,f func(name string)) {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(from,f)
	}
}

func (o *Once) doSlow(from string,f func(name string)) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f(from)
	}
}

func ChangeStage(name string) {
	fmt.Println("Hello",name)
}

func main(){
	var executeOnce Once
	executeOnce.Do("HaHaHaHaHaHa",ChangeStage)
}