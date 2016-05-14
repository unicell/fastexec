package fastexec

import (
	"bytes"
	"testing"
)

type FakeJob struct {
	data []byte
}

// Execute a Job with associated data
func (j *FakeJob) Execute() error {
	j.data = bytes.ToUpper(j.data)
	return nil
}

// Get input data from a Job
func (j FakeJob) GetData() []byte {
	return j.data
}

// Get output result from a Job
func (j FakeJob) GetResult() []byte {
	return j.data
}

func TestWorker(t *testing.T) {
	testCases := []struct {
		about    string
		jobs     []FakeJob
		expected [][]byte
	}{
		{
			about:    "fake single job execution",
			jobs:     []FakeJob{FakeJob{data: []byte("fakedata")}},
			expected: [][]byte{[]byte("FAKEDATA")},
		},
		{
			about: "fake multiple job execution",
			jobs: []FakeJob{
				FakeJob{data: []byte("fakedata1")},
				FakeJob{data: []byte("fakedata2")},
			},
			expected: [][]byte{
				[]byte("FAKEDATA1"),
				[]byte("FAKEDATA2")},
		},
	}

	for _, c := range testCases {
		in := make(chan Job, 100)
		out := make(chan Job, 100)

		for _, j := range c.jobs {
			// make sure the value is copied on stack
			j := j
			in <- &j
		}

		close(in)
		Worker(in, out)

		if len(out) != len(c.jobs) {
			t.Errorf("testcase: %s, expected processed %d got %d", c.about, len(c.jobs), len(out))
		}

		close(out)
		i := 0
		for j := range out {
			if !bytes.Equal(j.GetResult(), c.expected[i]) {
				t.Errorf("testcase: %s, expected result %v got %v", c.about, string(c.expected[i]), string(j.GetResult()))
			}
			i++
		}
	}
}
