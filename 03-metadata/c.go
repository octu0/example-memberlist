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
  meta  MyMetaData
}
func (d *MyDelegate) NodeMeta(limit int) []byte {
  return d.meta.Bytes()
}
func (d *MyDelegate) LocalState(join bool) []byte {
  // not use, noop
  return []byte("")
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
  // not use
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
  // not use, noop
  return nil
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
  // not use
}

type MyMetaData struct {
  Region   string   `json:"region"`
  Zone     string   `json:"zone"`
  ShardId  uint16   `json:"shard-id"`
  Weight   uint64   `json:"weight"`
}
func (m MyMetaData) Bytes() []byte {
  data, err := json.Marshal(m)
  if err != nil {
    return []byte("")
  }
  return data
}
func ParseMyMetaData(data []byte) (MyMetaData, bool) {
  meta := MyMetaData{}
  if err := json.Unmarshal(data, &meta); err != nil {
    return meta, false
  }
  return meta, true
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
