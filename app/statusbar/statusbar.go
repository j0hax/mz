package statusbar

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rivo/tview"
)

/*
StatusBar is a special kind of tview.InputField.

It is used to display application status, such as loading.
*/
type StatusBar struct {
	*tview.InputField
	app         *tview.Application
	loadingDone chan bool
	mutex       sync.Mutex
}

func NewStatusBar(application *tview.Application) *StatusBar {
	s := StatusBar{
		InputField:  tview.NewInputField(),
		app:         application,
		loadingDone: make(chan bool),
	}

	s.SetDisabled(true)

	return &s
}

func (f *StatusBar) setLabel(s string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.SetLabelWidth(len(s) + 1)
	f.SetLabel(s)
}

func (f *StatusBar) setMessage(items ...string) {
	message := strings.Join(items, " ")

	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.SetText(message)
}

// StartLoading displays a loading animation with a specified message.
//
// The animation stops as soon as DoneLoading is called.
func (f *StatusBar) StartLoading(message string) {
	f.setLabel("Loading")
	go func() {
		count := 0
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-f.loadingDone:
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
	case f.loadingDone <- true:
		return
	default:
		panic("Did not call StartLoading()!")
	}
}
