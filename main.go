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
  iterations := flag.Int64("n", -1, "number of search iterations (-1 means 'forever')")
  // threads := flag.Int("threads", 1, "number of concurrent threads to run")
  flag.Parse()
  log := make(chan string)
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else {
    s := kember.Searcher{ log, *start, *iterations, 0, *start }
    go kember.Search(&s)
    var msg string
    for ;true; {
      msg = <- log
      fmt.Printf("%8.1fM) [%v] %v\n", float64(s.I) / 1000000.0, time.Now().Format("2006-01-02 15:04:05 -0700 MST"), msg)
      if msg == "finished" {
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
