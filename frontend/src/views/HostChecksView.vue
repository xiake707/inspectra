<!--
  主机检查指标展示组件（页面级）
  ============================================================
  接入方式（不改动任何现有文件，仅在你的 App.vue 里添加即可）：
    1) import HostChecksView from "./views/HostChecksView.vue";
    2) 在 <main> 中渲染：
         <HostChecksView :items="hosts[0]?.items" title="主机检查" />
    3) 接入真实数据时把 :items 换成后端返回的 items 数组即可，
       组件会自动区分两类卡片：
         · 使用率类（值中含 %，如 CPU / Memory / Disk）：数字 + 进度条
         · 服务状态类（如 Docker / Nginx / MySQL ...）：active/inactive 状态色卡片

  当前展示的是组件内置的静态 mock 数据（mockCheckItems），
  后续修改数据只需编辑下方的 mockCheckItems 数组。
-->
<script setup>
import { computed } from "vue";
import {
  Cpu,
  Odometer,
  FolderOpened,
  Box,
  SetUp,
  Collection,
  Connection,
  Aim,
  Timer,
  Monitor,
  InfoFilled,
} from "@element-plus/icons-vue";
import SectionHeading from "../components/SectionHeading.vue";

// ================= 静态 mock 数据（接入真实数据时替换） =================
/*const mockCheckItems = [
  { name: "CPU", value: "2%", status: "OK" },
  { name: "Memory", value: "20% (398M/1960M)", status: "OK" },
  {
    name: "Disk",
    value: "32%",
    status: "OK",
    detail: "系统磁盘总使用率: 32%",
  },
  { name: "Docker", value: "运行中(0容器)", status: "OK" },
  { name: "Nginx", value: "failed\nunknown", status: "WARN" },
  {
    name: "MySQL",
    value: "已停止",
    status: "FAIL",
    detail: "服务 mysqld 未运行",
  },
  { name: "Network", value: "可达(68.8ms)", status: "OK" },
  {
    name: "DNS",
    value: "失败",
    status: "FAIL",
    detail:
      "bash: line 1: nslookup: command not found\nbash: line 1: dig: command not found\nFAIL\n",
  },
  { name: "TimeSync", value: "chronyd 正常", status: "OK" },
];
*/

const props = defineProps({
  selectedHost: { type: Object, default: null },
  title: { type: String, default: "主机检查" },
  description: {
    type: String,
    default:
      "单台主机的巡检指标：资源类（CPU / Memory / Disk）展示使用率进度条，服务类展示 active / inactive 运行状态。",
  },
});
const emit = defineEmits(["open-trend"]);

const displayItems = computed(() =>  props.selectedHost?.items ?? []);

const hostName = computed(
  () => props.selectedHost?.host ?? props.selectedHost?.Host ?? "未选择主机",
);
const hostPort = computed(() => props.selectedHost?.port ?? props.selectedHost?.Port ?? "22");
const hostCheckedAt = computed(
  () => props.selectedHost?.CheckTime ?? props.selectedHost?.checkedAt ?? "静态演示数据",
);

function openTrend(item) {
  const host = props.selectedHost?.host ?? props.selectedHost?.Host ?? "";
  emit("open-trend", { host, metric: item.name });
}

// 指标图标映射（未匹配到的一律回退到 Monitor）
const ICONS = {
  CPU: Cpu,
  Memory: Odometer,
  Disk: FolderOpened,
  Docker: Box,
  Nginx: SetUp,
  MySQL: Collection,
  Network: Connection,
  DNS: Aim,
  TimeSync: Timer,
};
const FALLBACK_ICON = Monitor;
const RESOURCE_NAMES = new Set(["CPU", "Memory", "Disk"]);

// 从 value 中解析百分比，例如 "20% (398M/1960M)" -> 20；解析不到返回 null
function parsePercent(value) {
  const match = String(value ?? "").match(/(\d+(?:\.\d+)?)\s*%/);
  return match ? Number(match[1]) : null;
}

// 是否“使用率类”指标：值里带 % 就按使用率卡片渲染
const isUsage = (item) => parsePercent(item.value) !== null;

// 使用率进度条颜色：<60 绿 / 60~85 橙 / >85 红（可按需改为跟随 status）
function barColor(percent) {
  if (percent == null) return "var(--muted)";
  if (percent > 85) return "#c75454";
  if (percent >= 60) return "var(--orange)";
  return "var(--green)";
}

// el-tag 类型映射
function tagType(status) {
  const s = String(status ?? "").toUpperCase();
  if (s === "WARN") return "warning";
  if (s === "FAIL") return "danger";
  return "success";
}

// 服务状态推导：OK -> active(绿) / FAIL -> inactive(红) / WARN -> warn(橙)
// 没有 status 字段时，根据 value 文本兜底推断
function serviceState(item) {
  const s = String(item.status ?? "").toUpperCase();
  if (s === "FAIL") return { key: "inactive", label: "inactive" };
  if (s === "WARN") return { key: "warn", label: "warn" };
  if (s === "OK" || s === "PASS") return { key: "active", label: "active" };
  const v = String(item.value ?? "");
  if (/运行|正常|可达|active|running|up/i.test(v))
    return { key: "active", label: "active" };
  if (/停止|失败|failed|inactive|down|error/i.test(v))
    return { key: "inactive", label: "inactive" };
  return { key: "warn", label: "warn" };
}

function iconOf(item) {
  return ICONS[item.name] || FALLBACK_ICON;
}

// 预组装两类卡片，模板里直接取用，避免重复调用解析函数
const usageCards = computed(() =>
  displayItems.value
    .filter((item) => RESOURCE_NAMES.has(item.name) && isUsage(item))
    .map((item) => {
      const percent = parsePercent(item.value);
      return {
        item,
        percent,
        tag: tagType(item.status),
        color: barColor(percent),
      };
    }),
);

const serviceCards = computed(() =>
  displayItems.value
    .filter((item) => !RESOURCE_NAMES.has(item.name) || !isUsage(item))
    .map((item) => ({ item, state: serviceState(item), tag: tagType(item.status) })),
);

// 汇总统计
const okCount = computed(
  () => displayItems.value.filter((i) => String(i.status).toUpperCase() === "OK").length,
);
const warnCount = computed(
  () => displayItems.value.filter((i) => String(i.status).toUpperCase() === "WARN").length,
);
const failCount = computed(
  () => displayItems.value.filter((i) => String(i.status).toUpperCase() === "FAIL").length,
);
</script>

<template>
  <section class="page-view">
    <SectionHeading
      eyebrow="HOST CHECKS"
      :title="props.selectedHost ? `${hostName} 主机检查` : title"
      :description="description"
    >
      <div class="check-summary">
        <span class="summary-pill ok"><i></i>正常 {{ okCount }}</span>
        <span class="summary-pill warn"><i></i>警告 {{ warnCount }}</span>
        <span class="summary-pill fail"><i></i>异常 {{ failCount }}</span>
        <span class="summary-pill total"><i></i>共 {{ displayItems.length }} 项</span>
      </div>
    </SectionHeading>

    <article class="host-identity panel">
      <div class="host-identity-main">
        <span class="host-identity-icon"><el-icon><Monitor /></el-icon></span>
        <div>
          <span class="panel-kicker">TARGET HOST</span>
          <strong>{{ hostName }}</strong>
        </div>
      </div>
      <div class="host-identity-meta">
        <span><small>SSH 端口</small><b>{{ hostPort }}</b></span>
        <span><small>最近检查</small><b>{{ hostCheckedAt }}</b></span>
      </div>
    </article>

    <!-- 使用率类指标：大数字 + 进度条 -->
    <div v-if="usageCards.length" class="metric-grid resource-metric-list">
      <article
        v-for="card in usageCards"
        :key="'usage-' + card.item.name"
        class="metric-card usage-card resource-card-clickable"
        @click="openTrend(card.item)" 
      >
        <div class="resource-card-title">
          <span class="metric-icon usage">
            <el-icon><component :is="iconOf(card.item)" /></el-icon>
          </span>
          <span class="metric-name">{{ card.item.name }}</span>
        </div>
        <strong class="usage-percent">{{ card.percent ?? "—" }}<small>%</small></strong>
        <el-progress
          :percentage="card.percent ?? 0"
          :stroke-width="9"
          :color="card.color"
          :trail-color="'#eef1f2'"
          :show-text="false"
        />
        <span class="resource-card-value">{{ card.item.value }}</span>
        <el-tag size="small" effect="light" :type="card.tag">{{ card.item.status }}</el-tag>
      </article>
    </div>

    <!-- 服务状态类指标：active / inactive / warn 状态色卡片 -->
    <div v-if="serviceCards.length" class="metric-grid">
      <article
        v-for="card in serviceCards"
        :key="'service-' + card.item.name"
        class="metric-card service-card"
        :class="'state-' + card.state.key"
      >
        <div class="metric-head">
          <span class="metric-icon service">
            <el-icon><component :is="iconOf(card.item)" /></el-icon>
          </span>
          <span class="metric-name">{{ card.item.name }}</span>
          <span class="state-chip" :class="'state-' + card.state.key">
            <i></i>{{ card.state.label }}
          </span>
        </div>
        <div class="service-value" :class="'state-' + card.state.key">
          {{ card.item.value }}
        </div>
        <div v-if="card.item.detail" class="service-detail">
          {{ card.item.detail }}
        </div>
      </article>
    </div>

    <div class="form-note">
      <el-icon><InfoFilled /></el-icon>
      <div>
        <b>注意检查时间</b>
      </div>
    </div>
  </section>
</template>
