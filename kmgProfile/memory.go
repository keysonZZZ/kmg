package kmgProfile
import (
"runtime"
"runtime/debug"
)

func GcAndGetMemoryUsage()uint64{
	Gc()
	memStat:=&runtime.MemStats{}
	runtime.ReadMemStats(memStat)
	return memStat.Alloc
}

func GetMemeoryUsage()uint64{
	memStat:=&runtime.MemStats{}
	runtime.ReadMemStats(memStat)
	return memStat.Alloc
}

func Gc(){
	debug.FreeOSMemory()
}