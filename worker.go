package fastexec

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/juju/ratelimit"
)

var pwg, cwg sync.WaitGroup
var pending, done chan Job

var tb *ratelimit.Bucket

func Worker(in <-chan Job, out chan<- Job) {
	for j := range in {
		if Config.ratelimit > 0 {
			wait := tb.Take(int64(Config.chunks))
			if wait > 0 {
				glog.V(4).Infof("ratelimited - sleeping - %v", wait)
				time.Sleep(wait)
			}
		}
		err := j.Execute()
		if err != nil {
			glog.Errorf("Error executing %v", err)
			continue
		}
		out <- j
	}
}

func StateMonitor(in <-chan Job) {
	glog.V(2).Infof("--> state monitor")
	for j := range in {
		fmt.Print(string(j.GetResult()))
		cwg.Done()
	}
}

// init channels for workers
// optionally enable ratelimit if enabled
func initWorkers() {
	pending, done = make(chan Job), make(chan Job)

	// initialize ratelimit control if non-zero value specified
	if Config.ratelimit > 0 {
		quantum := int64(Config.chunks)
		freqency := float64(Config.ratelimit) / float64(Config.chunks)
		interval := time.Duration(float64(1) / float64(freqency) * 1e9)
		glog.V(2).Infof("--> ratelimit - rate %d quantum %d interval %v", Config.ratelimit, quantum, interval)

		tb = ratelimit.NewBucketWithQuantum(interval, Config.ratelimit, quantum)
	}
}

func startWorkers(args []string, pending chan Job, done chan Job) {
	// launch workers
	for i := 0; i < Config.workers; i++ {
		go Worker(pending, done)
	}
	go StateMonitor(done)
}

func StartWorkers() {
	initWorkers()
	startWorkers(Config.args, pending, done)
}

func WaitToFinish() {
	pwg.Wait()
	cwg.Wait()
}
