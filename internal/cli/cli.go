package cli

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/taurusxin/tcping-go/internal/filter"
	"github.com/taurusxin/tcping-go/internal/stats"
)

const version = "1.3.0"

// Run executes the tcping CLI.
func Run() {
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

	showHelp := flag.BoolP("help", "h", false, "Show help")
	showVersion := flag.BoolP("version", "v", false, "Show version")
	flag.IntVarP(&count, "count", "c", 4, "Number of probes")
	flag.DurationVarP(&timeoutDuration, "timeout", "s", 2*time.Second, "Timeout")
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

	stopped = false
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ip := ""

	if net.ParseIP(hostname) == nil {
		ips, err := net.LookupIP(hostname)
		if err != nil {
			fmt.Printf("Failed to resolve %s: %s\n", hostname, err)
			os.Exit(1)
		}
		record := "A"
		if ipv6 {
			record = "AAAA"
		}
		ip, err = filter.IP(ips, ipv6)
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
		probe := 0
		for i := 0; infinite || i < count; i++ {
			start := time.Now()
			conn, err := net.DialTimeout("tcp", address, timeoutDuration)
			duration := time.Since(start)

			if err != nil {
				attemptCount++
				probe++
				fmt.Printf("[%d] Connection to %s failed: %s\n", probe, address, "timeout")
			} else if duration < 500*time.Microsecond {
				conn.Close()
				if !infinite {
					count++
				}
				continue
			} else {
				attemptCount++
				probe++
				successCount++
				successDelay = append(successDelay, duration)
				fmt.Printf("[%d] Reply from %s: time=%.3fms\n", probe, address, float64(duration)/float64(time.Millisecond))
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
		minDelay = stats.Min(successDelayMs)
		maxDelay = stats.Max(successDelayMs)
		avgDelay = stats.Avg(successDelayMs)
	}
	if attemptCount > 0 {
		successRate = float64(successCount) / float64(attemptCount) * 100
	}
	fmt.Printf("Test finished, success %d/%d (%.2f%%)\nmin = %.3fms, max = %.3fms, avg = %.3fms\n", successCount, attemptCount, successRate, minDelay, maxDelay, avgDelay)
}
