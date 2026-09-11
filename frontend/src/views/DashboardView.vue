<script setup>
import { computed, ref } from "vue";
import {
    Close,
  CircleCheck,
  CircleClose,
  Warning,
  Clock,
  Plus,
  TrendCharts, Monitor,
} from "@element-plus/icons-vue";
import SectionHeading from "../components/SectionHeading.vue";
import StatCard from "../components/StatCard.vue";
//接受变量，设置时间名称
const props = defineProps({ hosts: { type: Array, default: () => [] } });
const emit = defineEmits(["navigate", "open-checks"]);
const hoveredHost = ref(null);
const getStatus = (host) => String(host.Status ?? host.status ?? "").toUpperCase();
//统计pass的主机数据
const passHostArr = computed(() => {
    let arr = []
    props.hosts.forEach((host) => {
        if (getStatus(host) === "PASS"){
           arr.push(host) 
        }
    });
    return arr
});
//统计fail的主机数据
const failHostArr = computed(() => {
    let arr = []
    props.hosts.forEach((host) => {
        if (getStatus(host) === "FAIL"){
           arr.push(host) 
        }
    });
    return arr
});
//找到引起fail的指标
const getFailItem = (host) => {
    let failDetail = ""
    host.items?.forEach((item) => {
       if (getStatus(item) === "FAIL") {
            failDetail = item.name + ": " +item.detail
       }
    })

    return failDetail
}
    
//统计warn的主机数据
const warnHostArr = computed(() => {
    let arr = []
    props.hosts.forEach((host) => {
        if (getStatus(host) === "WARN"){
           arr.push(host) 
        }
    });
    return arr
});
//统计unconnection的主机
const unconHostArr = computed(() =>{
    let arr=[]
    props.hosts.forEach((host) => {
        if (getStatus(host) === "UNCON"){
           arr.push(host) 
        }
    });
    return arr
});
//计算检查时间
let checktime = computed(() => {
  return String(props.hosts[0]?.CheckTime ?? "09:38");
});
let checktimeT = computed(() => {
  return String(props.hosts[0]?.CheckTime ?? "09:38").slice(11, 19);
});
//统计每种状态的数量
const itemsLenth = computed(() => {
  let length = 0;
  props.hosts?.forEach((host) => {
    if (getStatus(host) !== "UNCON") length = host.items?.length ?? length;
  });
  return length;
});
const passCount = computed(() => {
  return props.hosts?.filter((host) => getStatus(host) === "PASS").length || 0;
});
const failCount = computed(() => {
  return props.hosts?.filter((host) => getStatus(host) === "FAIL").length || 0;
});
const unconCount = computed(() => {
  return props.hosts?.filter((host) => getStatus(host) === "UNCON").length || 0;
});
const warnCount = computed(() => {
  return props.hosts?.filter((host) => getStatus(host) === "WARN").length || 0;
});
//计算百分比
const toPercent = (count) => {
  const total = props.hosts.length;
  return total === 0 ? "0%" : `${((count / total) * 100).toFixed(1)}%`;
};
const passPercent = computed(() => {
  return toPercent(passCount.value);
});
const warnPercent = computed(() => toPercent(warnCount.value));
const failPercent = computed(() => toPercent(failCount.value));
const unconPercent = computed(() => toPercent(unconCount.value));
</script>
<template>
  <section class="page-view">
    <SectionHeading
      eyebrow="SATURDAY · 23 AUG 2026"
      title="早上好，Xwang"
      description="这里是你的 Linux 主机健康概览。"
    >
      <el-button
        type="primary"
        :icon="Plus"
        @click="emit('navigate', 'hosts', true)"
        >添加主机</el-button
      >
    </SectionHeading>
    <div class="stat-grid">
      <StatCard
        label="托管主机"
        :value="hosts.length || 暂无托管主机"
        detail="已托管主机"
        tone="blue"
        :icon="TrendCharts"
      />
      <StatCard
        label="健康状态"
        :value="`${passPercent}`"
        :detail="`${passCount || 0} 台主机运行正常`"
        tone="green"
        :icon="CircleCheck"
      />
      <StatCard
        label="待处理告警"
        :value="warnCount"
        detail="需要关注 WARN"
        tone="orange"
        :icon="Warning"
      />
      <StatCard
        label="断联的主机"
        :value="unconCount"
        detail="失去联系的主机"
        tone="red"
        :icon="CircleClose"
      />
      <StatCard
        label="服务异常主机"
        :value="failCount"
        detail="服务异常掉线"
        tone="orange"
        :icon="Warning"
      />
      <StatCard
        label="最近巡检"
        :value="checktimeT || '09:38'"
        :detail="
          `${checktime} ~ ${itemsLenth} 项指标` || '2026-08-23 · 18 项指标'
        "
        tone="gray"
        :icon="Clock"
      />
    </div>
    <div class="overview-grid">
        <article class="panel">
            <div class="panel-header">
                <div>
                    <span class="panel-kicker">HEALTH OVERVIEW</span>
                    <h2>集群健康状态</h2>
                </div>
                <el-button link type="primary" @click="emit('navigate', 'hosts')"
                                >查看全部 →</el-button
                            >
            </div>
                    <div class="health-layout">
                        <div class="donut" :style="{ '--pass-percent': passPercent }">
                            <div>
                                <strong>{{ passPercent }}</strong>
                                <small>健康率</small>
                            </div>
                        </div>
                        <div class="legend-list">
                            <div>
                                <i class="legend-dot green-dot"></i>
                                <span>PASS</span>
                                <b>{{ passCount}}</b>
                                <small>{{ passPercent }}</small>
                            </div>
                            <div>
                                <i class="legend-dot orange-dot"></i>
                                <span>WARN</span>
                                <b>{{ warnCount }}</b>
                                <small>{{warnPercent}}</small>
                            </div>
                            <div>
                                <i class="legend-dot red-dot"></i>
                                <span>FAIL</span>
                                <b>{{failCount}}</b>
                                <small>{{ failPercent }}</small>
                            </div>
                            <div>
                                <i class="legend-dot red-dot"></i>
                                <span>UNCON</span>
                                <b>{{unconCount}}</b>
                                <small>{{ unconPercent }}</small>
                            </div>
                        </div>
                    </div>
                    <div class="activity-list">
                        <div
                            class="activity-row is-clickable"
                            v-for="(host,index) in passHostArr"
                            :key="host.host"
                            :class="{ 'is-hovered': hoveredHost === host }"
                            @click="emit('open-checks', host)"
                            @mouseenter="hoveredHost = host"
                            @mouseleave="hoveredHost = null"
                        >
                            <span class="activity-mark pass"><CircleCheck /></span>
                            <div>
                                <b>{{host.host}}</b>
                                <small>服务及资源充足稳定</small>
                            </div>
                            <time>{{host.CheckTime.slice(11,16)}}</time>
                        </div>
                    </div>

        </article>
      <article class="panel uncon-panel">
        <div class="panel-header">
          <div>
            <span class="panel-kicker">DON'T CONNECTION</span>
            <h2 style="color: red;">UNCON</h2>
          </div>
          <span class="muted">失去联系的节点</span>
        </div>
        <div class="activity-list">
          <div class="activity-row" v-for="(host,index) in unconHostArr" :key="host.host">
            <span class="activity-mark unconnection"><CircleClose /></span>
            <div>
                <b>{{host.host}}</b>
              <small>{{host.error}}</small>
            </div>
            <time>{{host.CheckTime.slice(11,16)}}</time>
          </div>
        </div>
      </article>
      <article class="panel warn-panel">
        <div class="panel-header">
          <div>
            <span class="panel-kicker">RESOURCE ALERT</span>
            <h2 style="color: #e5ae91;">WARN</h2>
          </div>
          <span class="muted">资源告警</span>
        </div>
        <div class="activity-list">
          <div
            class="activity-row is-clickable"
            v-for="(host,index) in warnHostArr"
            :key="host.host"
            :class="{ 'is-hovered': hoveredHost === host }"
            @click="emit('open-checks', host)"
            @mouseenter="hoveredHost = host"
            @mouseleave="hoveredHost = null"
          >
            <span class="activity-mark warning"><Warning /></span>
            <div>
                <b>{{host.host}}</b>
              <small>{{host.error}}</small>
            </div>
            <time>{{host.CheckTime.slice(11,16)}}</time>
          </div>
        </div>
      </article>
      <article class="panel fail-panel">
        <div class="panel-header">
          <div>
            <span class="panel-kicker">SERVICE FIAL</span>
            <h2 style="color: #c05f29;">FAIL</h2>
          </div>
          <span class="muted">服务异常</span>
        </div>
        <div class="activity-list">
          <div
            class="activity-row is-clickable"
            v-for="(host,index) in failHostArr"
            :key="host.host"
            :class="{ 'is-hovered': hoveredHost === host }"
            @click="emit('open-checks', host)"
            @mouseenter="hoveredHost = host"
            @mouseleave="hoveredHost = null"
          >
            <span class="activity-mark fail"><Close /></span>
            <div>
                <b>{{host.host}}</b>
                <small>{{getFailItem(host)}}</small>
            </div>
            <time>{{host.CheckTime.slice(11,16)}}</time>
          </div>
        </div>
      </article>
    </div>
    <article class="panel quick-panel">
      <div>
        <span class="panel-kicker">QUICK ACCESS</span>
        <h2>继续工作</h2>
      </div>
      <div class="quick-actions">
        <button class="quick-action" @click="emit('navigate', 'trends')">
          <TrendCharts />
          <span>
            <b>查看趋势分析</b>
            <small>对比 15 天资源变化</small> </span
          >→
        </button>
        <button class="quick-action" @click="emit('navigate', 'hosts', true)">
          <Plus />
          <span>
            <b>添加新的主机</b>
            <small>录入 SSH 连接配置</small> </span
          >→
        </button>
      </div>
    </article>
  </section>
</template>

<style scoped>
.activity-row.is-clickable {
  cursor: pointer;
  border-radius: 6px;
  padding-left: 10px;
  padding-right: 10px;
  transition:
    background-color 160ms ease,
    box-shadow 160ms ease,
    transform 160ms ease;
}

.activity-row.is-clickable:hover,
.activity-row.is-clickable.is-hovered {
  background: #fff7ef;
  box-shadow: inset 3px 0 0 #c05f29;
  transform: translateX(3px);
}
</style>
