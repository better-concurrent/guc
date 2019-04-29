package main

import (
	"fmt"
	"github.com/better-concurrent/guc/util"
	"testing"
	"time"
)

func TestRuntimeLink(t *testing.T) {
	fmt.Println("now test " + t.Name())
	if util.SyncRuntimeCanSpin(0) {
		t.Error("SyncRuntimeCanSpin failed!")
	}
	// nanoTime
	fmt.Println(util.SyncRuntimeNanoTime())
	// do spin
	util.SyncRuntimeDoSpin()

	// sema
	var sema uint32 = 0
	c := make(chan string, 1)
	go func() {
		util.SyncRuntimeSemacquire(&sema)
		c <- "ok"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		util.SyncRuntimeSemrelease(&sema, false)
	}()

	select {
	case ok := <-c:
		fmt.Println("sema " + ok)
	case <-time.After(2 * time.Second):
		t.Error("sema error!")
	}
}
