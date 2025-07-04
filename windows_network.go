//go:build windows
// +build windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 可疑端口列表
var suspiciousPorts = map[int]string{
	22:    "SSH",
	23:    "Telnet",
	445:   "SMB",
	1433:  "MSSQL",
	3306:  "MySQL",
	3389:  "RDP",
	4444:  "Metasploit",
	5432:  "PostgreSQL",
	5900:  "VNC",
	6379:  "Redis",
	27017: "MongoDB",
}

// 网络连接分析结果
type NetworkAnalysis struct {
	LocalAddr     string
	RemoteAddr    string
	State         string
	ProcessName   string
	ProcessID     int32
	ListeningPort int
	Protocol      string
}

// 分析网络连接
func analyzeNetworkConnections() {
	fmt.Println("=== 网络连接分析 ===")

	// 获取所有网络连接
	conns, err := net.Connections("all")
	if err != nil {
		fmt.Printf("获取网络连接失败: %v\n", err)
		return
	}

	// 分析每个连接
	for _, conn := range conns {
		// 获取进程信息
		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue
		}

		name, _ := proc.Name()
		localPort := conn.Laddr.Port
		remotePort := conn.Raddr.Port

		// 检查可疑端口
		if service, ok := suspiciousPorts[int(localPort)]; ok {
			fmt.Printf("[警告] 发现可疑端口监听:\n")
			fmt.Printf("端口: %d (%s)\n", localPort, service)
			fmt.Printf("进程: %s (PID: %d)\n", name, conn.Pid)
			fmt.Printf("状态: %s\n\n", conn.Status)
		}

		// 检查可疑远程连接
		if service, ok := suspiciousPorts[int(remotePort)]; ok {
			fmt.Printf("[警告] 发现可疑远程连接:\n")
			fmt.Printf("远程地址: %s:%d (%s)\n", conn.Raddr.IP, remotePort, service)
			fmt.Printf("本地地址: %s:%d\n", conn.Laddr.IP, localPort)
			fmt.Printf("进程: %s (PID: %d)\n", name, conn.Pid)
			fmt.Printf("状态: %s\n\n", conn.Status)
		}
	}
}

// 分析网络接口
func analyzeNetworkInterfaces() {
	fmt.Println("=== 网络接口分析 ===")

	// 获取所有网络接口
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("获取网络接口失败: %v\n", err)
		return
	}

	for _, iface := range ifaces {
		fmt.Printf("接口: %s\n", iface.Name)
		fmt.Printf("MAC地址: %s\n", iface.HardwareAddr)
		fmt.Printf("状态: %v\n", iface.Flags)

		// 获取IP地址
		if len(iface.Addrs) > 0 {
			for _, addr := range iface.Addrs {
				fmt.Printf("IP地址: %s\n", addr.Addr)
			}
		}
		fmt.Println()
	}
}

// 分析网络流量
func analyzeNetworkTraffic() {
	fmt.Println("=== 网络流量分析 ===")

	// 获取网络IO计数器
	ioStats, err := net.IOCounters(true)
	if err != nil {
		fmt.Printf("获取网络流量统计失败: %v\n", err)
		return
	}

	for _, io := range ioStats {
		fmt.Printf("接口: %s\n", io.Name)
		fmt.Printf("发送字节: %d\n", io.BytesSent)
		fmt.Printf("接收字节: %d\n", io.BytesRecv)
		fmt.Printf("发送包数: %d\n", io.PacketsSent)
		fmt.Printf("接收包数: %d\n", io.PacketsRecv)
		fmt.Printf("错误数: %d\n", io.Errin+io.Errout)
		fmt.Printf("丢包数: %d\n\n", io.Dropin+io.Dropout)
	}
}

// 将GBK编码转换为UTF-8 (在windows_network.go中重复定义以避免依赖)
func gbkToUTF8Net(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return string(data), err // 如果转换失败，返回原始数据
	}
	return string(utf8Data), nil
}

// 分析防火墙规则
func analyzeFirewallRules() {
	fmt.Println("=== 防火墙规则分析 ===")

	// 获取防火墙规则
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name=all")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("获取防火墙规则失败: %v\n", err)
		return
	}

	// 转换编码
	utf8Output, convErr := gbkToUTF8Net(output)
	var outputStr string
	if convErr == nil {
		outputStr = utf8Output
	} else {
		outputStr = string(output)
	}

	// 分析输出
	rules := strings.Split(outputStr, "\r\n\r\n")
	for _, rule := range rules {
		if strings.Contains(rule, "允许") && strings.Contains(rule, "入站") {
			fmt.Printf("[注意] 发现入站允许规则:\n%s\n\n", rule)
		}
	}
}

// 检查DNS设置
func checkDNSSettings() {
	fmt.Println("=== DNS设置检查 ===")

	// 获取DNS服务器设置
	cmd := exec.Command("ipconfig", "/all")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("获取DNS设置失败: %v\n", err)
		return
	}

	// 转换编码
	utf8Output, convErr := gbkToUTF8Net(output)
	var outputStr string
	if convErr == nil {
		outputStr = utf8Output
	} else {
		outputStr = string(output)
	}

	// 分析输出
	if strings.Contains(outputStr, "DNS 服务器") {
		dnsServers := strings.Split(outputStr, "DNS 服务器")
		for i := 1; i < len(dnsServers); i++ {
			fmt.Printf("DNS服务器配置:\n%s\n", dnsServers[i])
		}
	}
}