package kember

import (
	"io"
	"time"
	"fmt"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

func Valid(hash string) bool {
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

func Search(curr string, iterations int) {
	log(0, "search(%v, %v)!", curr, iterations)
	h := md5.New()
	io.WriteString(h, curr)
	for i := 0; i < 20; i++ {
		curr = string(increment([]rune(curr)))
		fmt.Println(curr)
	}
	i := 0
	for ; iterations < 0 || i < iterations; i++ {
		if i % 1000000 == 0 {
			log(i, "...")
		}
		h.Reset()
		io.WriteString(h, curr)
		sum := h.Sum(nil)
		hash := hex.EncodeToString(sum[0:16])
		if curr == hash {
			log(i, "MATCH! " + hash)
		}
		curr = string(increment([]rune(curr)))
	}
	log(i, "finished")
}

func log(i int, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Printf("%d) [%v] %v\n", i, time.Now().Format("2006-01-02 15:04:05 -0700 MST"), msg)
}

func increment(runes []rune) []rune {
	runeCount := len(runes)
	pos := runeCount - 1
	for ; pos >= 0 && runes[pos] == 'f'; pos-- {}
	if pos < 0 {
		log(-1, "OVERFLOW!")
		return []rune("00000000000000000000000000000000")
	}
	for i := pos; i < runeCount; i++ {
		runes[i] = next(runes[i])
	}
	return runes
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
