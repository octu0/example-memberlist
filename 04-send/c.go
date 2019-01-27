package main

import(
  "log"
  "os"
  "os/signal"
  "syscall"
  "context"
  "encoding/json"
)

type MyDelegate struct {
  msgCh  chan []byte
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
  d.msgCh <- msg
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
  // not use, noop
  return []byte("")
}
func (d *MyDelegate) LocalState(join bool) []byte {
  // not use, noop
  return []byte("")
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  // not use, noop
  return nil
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
  // not use
}

type MyMessage struct {
  Key    string  `json:"key"`
  Value  uint64  `json:"value"`
}
func (m *MyMessage) Bytes() []byte {
  data, err := json.Marshal(m)
  if err != nil {
    return []byte("")
  }
  return data
}
func ParseMyMessage(data []byte) (*MyMessage, bool) {
  msg := new(MyMessage)
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
