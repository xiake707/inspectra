<script setup>
import {computed, ref} from 'vue'
import {Search, Plus, MoreFilled, VideoPlay} from '@element-plus/icons-vue'
import SectionHeading from '../components/SectionHeading.vue'

const props = defineProps({
  hosts: {type: Array, default: () => []},
  checking: {type: Boolean, default: false},
});
const emit = defineEmits(['add', 'start-check']);
const query = ref('')
const filteredHosts = computed(() => props.hosts.filter(h => `${h.Host || ''}`.toLowerCase().includes(query.value.toLowerCase())))
</script>
<template>
  <section class="page-view">
    <SectionHeading eyebrow="HOST CONFIGURATION" title="主机管理" description="维护巡检目标与 SSH 连接配置。">
      <el-button type="primary" :icon="Plus" @click="emit('add')">添加主机</el-button>
    </SectionHeading>
    <article class="panel check-action-panel">
      <div class="check-action-copy">
        <span class="panel-kicker">HOST INSPECTION</span>
        <h2>开始检查</h2>
        <p>检查已配置主机的资源、服务和网络状态，并更新最新巡检记录。</p>
      </div>
      <el-button
        type="primary"
        :icon="VideoPlay"
        :loading="checking"
        :disabled="checking"
        @click="emit('start-check')"
      >
        {{ checking ? '检查中...' : '开始检查' }}
      </el-button>
    </article>
    <article class="panel hosts-panel">
      <div class="panel-header">
        <div><span class="panel-kicker">HOST_CONFIG</span>
          <h2>已配置主机
            <el-tag size="small" round>{{ hosts.length }}</el-tag>
          </h2>
        </div>
        <el-input v-model="query" :prefix-icon="Search" placeholder="搜索主机地址" clearable style="width: 220px"/>
      </div>
      <el-table :data="filteredHosts" class="host-table" empty-text="没有匹配的主机">
        <el-table-column label="主机地址" min-width="180">
          <template #default="{ row }"><b>{{ row.Host }}</b><small>Linux server</small></template>
        </el-table-column>
        <el-table-column prop="User" label="SSH 用户" width="120"/>
        <el-table-column prop="Port" label="端口" width="90"/>
        <el-table-column label="认证方式" width="130">
          <template #default="{ row }">{{ row.KeyFile ? 'SSH 私钥' : '密码' }}</template>
        </el-table-column>
        <el-table-column label="最近状态" width="130">
          <template #default="{ row }">
            <el-tag :type="row.status === 'PASS' ? 'success' : row.status === 'WARN' ? 'warning' : 'info'"
                    effect="light">{{ row.status || '待巡检' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="最近巡检" width="130">
          <template #default="{ row }">{{ row.checkedAt || '暂无记录' }}</template>
        </el-table-column>
        <el-table-column label="" width="55">
          <template #default>
            <el-button text :icon="MoreFilled" aria-label="更多操作"/>
          </template>
        </el-table-column>
      </el-table>
    </article>
    <div class="form-note">
      <el-icon>
        <InfoFilled/>
      </el-icon>
      <div><b>数据尚未连接</b>
        <p>当前列表来自 mock 数据。接入后端时，`GET /api/hosts` 的返回值会自动填充到这里。</p></div>
    </div>
  </section>
</template>
