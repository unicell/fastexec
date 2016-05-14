package fastexec

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestExecJobs(t *testing.T) {
	testCases := []struct {
		about    string
		job      ExecJob
		err      error
		expected []byte
	}{
		{
			about: "simple cat",
			job: ExecJob{
				args: []string{"cat"},
				data: []byte("A\n"),
			},
			err:      nil,
			expected: []byte("A\n"),
		},
		{
			about: "sort with argucment",
			job: ExecJob{
				args: []string{"sort", "-n"},
				data: []byte("101\n102\n999\n-9\n"),
			},
			err:      nil,
			expected: []byte("-9\n101\n102\n999\n"),
		},
		{
			about: "simple cat1",
			job: ExecJob{
				args: []string{"non-existing-command"},
				data: []byte("A\n"),
			},
			err:      &os.PathError{"write", "|1", errors.New("bad file descriptor")},
			expected: []byte(""),
		},
	}

	for _, c := range testCases {
		err := c.job.execute()
		switch err.(type) {
		case error:
			if err.Error() != c.err.Error() {
				t.Errorf("testcase: %s, expected - %T got - %T", c.about, c.err, err)
			}
		default:
			if !bytes.Equal(c.expected, c.job.GetResult()) {
				t.Errorf("testcase: %s, expected - %v got - %v", c.about, c.expected, c.job.GetResult())
			}
		}
	}
}
