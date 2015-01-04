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
	Iterations int64
	I int64
	Curr string
}

type Status int
const (
	TICK Status = iota
	MATCH
	DONE
)

type StatusUpdate struct {
	Status Status
	I int64
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
	gs.I = int64(0)
	runes := []rune(gs.Curr)
	h := md5.New()
	for ; gs.Iterations < 0 || gs.I < gs.Iterations; gs.I++ {
		if gs.I % 10000000 == 0 {
			gs.Log <- StatusUpdate{TICK, gs.I, gs.Curr}
		}
		h.Reset()
		io.WriteString(h, gs.Curr)
		sum := h.Sum(nil)
		hash := hex.EncodeToString(sum[0:16])
		if gs.Curr == hash {
			gs.Log <- StatusUpdate{MATCH, gs.I, gs.Curr}
		}
		increment(runes)
		gs.Curr = string(runes)
	}
	gs.Log <- StatusUpdate{DONE, gs.I, gs.Curr}
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
