package main

import(
  "log"
  "os"
  "os/signal"
  "syscall"
  "context"
  "encoding/json"
  "github.com/hashicorp/memberlist"
)

type MyDelegate struct {
  msgCh      chan []byte
  broadcasts *memberlist.TransmitLimitedQueue
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
  d.msgCh <- msg
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  return d.broadcasts.GetBroadcasts(overhead, limit)
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
  // not use, noop
  return []byte("")
}
func (d *MyDelegate) LocalState(join bool) []byte {
  // not use, noop
  return []byte("")
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
  // not use
}

type MyEventDelegate struct {
  Num int
}
func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
  d.Num += 1
}
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
  d.Num -= 1
}
func (d *MyEventDelegate) NotifyUpdate(node *memberlist.Node) {
  // skip
}

type MyBroadcastMessage struct {
  Key    string  `json:"key"`
  Value  uint64  `json:"value"`
}
func (m MyBroadcastMessage) Invalidates(other memberlist.Broadcast) bool {
  return false
}
func (m MyBroadcastMessage) Finished() {
  // nop
}
func (m MyBroadcastMessage) Message() []byte {
  data, err := json.Marshal(m)
  if err != nil {
    return []byte("")
  }
  return data
}

func ParseMyBroadcastMessage(data []byte) (*MyBroadcastMessage, bool) {
  msg := new(MyBroadcastMessage)
  if err := json.Unmarshal(data, &msg); err != nil {
    return nil, false
  }
  return msg, true
}

func wait_signal(cancel context.CancelFunc){
  signal_chan := make(chan os.Signal, 1)
  signal.Notify(signal_chan, syscall.SIGINT)
  for {
    select {
    case s := <-signal_chan:
      log.Printf("signal %s happen", s.String())
      cancel()
    }
  }
}
