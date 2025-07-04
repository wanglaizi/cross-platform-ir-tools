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
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 重要的注册表路径
var criticalRegPaths = []string{
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\RunOnce",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\RunServices",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies",
	"SYSTEM\\CurrentControlSet\\Services",
	"SYSTEM\\CurrentControlSet\\Control\\SafeBoot",
	"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon",
	"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders",
}

// 检查注册表项
func checkRegistry() {
	fmt.Println("=== 注册表检查 ===")

	// 检查HKEY_LOCAL_MACHINE
	fmt.Println("\n检查HKEY_LOCAL_MACHINE:")
	for _, path := range criticalRegPaths {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.READ)
		if err != nil {
			fmt.Printf("无法打开注册表项 %s: %v\n", path, err)
			continue
		}
		defer key.Close()

		// 获取所有值
		values, err := key.ReadValueNames(0)
		if err != nil {
			fmt.Printf("无法读取值: %v\n", err)
			continue
		}

		fmt.Printf("\n[%s]\n", path)
		for _, name := range values {
			val, _, err := key.GetStringValue(name)
			if err == nil {
				fmt.Printf("%s = %s\n", name, val)
			}
		}
	}

	// 检查HKEY_CURRENT_USER
	fmt.Println("\n检查HKEY_CURRENT_USER:")
	for _, path := range criticalRegPaths {
		key, err := registry.OpenKey(registry.CURRENT_USER, path, registry.READ)
		if err != nil {
			continue
		}
		defer key.Close()

		values, err := key.ReadValueNames(0)
		if err != nil {
			continue
		}

		fmt.Printf("\n[%s]\n", path)
		for _, name := range values {
			val, _, err := key.GetStringValue(name)
			if err == nil {
				fmt.Printf("%s = %s\n", name, val)
			}
		}
	}
}

// 将GBK编码转换为UTF-8 (在windows_registry.go中重复定义以避免依赖)
func gbkToUTF8Reg(data []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	utf8Data, err := io.ReadAll(reader)
	if err != nil {
		return string(data), err // 如果转换失败，返回原始数据
	}
	return string(utf8Data), nil
}

// 检查系统文件完整性
func checkSystemFileIntegrity() {
	fmt.Println("\n=== 系统文件完整性检查 ===")

	// 检查系统关键文件
	criticalFiles := []string{
		"C:\\Windows\\System32\\ntoskrnl.exe",
		"C:\\Windows\\System32\\winlogon.exe",
		"C:\\Windows\\System32\\services.exe",
		"C:\\Windows\\System32\\lsass.exe",
		"C:\\Windows\\System32\\svchost.exe",
		"C:\\Windows\\System32\\csrss.exe",
	}

	for _, file := range criticalFiles {
		// 检查文件是否存在
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("警告: 文件不存在 - %s\n", file)
			continue
		}

		// 获取文件数字签名信息
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-AuthenticodeSignature '%s' | Select-Object -Property Status,SignerCertificate", file))
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("无法验证文件签名 %s: %v\n", file, err)
			continue
		}

		fmt.Printf("\n文件: %s\n", file)
		// 转换编码
		utf8Output, convErr := gbkToUTF8Reg(output)
		if convErr == nil {
			fmt.Printf("签名信息:\n%s\n", utf8Output)
		} else {
			fmt.Printf("签名信息:\n%s\n", string(output))
		}

		// 获取文件属性
		fileInfo, err := os.Stat(file)
		if err == nil {
			fmt.Printf("大小: %d 字节\n", fileInfo.Size())
			fmt.Printf("修改时间: %v\n", fileInfo.ModTime())
		}
	}
}

// 检查可疑文件
func checkSuspiciousFiles() {
	fmt.Println("\n=== 可疑文件检查 ===")

	// 检查常见的恶意软件位置
	suspiciousPaths := []string{
		os.Getenv("TEMP"),
		os.Getenv("APPDATA"),
		os.Getenv("LOCALAPPDATA"),
		"C:\\Windows\\Temp",
	}

	// 可疑文件扩展名
	suspiciousExts := []string{
		".exe", ".dll", ".bat", ".cmd", ".ps1", ".vbs", ".js",
	}

	for _, path := range suspiciousPaths {
		fmt.Printf("\n检查目录: %s\n", path)
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				for _, suspiciousExt := range suspiciousExts {
					if ext == suspiciousExt {
						// 检查文件创建时间
						if time.Since(info.ModTime()) < 24*time.Hour {
							fmt.Printf("发现可疑文件: %s\n", path)
							fmt.Printf("大小: %d 字节\n", info.Size())
							fmt.Printf("修改时间: %v\n", info.ModTime())
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("检查目录出错 %s: %v\n", path, err)
		}
	}
}