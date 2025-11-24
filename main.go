package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	var (
		port            int
		count           int
		timeoutDuration time.Duration
		infinite        bool
		ipv6            bool
		fast            bool

		successCount int
		successDelay []time.Duration
		attemptCount int
		stopped      bool
	)

	version := "1.2.0"

	showHelp := flag.BoolP("help", "h", false, "Show help")
	showVersion := flag.BoolP("version", "v", false, "Show version")
	flag.IntVarP(&count, "count", "c", 4, "Number of probes, default 4")
	flag.DurationVarP(&timeoutDuration, "timeout", "s", 2*time.Second, "Timeout, default 2s")
	flag.BoolVarP(&infinite, "infinite", "t", false, "Infinite probes")
	flag.BoolVarP(&ipv6, "ipv6", "6", false, "Use IPv6; requires domain name")
	flag.BoolVarP(&fast, "fast", "f", false, "Fast mode; reduce delay between successful probes")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("tcping v%s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	hostname := args[0]
	port = 80
	if len(args) >= 2 {
		p, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Port must be an integer")
			os.Exit(1)
		}
		port = p
	}
	if port < 1 || port > 65535 {
		fmt.Println("Port must be between 1 and 65535")
		os.Exit(1)
	}

	// 设置信号捕获
	stopped = false
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ip := ""

	if net.ParseIP(hostname) == nil {
		// 解析域名
		ips, err := net.LookupIP(hostname)
		if err != nil {
			fmt.Printf("Failed to resolve %s: %s\n", hostname, err)
			os.Exit(1)
		}
		record := "A"
		if ipv6 {
			record = "AAAA"
		}
		ip, err = filterIP(ips, ipv6)
		if err != nil {
			fmt.Printf("No %s record found for %s\n", record, hostname)
			os.Exit(1)
		}
		fmt.Printf("Using %s %s record: %s\n", hostname, record, ip)
	} else {
		ip = hostname
	}

	address := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

	if infinite {
		count = -1
	}

	done := make(chan bool, 1)
	go func() {
		for i := 0; infinite || i < count; i++ {
			start := time.Now()
			conn, err := net.DialTimeout("tcp", address, timeoutDuration)
			duration := time.Since(start)
			attemptCount++

			fmt.Printf("[%d] ", i+1)
			if err != nil {
				fmt.Printf("Connection to %s failed: %s\n", address, "timeout")
			} else {
				successCount++
				successDelay = append(successDelay, duration)
				fmt.Printf("Reply from %s: time=%s\n", address, fmt.Sprintf("%.3fms", float64(duration)/float64(time.Millisecond)))
				conn.Close()
			}

			if !infinite && attemptCount >= count {
				break
			}
			if fast {
				time.Sleep(150 * time.Millisecond)
			} else {
				time.Sleep(1 * time.Second)
			}
		}
		done <- true
	}()

	select {
	case <-sigChan:
		fmt.Println("\nTest interrupted by user")
		stopped = true
	case <-done:
	}

	if !stopped {
		fmt.Println()
	}

	successDelayMs := make([]float64, len(successDelay))
	for i, delay := range successDelay {
		successDelayMs[i] = float64(delay) / float64(time.Millisecond)
	}

	minDelay := 0.0
	maxDelay := 0.0
	avgDelay := 0.0
	successRate := 0.0

	if successCount > 0 {
		minDelay = float64_min(successDelayMs)
		maxDelay = float64_max(successDelayMs)
		avgDelay = float64_avg(successDelayMs)
	}
	if attemptCount > 0 {
		successRate = float64(successCount) / float64(attemptCount) * 100
	}
	fmt.Printf("Test finished, success %d/%d (%.2f%%)\nmin = %.3fms, max = %.3fms, avg = %.3fms\n", successCount, attemptCount, successRate, minDelay, maxDelay, avgDelay)
}

func filterIP(ips []net.IP, ipv6 bool) (string, error) {
	if ipv6 {
		for _, ip := range ips {
			if ip.To16() != nil && ip.To4() == nil {
				return ip.String(), nil
			}
		}
	} else {
		for _, ip := range ips {
			if ip.To4() != nil && ip.To16() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no suitable IP address found")
}

func float64_min(array []float64) float64 {
	min := array[0]
	for _, value := range array {
		if value < min {
			min = value
		}
	}
	return min
}

func float64_max(array []float64) float64 {
	max := array[0]
	for _, value := range array {
		if value > max {
			max = value
		}
	}
	return max
}

func float64_avg(array []float64) float64 {
	sum := 0.0
	for _, value := range array {
		sum += value
	}
	return sum / float64(len(array))
}
