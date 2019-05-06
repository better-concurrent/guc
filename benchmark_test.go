package guc

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

type taskConfig struct {
	readGCount  int
	writeGCount int
	countPerG   int
}

func runGuc(config taskConfig) {
	cmap := NewConcurrentHashMap(64, 16)
	endCh := make(chan int)
	value0 := valueObject{v: "v"}
	cmap.Store(keyObject2{i: -1}, value0)
	for i := 0; i < config.writeGCount; i++ {
		ii := i
		go func() {
			begin := config.countPerG
			beginTime := time.Now()
			for n := 0; n < begin+config.countPerG; n++ {
				key := keyObject2{i: n}
				cmap.Store(key, value0)
			}
			endTime := time.Now()
			fmt.Printf("guc.map %d write cost is %d\n", ii,
				endTime.Sub(beginTime).Nanoseconds()/1e6)
			endCh <- 1
		}()
	}
	for i := 0; i < config.readGCount; i++ {
		ii := i
		go func() {
			begin := config.countPerG
			beginTime := time.Now()
			for n := 0; n < begin+config.countPerG; n++ {
				key := keyObject2{i: n}
				cmap.Load(key)
			}
			endTime := time.Now()
			fmt.Printf("guc.map %d read cost is %d\n", ii,
				endTime.Sub(beginTime).Nanoseconds()/1e6)
			endCh <- 1
		}()
	}
	total := config.writeGCount + config.readGCount
	for i := 0; i < total; i++ {
		<-endCh
	}
}

func runSyncMap(config taskConfig) {
	smap := new(sync.Map)
	endCh := make(chan int)
	value0 := valueObject{v: "v"}
	for i := 0; i < config.writeGCount; i++ {
		ii := i
		go func() {
			begin := config.countPerG
			beginTime := time.Now()
			for n := 0; n < begin+config.countPerG; n++ {
				key := keyObject2{i: n}
				smap.Store(key, value0)
			}
			endTime := time.Now()
			fmt.Printf("sync.map %d write cost is %d\n", ii,
				endTime.Sub(beginTime).Nanoseconds()/1e6)
			endCh <- 1
		}()
	}
	for i := 0; i < config.readGCount; i++ {
		ii := i
		go func() {
			begin := config.countPerG
			beginTime := time.Now()
			for n := 0; n < begin+config.countPerG; n++ {
				key := keyObject2{i: n}
				smap.Load(key)
			}
			endTime := time.Now()
			fmt.Printf("sync.map %d read cost is %d\n", ii,
				endTime.Sub(beginTime).Nanoseconds()/1e6)
			endCh <- 1
		}()
	}
	total := config.writeGCount + config.readGCount
	for i := 0; i < total; i++ {
		<-endCh
	}
}

func TestConcurrentHashMap(t *testing.T) {
	runtime.GOMAXPROCS(4)
	readGCount := 4
	writeGCount := 4
	countPerG := 1024 * 32
	config := taskConfig{readGCount: readGCount, writeGCount: writeGCount, countPerG: countPerG}
	runGuc(config)
	runSyncMap(config)
}
