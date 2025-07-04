# 系统应急响应工具集

[![Build and Release](https://github.com/wanglaizi/cross-platform-ir-tools/actions/go.yml/badge.svg)](https://github.com/wanglaizi/cross-platform-ir-tools/actions/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wanglaizi/cross-platform-ir-tools)](https://goreportcard.com/report/github.com/wanglaizi/cross-platform-ir-tools)
[![codecov](https://codecov.io/gh/wanglaizi/cross-platform-ir-tools/branch/main/graph/badge.svg)](https://codecov.io/gh/wanglaizi/cross-platform-ir-tools)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

一个应急响应工具集，主要专为 Windows 系统设计，提供全面的安全分析、事件响应和系统取证功能。同时包含一个 Linux 系统的应急响应脚本作为补充工具。

## 🚀 功能特性

## Windows应急响应工具

本项目主要提供Windows平台的应急响应功能，同时包含一个Linux系统的应急响应脚本作为补充。

### Windows平台功能
### 功能介绍

这是一个用Go语言开发的Windows系统应急响应工具，集成了以下功能：

1. 基础应急响应检查 (-ir)
   - 系统信息收集
   - CPU使用情况
   - 内存使用情况
   - 磁盘信息
   - 网络连接
   - 进程信息
   - 自启动项
   - 计划任务

2. 注册表和文件完整性检查 (-reg)
   - 关键注册表项检查
   - 系统文件完整性验证
   - 可疑文件检测
   - 数字签名验证

3. 内存和进程行为分析 (-mem)
   - 系统内存使用分析
   - 可疑进程识别
   - 进程行为监控
   - DLL加载检查

4. 系统日志分析 (-log)
   - 系统日志分析（启动、关机、错误事件）
   - 安全日志分析（登录、账户管理、策略更改）
   - 应用程序日志分析（程序错误、服务失败）
   - PowerShell日志分析（执行策略、脚本执行）
   - IIS和防火墙日志分析

5. 网络安全分析 (-net)
   - 可疑网络连接检测
   - 网络接口分析
   - 网络流量统计
   - 防火墙规则审计
   - DNS设置检查

6. 系统安全基线检查 (-baseline)
   - 密码策略检查
   - 用户账户审计
   - 系统服务状态
   - 系统补丁检查
   - 审计策略配置
   - 文件系统权限
   - 共享设置检查
   - UAC配置检查
   - Windows Defender状态

### Linux应急响应脚本

项目还包含一个Linux系统的应急响应脚本 (`linux_forensics.sh`)，提供以下功能：

1. 基础系统检查 (-b, --basic)
   - 系统信息收集
   - CPU和内存信息
   - 磁盘和网络信息

2. 内存分析 (-m, --memory)
   - 内存使用详情
   - 内存泄漏检测
   - 进程行为监控

3. 安全检查 (-s, --security)
   - 文件完整性检查
   - 用户安全检查
   - 服务和端口检查

4. 日志分析 (-l, --log)
   - 系统日志分析
   - 安全日志分析
   - 应用日志分析

5. 网络分析 (-n, --network)
   - 网络接口分析
   - 网络连接分析
   - 防火墙配置检查

6. 安全基线检查 (-c, --baseline)
   - 密码策略检查
   - 系统更新检查
   - SSH配置检查

## 📦 安装使用

### 从 Release 下载

1. 访问 [Releases 页面](https://github.com/wanglaizi/cross-platform-ir-tools/releases)
2. 下载适合您系统的版本：
   - `incident_response_windows_amd64.exe`: Windows x64 版本
   - `incident_response_windows_386.exe`: Windows x86 版本

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/wanglaizi/cross-platform-ir-tools.git
cd incident_response

# 安装依赖
go mod download

# 编译
go build -o incident_response .

# Windows 编译
go build -o incident_response.exe .
```

## 使用方法

### Windows工具使用

本工具需要管理员权限才能运行。支持以下命令行参数：

```bash
# 运行所有检查
incident_response.exe -all

# 只运行基础应急响应检查
incident_response.exe -ir

# 只运行注册表和文件完整性检查
incident_response.exe -reg

# 只运行内存和进程行为分析
incident_response.exe -mem

# 只运行系统日志分析
incident_response.exe -log

# 只运行网络安全分析
incident_response.exe -net

# 只运行系统安全基线检查
incident_response.exe -baseline

# 组合使用多个检查项
incident_response.exe -ir -net -baseline

# 禁用报告生成
incident_response.exe -all -report=false
```

### Linux脚本使用

Linux应急响应脚本需要root权限运行。支持以下命令行参数：

```bash
# 显示帮助信息
./linux_forensics.sh -h

# 执行所有检查
./linux_forensics.sh -a

# 执行基础系统检查
./linux_forensics.sh -b

# 执行内存分析
./linux_forensics.sh -m

# 执行安全检查
./linux_forensics.sh -s

# 执行日志分析
./linux_forensics.sh -l

# 执行网络分析
./linux_forensics.sh -n

# 执行安全基线检查
./linux_forensics.sh -c

# 生成HTML报告
./linux_forensics.sh -r

# 指定输出文件
./linux_forensics.sh -a -o custom_report.html
```

## 报告说明

### Windows工具报告

Windows工具会在`reports`目录下生成HTML格式的检查报告，包含：

- 检查时间和系统信息
- 问题统计（严重/警告/信息）
- 详细检查结果
  - 基础系统信息
  - 安全配置状态
  - 发现的问题
  - 处理建议

报告文件名格式：`report_YYYYMMDD_HHMMSS.html`

### Linux脚本报告

Linux脚本可以生成HTML格式的检查报告，包含：

- 检查时间和系统信息
- 安全状态统计
- 基础系统信息
- 安全检查结果

报告文件名格式：`forensics_report_YYYYMMDD_HHMMSS.html`（可通过 -o 参数自定义）

## 输出说明

### Windows工具输出

Windows工具会将检查结果输出到控制台，包括：

- 系统基本信息
- CPU和内存使用情况
- 磁盘信息和使用状况
- 网络连接和端口信息
- 进程列表和资源使用情况
- 自启动项和计划任务
- 注册表关键项检查结果
- 系统文件完整性验证结果
- 可疑文件检测结果
- 内存分析结果
- 进程行为监控结果

### Linux脚本输出

Linux脚本会将检查结果输出到控制台，包括：

- 系统基本信息（主机名、内核版本、运行时间）
- CPU和内存使用情况
- 磁盘和网络信息
- 进程和服务状态
- 安全配置检查结果
- 日志分析结果
- 网络连接和防火墙状态
- 用户和权限检查结果


## 🏗️ 开发指南

### 环境要求

**Windows工具开发：**
- Go 1.20 或更高版本
- Git
- Windows 管理员权限（用于测试）

**Linux脚本使用：**
- Bash shell
- Linux系统（支持systemd的发行版）
- root权限

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/wanglaizi/cross-platform-ir-tools.git
cd incident_response

# 安装依赖
go mod download

# 运行测试
go test -v ./...

# 代码格式化
go fmt ./...

# 代码检查
golangci-lint run
```

### 项目结构

```
.
├── .github/
│   └── workflows/
│       └── go.yml          # GitHub Actions CI/CD 配置
├── docs/
│   └── CI_TROUBLESHOOTING.md # CI/CD 故障排除文档
├── reports/                # 生成的报告目录
├── go.mod                  # Go 模块文件
├── go.sum                  # Go 依赖校验文件
├── main.go                 # 主程序入口
├── main_windows.go         # Windows 特定主程序
├── windows_baseline.go     # Windows 基线检查
├── windows_ir.go           # Windows 事件响应
├── windows_log.go          # Windows 日志分析
├── windows_memory.go       # Windows 内存分析
├── windows_network.go      # Windows 网络分析
├── windows_registry.go     # Windows 注册表检查
├── windows_report.go       # Windows 报告生成
├── linux_forensics.sh      # Linux 应急响应脚本（补充工具）
├── CONTRIBUTING.md         # 贡献指南
├── LICENSE                 # 许可证文件
└── README.md               # 项目说明文档
```

## 🔄 CI/CD 流水线

本项目使用 GitHub Actions 实现完整的 CI/CD 流水线：

### 工作流程

1. **测试** (`test`)
   - 运行单元测试
   - 生成代码覆盖率报告

2. **Windows 平台构建** (`build`)
   - 支持 Windows x64 (amd64)
   - 支持 Windows x86 (386)
   - 生成构建产物

3. **自动发布** (`release`)
   - 当推送标签时自动触发
   - 生成 changelog
   - 创建 GitHub Release
   - 上传 Windows 平台的二进制文件

### 触发条件

- **推送到 main/develop 分支**: 运行完整的 CI 流程
- **Pull Request**: 运行代码检查、测试和构建
- **推送标签 (v*)**: 运行完整流程并自动发布

### 发布新版本

```bash
# 创建并推送标签
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions 将自动：
# 1. 运行所有检查和测试
# 2. 构建多平台二进制文件
# 3. 创建 GitHub Release
# 4. 上传构建产物
```

## 🛡️ 安全考虑

- 所有工具都需要提升权限运行
- 不会收集或传输敏感数据
- 所有操作都在本地执行
- 遵循最小权限原则
- 定期进行安全扫描

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🤝 贡献指南

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `go fmt` 格式化代码
- 通过 `golangci-lint` 检查
- 添加适当的注释和文档
- 编写单元测试

## 📞 支持

如果您遇到问题或有建议，请：

1. 查看 [Issues](https://github.com/wanglaizi/cross-platform-ir-tools/issues)
2. 创建新的 Issue
3. 联系维护者

## 🙏 致谢

感谢以下开源项目：

- [gopsutil](https://github.com/shirou/gopsutil) - 系统信息获取
- [golang.org/x/sys](https://golang.org/x/sys) - 系统调用接口
- [golangci-lint](https://github.com/golangci/golangci-lint) - 代码质量检查

---

**⚠️ 免责声明**: 本工具仅用于合法的系统管理和安全分析目的。使用者需确保遵守相关法律法规。

## ⚠️ 注意事项

### Windows工具注意事项

1. **权限要求**: 必须以管理员权限运行
2. **系统资源**: 部分功能可能会占用较多系统资源，建议在系统负载较低时运行
3. **执行时间**: 文件完整性检查和日志分析可能需要较长时间，请耐心等待
4. **结果确认**: 如果发现异常，建议进一步分析和确认
5. **定期检查**: 建议定期运行安全基线检查，及时发现安全隐患

### Linux脚本注意事项

1. **权限要求**: 必须以root权限运行
2. **系统兼容性**: 主要针对支持systemd的Linux发行版设计
3. **执行时间**: 某些检查项可能需要较长时间，特别是文件完整性检查
4. **网络依赖**: 部分功能需要网络连接（如系统更新检查）
5. **日志文件**: 确保相关日志文件存在且可读

## 📚 主要依赖

- **gopsutil**: 系统信息获取库
- **golang.org/x/sys**: 系统调用接口
- **标准库**: 使用 Go 标准库实现核心功能