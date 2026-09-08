package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"inspector/checks"
	"inspector/config"
	mysql "inspector/my_mysql"
	"inspector/reporter"
	"inspector/sshutil"
)

func main() {
	// 配置文件路径
	fmt.Println("============================================")
	fmt.Println("  🔍 go-inspector — 多主机批量巡检工具")
	fmt.Println("============================================")
	fmt.Println("📄 读取配置: hosts.yaml")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ 加载配置失败: %v\n", err)
		os.Exit(1)
	}

	if len(*cfg) == 0 {
		fmt.Fprintln(os.Stderr, "❌ 配置中没有服务器，请在 hosts.yaml 中添加服务器列表")
		os.Exit(1)
	}

	fmt.Printf("🌐 共 %d 台主机，开始并发巡检...\n\n", len(*cfg))

	// 同一轮巡检使用同一个时间，作为各表之间的关联 ID。
	checkTime := time.Now().Truncate(time.Second)

	// 并发巡检
	results := make([]checks.HostResult, len(*cfg))
	var wg sync.WaitGroup

	for i, srv := range *cfg {
		wg.Add(1)
		go func(idx int, hostCfg config.HostConfig) {
			defer wg.Done()
			results[idx] = inspectHost(hostCfg, checkTime)
		}(i, srv)
	}

	wg.Wait()
	if err = mysql.SetLastTime(checkTime); err != nil {
		fmt.Fprintf(os.Stderr, "写入最近巡检时间失败：%v\n", err)
		os.Exit(1)
	}
	//将数据永久存储在Mysql中。
	if err = mysql.SetData(&results); err != nil {
		fmt.Fprintf(os.Stderr, "检查数据吸入数据库失败：%v\n", err)
		os.Exit(1)
	}

	// 输出表格
	fmt.Println()
	reporter.PrintTable(results)

	// 生成报告
	fmt.Println()
	jsonPath := "report.json"
	if err := reporter.GenerateJSON(results, jsonPath); err != nil {
		fmt.Printf("⚠ JSON 报告生成失败: %v\n", err)
	} else {
		fmt.Printf("📄 JSON 报告已生成: %s\n", jsonPath)
	}

	htmlPath := "report.html"
	if err := reporter.GenerateHTML(results, htmlPath); err != nil {
		fmt.Printf("⚠ HTML 报告生成失败: %v\n", err)
	} else {
		fmt.Printf("📄 HTML 报告已生成: %s\n", htmlPath)
	}

	// 统计
	passCount, warnCount, failCount := 0, 0, 0
	for _, r := range results {
		switch checks.OverallStatus(r.Items) {
		case "PASS":
			passCount++
		case "WARN":
			warnCount++
		case "FAIL":
			failCount++
		}
	}
	fmt.Println()
	fmt.Printf("✅ 巡检完成: PASS=%d | WARN=%d | FAIL=%d\n", passCount, warnCount, failCount)
}

// inspectHost 对单台主机执行巡检
func inspectHost(hostCfg config.HostConfig, checkTime time.Time) checks.HostResult {
	result := checks.HostResult{
		Host:      hostCfg.Host,
		CheckTime: checkTime,
	}

	// 连接 SSH
	client, err := sshutil.NewClient(
		hostCfg.Host,
		hostCfg.Port,
		hostCfg.User,
		hostCfg.Password,
		hostCfg.KeyFile,
	)
	if err != nil {
		result.Error = err.Error()
		//连接主机失败添加失败状态
		result.Status = checks.StatusUNCON
		return result
	}
	defer client.Close()

	// 执行所有检查
	for _, check := range checks.AllChecks() {
		item := check.Fn(client)
		// 如果检查函数没有设置 Name，使用注册名称
		if item.Name == "" {
			item.Name = check.Name
		}
		result.Items = append(result.Items, item)

		// 小延迟避免过快连续执行
		time.Sleep(50 * time.Millisecond)
	}
	result.Status = checks.OverallStatus(result.Items)
	return result
}
