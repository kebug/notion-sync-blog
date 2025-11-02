<template>
  <n-card>
    <template #header>
      <n-h2>欢迎使用 Notion Sync Blog</n-h2>
    </template>
    <n-p>
      这是一个将 Notion 内容同步并转换为静态博客的工具。
    </n-p>
    <n-space>
      <n-button type="primary" @click="$router.push('/projects')">
        开始使用
      </n-button>
    </n-space>
  </n-card>
</template>

<script setup lang="ts">
import { NCard, NH2, NP, NSpace, NButton } from 'naive-ui'
</script>

