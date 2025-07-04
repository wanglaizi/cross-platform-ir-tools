//go:build windows
// +build windows

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func checkError(err error) {
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	}
}

func getSystemInfo() {
	fmt.Println("=== 系统信息 ===")
	hostInfo, err := host.Info()
	checkError(err)
	if err == nil {
		fmt.Printf("主机名: %s\n", hostInfo.Hostname)
		fmt.Printf("操作系统: %s\n", hostInfo.OS)
		fmt.Printf("平台: %s\n", hostInfo.Platform)
		fmt.Printf("平台版本: %s\n", hostInfo.PlatformVersion)
		fmt.Printf("内核版本: %s\n", hostInfo.KernelVersion)
		fmt.Printf("启动时间: %s\n", time.Unix(int64(hostInfo.BootTime), 0))
	}
}

func getCPUInfo() {
	fmt.Println("\n=== CPU信息 ===")
	cpuInfo, err := cpu.Info()
	checkError(err)
	if err == nil {
		for _, info := range cpuInfo {
			fmt.Printf("CPU型号: %s\n", info.ModelName)
			fmt.Printf("核心数: %d\n", info.Cores)
			fmt.Printf("频率: %.2f MHz\n", info.Mhz)
		}
	}

	percentages, err := cpu.Percent(time.Second, true)
	checkError(err)
	if err == nil {
		for i, percentage := range percentages {
			fmt.Printf("CPU%d使用率: %.2f%%\n", i, percentage)
		}
	}
}

func getMemoryInfo() {
	fmt.Println("\n=== 内存信息 ===")
	virtual, err := mem.VirtualMemory()
	checkError(err)
	if err == nil {
		fmt.Printf("总内存: %.2f GB\n", float64(virtual.Total)/(1024*1024*1024))
		fmt.Printf("可用内存: %.2f GB\n", float64(virtual.Available)/(1024*1024*1024))
		fmt.Printf("内存使用率: %.2f%%\n", virtual.UsedPercent)
	}
}

func getDiskInfo() {
	fmt.Println("\n=== 磁盘信息 ===")
	partitions, err := disk.Partitions(true)
	checkError(err)
	if err == nil {
		for _, partition := range partitions {
			fmt.Printf("\n分区: %s\n", partition.Device)
			fmt.Printf("挂载点: %s\n", partition.Mountpoint)
			fmt.Printf("文件系统: %s\n", partition.Fstype)

			usage, err := disk.Usage(partition.Mountpoint)
			if err == nil {
				fmt.Printf("总空间: %.2f GB\n", float64(usage.Total)/(1024*1024*1024))
				fmt.Printf("已用空间: %.2f GB\n", float64(usage.Used)/(1024*1024*1024))
				fmt.Printf("可用空间: %.2f GB\n", float64(usage.Free)/(1024*1024*1024))
				fmt.Printf("使用率: %.2f%%\n", usage.UsedPercent)
			}
		}
	}
}

func getNetworkInfo() {
	fmt.Println("\n=== 网络信息 ===")
	interfaces, err := net.Interfaces()
	checkError(err)
	if err == nil {
		for _, iface := range interfaces {
			fmt.Printf("\n网卡名称: %s\n", iface.Name)
			fmt.Printf("MAC地址: %s\n", iface.HardwareAddr)
			fmt.Printf("状态: %v\n", iface.Flags)

			if len(iface.Addrs) > 0 {
				for _, addr := range iface.Addrs {
					fmt.Printf("IP地址: %s\n", addr.Addr)
				}
			}
		}
	}

	conns, err := net.Connections("all")
	checkError(err)
	if err == nil {
		fmt.Printf("\n活动连接数: %d\n", len(conns))
		for i, conn := range conns {
			if i >= 5 { // 只显示前5个连接
				break
			}
			fmt.Printf("本地地址: %s:%d\n", conn.Laddr.IP, conn.Laddr.Port)
			if conn.Raddr.IP != "" {
				fmt.Printf("远程地址: %s:%d\n", conn.Raddr.IP, conn.Raddr.Port)
			}
			fmt.Printf("状态: %s\n\n", conn.Status)
		}
	}
}

func getProcessInfo() {
	fmt.Println("\n=== 进程信息 ===")
	processes, err := process.Processes()
	checkError(err)
	if err == nil {
		fmt.Printf("总进程数: %d\n", len(processes))
		fmt.Println("\nCPU使用率最高的进程:")

		type ProcessInfo struct {
			pid     int32
			name    string
			cpu     float64
			memory  float32
			cmdline string
		}

		var processInfos []ProcessInfo
		for _, p := range processes {
			name, _ := p.Name()
			cpu, _ := p.CPUPercent()
			mem, _ := p.MemoryPercent()
			cmd, _ := p.Cmdline()
			processInfos = append(processInfos, ProcessInfo{
				pid:     p.Pid,
				name:    name,
				cpu:     cpu,
				memory:  mem,
				cmdline: cmd,
			})
		}

		// 按CPU使用率排序
		for i := 0; i < len(processInfos)-1; i++ {
			for j := 0; j < len(processInfos)-i-1; j++ {
				if processInfos[j].cpu < processInfos[j+1].cpu {
					processInfos[j], processInfos[j+1] = processInfos[j+1], processInfos[j]
				}
			}
		}

		// 显示前5个进程
		for i := 0; i < 5 && i < len(processInfos); i++ {
			fmt.Printf("PID: %d\n", processInfos[i].pid)
			fmt.Printf("名称: %s\n", processInfos[i].name)
			fmt.Printf("CPU使用率: %.2f%%\n", processInfos[i].cpu)
			fmt.Printf("内存使用率: %.2f%%\n", processInfos[i].memory)
			fmt.Printf("命令行: %s\n\n", processInfos[i].cmdline)
		}
	}
}

// 将GBK编码转换为UTF-8 (在windows_ir.go中重复定义以避免依赖)
func gbkToUTF8IR(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return string(data), err // 如果转换失败，返回原始数据
	}
	return string(utf8Data), nil
}

func getAutoRuns() {
	fmt.Println("\n=== 自启动项检查 ===")
	// 检查注册表自启动项
	cmd := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	output, err := cmd.Output()
	if err == nil {
		fmt.Println("系统自启动项:")
		utf8Output, convErr := gbkToUTF8IR(output)
		if convErr == nil {
			fmt.Println(utf8Output)
		} else {
			fmt.Println(string(output))
		}
	}

	cmd = exec.Command("reg", "query", "HKEY_CURRENT_USER\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run")
	output, err = cmd.Output()
	if err == nil {
		fmt.Println("用户自启动项:")
		utf8Output, convErr := gbkToUTF8IR(output)
		if convErr == nil {
			fmt.Println(utf8Output)
		} else {
			fmt.Println(string(output))
		}
	}

	// 检查启动文件夹
	startupPath := filepath.Join(os.Getenv("APPDATA"), "Microsoft\\Windows\\Start Menu\\Programs\\Startup")
	fmt.Printf("\n启动文件夹 (%s) 内容:\n", startupPath)
	files, err := os.ReadDir(startupPath)
	if err == nil {
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}

func getScheduledTasks() {
	fmt.Println("\n=== 计划任务检查 ===")
	cmd := exec.Command("schtasks", "/query", "/fo", "LIST")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8IR(output)
		var outputStr string
		if convErr == nil {
			outputStr = utf8Output
		} else {
			outputStr = string(output)
		}
		tasks := strings.Split(outputStr, "\n")
		for _, task := range tasks {
			if strings.Contains(task, "TaskName:") {
				fmt.Println(task)
			}
		}
	}
}