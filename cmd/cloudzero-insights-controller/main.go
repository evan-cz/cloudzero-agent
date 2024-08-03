package main

import (
	"fmt"

	"github.com/cloudzero/cloudzero-insights-controller/pkg/build"
)

func main() {
	fmt.Printf("CloudZero Insights Controller %s\n", build.GetVersion())
}
