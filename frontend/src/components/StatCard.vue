<script setup>
defineProps({
  label: String,
  value: [String, Number],
  detail: String,
  tone: {type: String, default: 'blue'},
  icon: Object
})
</script>
<template>
  <article class="stat-card">
    <div class="stat-top">
        <span>{{ label }}</span>
        <span :class="['stat-icon', tone]">
            <component :is="icon"/>
        </span>
    </div>
    <strong>{{ value }}</strong><small>{{ detail }}</small></article>
</template>
