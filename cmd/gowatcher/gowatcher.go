package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atotto/clipboard/watcher"
)

type myListener struct{}

func main() {
	w := watcher.NewClipboardObserver(500 * time.Millisecond)
	l := myListener{}
	w.AddListener(&l)

	go w.Start()
	defer w.Stop()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func (m *myListener) OnChange(text string) {
	fmt.Println("New clip ", text)
}
