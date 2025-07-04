//go:build windows
// +build windows

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shirou/gopsutil/v3/host"
)

func runIncidentResponse() {
	fmt.Println("Windows系统应急响应工具 v1.0")
	getSystemInfo()
	getCPUInfo()
	getMemoryInfo()
	getDiskInfo()
	getNetworkInfo()
	getProcessInfo()
	getAutoRuns()
	getScheduledTasks()
}

// isAdmin 检查程序是否以管理员权限运行
func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

// OutputCapture 结构体用于捕获输出
type OutputCapture struct {
	buffer bytes.Buffer
	oldStdout *os.File
}

// 开始捕获输出
func (oc *OutputCapture) Start() error {
	// 保存原始stdout
	oc.oldStdout = os.Stdout

	// 创建管道
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	// 设置新的stdout
	os.Stdout = w

	// 启动goroutine来读取输出
	go func() {
		io.Copy(&oc.buffer, r)
	}()

	return nil
}

// 停止捕获输出
func (oc *OutputCapture) Stop() string {
	// 恢复原始stdout
	os.Stdout = oc.oldStdout

	// 返回捕获的输出
	return oc.buffer.String()
}

func main() {
	// 检查管理员权限
	if !isAdmin() {
		fmt.Println("错误：此工具需要管理员权限运行")
		os.Exit(1)
	}

	// 解析命令行参数
	var (
		runAll      = flag.Bool("all", false, "运行所有检查")
		runIR       = flag.Bool("ir", false, "运行基础应急响应检查")
		runReg      = flag.Bool("reg", false, "运行注册表和文件完整性检查")
		runMemory   = flag.Bool("mem", false, "运行内存和进程行为分析")
		runLog      = flag.Bool("log", false, "运行系统日志分析")
		runNet      = flag.Bool("net", false, "运行网络安全分析")
		runBaseline = flag.Bool("baseline", false, "运行系统安全基线检查")
		genReport   = flag.Bool("report", true, "生成HTML格式检查报告")
	)

	flag.Parse()

	// 如果没有指定任何参数，显示帮助信息
	if !*runAll && !*runIR && !*runReg && !*runMemory && !*runLog && !*runNet && !*runBaseline {
		flag.Usage()
		os.Exit(1)
	}

	if *runAll || *runIR {
		fmt.Println("\n[+] 开始基础应急响应检查...")
		runIncidentResponse()
	}

	if *runAll || *runReg {
		fmt.Println("\n[+] 开始注册表和文件完整性检查...")
		checkRegistry()
		checkSystemFileIntegrity()
		checkSuspiciousFiles()
	}

	if *runAll || *runMemory {
		fmt.Println("\n[+] 开始内存和进程行为分析...")
		analyzeMemory()
		monitorProcessBehavior()
	}

	if *runAll || *runLog {
		fmt.Println("\n[+] 开始系统日志分析...")
		analyzeSystemLogs()
		analyzeSecurityLogs()
		analyzeApplicationLogs()
		analyzePowerShellLogs()
		analyzeLogFiles()
	}

	if *runAll || *runNet {
		fmt.Println("\n[+] 开始网络安全分析...")
		analyzeNetworkConnections()
		analyzeNetworkInterfaces()
		analyzeNetworkTraffic()
		analyzeFirewallRules()
		checkDNSSettings()
	}

	if *runAll || *runBaseline {
		fmt.Println("\n[+] 开始系统安全基线检查...")
		checkPasswordPolicy()
		checkUserAccounts()
		checkSystemServices()
		checkSystemPatches()
		checkAuditPolicy()
		checkFileSystemPermissions()
		checkShareSettings()
		checkUACSettings()
		checkWindowsDefender()
	}

	// 生成报告时获取系统信息

	// 如果需要生成报告
	if *genReport {
		// 创建检查结果列表
		var results []CheckResult

		// 添加系统信息作为基本信息
		hostInfo, _ := host.Info()
		sysInfo := fmt.Sprintf("主机名: %s\n操作系统: %s\n平台: %s %s\n", 
			hostInfo.Hostname, hostInfo.OS, hostInfo.Platform, hostInfo.PlatformVersion)

		// 生成报告
		if err := generateReport(results, sysInfo); err != nil {
			fmt.Printf("生成报告失败: %v\n", err)
		}
	}
}