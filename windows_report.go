//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
	"time"
	"html/template"
	"bytes"
	"path/filepath"
)

// 报告数据结构
type Report struct {
	Timestamp     string
	SystemInfo    string
	CheckResults  []CheckResult
	TotalIssues   int
	CriticalCount int
	WarningCount  int
	InfoCount     int
}

// 检查结果结构
type CheckResult struct {
	Category    string
	Description string
	Severity    string
	Status      string
	Details     string
}

// HTML模板
var reportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Windows系统应急响应报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; border-radius: 5px; }
        .summary { margin: 20px 0; }
        .issue { margin: 10px 0; padding: 10px; border-radius: 5px; }
        .critical { background-color: #ffe6e6; border-left: 5px solid #ff0000; }
        .warning { background-color: #fff3e6; border-left: 5px solid #ff9900; }
        .info { background-color: #e6f3ff; border-left: 5px solid #0066cc; }
        .status-ok { color: green; }
        .status-error { color: red; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Windows系统应急响应报告</h1>
        <p>生成时间: {{.Timestamp}}</p>
    </div>

    <div class="summary">
        <h2>系统信息</h2>
        <pre>{{.SystemInfo}}</pre>
        
        <h2>检查结果统计</h2>
        <p>总问题数: {{.TotalIssues}}</p>
        <p>严重问题: {{.CriticalCount}}</p>
        <p>警告: {{.WarningCount}}</p>
        <p>信息: {{.InfoCount}}</p>
    </div>

    <div class="results">
        <h2>详细检查结果</h2>
        {{range .CheckResults}}
        <div class="issue {{.Severity}}">
            <h3>{{.Category}}</h3>
            <p><strong>描述:</strong> {{.Description}}</p>
            <p><strong>严重程度:</strong> {{.Severity}}</p>
            <p><strong>状态:</strong> <span class="status-{{if eq .Status "正常"}}ok{{else}}error{{end}}">{{.Status}}</span></p>
            {{if .Details}}
            <pre>{{.Details}}</pre>
            {{end}}
        </div>
        {{end}}
    </div>
</body>
</html>
`

// 生成报告
func generateReport(results []CheckResult, sysInfo string) error {
	// 创建报告目录
	reportDir := "reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("创建报告目录失败: %v", err)
	}

	// 统计问题数量
	var criticalCount, warningCount, infoCount int
	for _, result := range results {
		switch result.Severity {
		case "critical":
			criticalCount++
		case "warning":
			warningCount++
		case "info":
			infoCount++
		}
	}

	// 准备报告数据
	report := Report{
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		SystemInfo:    sysInfo,
		CheckResults:  results,
		TotalIssues:   len(results),
		CriticalCount: criticalCount,
		WarningCount:  warningCount,
		InfoCount:     infoCount,
	}

	// 解析模板
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return fmt.Errorf("解析报告模板失败: %v", err)
	}

	// 生成报告内容
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, report); err != nil {
		return fmt.Errorf("生成报告内容失败: %v", err)
	}

	// 保存报告文件
	reportPath := filepath.Join(reportDir, fmt.Sprintf("report_%s.html", time.Now().Format("20060102_150405")))
	if err := os.WriteFile(reportPath, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("保存报告文件失败: %v", err)
	}

	fmt.Printf("报告已生成: %s\n", reportPath)
	return nil
}

// 添加检查结果
func addCheckResult(results *[]CheckResult, category, description, severity, status, details string) {
	*results = append(*results, CheckResult{
		Category:    category,
		Description: description,
		Severity:    severity,
		Status:      status,
		Details:     details,
	})
}