package kember

import (
	"io"
	"strings"
	"crypto/md5"
	"encoding/hex"
)

type Searcher struct {
	Log chan StatusUpdate
	Start string
	Iterations uint64
}

type Status uint
const (
	TICK Status = iota
	MATCH
	DONE
)

type StatusUpdate struct {
	Status Status
	I uint64
	Curr string
}

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

func Search(gs *Searcher) {
	i := uint64(0)
	curr := gs.Start
	runes := []rune(curr)
	h := md5.New()
	for ; gs.Iterations == 0 || i < gs.Iterations; i++ {
		if i % 10000000 == 0 {
			gs.Log <- StatusUpdate{TICK, i, curr}
		}
		h.Reset()
		io.WriteString(h, curr)
		sum := h.Sum(nil)
		hash := hex.EncodeToString(sum[0:16])
		if curr == hash {
			gs.Log <- StatusUpdate{MATCH, i, curr}
		}
		increment(runes)
		curr = string(runes)
	}
	gs.Log <- StatusUpdate{DONE, i, curr}
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
