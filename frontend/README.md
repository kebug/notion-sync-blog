# Notion Sync Blog Frontend

前端基于 Vue 3 + TypeScript + Naive UI 构建的现代化管理界面。

## 技术栈

- **Vue 3** - 组合式 API
- **TypeScript** - 类型安全
- **Vite** - 快速构建工具
- **Naive UI** - Vue 3 组件库
- **Vue Router** - 路由管理
- **Pinia** - 状态管理
- **Axios** - HTTP 客户端

## 项目结构

```
src/
├── api/              # API 接口
│   ├── client.ts     # HTTP 客户端
│   ├── projects.ts   # 项目相关 API
│   ├── documents.ts  # 文档相关 API
│   ├── sync.ts       # 同步相关 API
│   └── logs.ts       # 日志相关 API
├── components/       # 可复用组件
│   ├── DocumentList.vue
│   └── SyncLogs.vue
├── composables/      # 组合式函数
│   ├── useLoading.ts
│   └── useError.ts
├── router/          # 路由配置
│   └── index.ts
├── styles/          # 全局样式
│   └── global.css
├── types/           # 类型定义
│   └── index.ts
├── utils/           # 工具函数
│   ├── date.ts
│   └── format.ts
├── views/           # 页面组件
│   ├── Projects.vue
│   └── ProjectDetail.vue
├── App.vue         # 根组件
└── main.ts         # 入口文件
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发环境运行

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

### 类型检查

```bash
npm run type-check
```

## 功能特性

### 项目管理
- 创建、查看、删除项目
- 项目列表分页显示
- 响应式设计

### 文档管理
- 查看项目下的所有文档
- 文档状态显示（草稿/已发布）
- 文档内容路径显示

### 同步功能
- 手动同步数据库
- 转换为 Astro 格式
- 实时状态反馈

### 同步日志
- 查看所有同步操作记录
- 状态筛选（成功/失败/处理中）
- 错误信息显示

## 开发指南

### 组件开发

使用 `<script setup>` 语法：

```vue
<script setup lang="ts">
import { ref } from 'vue'
import type { Project } from '@/types'

const projects = ref<Project[]>([])
</script>
```

### API 调用

在 `api/` 目录下创建 API 模块：

```typescript
// api/projects.ts
import { apiClient } from './client'
import type { Project } from '@/types'

export const projectApi = {
  async getProjects() {
    return apiClient.get<Project[]>('/projects')
  },
}
```

### 样式规范

- 使用 CSS 变量
- 响应式设计（移动端优先）
- 遵循 Naive UI 设计规范

## 环境变量

```env
# .env
VITE_API_URL=http://localhost:8080/api
```

## 浏览器支持

- Chrome >= 90
- Firefox >= 88
- Safari >= 14
- Edge >= 90

## 许可证

MIT License
