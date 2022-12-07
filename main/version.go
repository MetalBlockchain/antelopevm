package main

import "fmt"

var (
	// GitCommit is set by the build script
	GitCommit string
	// Version is the version of Antelope
	Version string = "v0.0.1"
)

func init() {
	if len(GitCommit) != 0 {
		Version = fmt.Sprintf("%s@%s", Version, GitCommit)
	}
}
