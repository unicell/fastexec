package fastexec

import (
	"bytes"
	"strings"
	"testing"
)

func TestDataChunks(t *testing.T) {
	testCases := []struct {
		args         []string
		data         string
		chunks       int
		expectedData [][]byte
	}{
		{
			args:   []string{"test", "arg1"},
			data:   "A\nB\nC\nD\n",
			chunks: 1,
			expectedData: [][]byte{
				[]byte("A\n"), []byte("B\n"), []byte("C\n"), []byte("D\n"),
			},
		},
		{
			args:   []string{"test", "arg1"},
			data:   "A\nB\nC\nD\n",
			chunks: 2,
			expectedData: [][]byte{
				[]byte("A\nB\n"), []byte("C\nD\n"),
			},
		},
		{
			args:   []string{"test", "arg1"},
			data:   "A\nB\nC\nD\n",
			chunks: 3,
			expectedData: [][]byte{
				[]byte("A\nB\nC\n"), []byte("D\n"),
			},
		},
		{
			args:   []string{"test", "arg1"},
			data:   "A\nB\nC\nD\n",
			chunks: 4,
			expectedData: [][]byte{
				[]byte("A\nB\nC\nD\n"),
			},
		},
	}

	for _, c := range testCases {
		Config.chunks = c.chunks
		p := make(chan *Job, 100)
		initDataChunks(c.args, strings.NewReader(c.data), p)
		close(p)

		if len(c.expectedData) != len(p) {
			t.Errorf("Number of data chunks didn't match:\n\t%v\n\t%v", len(c.expectedData), len(p))
		}

		i := 0
		for j := range p {
			if !bytes.Equal(c.expectedData[i], j.data) {
				t.Errorf("Content of data chunks didn't match:\n\t%v\n\t%v", string(c.expectedData[i]), string(j.data))
			}
			i++
		}
	}
}
