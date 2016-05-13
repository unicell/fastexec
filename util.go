package fastexec

import (
	"bufio"
	"bytes"
	"io"

	"github.com/golang/glog"
)

func AddNewJob(job *Job, pending chan<- *Job) {
	glog.V(2).Infof("--> data - submitting chunk...")
	glog.V(4).Infof("\n>\n%s<\n", string(job.data))
	pwg.Add(1)
	cwg.Add(1)
	pending <- job
}

func initDataChunks(args []string, r io.Reader, pending chan *Job) {
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

// init data pool by reading from reader input and break it into chunks
func InitDataChunks(r io.Reader) {
	initDataChunks(Config.args, r, pending)
}
