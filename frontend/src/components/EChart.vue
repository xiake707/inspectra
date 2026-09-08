<script setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue'
import * as echarts from 'echarts'

const props = defineProps({option: {type: Object, required: true}, height: {type: String, default: '320px'}})
const chartRef = ref(null);
let chart

function render() {
  if (!chart) return;
  chart.setOption(props.option, true);
  chart.resize()
}

function resize() {
  chart?.resize()
}

onMounted(async () => {
  await nextTick();
  chart = echarts.init(chartRef.value);
  render();
  window.addEventListener('resize', resize)
})
watch(() => props.option, render, {deep: true})
onBeforeUnmount(() => {
  window.removeEventListener('resize', resize);
  chart?.dispose()
})
</script>
<template>
  <div ref="chartRef" class="echart" :style="{ height }"></div>
</template>
