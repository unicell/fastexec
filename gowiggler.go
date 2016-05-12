package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"github.com/golang/glog"
)

var Config struct {
	chunks  int
	workers int
}

// Job contains boths command and data to apply command on
type Job struct {
	args   []string
	data   []byte
	result []byte
}

var pwg, cwg sync.WaitGroup
var pending, done chan *Job

func (j *Job) Execute() error {
	defer pwg.Done()
	glog.V(4).Infof("--> executor - cmd - %s", j.args)
	glog.V(4).Infof("--> executor - data\n>>\n%s<<\n", string(j.data))

	cmd := exec.Command(j.args[0], j.args[1:]...)
	cmdIn, _ := cmd.StdinPipe()
	cmdOut, _ := cmd.StdoutPipe()

	cmd.Start()
	_, err := cmdIn.Write(j.data)
	if err != nil {
		glog.Errorf("Error writing to pipe %v", err)
		return err
	}
	cmdIn.Close()

	cmdBytes, err := ioutil.ReadAll(cmdOut)
	cmd.Wait()
	if err != nil {
		glog.Errorf("Error reading from pipe %v", err)
		return err
	}

	j.result = cmdBytes
	glog.V(4).Infof("--> executor - result\n>>\n%s<<\n", string(cmdBytes))

	return nil
}

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

func AddNewJob(job *Job, pending chan<- *Job) {
	glog.V(2).Infof("--> data - submitting chunk...")
	glog.V(4).Infof("\n>\n%s<\n", string(job.data))
	pwg.Add(1)
	cwg.Add(1)
	pending <- job
}

// init data pool by reading from reader input and break it into chunks
func InitDataChunks(args []string, r io.Reader, pending chan *Job) {
	glog.V(2).Infof("--> args - %s", args)
	count := 0
	scanner := bufio.NewScanner(r)
	var buf *bytes.Buffer
	for i := 0; scanner.Scan(); i++ {
		if i%Config.chunks == 0 {
			if buf != nil {
				AddNewJob(&Job{args: args, data: buf.Bytes()}, pending)
				count++
			}
			buf = new(bytes.Buffer)
		}
		token := scanner.Text()
		glog.V(6).Infof("--> writing to buf: %s", token)
		buf.WriteString(token)
		buf.WriteString("\n")
	}

	// handling last chunk
	if buf != nil {
		AddNewJob(&Job{args: args, data: buf.Bytes()}, pending)
		count++
	}
}

func StartWorkers(args []string, pending chan *Job, done chan *Job) {
	// launch work pools
	for i := 0; i < Config.workers; i++ {
		go Worker(pending, done)
	}
	go StateMonitor(done)
}

func init() {
	flag.IntVar(&Config.chunks, "chunks", 1, "size of data chunk for one job")
	flag.IntVar(&Config.workers, "workers", 1, "num of workers")
}

func main() {
	flag.Parse()
	args := flag.Args()
	glog.V(2).Infof("--> cmd to run: %s", args)

	pending, done = make(chan *Job), make(chan *Job)
	StartWorkers(args, pending, done)
	InitDataChunks(args, os.Stdin, pending)

	pwg.Wait()
	cwg.Wait()
}
