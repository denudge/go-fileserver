package main

import "regexp"

var (
	folderRegexp = regexp.MustCompile("[a-zA-Z0-9_-]+")
	fileRegexp   = regexp.MustCompile("[a-zA-Z0-9_.-]+")

	port       = 8080
	rootDir    = "/tmp"
	bufferSize = 8092
	debug      bool
)
