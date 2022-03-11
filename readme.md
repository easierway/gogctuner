# GOGCTuner [beta]

idea is from this article [How We Saved 70K Cores Across 30 Mission-Critical Services (Large-Scale, Semi-Automated Go GC Tuning @Uber) ](https://eng.uber.com/how-we-saved-70k-cores-across-30-mission-critical-services/)

** This version is updated by Chao.Cai (chaocai@icloud.com) to reduce the GC effect.

## How to use this lib?

just call NewTuner when initializing app :

```go
func initProcess() {
	var (
		inCgroup = false //true: inside docker; false: outsie docker (VM)
		// highPercent & lowPercent is to control GC tigger and scope
		highPercent = 70
		lowPercent = 30 
	)
	go NewTuner(inCgroup, percent)
}
```

## Current Status

more tests are needed

