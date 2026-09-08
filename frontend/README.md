# Go Inspector Frontend

Go Inspector 的 Vue 3 前端工程，使用 Vite、Element Plus、ECharts 和 Axios。当前默认使用 mock 数据，页面可以完整运行，但不会连接 Go API、MySQL 或 Linux 主机。

## 技术栈

- Vue 3 + `<script setup>`：页面与组件化开发
- Vite：开发服务器和生产构建
- Element Plus：表单、表格、弹窗、按钮和反馈组件
- ECharts：资源趋势折线图和主机对比柱状图
- Axios：统一 HTTP 客户端，已预留 API 层

## 目录结构

```text
frontend/
├── index.html
├── package.json
├── vite.config.js
├── src/
│   ├── main.js                 Vue 入口和 Element Plus 注册
│   ├── App.vue                 布局、导航、页面状态和全局数据加载
│   ├── styles.css              全局视觉变量、布局和响应式样式
│   ├── components/
│   │   ├── EChart.vue           ECharts 生命周期封装
│   │   ├── HostFormDialog.vue   添加主机表单与提交事件
│   │   ├── SectionHeading.vue   页面标题复用组件
│   │   └── StatCard.vue         总览统计卡片
│   ├── views/
│   │   ├── DashboardView.vue    系统总览
│   │   ├── TrendsView.vue       趋势分析
│   │   └── HostsView.vue        主机管理
│   ├── services/api.js          Axios 实例和所有后端请求边界
│   ├── mock/hosts.js            演示主机与图表数据
│   └── utils/hostValidation.js  host_config 输入规则
└── dist/                        npm run build 生成的 Nginx 静态文件
```

## 启动和构建

```bash
npm install
npm run dev
npm run build
npm run preview
```

开发地址默认是 `http://localhost:5173`。生产构建后的 `dist/` 可以作为 Nginx 静态根目录。

## 接入后端

所有请求集中在 `src/services/api.js`。把 `const USE_MOCK = true` 改成 `false` 后，页面会使用以下接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/hosts` | 主机列表及最近巡检状态 |
| `GET` | `/api/hosts/{host}/metrics?days=15&metric=cpu,memory` | 指定主机历史指标 |
| `GET` | `/api/metrics?days=15&metric=cpu&hosts=all` | 跨主机指标比较 |
| `POST` | `/api/hosts` | 新增 `host_config` 配置 |

推荐趋势响应：

```json
{
  "xAxis": ["2026-08-09", "2026-08-10"],
  "series": [
    { "name": "CPU", "data": [31, 34] },
    { "name": "Memory", "data": [46, 47] }
  ]
}
```

## 修改图表数据

修改 `src/mock/hosts.js`：`mockTrend.xAxis` 是趋势图 x 轴日期，`mockTrend.series[].data` 是趋势图 y 轴 CPU/Memory 百分比；`mockComparison.xAxis` 是主机对比图 x 轴地址，`mockComparison.series[].data` 是对应 y 轴指标值。默认趋势数据有 15 个参考值，数组长度应保持一致。

## 添加主机校验

`src/utils/hostValidation.js` 对应现有 `host_config` 表字段：`Host` 是合法 IPv4 且最多 15 个字符；`User` 必填，字母或下划线开头，仅允许字母、数字、`_`、`.`、`-`，最多 20 位；`Port` 为 1-65535 的整数；`Password` 最多 20 位；`KeyFile` 最多 30 位且必须是绝对路径或 `$HOME/` 开头。

表单目前通过 `POST /api/hosts` 的事件边界生成 payload，但 mock 模式不会发出网络请求。

## Nginx 部署

将 `dist/` 内容复制到 Nginx 静态目录，并配置 `location / { try_files $uri $uri/ /index.html; }` 与 `location /api/ { proxy_pass http://goapi:8080; }`。本次没有修改现有 Nginx 配置，也没有连接后端。
