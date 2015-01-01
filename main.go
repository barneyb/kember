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
  flag.Parse()
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else {
    kember.Search(*start, *iterations)
  }
}

func randHash() string {
  h := md5.New()
  io.WriteString(h, time.Now().Format(time.RFC3339Nano))
  sum := h.Sum(nil)
  return hex.EncodeToString(sum[0:16])
}
