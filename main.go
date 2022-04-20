package main

import (
	"os/exec"
)

func CopyGeneratedProtoFilesToMount() {
	cmd := exec.Command("cp", "--recursive", "proto/", "generated_proto/")
	cmd.Run()
}

func main() {
	CopyGeneratedProtoFilesToMount()
}
