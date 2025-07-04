//go:build windows
// +build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"bytes"
)

// 安全基线检查项
type BaselineCheck struct {
	Name        string
	Description string
	Status      string
	Severity    string
	Remediation string
}

// 将GBK编码转换为UTF-8
func gbkToUTF8(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return string(data), err // 如果转换失败，返回原始数据
	}
	return string(utf8Data), nil
}

// 检查密码策略
func checkPasswordPolicy() {
	fmt.Println("=== 密码策略检查 ===")

	cmd := exec.Command("net", "accounts")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("当前密码策略:\n%s\n", utf8Output)
		} else {
			fmt.Printf("当前密码策略:\n%s\n", string(output))
		}
	}

	// 检查密码复杂度要求
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Netlogon\Parameters`, registry.READ)
	if err == nil {
		defer key.Close()
		if val, _, err := key.GetIntegerValue("RequireStrongKey"); err == nil && val == 0 {
			fmt.Println("[警告] 未启用强密码要求")
		}
	}
}

// 检查用户账户设置
func checkUserAccounts() {
	fmt.Println("\n=== 用户账户检查 ===")

	// 检查管理员组成员
	cmd := exec.Command("net", "localgroup", "Administrators")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("管理员组成员:\n%s\n", utf8Output)
		} else {
			fmt.Printf("管理员组成员:\n%s\n", string(output))
		}
	}

	// 检查来宾账户状态
	cmd = exec.Command("net", "user", "Guest")
	output, err = cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			if !strings.Contains(utf8Output, "Account active               No") {
				fmt.Println("[警告] Guest账户未禁用")
			}
		} else {
			if !strings.Contains(string(output), "Account active               No") {
				fmt.Println("[警告] Guest账户未禁用")
			}
		}
	}
}

// 检查系统服务
func checkSystemServices() {
	fmt.Println("\n=== 系统服务检查 ===")

	// 检查关键服务状态
	criticalServices := []string{
		"Windows Defender",
		"Windows Firewall",
		"Windows Update",
		"Remote Registry",
	}

	for _, service := range criticalServices {
		cmd := exec.Command("sc", "query", service)
		output, err := cmd.Output()
		if err == nil {
			utf8Output, convErr := gbkToUTF8(output)
			if convErr == nil {
				if strings.Contains(utf8Output, "RUNNING") {
					fmt.Printf("%s: 运行中\n", service)
				} else {
					fmt.Printf("[警告] %s: 未运行\n", service)
				}
			} else {
				if strings.Contains(string(output), "RUNNING") {
					fmt.Printf("%s: 运行中\n", service)
				} else {
					fmt.Printf("[警告] %s: 未运行\n", service)
				}
			}
		}
	}
}

// 检查系统补丁
func checkSystemPatches() {
	fmt.Println("\n=== 系统补丁检查 ===")

	cmd := exec.Command("wmic", "qfe", "list", "brief")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("已安装的补丁:\n%s\n", utf8Output)
		} else {
			fmt.Printf("已安装的补丁:\n%s\n", string(output))
		}
	}
}

// 检查系统审计策略
func checkAuditPolicy() {
	fmt.Println("\n=== 审计策略检查 ===")

	cmd := exec.Command("auditpol", "/get", "/category:*")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("当前审计策略:\n%s\n", utf8Output)
		} else {
			fmt.Printf("当前审计策略:\n%s\n", string(output))
		}
	}
}

// 检查文件系统权限
func checkFileSystemPermissions() {
	fmt.Println("\n=== 文件系统权限检查 ===")

	// 检查系统关键目录权限
	criticalPaths := []string{
		"C:\\Windows\\System32",
		"C:\\Windows\\System32\\config",
		"C:\\Program Files",
		"C:\\Program Files (x86)",
	}

	for _, path := range criticalPaths {
		cmd := exec.Command("icacls", path)
		output, err := cmd.Output()
		if err == nil {
			utf8Output, convErr := gbkToUTF8(output)
			if convErr == nil {
				fmt.Printf("%s 权限:\n%s\n", path, utf8Output)
			} else {
				fmt.Printf("%s 权限:\n%s\n", path, string(output))
			}
		}
	}
}

// 检查共享设置
func checkShareSettings() {
	fmt.Println("\n=== 共享设置检查 ===")

	cmd := exec.Command("net", "share")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("当前共享:\n%s\n", utf8Output)
		} else {
			fmt.Printf("当前共享:\n%s\n", string(output))
		}
	}
}

// 检查UAC设置
func checkUACSettings() {
	fmt.Println("\n=== UAC设置检查 ===")

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System`, registry.READ)
	if err == nil {
		defer key.Close()

		if val, _, err := key.GetIntegerValue("EnableLUA"); err == nil {
			if val == 0 {
				fmt.Println("[警告] UAC已禁用")
			} else {
				fmt.Println("UAC已启用")
			}
		}
	}
}

// 检查Windows Defender设置
func checkWindowsDefender() {
	fmt.Println("\n=== Windows Defender检查 ===")

	cmd := exec.Command("powershell", "Get-MpPreference")
	output, err := cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("Windows Defender配置:\n%s\n", utf8Output)
		} else {
			fmt.Printf("Windows Defender配置:\n%s\n", string(output))
		}
	}

	cmd = exec.Command("powershell", "Get-MpComputerStatus")
	output, err = cmd.Output()
	if err == nil {
		utf8Output, convErr := gbkToUTF8(output)
		if convErr == nil {
			fmt.Printf("Windows Defender状态:\n%s\n", utf8Output)
		} else {
			fmt.Printf("Windows Defender状态:\n%s\n", string(output))
		}
	}
}