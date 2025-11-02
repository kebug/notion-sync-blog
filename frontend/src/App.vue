<template>
  <n-config-provider>
    <n-layout>
      <n-layout-header bordered>
        <n-space justify="space-between" style="padding: 12px 24px;">
          <n-h2 style="margin: 0;">Notion Sync Blog</n-h2>
          <n-space>
            <n-button @click="$router.push('/')">首页</n-button>
            <n-button @click="$router.push('/projects')">项目</n-button>
          </n-space>
        </n-space>
      </n-layout-header>
      <n-layout-content style="padding: 24px;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NLayout, NLayoutHeader, NLayoutContent, NSpace, NButton, NH2 } from 'naive-ui'
</script>

