package lib

import "time"

const (
	totalFile     = 3000
	contentLength = 5000
	TimeDuration  = time.Second * 3
)

var tempPath = "/tmp/context-cancelation-pipeline"
