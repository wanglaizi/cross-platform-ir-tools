# 贡献指南

感谢您对系统应急响应工具集项目的关注！我们欢迎各种形式的贡献。

## 🤝 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议：

1. 首先搜索 [现有 Issues](https://github.com/username/incident_response/issues) 确认问题未被报告
2. 创建新的 Issue，请包含：
   - 清晰的问题描述
   - 重现步骤（如果是 bug）
   - 期望的行为
   - 系统环境信息
   - 相关的错误日志或截图

### 提交代码

1. **Fork 仓库**
   ```bash
   git clone https://github.com/your-username/incident_response.git
   cd incident_response
   ```

2. **创建特性分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或者修复分支
   git checkout -b fix/your-fix-name
   ```

3. **进行开发**
   - 遵循项目的代码规范
   - 添加必要的测试
   - 更新相关文档

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   # 或者
   git commit -m "fix: fix your bug description"
   ```

5. **推送分支**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 提供清晰的 PR 描述
   - 关联相关的 Issues
   - 确保 CI 检查通过

## 📋 代码规范

### Go 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查代码
- 通过 `golangci-lint` 检查
- 函数和方法需要添加注释
- 导出的类型、变量、常量需要添加注释

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型 (type):**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式化（不影响代码逻辑）
- `refactor`: 代码重构
- `test`: 添加或修改测试
- `chore`: 构建过程或辅助工具的变动
- `perf`: 性能优化
- `ci`: CI/CD 相关更改

**示例:**
```
feat(windows): add registry integrity check

fix(linux): resolve memory leak in process monitoring

docs: update installation instructions

test(network): add unit tests for network analysis
```

## 🧪 测试

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行测试并生成覆盖率报告
go test -v -race -coverprofile=coverage.out ./...

# 查看覆盖率报告
go tool cover -html=coverage.out
```

### 编写测试

- 为新功能添加单元测试
- 测试文件命名为 `*_test.go`
- 测试函数命名为 `TestXxx`
- 使用表驱动测试处理多个测试用例
- 模拟外部依赖（文件系统、网络等）

### 测试示例

```go
func TestProcessAnalysis(t *testing.T) {
    tests := []struct {
        name     string
        input    ProcessInfo
        expected bool
    }{
        {
            name: "normal process",
            input: ProcessInfo{Name: "notepad.exe", CPU: 1.0},
            expected: false,
        },
        {
            name: "suspicious process",
            input: ProcessInfo{Name: "malware.exe", CPU: 90.0},
            expected: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsSuspiciousProcess(tt.input)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## 🔍 代码审查

### 审查清单

- [ ] 代码遵循项目规范
- [ ] 添加了必要的测试
- [ ] 测试覆盖率满足要求
- [ ] 文档已更新
- [ ] 没有引入安全漏洞
- [ ] 性能影响可接受
- [ ] 向后兼容性

### 审查重点

1. **安全性**: 检查是否有安全漏洞，特别是权限提升、文件操作等
2. **性能**: 避免不必要的资源消耗
3. **错误处理**: 确保适当的错误处理和日志记录
4. **跨平台兼容性**: 确保代码在不同平台上正常工作
5. **代码可读性**: 代码应该清晰易懂

## 🚀 发布流程

### 版本号规范

使用 [Semantic Versioning](https://semver.org/)：

- `MAJOR.MINOR.PATCH` (例如: 1.2.3)
- `MAJOR`: 不兼容的 API 更改
- `MINOR`: 向后兼容的功能添加
- `PATCH`: 向后兼容的 bug 修复

### 发布步骤

1. 更新版本号和 CHANGELOG
2. 创建 release 分支
3. 运行完整测试套件
4. 创建并推送标签
5. GitHub Actions 自动构建和发布

```bash
# 创建标签
git tag v1.2.3
git push origin v1.2.3
```

## 📚 开发环境

### 必需工具

- Go 1.20+
- Git
- golangci-lint
- 管理员/root 权限（用于测试）

### 推荐工具

- VS Code 或 GoLand
- Go 扩展/插件
- Git 客户端

### 环境设置

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装 pre-commit hooks（可选）
go install github.com/pre-commit/pre-commit@latest
```

## 🆘 获取帮助

如果您在贡献过程中遇到问题：

1. 查看项目文档和 Issues
2. 在 Discussions 中提问
3. 联系维护者

## 📄 许可证

通过贡献代码，您同意您的贡献将在 MIT 许可证下发布。

---

再次感谢您的贡献！🎉