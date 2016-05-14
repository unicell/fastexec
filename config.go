package fastexec

import (
	"flag"

	"github.com/golang/glog"
)

var Config struct {
	chunks    int
	workers   int
	args      []string
	ratelimit int64
}

func InitFlags() {
	flag.IntVar(&Config.chunks, "chunks", 1, "size of data chunk for one job")
	flag.IntVar(&Config.workers, "workers", 1, "num of workers")
	flag.Int64Var(&Config.ratelimit, "ratelimit", -1, "data processing ratelimit / sec")
	flag.Parse()

	Config.args = flag.Args()
	glog.V(2).Infof("--> cmd to run: %s", Config.args)
}
