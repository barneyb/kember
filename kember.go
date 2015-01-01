package main

import (
	"time"
	"bytes"
	"fmt"
	"flag"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func main() {
	start := flag.String("start", "00000000000000000000000000000000", "hash to start searching from")
	iterations := flag.Int("n", -1, "number of search iterations (-1 means 'forever')")
	flag.Parse()
	if ! valid(*start) {
		fmt.Println("The starting hash is invalid.")
	} else {
		Search(*start, *iterations)
	}
}

func valid(hash string) bool {
	if len(hash) != 32 {
		return false
	}
	for _, runeValue := range strings.ToLower(hash) {
		if ((runeValue < 'a' || runeValue > 'f') && (runeValue < '0' || runeValue > '9')) {
	        return false
	    }
	}
	return true
}

func Search(start string, iterations int) {
	fmt.Printf("search(%v, %v)!\n", start, iterations)
	curr, _ := hex.DecodeString(start)
	fmt.Println(curr)
	for i := 0; i < 70; i++ {
		increment(&curr)
		fmt.Println(curr)
	}
	for i := 0; iterations < 0 || i < iterations; i++ {
		if i % 1000000 == 0 {
			fmt.Printf("%d) [%v]\n", i, time.Now())
		}
		sum := md5.Sum(curr)
		if bytes.Equal(curr, sum[0:16]) {
			fmt.Printf("%d) [%v] %v\n", i, time.Now(), hex.EncodeToString(sum[0:16]))
		}
		increment(&curr)
	}
}

func increment(curr *[]byte) {
}
