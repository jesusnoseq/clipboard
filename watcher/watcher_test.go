package watcher

import (
	"testing"
	"time"

	"github.com/atotto/clipboard"
)

type testListener struct {
	changes chan string
}

func TestWatcher(t *testing.T) {
	expectedOne := "expected first"
	expectedTwo := "expected second"

	//  clean clipboard
	clipboard.WriteAll("")

	w := NewClipboardObserver(50 * time.Millisecond)
	m := testListener{
		changes: make(chan string, 2),
	}
	w.AddListener(&m)

	go w.Start()

	clipboard.WriteAll(expectedOne)
	wait()
	clipboard.WriteAll(expectedTwo)
	wait()
	clip := <-m.changes
	if clip != expectedOne {
		t.Errorf("Error - got %s expected %s", clip, expectedOne)
	}

	clip = <-m.changes
	if clip != expectedTwo {
		t.Errorf("Error - got %s expected %s", clip, expectedTwo)
	}
}

func wait() {
	time.Sleep(100 * time.Millisecond)
}

func (m *testListener) OnChange(text string) {
	m.changes <- text
}
