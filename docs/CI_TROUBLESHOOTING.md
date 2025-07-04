# GitHub Actions CI/CD 故障排除指南

本文档提供了系统应急响应工具集项目中 GitHub Actions 工作流常见问题的解决方案。

## 🔧 常见问题及解决方案

### 1. golangci-lint 检查失败

**问题描述**: golangci-lint 报告代码格式或质量问题

**解决方案**:
```bash
# 本地运行 golangci-lint 检查
golangci-lint run

# 自动修复可修复的问题
golangci-lint run --fix

# 格式化代码
go fmt ./...

# 整理导入
go mod tidy
```

**常见错误**:
- `fmt.Println arg list ends with redundant newline`: 移除 fmt.Println 中多余的换行符
- `unused import`: 移除未使用的导入包
- `ineffassign`: 修复无效的赋值

### 2. 测试失败

**问题描述**: 测试步骤失败或没有测试文件

**解决方案**:
```bash
# 运行测试
go test -v ./...

# 运行测试并生成覆盖率
go test -v -coverprofile=coverage.out ./...

# 查看覆盖率报告
go tool cover -html=coverage.out
```

**注意事项**:
- 确保至少有一个 `*_test.go` 文件
- 测试函数必须以 `Test` 开头
- 避免在测试中使用需要管理员权限的功能

### 3. 构建失败

**问题描述**: 多平台构建失败

**解决方案**:
```bash
# 检查不同平台的构建
GOOS=linux GOARCH=amd64 go build .
GOOS=windows GOARCH=amd64 go build .
GOOS=darwin GOARCH=amd64 go build .

# 检查构建标签
go build -tags linux .
go build -tags windows .
```

**常见问题**:
- 平台特定的代码没有正确使用构建标签
- 依赖包不支持某些平台
- CGO 相关问题

### 4. 依赖问题

**问题描述**: 依赖下载或版本冲突

**解决方案**:
```bash
# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download

# 更新依赖
go get -u ./...

# 整理依赖
go mod tidy
```

### 5. 安全扫描问题

**问题描述**: Gosec 安全扫描报告问题

**解决方案**:
```bash
# 本地运行 gosec
gosec ./...

# 生成详细报告
gosec -fmt json -out gosec-report.json ./...
```

**常见安全问题**:
- G204: 使用变量启动子进程
- G301/G302: 文件权限问题
- G304: 文件包含漏洞

## 🛠️ 本地开发环境设置

### 安装必要工具

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 gosec
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# 验证安装
golangci-lint version
gosec -version
```

### 预提交检查脚本

创建 `scripts/pre-commit.sh`:
```bash
#!/bin/bash
set -e

echo "运行代码格式化..."
go fmt ./...

echo "整理依赖..."
go mod tidy

echo "运行代码检查..."
golangci-lint run

echo "运行测试..."
go test -v ./...

echo "运行安全扫描..."
gosec ./...

echo "所有检查通过！"
```

## 📋 CI/CD 配置优化建议

### 1. 使用缓存加速构建

```yaml
- name: Cache Go modules
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### 2. 并行执行任务

```yaml
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.19, 1.20, 1.21]
        os: [ubuntu-latest, windows-latest, macos-latest]
```

### 3. 条件执行

```yaml
- name: Run security scan
  if: github.event_name == 'push' && github.ref == 'refs/heads/main'
  run: gosec ./...
```

### 4. 错误处理

```yaml
- name: Run tests
  run: go test -v ./...
  continue-on-error: true

- name: Upload test results
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: test-results
    path: test-results.xml
```

## 🔍 调试技巧

### 1. 启用调试日志

```yaml
env:
  ACTIONS_STEP_DEBUG: true
  ACTIONS_RUNNER_DEBUG: true
```

### 2. 使用 tmate 进行远程调试

```yaml
- name: Setup tmate session
  if: failure()
  uses: mxschmitt/action-tmate@v3
```

### 3. 保存构建产物

```yaml
- name: Upload build artifacts
  if: always()
  uses: actions/upload-artifact@v3
  with:
    name: build-logs
    path: |
      *.log
      coverage.out
      gosec.sarif
```

## 📚 参考资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [golangci-lint 配置](https://golangci-lint.run/usage/configuration/)
- [Go 测试最佳实践](https://golang.org/doc/tutorial/add-a-test)
- [Gosec 安全规则](https://securecodewarrior.github.io/gosec/)

## 🆘 获取帮助

如果遇到无法解决的问题：

1. 检查 [GitHub Actions 日志](https://github.com/username/incident_response/actions)
2. 查看 [项目 Issues](https://github.com/username/incident_response/issues)
3. 参考 [Go 官方文档](https://golang.org/doc/)
4. 联系项目维护者

---

**提示**: 定期更新 CI/CD 配置和工具版本，以获得最佳的构建体验和安全性。