package watcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

// Listener are notified
type Listener interface {
	OnChange(text string)
}

type clipboardObserver struct {
	listeners     sync.Map
	watchInterval time.Duration
	current       string
	stopedChan    chan struct{}
	stopChan      chan struct{}
}

// NewClipboardObserver is the constructor to create a clipboardObserver
func NewClipboardObserver(watchInterval time.Duration) *clipboardObserver {
	return &clipboardObserver{
		listeners:     sync.Map{},
		stopChan:      make(chan struct{}),
		stopedChan:    make(chan struct{}),
		watchInterval: watchInterval,
	}
}

// AddListener adds a listener to notify changes
func (c *clipboardObserver) AddListener(l Listener) {
	c.listeners.Store(l, struct{}{})
}

// RemoveListener removes a listener
func (c *clipboardObserver) RemoveListener(l Listener) {
	c.listeners.Delete(l)
}

// notify changes to listeners
func (c *clipboardObserver) notify(text string) {
	c.listeners.Range(func(k interface{}, v interface{}) bool {
		k.(Listener).OnChange(text)
		return true
	})
}

// check if there is something to notify
func (c *clipboardObserver) check() error {
	newText, err := clipboard.ReadAll()
	if err != nil {
		return err
	}
	if c.current != newText {
		c.current = newText
		c.notify(newText)
	}
	return nil
}

// Start watch for changes in clipboard
func (c *clipboardObserver) Start() error {
	for {
		select {
		case <-c.stopChan:
			c.stopedChan <- struct{}{}
			return nil
		case <-time.After(c.watchInterval):
			fmt.Println("Executing at ", time.Now())
			err := c.check()
			if err != nil {
				return err
			}
		}
	}
}

// Stop the observer gorutine
func (c *clipboardObserver) Stop() {
	c.stopChan <- struct{}{}
	<-c.stopedChan
}
