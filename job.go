package fastexec

import (
	"io/ioutil"
	"os/exec"

	"github.com/golang/glog"
)

type Job interface {
	Execute() error
	GetData() []byte
	GetResult() []byte
}

// Job contains boths command and data to apply command on
type ExecJob struct {
	args   []string
	data   []byte
	result []byte
}

func (j *ExecJob) execute() error {
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

// Execute a Job with associated data
func (j *ExecJob) Execute() error {
	defer pwg.Done()
	return j.execute()
}

// Get input data from a Job
func (j ExecJob) GetData() []byte {
	return j.data
}

// Get output result from a Job
func (j ExecJob) GetResult() []byte {
	return j.result
}
