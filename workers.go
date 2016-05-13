package fastexec

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
)

var pwg, cwg sync.WaitGroup
var pending, done chan *Job

func Worker(in <-chan *Job, out chan<- *Job) {
	for j := range in {
		err := j.Execute()
		if err != nil {
			glog.Errorf("Error executing %v", err)
			continue
		}
		out <- j
	}
}

func StateMonitor(in <-chan *Job) {
	glog.V(2).Infof("--> state monitor")
	for j := range in {
		fmt.Print(string(j.result))
		cwg.Done()
	}
}

func initWorkers() {
	pending, done = make(chan *Job), make(chan *Job)
}

func startWorkers(args []string, pending chan *Job, done chan *Job) {
	// launch work pools
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
