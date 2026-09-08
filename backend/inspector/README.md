# Go-Inspector 🔍

多主机批量巡检工具：通过 SSH 并发登录多台服务器，自动检查系统运行状况，并生成终端表格、JSON、HTML 三种形式的巡检报告。

## 功能特性

- 🚀 并发巡检：多台主机同时检查，结果按主机索引汇总，互不干扰
- 🔐 灵活认证：支持密码、指定密钥文件，自动回退到 `~/.ssh/id_rsa` / `~/.ssh/id_ed25519`
- 📋 九大检查项：CPU、内存、磁盘、Docker、Nginx、MySQL、网络、DNS、时间同步
- 📊 三种报告：终端彩色表格、JSON 结构化报告、自包含 HTML 网页报告
- ⚙️ 配置灵活：支持纯 IP 列表和完整主机配置两种 YAML 格式，支持全局默认值

## 环境要求

- Go 1.21+
- 目标主机可被 SSH 访问（默认端口 22）

## 快速开始

```bash
# 编译
go build -o go-inspector .

# 运行（默认读取 hosts.yaml）
./go-inspector

# 指定其他配置文件
./go-inspector /path/to/hosts.yaml
```

运行结束后会在当前目录生成：

- `report.json` — 结构化 JSON 报告
- `report.html` — 可视化 HTML 报告（浏览器直接打开即可）

## 配置说明

配置文件为 YAML 格式，支持两种写法：

### 格式一：简洁写法（纯 IP 列表）

```yaml
global:
  user: root      # 默认用户名
  port: 22        # 默认 SSH 端口

servers:
  - 192.168.10.10
  - 192.168.10.20
```

### 格式二：完整写法（逐台覆盖全局配置）

```yaml
global:
  user: root
  port: 22

servers:
  - host: 192.168.10.40
    user: admin       # 覆盖全局用户名
    port: 2222        # 覆盖全局端口
    password: "p@ssw0rd"
  - host: 192.168.10.50
    user: root
    key_file: "~/.ssh/my_key"   # 支持 ~ 展开
```

### 认证优先级

1. 主机配置中指定的 `key_file`
2. 主机/全局配置中的 `password`
3. 自动尝试 `~/.ssh/id_rsa`
4. 自动尝试 `~/.ssh/id_ed25519`

> ⚠️ 注意：密码以明文形式保存在配置文件中，请妥善保管配置文件权限。生产环境建议使用 SSH 密钥认证。

## 巡检项与判定阈值

| 检查项 | 说明 | WARN | FAIL |
| --- | --- | --- | --- |
| CPU | `top` / `uptime` 负载 | > 80% | > 90%（或负载 > 8） |
| Memory | `free -m` 内存使用率 | > 85% | > 95% |
| Disk | `df` 所有挂载点取最差 | > 85% | > 95% |
| Docker | 是否安装、服务状态、容器数 | 服务非 active | 未运行 |
| Nginx | 是否安装、服务/进程状态 | 状态未知 | 已停止 |
| MySQL | mysqld/mysql/mariadb 服务或进程 | 状态未知 | 未运行 |
| Network | ping 外网（8.8.8.8 / 114.114.114.114） | - | 全部不通 |
| DNS | nslookup / dig baidu.com | 结果异常 | 解析失败 |
| TimeSync | chronyd / ntpd / timesyncd / timedatectl | 未检测到同步服务 | - |

总体状态规则：任一 FAIL → `FAIL`；否则任一 WARN → `WARN`；全部 OK → `PASS`。

## 目录结构

```
go-inspector/
├── main.go               # 入口：并发调度、报告生成、结果统计
├── config/
│   └── config.go         # YAML 配置加载与默认值填充
├── sshutil/
│   └── ssh.go            # SSH 客户端封装（认证、命令执行）
├── checks/
│   └── checks.go         # 九大检查项与阈值判定
├── reporter/
│   └── reporter.go       # 终端表格 / JSON / HTML 报告
├── templates/            # 预留模板目录（当前 HTML 内嵌在代码中）
└── hosts.yaml            # 示例配置
```

## 扩展新检查项

在 `checks/checks.go` 中实现一个 `func CheckXxx(client *sshutil.Client) checks.Item`，然后在 `AllChecks()` 注册表中添加即可，无需改动其他代码。

## 已知限制

- SSH 主机密钥未校验（使用 `InsecureIgnoreHostKey`），存在中间人攻击风险，公网环境请谨慎使用
- 远程命令未设置超时，目标主机命令卡住时会一直等待
- 未提供 `go.sum`，建议运行 `go mod tidy` 后提交以锁定依赖版本
