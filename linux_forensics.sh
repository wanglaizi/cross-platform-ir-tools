#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 帮助信息
show_help() {
    echo -e "${GREEN}Linux系统应急响应工具集 v2.0${NC}"
    echo -e "用法: $0 [选项]"
    echo -e "选项:"
    echo -e "  -h, --help\t\t显示帮助信息"
    echo -e "  -a, --all\t\t执行所有检查"
    echo -e "  -b, --basic\t\t执行基础系统检查"
    echo -e "  -m, --memory\t\t执行内存分析"
    echo -e "  -s, --security\t\t执行安全检查"
    echo -e "  -l, --log\t\t执行日志分析"
    echo -e "  -n, --network\t\t执行网络分析"
    echo -e "  -c, --baseline\t执行安全基线检查"
    echo -e "  -r, --report\t\t生成HTML报告"
    echo -e "  -o, --output <file>\t指定输出文件"
    exit 0
}

# 检查是否以root权限运行
check_root() {
    if [ "$(id -u)" != "0" ]; then
        echo -e "${RED}错误: 此脚本需要root权限运行${NC}"
        exit 1
    fi
}

# 基础系统检查模块
run_basic_check() {
    echo -e "\n${BLUE}[+] 执行基础系统检查${NC}"
    
    # 系统信息
    echo -e "\n${GREEN}=== 系统信息 ===${NC}"
    echo -e "${YELLOW}主机名:${NC} $(hostname)"
    echo -e "${YELLOW}内核版本:${NC} $(uname -r)"
    echo -e "${YELLOW}操作系统:${NC} $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
    echo -e "${YELLOW}系统时间:${NC} $(date)"
    echo -e "${YELLOW}运行时间:${NC} $(uptime)"
    
    # CPU信息
    echo -e "\n${GREEN}=== CPU信息 ===${NC}"
    echo -e "${YELLOW}CPU型号:${NC}"
    cat /proc/cpuinfo | grep 'model name' | uniq
    echo -e "${YELLOW}CPU核心数:${NC} $(nproc)"
    echo -e "${YELLOW}CPU使用率:${NC}"
    top -bn1 | grep '%Cpu'
    
    # 内存信息
    echo -e "\n${GREEN}=== 内存信息 ===${NC}"
    free -h
    
    # 磁盘信息
    echo -e "\n${GREEN}=== 磁盘信息 ===${NC}"
    df -h
    
    # 网络信息
    echo -e "\n${GREEN}=== 网络信息 ===${NC}"
    ip addr
    netstat -tuln
}

# 内存分析模块
run_memory_analysis() {
    echo -e "\n${BLUE}[+] 执行内存分析${NC}"
    
    # 内存使用详情
    echo -e "\n${GREEN}=== 内存使用详情 ===${NC}"
    vmstat 1 5
    echo -e "\n${YELLOW}内存使用TOP 10进程:${NC}"
    ps aux | sort -rn -k4 | head -n 10
    
    # 内存泄漏检测
    echo -e "\n${GREEN}=== 内存泄漏检测 ===${NC}"
    echo -e "${YELLOW}持续增长的进程:${NC}"
    ps aux | awk '$6>50000' | sort -rn -k6
    
    # 进程行为监控
    echo -e "\n${GREEN}=== 进程行为监控 ===${NC}"
    echo -e "${YELLOW}异常进程:${NC}"
    ps aux | awk '$3>50.0 || $4>50.0'
    echo -e "\n${YELLOW}僵尸进程:${NC}"
    ps aux | awk '$8~"Z"'
}

# 安全检查模块
run_security_check() {
    echo -e "\n${BLUE}[+] 执行安全检查${NC}"
    
    # 文件完整性检查
    echo -e "\n${GREEN}=== 文件完整性检查 ===${NC}"
    echo -e "${YELLOW}SUID文件:${NC}"
    find / -type f -perm -4000 2>/dev/null
    
    # 用户安全
    echo -e "\n${GREEN}=== 用户安全检查 ===${NC}"
    echo -e "${YELLOW}特权用户:${NC}"
    awk -F: '$3==0' /etc/passwd
    echo -e "\n${YELLOW}最近用户活动:${NC}"
    last | head -n 5
    
    # 服务检查
    echo -e "\n${GREEN}=== 服务检查 ===${NC}"
    echo -e "${YELLOW}运行的服务:${NC}"
    systemctl list-units --type=service --state=running
    
    # 开放端口
    echo -e "\n${GREEN}=== 端口检查 ===${NC}"
    echo -e "${YELLOW}开放的端口:${NC}"
    netstat -tulnp
}

# 日志分析模块
run_log_analysis() {
    echo -e "\n${BLUE}[+] 执行日志分析${NC}"
    
    # 系统日志
    echo -e "\n${GREEN}=== 系统日志分析 ===${NC}"
    echo -e "${YELLOW}系统错误:${NC}"
    journalctl -p 3 -xb | tail -n 10
    
    # 安全日志
    echo -e "\n${GREEN}=== 安全日志分析 ===${NC}"
    echo -e "${YELLOW}认证失败:${NC}"
    grep -i "failed" /var/log/auth.log 2>/dev/null | tail -n 10
    
    # 应用日志
    echo -e "\n${GREEN}=== 应用日志分析 ===${NC}"
    if [ -f "/var/log/apache2/error.log" ]; then
        echo -e "${YELLOW}Apache错误日志:${NC}"
        tail -n 10 /var/log/apache2/error.log
    fi
    if [ -f "/var/log/nginx/error.log" ]; then
        echo -e "${YELLOW}Nginx错误日志:${NC}"
        tail -n 10 /var/log/nginx/error.log
    fi
}

# 网络分析模块
run_network_analysis() {
    echo -e "\n${BLUE}[+] 执行网络分析${NC}"
    
    # 网络接口
    echo -e "\n${GREEN}=== 网络接口分析 ===${NC}"
    ip -s link
    
    # 网络连接
    echo -e "\n${GREEN}=== 网络连接分析 ===${NC}"
    echo -e "${YELLOW}活动连接:${NC}"
    netstat -antp | grep ESTABLISHED
    
    # 防火墙配置
    echo -e "\n${GREEN}=== 防火墙配置 ===${NC}"
    if command -v iptables >/dev/null 2>&1; then
        echo -e "${YELLOW}iptables规则:${NC}"
        iptables -L -n
    fi
    if command -v ufw >/dev/null 2>&1; then
        echo -e "${YELLOW}UFW状态:${NC}"
        ufw status
    fi
}

# 安全基线检查模块
run_baseline_check() {
    echo -e "\n${BLUE}[+] 执行安全基线检查${NC}"
    
    # 密码策略
    echo -e "\n${GREEN}=== 密码策略检查 ===${NC}"
    grep -i "^password" /etc/security/pwquality.conf 2>/dev/null
    
    # 系统更新
    echo -e "\n${GREEN}=== 系统更新检查 ===${NC}"
    if command -v apt >/dev/null 2>&1; then
        apt list --upgradable 2>/dev/null
    elif command -v yum >/dev/null 2>&1; then
        yum check-update
    fi
    
    # SSH配置
    echo -e "\n${GREEN}=== SSH配置检查 ===${NC}"
    grep -i "^PermitRootLogin\|^PasswordAuthentication" /etc/ssh/sshd_config
}

# 生成HTML报告
generate_report() {
    local report_file="$1"
    echo -e "\n${BLUE}[+] 生成HTML报告${NC}"
    
    # 创建HTML报告
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Linux系统安全检查报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #2c3e50; }
        h2 { color: #34495e; margin-top: 20px; }
        .section { margin: 10px 0; padding: 10px; border: 1px solid #bdc3c7; }
        .warning { color: #c0392b; }
        .info { color: #2980b9; }
    </style>
</head>
<body>
    <h1>Linux系统安全检查报告</h1>
    <div class="section">
        <h2>检查时间</h2>
        <p>$(date)</p>
    </div>
    <div class="section">
        <h2>系统信息</h2>
        <p>主机名: $(hostname)</p>
        <p>内核版本: $(uname -r)</p>
        <p>操作系统: $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)</p>
    </div>
    <div class="section">
        <h2>安全状态</h2>
        <p>特权用户数量: $(awk -F: '\$3==0' /etc/passwd | wc -l)</p>
        <p>开放端口数量: $(netstat -tuln | grep LISTEN | wc -l)</p>
        <p>SUID文件数量: $(find / -type f -perm -4000 2>/dev/null | wc -l)</p>
    </div>
</body>
</html>
EOF
    
    echo -e "${GREEN}报告已生成: $report_file${NC}"
}

# 主函数
main() {
    local output_file=""
    local do_all=false
    local do_basic=false
    local do_memory=false
    local do_security=false
    local do_log=false
    local do_network=false
    local do_baseline=false
    local do_report=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                ;;
            -a|--all)
                do_all=true
                shift
                ;;
            -b|--basic)
                do_basic=true
                shift
                ;;
            -m|--memory)
                do_memory=true
                shift
                ;;
            -s|--security)
                do_security=true
                shift
                ;;
            -l|--log)
                do_log=true
                shift
                ;;
            -n|--network)
                do_network=true
                shift
                ;;
            -c|--baseline)
                do_baseline=true
                shift
                ;;
            -r|--report)
                do_report=true
                shift
                ;;
            -o|--output)
                output_file="$2"
                shift 2
                ;;
            *)
                echo -e "${RED}错误: 未知选项 $1${NC}"
                show_help
                ;;
        esac
    done
    
    # 如果没有指定任何选项，显示帮助
    if [[ "$do_all" == "false" && "$do_basic" == "false" && "$do_memory" == "false" && \
          "$do_security" == "false" && "$do_log" == "false" && "$do_network" == "false" && \
          "$do_baseline" == "false" && "$do_report" == "false" ]]; then
        show_help
    fi
    
    # 检查root权限
    check_root
    
    # 设置默认输出文件
    if [[ -z "$output_file" ]]; then
        output_file="forensics_report_$(date +%Y%m%d_%H%M%S).html"
    fi
    
    # 执行选定的检查项
    if [[ "$do_all" == "true" || "$do_basic" == "true" ]]; then
        run_basic_check
    fi
    
    if [[ "$do_all" == "true" || "$do_memory" == "true" ]]; then
        run_memory_analysis
    fi
    
    if [[ "$do_all" == "true" || "$do_security" == "true" ]]; then
        run_security_check
    fi
    
    if [[ "$do_all" == "true" || "$do_log" == "true" ]]; then
        run_log_analysis
    fi
    
    if [[ "$do_all" == "true" || "$do_network" == "true" ]]; then
        run_network_analysis
    fi
    
    if [[ "$do_all" == "true" || "$do_baseline" == "true" ]]; then
        run_baseline_check
    fi
    
    if [[ "$do_all" == "true" || "$do_report" == "true" ]]; then
        generate_report "$output_file"
    fi
}

# 执行主函数
main "$@"