package main

import (
  "flag"
  "fmt"
  "io"
  "time"
  "crypto/md5"
  "encoding/hex"
  "github.com/barneyb/kember/impl"
)

type Worker struct {
  Searcher kember.Searcher
  Ticks uint64
}

type StatusUpdate struct {
  Worker Worker
  Status kember.Status
  I uint64
  Curr string
}

func main() {
  start := flag.String("start", randHash(), "hash to start searching from")
  iterations := flag.Uint64("n", 0, "number of search iterations (0 means 'forever')")
  // threads := flag.Uint("threads", 1, "number of concurrent threads to run")
  flag.Parse()
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else {
    updates := make(chan StatusUpdate)
    total := float64(0)
    workers := 0



      workers++
      log := make(chan kember.StatusUpdate)
      s := kember.Searcher{ log, 1 * 1000 * 1000, *start, *iterations }
      w := Worker{ s, 0 }
      go kember.Search(&s)
      go func() {
        for {
          su := <- log
          updates <- StatusUpdate{w, su.Status, su.I, su.Curr}
        }
      }()



    var msg string
    for workers > 0 {
      su := <- updates
      switch su.Status {
        case kember.TICK:
          msg = su.Curr
        case kember.MATCH:
          msg = fmt.Sprintf("%v == %v <-- MATCH!!!", su.Curr, su.Curr)
        case kember.DONE:
          workers--
          msg = "finished"
      }
      su.Worker.Ticks = su.I
      i := float64(su.I) / 1000000.0
      total += i
      fmt.Printf("%.7s %7.1fM / %7.1fM %s %s\n", su.Worker.Searcher.Start, i, total, time.Now().Format("2006-01-02T15:04:05-0700"), msg)
    }
  }
}

func randHash() string {
  h := md5.New()
  io.WriteString(h, time.Now().Format(time.RFC3339Nano))
  sum := h.Sum(nil)
  return hex.EncodeToString(sum[0:16])
}
