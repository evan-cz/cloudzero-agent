package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	cpuUsagePtr := flag.Float64("cpu", 0.2, "CPU usage as a fraction (e.g., 0.2 for 200m)")
	sleepIntervalPtr := flag.Duration("sleep", 0, "Sleep interval in milliseconds between CPU stress cycles (e.g., 100)")
	memUsagePtr := flag.Int("mem", 0, "Memory usage in MB")
	flag.Parse()

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	numGoroutines := int(float64(numCPU) * (*cpuUsagePtr))

	fmt.Printf("Starting CPU and memory stress with %d goroutines...\n", numGoroutines)

	ctx, cancel := context.WithCancel(context.Background())

	// Capture termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for i := 0; i < numGoroutines; i++ {
		go func(ctx context.Context) {
			mem := make([]byte, *memUsagePtr*1024*1024) // Allocate memory
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Simulate CPU work
					for j := 0; j < 1_000_000; j++ {
					}
					// Simulate memory usage
					for k := range mem {
						mem[k] = byte(k)
					}
					if *sleepIntervalPtr > 0 {
						time.Sleep(*sleepIntervalPtr)
					}
				}
			}
		}(ctx)
	}

	go func() {
		// Wait for termination signal
		<-quit
		fmt.Println("Termination signal received. Stopping CPU and memory stress...")
		cancel()
	}()

	// Run stress indefinitely
	fmt.Println("CPU and memory stress will run indefinitely. Press Ctrl+C to stop.")
	<-ctx.Done()
}
