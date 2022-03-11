// updated Uber's GOGCTuner to bring the effect as "ballast"
// chaocai2001@icloud.com
package gogctuner

import (
	"runtime"
	"runtime/debug"
	"time"
)

type finalizer struct {
	ch  chan time.Time
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

// default GOGC = 100
// FIXME, if use sets the GOGC at startup
// how to get this value correctly?
var previousGOGC = 100

// don't trigger err log on every failure
var failCounter = -1

func getCurrentPercentAndChangeGOGC() {
	memPercent, err := getUsage()

	if err != nil {
		failCounter++
		if failCounter%10 == 0 {
			println("failed to adjust GC", err.Error())
		}
		return
	}
	// hard_target = live_dataset + live_dataset * (GOGC / 100).
	// 	hard_target =  memoryLimitInPercent
	// 	live_dataset = memPercent
	//  so gogc = (hard_target - livedataset) / live_dataset * 100

	if memPercent < memoryBottomInPercent {
		previousGOGC = debug.SetGCPercent(int(smallGCPercent))
		return
	}

	newgogc := (memoryLimitInPercent - memPercent + memoryBottomInPercent) / memPercent * 100.0

	// if newgogc < 0, we have to use the previous gogc to determine the next
	if newgogc < 0 {
		newgogc = float64(previousGOGC) * memoryLimitInPercent / memPercent
	}

	previousGOGC = debug.SetGCPercent(int(newgogc))
}

func finalizerHandler(f *finalizerRef) {
	select {
	case f.parent.ch <- time.Time{}:
	default:
	}

	getCurrentPercentAndChangeGOGC()
	runtime.SetFinalizer(f, finalizerHandler)
}

// NewTuner
//   set useCgroup to true if your app is in docker
//   set highPercent and lowPercent to control the gc trigger

func NewTuner(useCgroup bool, highPercent float64, lowPercent float64) *finalizer {
	if useCgroup {
		getUsage = getUsageCGroup
	} else {
		getUsage = getUsageNormal
	}

	memoryLimitInPercent = highPercent

	f := &finalizer{
		ch: make(chan time.Time, 1),
	}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
	return f
}
