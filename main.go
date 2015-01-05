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

const tickFrequency = 10 * 1000 * 1000

type Worker struct {
  Searcher *kember.Searcher
  Ticks uint64
  Done bool
}

type StatusUpdate struct {
  Worker *Worker
  Status kember.Status
  Curr string
}

func main() {
  rh := randHash()
  start := flag.String("start", rh, "hash to start searching from")
  iterations := flag.Uint64("n", 0, "number of search iterations (0 means 'forever')")
  threads := flag.Int("w", 1, "number of concurrent workers to run")
  flag.Parse()
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else if *threads < 1 {
    fmt.Println("At least one thread must be used.")
  } else {
    updates := make(chan StatusUpdate)
    workers := make([]*Worker, 0, *threads)

    ts := *threads
    for i := 0; i < ts; i++ {
      log := make(chan kember.StatusUpdate)
      var st string
      if i == 0 || *start != rh {
        st = *start
      } else {
        st = randHash()
      }
      s := kember.Searcher{ log, tickFrequency, st, *iterations / uint64(ts) }
      w := Worker{ &s, 0, false }
      workers = append(workers, &w)
      go kember.Search(&s)
      go func() {
        for {
          su := <- log
          w.Ticks = su.I
          updates <- StatusUpdate{ &w, su.Status, su.Curr }
          if su.Status == kember.DONE {
            w.Done = true
            break
          }
        }
      }()
    }



    var msg string
    keepGoing := func() bool {
      for _, w := range workers {
        if ! w.Done {
          return true
        }
      }
      return false
    }
    totalTicks := func() uint64 {
      total := uint64(0)
      for _, w := range workers {
        total += w.Ticks
      }
      return total
    }
    lastTotal := uint64(0)
    for keepGoing() {
      su := <- updates
      total := totalTicks()
      // only tick on the aggregate freq
      if lastTotal == 0 || (total - lastTotal) >= tickFrequency || total < lastTotal {
        lastTotal = total
        switch su.Status {
          case kember.TICK:
            msg = su.Curr
          case kember.MATCH:
            msg = fmt.Sprintf("%v == %v <-- MATCH!!!", su.Curr, su.Curr)
          case kember.DONE:
            msg = "finished"
        }
        fmt.Printf("%.7s %7.1e / %7.1e %s %s\n", su.Worker.Searcher.Start, float64(su.Worker.Ticks), float64(total), time.Now().Format("2006-01-02T15:04:05-0700"), msg)
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
