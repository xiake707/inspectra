<script setup>
import {computed, ref} from 'vue'
import SectionHeading from '../components/SectionHeading.vue'
import EChart from '../components/EChart.vue'
import {mockComparison, /*mockTrend*/} from '../mock/hosts.js'

const props = defineProps({
  hostIP: {type: String, default: '10.10.10.10' },
  selectedMetric: {type: String, default: 'both'},
  trends: {type: Object, default:() => ({xAxis: [], series: []})},
});
const trendData = computed(() => ({
  xAxis: Array.isArray(props.trends?.XAxis) ? props.trends.XAxis : [],
  series: Array.isArray(props.trends?.Series)
    ? props.trends.Series.map((item) => ({
        name: item.Name,
        color: item.Color,
        data: Array.isArray(item.Datas) ? item.Datas : [],
      }))
    : [],
}));
const metric = ref(props.selectedMetric.toLowerCase() || 'both');
const days = ref(15);
const compareMetric = ref('CPU')
const selectedSeries = computed(() => trendData.value.series)
const trendOption = computed(() => {
  const series = metric.value === 'both' ? selectedSeries.value : selectedSeries.value.filter(s => s.name.toLowerCase() === metric.value);
  return {
    color: series.map(s => s.color),
    tooltip: {trigger: 'axis', valueFormatter: v => `${v}%`},
    legend: {bottom: 0, data: series.map(s => s.name)},
    grid: {left: 45, right: 20, top: 22, bottom: 38},
    xAxis: {type: 'category', data: trendData.value.xAxis.slice(-days.value)},
    yAxis: {type: 'value', max: 100, axisLabel: {formatter: '{value}%'}, splitLine: {lineStyle: {color: '#edf0f2'}}},
    series: series.map(s => ({
      name: s.name,
      type: 'line',
      smooth: true,
      data: s.data.slice(-days.value),
      symbolSize: 7
    }))
  }
})
const compareOption = computed(() => {
  const series = mockComparison.series.find(s => s.name === compareMetric.value);
  return {
    tooltip: {trigger: 'axis', valueFormatter: v => `${v}%`},
    grid: {left: 45, right: 20, top: 15, bottom: 40},
    xAxis: {type: 'category', data: mockComparison.xAxis},
    yAxis: {type: 'value', max: 100, axisLabel: {formatter: '{value}%'}},
    series: [{
      name: series.name,
      type: 'bar',
      barWidth: 48,
      itemStyle: {color: '#6ba4de', borderRadius: [5, 5, 0, 0]},
      data: series.data
    }]
  }
})
</script>
<template>
  <section class="page-view">
    <SectionHeading eyebrow="RESOURCE TRENDS"
                    :title="hostIP ? `${hostIP} · ${metric.toUpperCase()}` : '趋势分析'"
                    description="沿时间或主机维度观察资源使用率。"><span class="control-note"><i class="status-dot"></i>数据更新时间：今天 09:28</span>
    </SectionHeading>
    <article class="panel chart-panel">
      <div class="panel-header chart-header">
        <div><span class="panel-kicker">TIME SERIES</span>
          <h2>{{
              hostIP ? `${hostIP} 的 ${metric === 'both' ? '资源' : metric.toUpperCase()} 趋势` : '主机资源趋势'
            }}</h2></div>
        <div class="chart-controls"><label>主机
          <el-select :model-value="hostIP || '全部主机'" size="small">
            <el-option :label="hostIP || '全部主机'" :value="hostIP || '全部主机'"/>
          </el-select>
        </label><label>指标
          <el-select v-model="metric" size="small">
            <el-option label="CPU + Memory" value="both"/>
            <el-option label="CPU" value="cpu"/>
            <el-option label="Memory" value="memory"/>
            <el-option label="Disk" value="disk"/>
          </el-select>
        </label><label>周期
          <el-select v-model="days" size="small">
            <el-option label="15 天" :value="15"/>
            <el-option label="7 天" :value="7"/>
          </el-select>
        </label></div>
      </div>
      <div class="axis-callout"><span><b>xAxis</b> = 检查日期</span><span><b>yAxis</b>= 使用率 %</span>
      </div>
      <EChart :option="trendOption"/>
    </article>
    <article class="panel chart-panel">
      <div class="panel-header">
        <div><span class="panel-kicker">HOST COMPARISON</span>
          <h2>主机指标对比</h2></div>
        <el-select v-model="compareMetric" size="small">
          <el-option v-for="item in ['CPU', 'Memory', 'Disk']" :key="item" :label="`${item} 使用率`" :value="item"/>
        </el-select>
      </div>
      <div class="axis-callout">
        <span><b>xAxis</b> = <code>mockComparison.xAxis</code> · 主机地址</span><span><b>yAxis</b> = <code>mockComparison.series[].data</code> · 指标 %</span>
      </div>
      <EChart :option="compareOption" height="220px"/>
    </article>
  </section>
</template>
