package fastexec

import (
	"bufio"
	"bytes"
	"io"

	"github.com/golang/glog"
)

func addJob(job Job, pending chan<- Job) {
	glog.V(2).Infof("--> data - submitting chunk...")
	glog.V(4).Infof("\n>\n%s<\n", string(job.GetData()))
	pwg.Add(1)
	cwg.Add(1)

	pending <- job
}

func initJobPool(args []string, r io.Reader, pending chan Job) {
	glog.V(2).Infof("--> args - %s", args)
	count := 0
	scanner := bufio.NewScanner(r)
	var buf *bytes.Buffer
	for i := 0; scanner.Scan(); i++ {
		if i%Config.chunks == 0 {
			if buf != nil {
				j := &ExecJob{args: args, data: buf.Bytes()}
				addJob(j, pending)
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
		j := &ExecJob{args: args, data: buf.Bytes()}
		addJob(j, pending)
		count++
	}
}

// init job pool by reading from reader input and break it into chunks
func InitJobPool(r io.Reader) {
	initJobPool(Config.args, r, pending)
}
