package main

import (
  "runtime"
  "time"
  "io"
  "strings"
  "fmt"
  "flag"
  "crypto/md5"
  "encoding/hex"
)

const tickFrequency = 1 * 1000 * 1000

type UpdateType uint
const (
  TICK UpdateType = iota
  MATCH
  DONE
)

type WorkerState struct {
  Log chan StatusUpdate
  Start *string
  TicksRequested uint64
  TicksCompleted uint64
}
func (state *WorkerState) Done() bool {
  return state.TicksRequested != 0 && state.TicksCompleted >= state.TicksRequested
}

type StatusUpdate struct {
  Worker *WorkerState
  Status UpdateType
  Ticks uint64
  Curr *string
}

type Config struct {
  WorkerCount int
  TicksRequested uint64
  Start *string
}

func getConfig() *Config {
  rh := "<random>"
  start := flag.String("start", rh, "hash to start searching from")
  requested := flag.Uint64("n", 0, "number of search blocks (0 means 'forever')")
  workers := flag.Int("w", 1, "number of concurrent workers to run")
  flag.Parse()

  if *start == rh {
    start = nil
  }

  if start != nil && ! validHash(start) {
    panic("The starting hash is invalid.")
  }
  if *workers < 1 {
    panic("At least one thread must be used.")
  }
  if *workers > runtime.NumCPU() {
    panic(fmt.Sprintf("This machine only has %d CPUs available.", runtime.NumCPU()))
  }

  return &Config{ *workers, *requested, start }
}

func randHash() *string {
  h := md5.New()
  io.WriteString(h, time.Now().Format(time.RFC3339Nano))
  sum := h.Sum(nil)
  hash := hex.EncodeToString(sum[0:16])
  return &hash
}

func startWorkers(config *Config, updates chan StatusUpdate) *[]*WorkerState{
  workers := make([]*WorkerState, 0, config.WorkerCount)
  for i := 0; i < config.WorkerCount; i++ {
    var st *string
    if i > 0 || config.Start == nil {
      st = randHash() // subsequent workers always randomize
    } else {
      st = config.Start
    }
    w := WorkerState{ updates, st, config.TicksRequested, 0 }
    workers = append(workers, &w)
    go search(&w)
  }
  return &workers
}

func main() {
  config := getConfig()
  updates := make(chan StatusUpdate)
  workers := *startWorkers(config, updates)

  keepGoing := func() bool {
    for _, w := range workers {
      if ! w.Done() {
        return true
      }
    }
    return false
  }
  totalTicks := func() uint64 {
    total := uint64(0)
    for _, w := range workers {
      total += w.TicksCompleted
    }
    return total
  }
  lastTotal := uint64(0)
  startTime := time.Now().Unix()
  lastTick := startTime
  var lastHash string
  log := func(msg string) {
    total := totalTicks()
    tick := time.Now().Unix()
    perSecTick := float64(0)
    dt := tick - lastTick
    if dt > 0 {
      perSecTick = float64(total - lastTotal) / float64(dt) * tickFrequency
    }
    perSecAll := float64(0)
    dt = tick - startTime
    if dt > 0 {
      perSecAll = float64(total) / float64(dt) * tickFrequency
    }
    lastTotal = total
    lastTick = tick
    fmt.Printf("%s %7.1e %7.1e/s %7.1e/s %s\n", time.Now().Format("2006-01-02T15:04:05-0700"), float64(total) * tickFrequency, perSecTick, perSecAll, msg)
  }
  go func() {
    c := time.Tick(5 * time.Second)
    for range c {
      log(lastHash)
    }
  }()
  fmt.Printf("%-24s %-7s %-9s %-9s %s\n", "timestamp", "tests", "tick", "overall", "message")
  log("starting")
  for keepGoing() {
    su := <- updates
    lastHash = *su.Curr
    switch su.Status {
      case MATCH:
        log(fmt.Sprintf("md5(%v) == %v <-- MATCH!!!", su.Curr, su.Curr))
      case DONE:
        log(fmt.Sprintf("%.7s done", *su.Worker.Start))
    }
  }
  log("exiting")
}

func search(state *WorkerState) {
  curr := *state.Start
  runes := []rune(curr)
  h := md5.New()
  for ; state.TicksRequested == 0 || state.TicksCompleted < state.TicksRequested; state.TicksCompleted++ {
    state.Log <- StatusUpdate{ state, TICK, state.TicksCompleted, &curr }
    for j := 0; j < tickFrequency; j++ {
      h.Reset()
      io.WriteString(h, curr)
      sum := h.Sum(nil)
      hash := hex.EncodeToString(sum[0:16])
      if curr == hash {
        state.Log <- StatusUpdate{ state, MATCH, state.TicksCompleted, &curr }
      }
      increment(runes)
      curr = string(runes)
    }
  }
  state.Log <- StatusUpdate{ state, DONE, state.TicksCompleted, &curr }
}

func increment(runes []rune) {
  runeCount := len(runes)
  pos := runeCount - 1
  for ; pos > 0 && runes[pos] == 'f'; pos-- {}
  for i := pos; i < runeCount; i++ {
    runes[i] = next(runes[i])
  }
}

func next(curr rune) rune {
  if curr == '9' {
    return 'a'
  } else if curr == 'f' {
    return '0'
  } else {
    return curr + 1
  }
}

func validHash(hash *string) bool {
  if len(*hash) != 32 {
    return false
  }
  for _, runeValue := range strings.ToLower(*hash) {
    if ((runeValue < 'a' || runeValue > 'f') && (runeValue < '0' || runeValue > '9')) {
          return false
      }
  }
  return true
}
