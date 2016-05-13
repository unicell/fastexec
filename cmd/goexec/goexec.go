package main

import (
	"os"

	"github.com/unicell/fastexec"
)

func main() {
	fastexec.InitFlags()
	fastexec.StartWorkers()
	fastexec.InitDataChunks(os.Stdin)
	fastexec.WaitToFinish()
}
