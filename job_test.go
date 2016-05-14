package fastexec

import (
	"bytes"
	"testing"
)

func TestExecJobs(t *testing.T) {
	testCases := []struct {
		about    string
		job      ExecJob
		expected []byte
	}{
		{
			about: "simple cat",
			job: ExecJob{
				args: []string{"cat"},
				data: []byte("A\n"),
			},
			expected: []byte("A\n"),
		},
		{
			about: "sort with argucment",
			job: ExecJob{
				args: []string{"sort", "-n"},
				data: []byte("101\n102\n999\n-9\n"),
			},
			expected: []byte("-9\n101\n102\n999\n"),
		},
	}

	for _, c := range testCases {
		err := c.job.execute()
		if err != nil {
			t.Errorf("testcase: %s, got error %v", c.about, err)
		}
		if !bytes.Equal(c.expected, c.job.GetResult()) {
			t.Errorf("testcase: %s, expected %v got %v", c.about, c.expected, c.job.GetResult())
		}
	}
}
