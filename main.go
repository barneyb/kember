package main

import (
  "flag"
  "fmt"
  "github.com/barneyb/kember/impl"
)

func main() {
  start := flag.String("start", "00000000000000000000000000000000", "hash to start searching from")
  iterations := flag.Int("n", -1, "number of search iterations (-1 means 'forever')")
  flag.Parse()
  if ! kember.Valid(*start) {
    fmt.Println("The starting hash is invalid.")
  } else {
    kember.Search(*start, *iterations)
  }
}
