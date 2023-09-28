package statusbar

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rivo/tview"
)

type StatusBar struct {
	app   *tview.Application
	field *tview.InputField
	done  chan bool
	mutex sync.Mutex
}

func NewStatusBar(application *tview.Application) *StatusBar {
	i := tview.NewInputField()
	i.SetDisabled(true)
	return &StatusBar{
		field: i,
		app:   application,
	}
}

func (f *StatusBar) setLabel(s string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.field.SetLabelWidth(len(s) + 1)
	f.field.SetLabel(s)
}

func (f *StatusBar) setMessage(items ...string) {
	message := strings.Join(items, " ")

	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.field.SetText(message)
}

// StartLoading displays a loading animation with a specified message.
//
// The animation stops as soon as DoneLoading is called.
func (f *StatusBar) StartLoading(message string) {
	f.setLabel("Loading")
	f.done = make(chan bool)
	go func() {
		count := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-f.done:
				ticker.Stop()
				f.setLabel("Finished Loading")
				f.setMessage(message)
				f.app.Draw()
				return
			case <-ticker.C:
				lbl := fmt.Sprintf("%s%s", message, strings.Repeat(".", count%4))
				f.setMessage(lbl)
				f.app.Draw()
				count += 1
			}
		}
	}()
}

// DoneLoading is called after calling StartLoading to indicate that the loading process has finished.
func (f *StatusBar) DoneLoading() {
	select {
	case f.done <- true:
		return
	default:
		panic("Did not call StartLoading()!")
	}
}
