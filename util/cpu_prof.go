package util

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/burgesQ/webfmwk/log"
)

// enablePprof enable the profiling of the CPU
func enablePprof() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Errorf("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Errorf("could not start CPU profile: ", err)
	}
}

// enableMemProf enable the profiling of the memory
func enableMemProf() (f *os.File) {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Errorf("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Errorf("could not write memory profile: ", err)
	}

	return
}

// EnableProfiling start the CPU & memory profiling
func EnableProfiling() (f *os.File) {
	enablePprof()
	return enableMemProf()
}

// StopProfiling stop the CPU & memory profiling
func StopProfiling(mem *os.File) {
	pprof.StopCPUProfile()
	if e := mem.Close(); e != nil {
		log.Errorf("cannot stop memory profiling : %s", e.Error())
	}
}
