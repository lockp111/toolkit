package jsonrpc2

import (
	"sync"

	"github.com/lockp111/toolkit/log"
)

type notify struct {
	method string
	data   interface{}
}

// Notifier ...
type Notifier struct {
	sync.Mutex
	conn    *Conn
	sending chan *notify
	isopen  bool
}

// NewNotifier ...
func NewNotifier(conn *Conn) *Notifier {
	sending := make(chan *notify, 1024)

	notifier := &Notifier{
		conn:    conn,
		sending: sending,
		isopen:  true,
	}
	conn.OnClose(func() {
		notifier.Close()
	})

	go notifier.loop()

	return notifier
}

func (n *Notifier) loop() {
	for {
		notify, ok := <-n.sending
		if !ok {
			break
		}

		if err := n.conn.Notify(notify.method, notify.data); err != nil {
			log.WithError(err).Error("WS send error")
			break
		}
	}
}

// Close ...
func (n *Notifier) Close() {
	n.Lock()
	defer n.Unlock()

	n.closeLocked()

}

func (n *Notifier) closeLocked() {
	if n.isopen {
		close(n.sending)
		n.isopen = false
	}
}

// Notify ...
func (n *Notifier) Notify(method string, data interface{}) {
	n.Lock()
	defer n.Unlock()

	if n.isopen {
		select {
		case n.sending <- &notify{
			method: method,
			data:   data,
		}:
		default:
			log.WithFields(log.Fields{
				"conn": n.conn,
			}).Error("Sending channel is full, close conn.")
			n.closeLocked()
			go n.conn.Close()
		}
	}
}
