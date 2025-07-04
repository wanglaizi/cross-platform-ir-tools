//go:build windows
// +build windows

package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v3/mem"
)

// 进程行为监控结构
type ProcessBehavior struct {
	PID           int32
	Name          string
	CPUUsage      float64
	MemoryUsage   float32
	ThreadCount   int32
	HandleCount   int32
	ReadBytes     uint64
	WriteBytes    uint64
	NetworkUsage  uint64
	DLLs          []string
	FileAccesses  []string
}

// 获取进程详细信息
func getProcessDetails(pid int32) (*ProcessBehavior, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}

	behavior := &ProcessBehavior{PID: pid}

	// 获取进程名称
	name, err := proc.Name()
	if err == nil {
		behavior.Name = name
	}

	// 获取CPU使用率
	cpu, err := proc.CPUPercent()
	if err == nil {
		behavior.CPUUsage = cpu
	}

	// 获取内存使用率
	mem, err := proc.MemoryPercent()
	if err == nil {
		behavior.MemoryUsage = mem
	}

	// 获取线程数
	threads, err := proc.NumThreads()
	if err == nil {
		behavior.ThreadCount = threads
	}

	// 获取句柄数 (Windows特定功能，暂时跳过)
	// handles, err := proc.NumHandles()
	// if err == nil {
	//	behavior.HandleCount = handles
	// }

	// 获取IO统计
	io, err := proc.IOCounters()
	if err == nil {
		behavior.ReadBytes = io.ReadBytes
		behavior.WriteBytes = io.WriteBytes
	}

	// 获取网络使用情况
	conns, err := proc.Connections()
	if err == nil {
		behavior.NetworkUsage = uint64(len(conns))
	}

	// 获取加载的DLL (Windows平台可能不支持，暂时跳过)
	// dlls, err := proc.MemoryMaps(true)
	// if err == nil {
	//	for _, dll := range *dlls {
	//		behavior.DLLs = append(behavior.DLLs, dll.Path)
	//	}
	// }

	return behavior, nil
}

// 监控进程行为
func monitorProcessBehavior() {
	fmt.Println("=== 进程行为监控 ===")

	// 获取所有进程
	processes, err := process.Processes()
	if err != nil {
		fmt.Printf("获取进程列表失败: %v\n", err)
		return
	}

	// 监控高CPU和内存使用的进程
	for _, proc := range processes {
		behavior, err := getProcessDetails(proc.Pid)
		if err != nil {
			continue
		}

		// 检查是否为异常行为
		if behavior.CPUUsage > 50 || behavior.MemoryUsage > 50 {
			fmt.Printf("\n发现高资源使用进程:\n")
			fmt.Printf("PID: %d\n", behavior.PID)
			fmt.Printf("名称: %s\n", behavior.Name)
			fmt.Printf("CPU使用率: %.2f%%\n", behavior.CPUUsage)
			fmt.Printf("内存使用率: %.2f%%\n", behavior.MemoryUsage)
			fmt.Printf("线程数: %d\n", behavior.ThreadCount)
			fmt.Printf("句柄数: %d\n", behavior.HandleCount)
			fmt.Printf("读取字节: %d\n", behavior.ReadBytes)
			fmt.Printf("写入字节: %d\n", behavior.WriteBytes)
			fmt.Printf("网络连接数: %d\n", behavior.NetworkUsage)
			
			if len(behavior.DLLs) > 0 {
				fmt.Println("加载的DLL:")
				for _, dll := range behavior.DLLs {
					fmt.Printf("  - %s\n", dll)
				}
			}
		}
	}
}

// 内存分析
func analyzeMemory() {
	fmt.Println("\n=== 内存分析 ===")

	// 获取系统内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("获取内存信息失败: %v\n", err)
		return
	}

	fmt.Printf("内存使用率: %.2f%%\n", memInfo.UsedPercent)
	fmt.Printf("总物理内存: %.2f GB\n", float64(memInfo.Total)/(1024*1024*1024))
	fmt.Printf("可用物理内存: %.2f GB\n", float64(memInfo.Available)/(1024*1024*1024))
	fmt.Printf("已使用内存: %.2f GB\n", float64(memInfo.Used)/(1024*1024*1024))
	fmt.Printf("空闲内存: %.2f GB\n", float64(memInfo.Free)/(1024*1024*1024))

	// 分析大内存进程
	fmt.Println("\n内存使用TOP 10进程:")
	processes, _ := process.Processes()
	type ProcessMemInfo struct {
		pid      int32
		name     string
		memory   float32
		path     string
		cmdline  string
	}

	var processMemList []ProcessMemInfo
	for _, p := range processes {
		name, _ := p.Name()
		mem, _ := p.MemoryPercent()
		path, _ := p.Exe()
		cmd, _ := p.Cmdline()
		processMemList = append(processMemList, ProcessMemInfo{
			pid:     p.Pid,
			name:    name,
			memory:  mem,
			path:    path,
			cmdline: cmd,
		})
	}

	// 按内存使用排序
	for i := 0; i < len(processMemList)-1; i++ {
		for j := 0; j < len(processMemList)-i-1; j++ {
			if processMemList[j].memory < processMemList[j+1].memory {
				processMemList[j], processMemList[j+1] = processMemList[j+1], processMemList[j]
			}
		}
	}

	// 显示前10个进程
	for i := 0; i < 10 && i < len(processMemList); i++ {
		fmt.Printf("\nPID: %d\n", processMemList[i].pid)
		fmt.Printf("名称: %s\n", processMemList[i].name)
		fmt.Printf("内存使用率: %.2f%%\n", processMemList[i].memory)
		fmt.Printf("路径: %s\n", processMemList[i].path)
		fmt.Printf("命令行: %s\n", processMemList[i].cmdline)
	}
}