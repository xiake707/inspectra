<script setup>
defineProps({eyebrow: String, title: String, description: String})
</script>
<template>
  <div class="page-heading">
    <div><p class="eyebrow">{{ eyebrow }}</p>
      <h1>{{ title }}</h1>
      <p class="subheading">{{ description }}</p></div>
    <slot/>
  </div>
</template>
