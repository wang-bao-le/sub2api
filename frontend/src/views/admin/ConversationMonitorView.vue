<template>
  <AppLayout>
    <div class="space-y-6">
      <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <h1 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">{{ t('admin.conversations.title') }}</h1>
          <p class="mt-1 max-w-2xl text-sm leading-6 text-gray-500 dark:text-gray-400">{{ t('admin.conversations.description') }} 内容可能包含用户隐私和敏感数据。</p>
        </div>
        <button class="btn btn-secondary inline-flex items-center gap-2" :disabled="loading" @click="load"><Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />刷新</button>
      </div>

      <section class="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-800 sm:p-5">
        <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
          <div class="flex items-center gap-3"><span class="flex h-9 w-9 items-center justify-center rounded-lg bg-primary-50 text-primary-600 dark:bg-primary-500/10 dark:text-primary-400"><Icon name="shield" size="sm" /></span><div><h2 class="font-medium text-gray-900 dark:text-white">采集开关</h2><p class="text-xs text-gray-500 dark:text-gray-400">默认关闭，开启后覆盖全部用户和分组。</p></div></div>
          <div class="flex flex-wrap items-center gap-3"><span class="text-sm text-gray-600 dark:text-gray-300">{{ config.enabled ? '已启用' : '已停用' }}</span><Toggle v-if="configLoaded" v-model="config.enabled" /><button class="btn btn-primary" :disabled="configSaving" @click="saveConfig">{{ configSaving ? '保存中…' : '保存配置' }}</button></div>
        </div>
      </section>

      <section class="card overflow-hidden">
        <div class="flex flex-col gap-3 border-b border-gray-100 p-4 dark:border-dark-700 sm:flex-row sm:items-center sm:p-5"><div class="relative min-w-0 flex-1"><Icon name="search" size="sm" class="absolute left-3 top-2.5 text-gray-400" /><input v-model="keyword" class="form-input w-full pl-9" placeholder="搜索用户、模型、Prompt 或回答" @keyup.enter="applySearch" /></div><select v-model="status" class="form-input sm:w-36"><option value="">全部状态</option><option value="200">成功</option><option value="failed">失败</option></select><button class="btn btn-primary" @click="applySearch">筛选</button><button class="btn btn-secondary" @click="resetSearch">清空</button></div>
        <div v-if="error" class="m-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-500/10 dark:text-red-300">加载失败，请刷新重试。</div>
        <div v-else-if="loading" class="flex items-center justify-center py-20 text-sm text-gray-500">加载中…</div>
        <div v-else-if="items.length === 0" class="flex flex-col items-center justify-center px-6 py-20 text-center"><Icon name="search" size="lg" class="text-gray-300 dark:text-dark-500" /><p class="mt-3 text-sm font-medium text-gray-700 dark:text-gray-200">暂无对话记录</p><p class="mt-1 text-xs text-gray-500">请先开启采集，或调整筛选条件。</p></div>
        <div v-else class="overflow-x-auto"><table class="min-w-full text-left text-sm"><thead class="bg-gray-50 text-xs uppercase tracking-wide text-gray-500 dark:bg-dark-900/40 dark:text-gray-400"><tr><th class="px-5 py-3">时间 / 用户</th><th class="px-5 py-3">平台 / 模型</th><th class="px-5 py-3">Prompt</th><th class="px-5 py-3">回答</th><th class="px-5 py-3">状态</th><th class="px-5 py-3">操作</th></tr></thead><tbody class="divide-y divide-gray-100 dark:divide-dark-700"><tr v-for="item in items" :key="item.id" class="align-top hover:bg-gray-50/70 dark:hover:bg-dark-700/30"><td class="whitespace-nowrap px-5 py-4"><div class="font-medium text-gray-900 dark:text-white">{{ item.username || item.user_email || `用户 #${item.user_id}` }}</div><div class="mt-1 text-xs text-gray-500">{{ formatTime(item.created_at) }}</div></td><td class="px-5 py-4"><div class="font-medium text-gray-800 dark:text-gray-200">{{ item.requested_model || '-' }}</div><div class="mt-1 text-xs text-gray-500">{{ item.provider || '-' }} · {{ item.stream ? 'SSE' : 'HTTP' }}</div></td><td class="max-w-[280px] px-5 py-4"><p class="line-clamp-3 whitespace-pre-wrap text-gray-600 dark:text-gray-300">{{ item.prompt_content || '-' }}</p></td><td class="max-w-[280px] px-5 py-4"><p class="line-clamp-3 whitespace-pre-wrap text-gray-600 dark:text-gray-300">{{ item.response_content || '-' }}</p><span v-if="item.truncated" class="mt-1 inline-block text-xs text-amber-600">内容已截断</span></td><td class="whitespace-nowrap px-5 py-4"><span class="rounded-full px-2.5 py-1 text-xs font-medium" :class="item.complete ? 'bg-green-50 text-green-700 dark:bg-green-500/10 dark:text-green-300' : 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-300'">{{ item.complete ? '成功' : '未完成' }}</span><div class="mt-2 text-xs text-gray-500">{{ item.duration_ms }} ms</div></td><td class="whitespace-nowrap px-5 py-4"><button class="text-primary-600 hover:text-primary-700 dark:text-primary-400" @click="openDetail(item.id)">查看详情</button><button v-if="config.manual_delete_enabled" class="ml-3 text-red-600 hover:text-red-700" @click="remove(item.id)">删除</button></td></tr></tbody></table></div>
        <div v-if="total > pageSize" class="flex items-center justify-between border-t border-gray-100 px-5 py-3 text-sm dark:border-dark-700"><span class="text-gray-500">共 {{ total }} 条</span><div class="flex items-center gap-2"><button class="btn btn-secondary px-3 py-1.5" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button><span class="text-gray-600 dark:text-gray-300">{{ page }}</span><button class="btn btn-secondary px-3 py-1.5" :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button></div></div>
      </section>
    </div>
    <BaseDialog v-model:show="detailOpen" title="对话详情" size="xl"><div v-if="detail" class="space-y-5"><p class="rounded-lg bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:bg-amber-500/10 dark:text-amber-200">内容可能包含用户隐私和敏感数据，请妥善处理。凭证字段不会作为监控内容保存。</p><div class="grid gap-4 sm:grid-cols-2"><div><p class="text-xs text-gray-500">用户</p><p class="mt-1 text-sm text-gray-900 dark:text-white">{{ detail.username || detail.user_email || `#${detail.user_id}` }}</p></div><div><p class="text-xs text-gray-500">请求</p><p class="mt-1 break-all text-sm text-gray-900 dark:text-white">{{ detail.provider }} · {{ detail.requested_model }} · {{ detail.endpoint }}</p></div></div><div><h3 class="mb-2 text-sm font-semibold text-gray-900 dark:text-white">用户 Prompt</h3><pre class="max-h-72 overflow-auto whitespace-pre-wrap rounded-lg bg-gray-50 p-4 text-sm leading-6 text-gray-700 dark:bg-dark-900 dark:text-gray-200">{{ detail.prompt_content || '-' }}</pre></div><div><h3 class="mb-2 text-sm font-semibold text-gray-900 dark:text-white">模型回答</h3><pre class="max-h-96 overflow-auto whitespace-pre-wrap rounded-lg bg-gray-50 p-4 text-sm leading-6 text-gray-700 dark:bg-dark-900 dark:text-gray-200">{{ detail.response_content || '-' }}</pre></div></div><div v-else class="py-12 text-center text-sm text-gray-500">加载中…</div></BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import Toggle from '@/components/common/Toggle.vue'
import { deleteConversation, getConversation, getConversationConfig, listConversations, updateConversationConfig, type ConversationConfig, type ConversationRecord } from '@/features/conversations/api'
const { t } = useI18n(); const items=ref<ConversationRecord[]>([]); const detail=ref<ConversationRecord|null>(null); const detailOpen=ref(false); const loading=ref(false); const error=ref(false); const configSaving=ref(false); const configLoaded=ref(false); const keyword=ref(''); const status=ref(''); const page=ref(1); const pageSize=20; const total=ref(0); const config=ref<ConversationConfig>({enabled:false,capture_all_groups:true,max_prompt_bytes:65536,max_response_bytes:262144,manual_delete_enabled:true})
async function load(){ loading.value=true; error.value=false; try { const result=await listConversations({page:page.value,page_size:pageSize,keyword:keyword.value,status:status.value}); items.value=result.items; total.value=result.total } catch { error.value=true } finally { loading.value=false } }
async function loadConfig(){ try { config.value=await getConversationConfig() } catch {} finally { configLoaded.value=true } }
async function saveConfig(){ configSaving.value=true; try { config.value=await updateConversationConfig(config.value) } finally { configSaving.value=false } }
function applySearch(){ page.value=1; load() }; function resetSearch(){ keyword.value=''; status.value=''; applySearch() }; function changePage(value:number){ page.value=value; load() }
async function openDetail(id:number){ detailOpen.value=true; detail.value=null; try { detail.value=await getConversation(id) } catch { detailOpen.value=false } }
async function remove(id:number){ if (!window.confirm('删除这条对话记录？此操作不可恢复。')) return; await deleteConversation(id); await load() }
function formatTime(value:string){ return new Date(value).toLocaleString() }
onMounted(()=>{ loadConfig(); load() })
</script>
