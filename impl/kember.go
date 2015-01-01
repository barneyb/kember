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

func Search(curr string, iterations int64) {
	log(0, "search(%v, %v)!", curr, iterations)
	runes := []rune(curr)
	h := md5.New()
	i := int64(0)
	for ; iterations < 0 || i < iterations; i++ {
		if i % 1000000 == 0 {
			log(i, curr)
		}
		h.Reset()
		io.WriteString(h, curr)
		sum := h.Sum(nil)
		hash := hex.EncodeToString(sum[0:16])
		if curr == hash {
			log(i, "%v == %v <-- MATCH!!!", curr, hash)
		}
		increment(runes)
		curr = string(runes)
	}
	log(i, "finished")
}

func log(i int64, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Printf("%8.1fM) [%v] %v\n", float64(i) / 1000000.0, time.Now().Format("2006-01-02 15:04:05 -0700 MST"), msg)
}

func increment(runes []rune) {
	runeCount := len(runes)
	pos := runeCount - 1
	for ; pos >= 0 && runes[pos] == 'f'; pos-- {}
	if pos < 0 {
		log(-1, "OVERFLOW!")
		pos = 0
	}
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
