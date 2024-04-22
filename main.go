package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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

		successCount int
		attemptCount int
		stopped      bool
	)

	version := "1.0.0"

	// 设置命令行参数
	showHelp := flag.BoolP("help", "h", false, "显示帮助信息")
	showVersion := flag.BoolP("version", "v", false, "显示版本信息")
	flag.IntVarP(&port, "port", "p", 80, "端口，默认为80")
	flag.IntVarP(&count, "count", "c", 4, "测试次数，默认为4次")
	flag.DurationVarP(&timeoutDuration, "timeout", "s", 2*time.Second, "超时时间，默认为2秒")
	flag.BoolVarP(&infinite, "infinite", "t", false, "无限次测试")
	flag.BoolVarP(&ipv6, "ipv6", "6", false, "使用 IPv6，需搭配域名使用")

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

	if port < 1 || port > 65535 {
		fmt.Println("端口号必须在1-65535之间")
		os.Exit(1)
	}

	hostname := args[0]

	// 设置信号捕获
	stopped = false
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ip := ""

	if net.ParseIP(hostname) == nil {
		// 解析域名
		ips, err := net.LookupIP(hostname)
		if err != nil {
			fmt.Printf("解析 %s 失败: %s\n", hostname, err)
			os.Exit(1)
		}
		record := "A"
		if ipv6 {
			record = "AAAA"
		}
		ip, err = filterIP(ips, ipv6)
		if err != nil {
			fmt.Printf("找不到 %s 的 %s 记录\n", hostname, record)
			os.Exit(1)
		}
		fmt.Printf("使用 %s 的 %s 记录: %s\n", hostname, record, ip)
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
				fmt.Printf("测试到 %s 的连接失败: %s\n", address, "连接超时")
			} else {
				successCount++
				fmt.Printf("来自 %s 的响应: 时间=%s\n", address, fmt.Sprintf("%.3fms", float64(duration)/float64(time.Millisecond)))
				conn.Close()
			}

			if !infinite && attemptCount >= count {
				break
			}
			time.Sleep(1 * time.Second)
		}
		done <- true
	}()

	// 等待信号或测试完成
	select {
	case <-sigChan:
		fmt.Println("\n测试被用户中断")
		stopped = true
	case <-done:
	}

	// 打印统计结果
	if !stopped {
		fmt.Println()
	}
	fmt.Printf("测试完成，成功次数: %d/%d\n", successCount, attemptCount)
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
	return "", fmt.Errorf("找不到合适的IP地址")
}
