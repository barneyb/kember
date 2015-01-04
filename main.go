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

func main() {
  start := flag.String("start", randHash(), "hash to start searching from")
  iterations := flag.Uint64("n", 0, "number of search iterations (0 means 'forever')")
  // threads := flag.Uint("threads", 1, "number of concurrent threads to run")
  flag.Parse()
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else {
    log := make(chan kember.StatusUpdate)
    s := kember.Searcher{ log, *start, *iterations }
    go kember.Search(&s)
    var msg string
    var su kember.StatusUpdate
    for ;true; {
      su = <- log
      switch su.Status {
        case kember.TICK:
          msg = su.Curr
        case kember.MATCH:
          msg = fmt.Sprintf("%v == %v <-- MATCH!!!", su.Curr, su.Curr)
        case kember.DONE:
          msg = "finished"
      }
      fmt.Printf("%.7s %8.1fM) [%s] %s\n", s.Start, float64(su.I) / 1000000.0, time.Now().Format("2006-01-02 15:04:05 -0700 MST"), msg)
      if su.Status == kember.DONE {
        break
      }
    }
  }
}

func randHash() string {
  h := md5.New()
  io.WriteString(h, time.Now().Format(time.RFC3339Nano))
  sum := h.Sum(nil)
  return hex.EncodeToString(sum[0:16])
}
