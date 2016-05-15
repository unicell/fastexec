package fastexec

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"
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

func generateFakeJob(i int) FakeJob {
	msg := []string{"fake job ID:", strconv.Itoa(i)}
	data := strings.Join(msg, " ")
	j := FakeJob{data: []byte(data)}
	return j
}

type TimeRange struct {
	min time.Duration
	max time.Duration
}

func TestRateLimitedWorker(t *testing.T) {
	testCases := []struct {
		about     string
		jobNumber int
		chunks    int
		ratelimit int64
		t         TimeRange
	}{
		{
			about:     "fake job run without ratelimit",
			jobNumber: 1,
			chunks:    1,
			ratelimit: -1,
			t:         TimeRange{min: time.Duration(0), max: 100 * time.Microsecond},
		},
		{
			about:     "fake job run ratelimit 1/s chunk 1 - burst",
			jobNumber: 1,
			chunks:    1,
			ratelimit: 1,
			t:         TimeRange{min: time.Duration(0), max: 100 * time.Microsecond},
		},
		{
			about:     "fake job run ratelimit 2/s chunk 1 - burst",
			jobNumber: 4,
			chunks:    1,
			ratelimit: 2,
			t: TimeRange{
				min: 1*time.Second - 10*time.Millisecond,
				max: 1*time.Second + 10*time.Millisecond},
		},
		{
			about:     "fake job run ratelimit 1/s chunk 2",
			jobNumber: 2,
			chunks:    2,
			ratelimit: 1,
			t: TimeRange{
				min: 4*time.Second - 10*time.Millisecond,
				max: 4*time.Second + 10*time.Millisecond},
		},
	}

	for _, c := range testCases {
		Config.chunks = c.chunks
		Config.ratelimit = c.ratelimit
		initWorkers()

		in := make(chan Job, 100)
		out := make(chan Job, 100)

		for i := 0; i < c.jobNumber; i++ {
			j := generateFakeJob(i)
			in <- &j
		}

		close(in)
		start := time.Now()
		Worker(in, out)

		if len(out) != c.jobNumber {
			t.Errorf("testcase: %s, expected processed %d got %d", c.about, c.jobNumber, len(out))
		}

		close(out)
		i := 0
		for j := range out {
			expected := strings.Join([]string{"FAKE JOB ID:", strconv.Itoa(i)}, " ")
			if !bytes.Equal(j.GetResult(), []byte(expected)) {
				t.Errorf("testcase: %s, expected result %v got %v", c.about, expected, string(j.GetResult()))
			}
			i++
		}

		elapsed := time.Now().Sub(start)
		if elapsed < c.t.min || elapsed > c.t.max {
			t.Errorf("testcase: %s, expected time [%v, %v] got %v", c.about, c.t.min, c.t.max, elapsed)
		}
	}
}
