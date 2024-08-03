package main

import (
	"fmt"

	"github.com/cloudzero/cloudzero-insights-controller/pkg/build"
)

func main() {
	fmt.Printf("CloudZero Agent Validator %s\n", build.GetVersion())
}
