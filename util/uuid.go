package util

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"
)

//------------------------------------------------------------
// ID utils
//------------------------------------------------------------

var (
	_baseID int64
	_lastID int64
)

//------------------------------------------------------------
// Init
//------------------------------------------------------------

func init() {
	randGen := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	_baseID = randGen.Int63()
	_baseID = 0 // XXX let's keep it silly simple
	_lastID = _baseID
}

//------------------------------------------------------------
// Generators
//------------------------------------------------------------

// Super simple ID generator, counter-like atomically increment.
func NewID() int64 {
	return atomic.AddInt64(&_lastID, 1)
}

// Generates proper UUID. Can be somewhat slow.
// In case of error falls back to NewID.
func NewUUID() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		t := time.Now()
		return fmt.Sprintf("%s", t.UnixNano())
	}
	return strings.Trim(fmt.Sprintf("%s", out), "\n")
}
