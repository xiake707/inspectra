<script setup>
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { Refresh, House, TrendCharts, Grid, Monitor } from "@element-plus/icons-vue";
import DashboardView from "./views/DashboardView.vue";
import TrendsView from "./views/TrendsView.vue";
import HostsView from "./views/HostsView.vue";
import HostChecksView from "./views/HostChecksView.vue";
import HostFormDialog from "./components/HostFormDialog.vue";
import { api } from "./services/api.js";

const activeView = ref(location.hash.slice(1) || "overview");
const addHostOpen = ref(false);
const trends = ref(null);
const hostlist = ref([])
const hosts = ref([]);
const hostIP = ref("");
const selectedHost = ref(null);
const selectedMetric = ref("both");
const isChecking = ref(false);
//const checkResults = ref([]);
const isLoading = ref(false);
const views = {
  overview: "系统总览",
  trends: "趋势分析",
  hosts: "主机管理",
  checks: "主机检查",
};

const pageTitle = computed(() => views[activeView.value] || views.overview);

function formatCheckTime(value) {
  if (!value) return "无记录";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "无记录";
  const pad = (number) => String(number).padStart(2, "0");
  return `${String(date.getFullYear()).slice(-2)}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

const managedHosts = computed(() =>
  hostlist.value.map((host) => {
    const latest = hosts.value.find(
      (item) => (item.Host || item.host) === host.Host,
    );
    return {
      ...host,
      status: latest?.Status || latest?.status || "待巡检",
      checkedAt: formatCheckTime(
        latest?.CheckTime || latest?.checkTime || latest?.checktime,
      ),
    };
  }),
);

function navigate(view, openForm = false) {
  //跳转函数
  activeView.value = view;
  location.hash = view;
  if (openForm) addHostOpen.value = true;
}
function openHostChecks(host) {
  selectedHost.value = host;
  //selectedMetric.value = "both";
  activeView.value = "checks";
  location.hash = "checks";
}
async function openHostTrend({ host, metric }) {
  hostIP.value = host;
  selectedMetric.value = metric;
  activeView.value = "trends";
  location.hash = "trends";
  await loadTrends(host);
}

async function openHostList() {
//  hostIP.value = host;
//  selectedMetric.value = metric;
  activeView.value = "hosts";
  location.hash = "hosts";
  await loadHostList();
}
async function loadTrends(host) {
  isLoading.value = true;
  try {
    trends.value = await api.getItems(host);
  } catch (error) {
    ElMessage.error(error.message);
  } finally {
    isLoading.value = false;
  }
}
async function loadHostList(host) {
  isLoading.value = true;
  try {
    hostlist.value = await api.getHosts();
  } catch (error) {
    ElMessage.error(error.message);
  } finally {
    isLoading.value = false;
  }
}
async function loadHosts() {
  isLoading.value = true;
  try {
    hosts.value = await api.listHosts();
  } catch (error) {
    ElMessage.error(error.message);
  } finally {
    isLoading.value = false;
  }
}
async function handleHostCreated(payload) {
  try {
    const response = await api.createHost(payload);
    const message = response.data?.message || "主机信息添加成功";
    await loadHostList();
    addHostOpen.value = false;
    ElMessage.success(message);
  } catch (error) {
    const message = error.response
      ? error.response.data?.message ||
        error.response.data?.error ||
        `服务器返回 HTTP ${error.response.status}`
      : "无法连接后端服务，请确认后端已启动并监听 8080 端口";
    ElMessage.error(`主机信息添加失败：${message}`);
  }
}
async function handleStartCheck() {
  if (isChecking.value) return;
  isChecking.value = true;
  try {
    const response = await api.checkHosts();
    const message = response.data?.message || "主机检查已开始";
    ElMessage.success(message);
    await Promise.all([loadHosts(), loadHostList()]);
  } catch (error) {
    const message = error.response
      ? error.response.data?.message ||
        error.response.data?.error ||
        `服务器返回 HTTP ${error.response.status}`
      : "无法连接后端服务，请确认后端已启动并监听 8080 端口";
    ElMessage.error(`开始检查失败：${message}`);
  } finally {
    isChecking.value = false;
  }
}
onMounted(() => {
  window.addEventListener("hashchange", () => {
    activeView.value = location.hash.slice(1) || "overview";
  });
  loadHosts();
});
</script>
<template>
  <div class="app-shell">
    <aside class="sidebar">
      <button class="brand" aria-label="返回总览" @click="navigate('overview')">
        <span class="brand-mark">GI</span>
        <span>
          <strong>Go Inspector</strong>
          <small>主机巡检平台</small>
        </span>
      </button>
      <div class="workspace-label">工作区</div>
      <nav class="nav-list">
        <button
          :class="['nav-item', { 'is-active': activeView === 'overview' }]"
          @click="navigate('overview')"
        >
          <House />总览
        </button>
        <button
          :class="['nav-item', { 'is-active': activeView === 'trends' }]"
          @click="navigate('trends')"
        >
          <TrendCharts />趋势分析
        </button>
        <button
          :class="['nav-item', { 'is-active': activeView === 'hosts' }]"
          @click="openHostList"
        >
          <Grid />主机管理
        </button>
        <button
          :class="['nav-item', { 'is-active': activeView === 'checks' }]"
          @click="navigate('checks')"
        >
          <Monitor />主机检查
        </button>
      </nav>
      <div class="sidebar-bottom">
        <div class="connection-state">
          <span class="status-dot"></span>
          <span><b>数据层已连接</b></span>
        </div>
        <div class="sidebar-meta">v1.1 · Linux inspectra</div>
      </div>
    </aside>
    <main class="main-content">
      <header class="topbar">
        <div class="breadcrumb">
          <span>工作区</span>
          <b>/</b>
          <strong>{{ pageTitle }}</strong>
        </div>
        <div class="topbar-actions">
          <span class="mode-pill"><i></i> MOCK DATA</span>
          <el-button
            text
            circle
            :icon="Refresh"
            aria-label="刷新数据"
            @click="loadHosts"
          />
          <div class="avatar">XW</div>
        </div>
      </header>
      <div v-loading="isLoading" class="view-container">
        <DashboardView
          v-if="activeView === 'overview'"
          :hosts="hosts"
          @navigate="navigate"
          @open-checks="openHostChecks"
        />
        <TrendsView
          v-else-if="activeView === 'trends'"
          :hostIP="hostIP"
          :trends="trends"
          :selected-metric="selectedMetric"
        />
        <!-- 主机检查：调试用视图，默认渲染组件内置的静态 mock 数据；
             接入真实数据后改为 :items="hosts[0]?.items" -->
        <HostChecksView
          v-else-if="activeView === 'checks'"
          :selectedHost="selectedHost"
          @open-trend="openHostTrend"
        />
        <HostsView
          v-else
          :hosts="managedHosts"
          :checking="isChecking"
          @add="addHostOpen = true"
          @start-check="handleStartCheck"
        />
      </div>
    </main>
    <HostFormDialog v-model="addHostOpen" @submit="handleHostCreated" />
  </div>
</template>
