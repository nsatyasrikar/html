package html

import (
	"log"
	"sync"
)

type Warning struct{ Code, Message, Tag, Path string }
type WarningHandler func(Warning)

var warningMu sync.RWMutex
var warningHandler WarningHandler = func(w Warning) { log.Printf("html warning [%s] %s (%s)", w.Code, w.Message, w.Path) }

func SetWarningHandler(h WarningHandler) { warningMu.Lock(); warningHandler = h; warningMu.Unlock() }
func emitWarning(w Warning) {
	warningMu.RLock()
	h := warningHandler
	warningMu.RUnlock()
	if h != nil {
		h(w)
	}
}
