package checks

import (
	"fmt"
	"strconv"
	"strings"
	"inspector/sshutil"
	"time"
)

// Status 检查状态
const (
	StatusOK   = "OK"
	StatusWARN = "WARN"
	StatusFAIL = "FAIL"
	StatusUNCON = "UNCON"
)

// Item 单个检查项结果
type Item struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// HostResult 单台主机巡检结果
type HostResult struct {
	Host  string `json:"host"`
	CheckTime time.Time `json: "checktime"`
	Items []Item `json:"items"`
	Status string `json: "status"`
	Error string `json:"error,omitempty"`
}

// ---------------------------------------------------------------------------
// 阈值常量
// ---------------------------------------------------------------------------
const (
	CPUWarnThreshold    = 80.0 // CPU 使用率 > 80% 警告
	CPUCalcThreshold    = 90.0 // CPU 使用率 > 90% 失败
	MemWarnThreshold    = 85.0 // 内存使用率 > 85% 警告
	MemFailThreshold    = 95.0 // 内存使用率 > 95% 失败
	DiskWarnThreshold   = 85.0 // 磁盘使用率 > 85% 警告
	DiskFailThreshold   = 95.0 // 磁盘使用率 > 95% 失败
)

// ---------------------------------------------------------------------------
// 检查函数类型
// ---------------------------------------------------------------------------

// CheckFunc 检查函数签名
type CheckFunc func(client *sshutil.Client) Item

// AllChecks 所有检查项（按顺序）
func AllChecks() []struct {
	Name string
	Fn   CheckFunc
} {
	return []struct {
		Name string
		Fn   CheckFunc
	}{
		{"CPU", CheckCPU},
		{"Memory", CheckMemory},
		{"Disk", CheckDisk},
		{"Docker", CheckDocker},
		{"Nginx", CheckNginx},
		{"MySQL", CheckMySQL},
		{"Network", CheckNetwork},
		{"DNS", CheckDNS},
		{"TimeSync", CheckTimeSync},
	}
}

// ---------------------------------------------------------------------------
// CPU 检查
// 使用 /proc/stat 或 top -bn1 获取 CPU 使用率
// ---------------------------------------------------------------------------
func CheckCPU(client *sshutil.Client) Item {
	// 读取 /proc/stat 计算 CPU 使用率
	out, err := client.Run("cat /proc/stat | grep '^cpu '")
	if err != nil {
		return Item{Name: "CPU", Value: "N/A", Status: StatusFAIL, Detail: err.Error()}
	}
	// 简单方法：用 top 一次性获取
	out2, err := client.Run("top -bn1 | grep 'Cpu(s)' | head -1")
	if err == nil && out2 != "" {
		return parseCPUFromTop(out2)
	}

	// 回退：用 uptime 获取负载
	out3, err := client.Run("uptime")
	if err != nil {
		return Item{Name: "CPU", Value: "N/A", Status: StatusFAIL, Detail: out}
	}
	return parseLoadFromUptime(out3)
}

func parseCPUFromTop(output string) Item {
	// 格式: %Cpu(s):  5.2 us,  2.1 sy,  0.0 ni, 92.0 id,  0.5 wa, ...
	// 或者: Cpu(s):  5.2% us,  2.1% sy, ...
	output = strings.TrimSpace(output)

	// 提取 idle 百分比
	var idle float64
	if strings.Contains(output, "id") {
		// 尝试多种格式
		parts := strings.Fields(output)
		for i, p := range parts {
			// 查找 "id" 或 "id,"
			pClean := strings.TrimRight(p, ",")
			if pClean == "id" && i > 0 {
				valStr := strings.TrimRight(parts[i-1], ",")
				if v, err := strconv.ParseFloat(valStr, 64); err == nil {
					idle = v
					break
				}
			}
		}
	}

	if idle == 0 && len(output) > 0 {
		// 也许格式不同，尝试其他解析
		// 简单提取：找到 idle 前面的数字
		idx := strings.Index(output, "id")
		if idx > 0 {
			before := strings.TrimSpace(output[:idx])
			fields := strings.Fields(before)
			if len(fields) > 0 {
				valStr := strings.TrimRight(fields[len(fields)-1], ",")
				if v, err := strconv.ParseFloat(valStr, 64); err == nil && idle ==0 {
					idle = v
					status := StatusOK
					return Item{Name: "CPU", Value: fmt.Sprintf("%d%%",idle),Status: status}
				}
			}
		}
	}

	if idle == 0 {
		return Item{Name: "CPU", Value: output, Status: StatusWARN, Detail: "无法解析 CPU 使用率"}
	}

	usage := 100.0 - idle
	usageStr := fmt.Sprintf("%.0f%%", usage)
	status := StatusOK
	if usage > CPUCalcThreshold {
		status = StatusFAIL
	} else if usage > CPUWarnThreshold {
		status = StatusWARN
	}
	return Item{Name: "CPU", Value: usageStr, Status: status}
}

func parseLoadFromUptime(output string) Item {
	output = strings.TrimSpace(output)
	// 格式: 10:30:00 up 10 days,  2:30,  5 users,  load average: 0.52, 0.30, 0.25
	idx := strings.LastIndex(output, "load average:")
	if idx < 0 {
		return Item{Name: "CPU", Value: "N/A", Status: StatusWARN, Detail: "无法解析 uptime: " + output}
	}
	loadPart := strings.TrimSpace(output[idx+len("load average:"):])
	parts := strings.Split(loadPart, ",")
	if len(parts) < 3 {
		return Item{Name: "CPU", Value: loadPart, Status: StatusOK}
	}
	// 取 1 分钟负载
	load1 := strings.TrimSpace(parts[0])
	val, err := strconv.ParseFloat(load1, 64)
	if err != nil {
		return Item{Name: "CPU", Value: "load: " + loadPart, Status: StatusOK}
	}
	status := StatusOK
	if val > 8.0 {
		status = StatusFAIL
	} else if val > 4.0 {
		status = StatusWARN
	}
	return Item{Name: "CPU", Value: fmt.Sprintf("load %.2f", val), Status: status}
}

// ---------------------------------------------------------------------------
// Memory 检查
// ---------------------------------------------------------------------------
func CheckMemory(client *sshutil.Client) Item {
	out, err := client.Run("free -m | grep 'Mem:'")
	if err != nil {
		return Item{Name: "Memory", Value: "N/A", Status: StatusFAIL, Detail: err.Error()}
	}
	// 格式: Mem:           7975        1842        4502         278        1630        5574
	fields := strings.Fields(out)
	if len(fields) < 3 {
		return Item{Name: "Memory", Value: out, Status: StatusWARN, Detail: "无法解析"}
	}
	total, err1 := strconv.ParseFloat(fields[1], 64)
	used, err2 := strconv.ParseFloat(fields[2], 64)
	if err1 != nil || err2 != nil || total == 0 {
		return Item{Name: "Memory", Value: out, Status: StatusWARN, Detail: "解析数值失败"}
	}
	usage := (used / total) * 100.0
	usageStr := fmt.Sprintf("%.0f%% (%.0fM/%.0fM)", usage, used, total)
	status := StatusOK
	if usage > MemFailThreshold {
		status = StatusFAIL
	} else if usage > MemWarnThreshold {
		status = StatusWARN
	}
	return Item{Name: "Memory", Value: usageStr, Status: status}
}

// ---------------------------------------------------------------------------
// Disk 检查 - 检查所有挂载点，报告最差情况
// ---------------------------------------------------------------------------
func CheckDisk(client *sshutil.Client) Item {
	out, err := client.Run("df  --exclude-type=tmpfs --exclude-type=devtmpfs | grep -v '^Filesystem'")
	if err != nil {
		return Item{Name: "Disk", Value: "N/A", Status: StatusFAIL, Detail: err.Error()}
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return Item{Name: "Disk", Value: "N/A", Status: StatusWARN}
	}

	tatol := 0.0
	used := 0.0
	//worstMount := ""
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// df -h: Filesystem Size Used Avail Use% Mounted on
		tatolStr := strings.TrimSpace(fields[1])
		tatol_, err := strconv.ParseFloat(tatolStr, 64)
		if err != nil {
			continue
		}
		usedStr := strings.TrimSpace(fields[2])
		used_, err := strconv.ParseFloat(usedStr, 64)
		if err != nil {
			continue
		}
		tatol = tatol + tatol_
		used = used + used_
	}

	if tatol == 0 {
		return Item{Name: "Disk", Value: "N/A", Status: StatusWARN}
	}

	status := StatusOK
	usage := (used/tatol)*100

	if usage > DiskFailThreshold {
		status = StatusFAIL
	} else if usage > DiskWarnThreshold {
		status = StatusWARN
	}

	detail := fmt.Sprintf("系统磁盘总使用率: %.0f%%", usage)
	return Item{Name: "Disk", Value: fmt.Sprintf("%.0f%%", usage), Status: status, Detail: detail}
}

// ---------------------------------------------------------------------------
// Docker 检查
// ---------------------------------------------------------------------------
func CheckDocker(client *sshutil.Client) Item {
	// 先检查 docker 是否安装
	_, err := client.Run("which docker")
	if err != nil {
		return Item{Name: "Docker", Value: "未安装", Status: StatusOK, Detail: "docker 未安装，跳过"}
	}

	// 检查 docker 服务状态
	out, err := client.Run("systemctl is-active docker 2>/dev/null || echo 'inactive'")
	if err != nil {
		// 忽略错误，继续检查
	}
	out = strings.TrimSpace(out)

	if out != "active" {
		// 尝试检查 docker ps 是否可用
		psOut, psErr := client.Run("docker ps -q 2>/dev/null | wc -l")
		if psErr != nil {
			return Item{Name: "Docker", Value: "未运行", Status: StatusFAIL, Detail: "Docker daemon 未运行"}
		}
		containerCount := strings.TrimSpace(psOut)
		return Item{Name: "Docker", Value: fmt.Sprintf("运行中(%s容器)", containerCount), Status: StatusWARN, Detail: "Docker 服务状态非 active"}
	}

	// 检查运行容器数量
	psOut, _ := client.Run("docker ps -q 2>/dev/null | wc -l")
	containerCount := strings.TrimSpace(psOut)

	return Item{Name: "Docker", Value: fmt.Sprintf("运行中(%s容器)", containerCount), Status: StatusOK}
}

// ---------------------------------------------------------------------------
// Nginx 检查
// ---------------------------------------------------------------------------
func CheckNginx(client *sshutil.Client) Item {
	// 检查 nginx 是否安装
	_, err := client.Run("which nginx")
	if err != nil {
		return Item{Name: "Nginx", Value: "未安装", Status: StatusOK, Detail: "nginx 未安装，跳过"}
	}

	out, err := client.Run("systemctl is-active nginx 2>/dev/null || echo 'unknown'")
	out = strings.TrimSpace(out)

	if out == "active" {
		return Item{Name: "Nginx", Value: "运行中", Status: StatusOK}
	}
	if out == "inactive" {
		return Item{Name: "Nginx", Value: "已停止", Status: StatusFAIL}
	}
	// 未知状态，尝试检查进程
	_, err = client.Run("pgrep -x nginx")
	if err == nil {
		return Item{Name: "Nginx", Value: "运行中(进程)", Status: StatusOK}
	}
	return Item{Name: "Nginx", Value: out, Status: StatusWARN}
}

// ---------------------------------------------------------------------------
// MySQL 检查
// ---------------------------------------------------------------------------
func CheckMySQL(client *sshutil.Client) Item {
	// 支持 mysqld 和 mariadb
	for _, svc := range []string{"mysqld", "mysql", "mariadb"} {
		out, err := client.Run(fmt.Sprintf("systemctl is-active %s 2>/dev/null || echo ''", svc))
		if err == nil {
			out = strings.TrimSpace(out)
			if out == "active" {
				return Item{Name: "MySQL", Value: "运行中", Status: StatusOK, Detail: fmt.Sprintf("服务: %s", svc)}
			}
			if out == "inactive" {
				return Item{Name: "MySQL", Value: "已停止", Status: StatusFAIL, Detail: fmt.Sprintf("服务 %s 未运行", svc)}
			}
		}
	}

	// 尝试检查进程
	_, err := client.Run("pgrep -x 'mysqld|mariadbd'")
	if err == nil {
		return Item{Name: "MySQL", Value: "运行中(进程)", Status: StatusOK}
	}

	// 检查是否安装
	_, err = client.Run("which mysqld 2>/dev/null || which mariadbd 2>/dev/null")
	if err != nil {
		return Item{Name: "MySQL", Value: "未安装", Status: StatusOK, Detail: "MySQL/MariaDB 未安装"}
	}

	return Item{Name: "MySQL", Value: "未知", Status: StatusWARN}
}

// ---------------------------------------------------------------------------
// 网络检查 - ping 网关或外网
// ---------------------------------------------------------------------------
func CheckNetwork(client *sshutil.Client) Item {
	// 快速检查：ping 8.8.8.8 一次，超时 2 秒
	out, err := client.Run("ping -c 1 -W 2 8.8.8.8 2>&1")
	if err != nil {
		// 尝试 ping 114.114.114.114
		out2, err2 := client.Run("ping -c 1 -W 2 114.114.114.114 2>&1")
		if err2 != nil {
			return Item{Name: "Network", Value: "不可达", Status: StatusFAIL, Detail: "无法 ping 通外部网络"}
		}
		_ = out2
	}
	_ = out

	// 提取延迟
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			// 格式: 64 bytes from 8.8.8.8: icmp_seq=1 ttl=117 time=12.3 ms
			idx := strings.Index(line, "time=")
			if idx >= 0 {
				timePart := line[idx+5:]
				timePart = strings.TrimSpace(timePart)
				timePart = strings.Fields(timePart)[0]
				return Item{Name: "Network", Value: fmt.Sprintf("可达(%sms)", timePart), Status: StatusOK}
			}
		}
	}
	return Item{Name: "Network", Value: "可达", Status: StatusOK}
}

// ---------------------------------------------------------------------------
// DNS 检查
// ---------------------------------------------------------------------------
func CheckDNS(client *sshutil.Client) Item {
	// 使用 nslookup 或 dig
	out, err := client.Run("nslookup baidu.com 2>&1 || dig +short baidu.com 2>&1 || echo 'FAIL'")
	if err != nil || strings.Contains(out, "FAIL") {
		return Item{Name: "DNS", Value: "失败", Status: StatusFAIL, Detail: out}
	}

	out = strings.TrimSpace(out)
	if strings.Contains(out, "Address:") || strings.Contains(out, "A") || (len(out) > 0 && !strings.Contains(out, "server can't") && !strings.Contains(out, "NXDOMAIN")) {
		return Item{Name: "DNS", Value: "正常", Status: StatusOK}
	}

	// 简单判断：有 IP 地址就认为 OK
	if strings.Contains(out, ".") && len(out) < 50 {
		return Item{Name: "DNS", Value: "正常", Status: StatusOK}
	}

	return Item{Name: "DNS", Value: "异常", Status: StatusWARN, Detail: out}
}

// ---------------------------------------------------------------------------
// 时间同步检查
// ---------------------------------------------------------------------------
func CheckTimeSync(client *sshutil.Client) Item {
	// 先检查 chronyd
	out, err := client.Run("chronyc tracking 2>/dev/null || echo ''")
	if err == nil && strings.Contains(out, "Reference ID") {
		// 提取同步状态
		lines := strings.Split(out, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Leap status") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[1]) == "Normal" {
					return Item{Name: "TimeSync", Value: "chronyd 正常", Status: StatusOK}
				}
			}
		}
		return Item{Name: "TimeSync", Value: "chronyd 运行", Status: StatusOK}
	}

	// 检查 ntpd
	out2, err := client.Run("systemctl is-active ntpd 2>/dev/null || echo ''")
	out2 = strings.TrimSpace(out2)
	if out2 == "active" {
		return Item{Name: "TimeSync", Value: "ntpd 运行", Status: StatusOK}
	}

	// 检查 timesyncd
	out3, _ := client.Run("systemctl is-active systemd-timesyncd 2>/dev/null || echo ''")
	out3 = strings.TrimSpace(out3)
	if out3 == "active" {
		return Item{Name: "TimeSync", Value: "timesyncd 运行", Status: StatusOK}
	}

	// 检查 timedatectl
	out4, err := client.Run("timedatectl show -p NTP -p NTPSynchronized 2>/dev/null || echo ''")
	if err == nil && strings.Contains(out4, "NTPSynchronized=yes") {
		return Item{Name: "TimeSync", Value: "已同步", Status: StatusOK}
	}

	return Item{Name: "TimeSync", Value: "未同步", Status: StatusWARN, Detail: "未检测到时间同步服务"}
}

// ---------------------------------------------------------------------------
// 汇总结果状态
// ---------------------------------------------------------------------------

// OverallStatus 根据所有检查项计算总体状态
func OverallStatus(items []Item) string {
	if len(items) == 0 {
		return "FAIL" // 无检查项视为失败
	}
	hasFail := false
	hasWarn := false
	for _, item := range items {
		switch item.Status {
		case StatusFAIL:
			hasFail = true
		case StatusWARN:
			hasWarn = true
		}
	}
	if hasFail {
		return "FAIL"
	}
	if hasWarn {
		return "WARN"
	}
	return "PASS"
}
