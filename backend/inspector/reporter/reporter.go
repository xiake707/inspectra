package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"inspector/checks"
)

// ReportData 报告数据
type ReportData struct {
	GeneratedAt string              `json:"generated_at"`
	Hosts       []checks.HostResult `json:"hosts"`
}

// PrintTable 终端打印表格
func PrintTable(results []checks.HostResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// 表头
	fmt.Fprintln(w, "Host\tCPU\tMemory\tDisk\tDocker\tNginx\tMySQL\tNetwork\tDNS\tTimeSync\tStatus")
	fmt.Fprintln(w, "----\t---\t------\t----\t------\t-----\t-----\t-------\t---\t--------\t------")

	for _, hr := range results {
		if hr.Error != "" {
			fmt.Fprintf(w, "%s\tERROR: %s\n", hr.Host, hr.Error)
			continue
		}
		// 构建列
		cols := []string{hr.Host}
		for _, item := range hr.Items {
			val := item.Value
			if item.Status == checks.StatusFAIL {
				val = val + " ❌"
			} else if item.Status == checks.StatusWARN {
				val = val + " ⚠"
			}
			cols = append(cols, val)
		}
		cols = append(cols, checks.OverallStatus(hr.Items))
		fmt.Fprintln(w, strings.Join(cols, "\t"))
	}

	w.Flush()
}

// GenerateJSON 生成 JSON 报告
func GenerateJSON(results []checks.HostResult, path string) error {
	report := ReportData{
		GeneratedAt: time.Now().Format(time.RFC3339),
		Hosts:       results,
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GenerateHTML 生成 HTML 报告
func GenerateHTML(results []checks.HostResult, path string) error {
	report := ReportData{
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
		Hosts:       results,
	}
	html := buildHTML(report)
	return os.WriteFile(path, []byte(html), 0644)
}

// escapeHTML 简易 HTML 转义
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// buildHTML 构建 HTML 内容（使用 strings.Replacer 避免 % 冲突）
func buildHTML(data ReportData) string {
	var rows strings.Builder
	for _, hr := range data.Hosts {
		overall := checks.OverallStatus(hr.Items)
		rowClass := ""
		switch overall {
		case "FAIL":
			rowClass = "class='fail'"
		case "WARN":
			rowClass = "class='warn'"
		case "PASS":
			rowClass = "class='pass'"
		}

		rows.WriteString("<tr " + rowClass + "><td>" + escapeHTML(hr.Host) + "</td>")
		if hr.Error != "" {
			rows.WriteString("<td colspan='9' style='color:red'>错误: " + escapeHTML(hr.Error) + "</td>")
		} else {
			for _, item := range hr.Items {
				cellClass := ""
				switch item.Status {
				case "FAIL":
					cellClass = "cell-fail"
				case "WARN":
					cellClass = "cell-warn"
				case "OK":
					cellClass = "cell-ok"
				}
				rows.WriteString("<td class='" + cellClass + "'>" + escapeHTML(item.Value) + "</td>")
			}
		}
		rows.WriteString("<td class='cell-" + strings.ToLower(overall) + "'><strong>" + overall + "</strong></td></tr>\n")
	}

	// HTML 模板
	tmpl := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Go-Inspector 巡检报告</title>
<style>
  body { font-family: 'Segoe UI', Arial, sans-serif; margin: 20px; background: #f5f7fa; color: #333; }
  h1 { color: #2c3e50; border-bottom: 3px solid #3498db; padding-bottom: 10px; }
  .meta { color: #7f8c8d; margin-bottom: 20px; }
  table { border-collapse: collapse; width: 100%; background: #fff; box-shadow: 0 2px 8px rgba(0,0,0,0.1); border-radius: 4px; overflow: hidden; }
  th { background: #3498db; color: #fff; padding: 12px 8px; text-align: left; font-size: 13px; }
  td { padding: 10px 8px; border-bottom: 1px solid #ecf0f1; font-size: 13px; }
  tr:hover { background: #f0f6ff; }
  .pass { background: #eafaf1; }
  .warn { background: #fef9e7; }
  .fail { background: #fdedec; }
  .cell-ok { color: #27ae60; font-weight: bold; }
  .cell-warn { color: #e67e22; font-weight: bold; }
  .cell-fail { color: #e74c3c; font-weight: bold; }
  .cell-pass { color: #27ae60; }
  .summary { margin-top: 20px; padding: 15px; background: #fff; border-radius: 4px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
</style>
</head>
<body>
<h1>🔍 Go-Inspector 巡检报告</h1>
<p class="meta">生成时间: __GEN_TIME__</p>
<table>
<thead>
<tr>
  <th>Host</th>
  <th>CPU</th>
  <th>Memory</th>
  <th>Disk</th>
  <th>Docker</th>
  <th>Nginx</th>
  <th>MySQL</th>
  <th>Network</th>
  <th>DNS</th>
  <th>TimeSync</th>
  <th>Status</th>
</tr>
</thead>
<tbody>
__ROWS__
</tbody>
</table>
<div class="summary">
  <strong>统计:</strong>
  共 __TOTAL__ 台主机 |
  ✅ PASS: __PASS__ |
  ⚠ WARN: __WARN__ |
  ❌ FAIL: __FAIL__
</div>
<p style="margin-top:20px; color:#95a5a6; font-size:12px;">Powered by Go-Inspector</p>
</body>
</html>`

	// 使用 Replacer 避免 % 冲突
	r := strings.NewReplacer(
		"__GEN_TIME__", data.GeneratedAt,
		"__ROWS__", rows.String(),
		"__TOTAL__", fmt.Sprintf("%d", len(data.Hosts)),
		"__PASS__", fmt.Sprintf("%d", countByStatus(data.Hosts, "PASS")),
		"__WARN__", fmt.Sprintf("%d", countByStatus(data.Hosts, "WARN")),
		"__FAIL__", fmt.Sprintf("%d", countByStatus(data.Hosts, "FAIL")),
	)
	return r.Replace(tmpl)
}

func countByStatus(hosts []checks.HostResult, status string) int {
	n := 0
	for _, h := range hosts {
		if checks.OverallStatus(h.Items) == status {
			n++
		}
	}
	return n
}
