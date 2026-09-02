<template>
  <div class="min-h-screen bg-gray-50 text-gray-950 dark:bg-dark-950 dark:text-white">
    <header class="border-b border-gray-200 bg-white/95 dark:border-dark-800 dark:bg-dark-900/95">
      <div class="mx-auto flex max-w-[1440px] items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
        <RouterLink to="/home" class="flex items-center gap-3 text-lg font-semibold">
          <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-9 w-9 rounded-lg object-contain" />
          <span>{{ siteName }}</span>
        </RouterLink>
      </div>
    </header>

    <main style="max-width: calc(1440px + 3cm)" class="mx-auto grid gap-8 px-4 py-8 sm:px-6 lg:grid-cols-[calc(220px+1.5cm)_minmax(0,1fr)_220px] lg:px-8">
      <aside class="hidden border-r border-gray-200 pr-6 dark:border-dark-700 lg:block">
        <nav class="sticky top-6 space-y-7" aria-label="文档导航">
          <section v-for="group in groups" :key="group.title">
            <h2 class="mb-3 px-3 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">
              {{ group.title }}
            </h2>
            <div class="space-y-1">
              <RouterLink
                v-for="item in group.items"
                :key="item.slug"
                :to="{ name: 'DocsPage', params: { slug: item.slug } }"
                class="block rounded-xl px-4 py-2 text-sm transition-colors"
                :class="item.slug === currentSlug ? 'bg-gray-900 font-semibold text-white shadow-sm dark:bg-white dark:text-gray-900' : 'text-gray-600 hover:bg-gray-100 dark:text-dark-300 dark:hover:bg-dark-800'"
              >{{ item.title }}</RouterLink>
            </div>
          </section>
        </nav>
      </aside>

      <section class="min-w-0">
        <div class="mb-6 lg:hidden">
          <label for="docs-mobile-nav" class="sr-only">文档导航</label>
          <select id="docs-mobile-nav" class="form-input w-full" :value="currentSlug" @change="navigateFromSelect">
            <option v-for="item in allItems" :key="item.slug" :value="item.slug">{{ item.title }}</option>
          </select>
        </div>

        <div v-if="loading" class="card flex min-h-[420px] items-center justify-center">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
        </div>
        <article v-else-if="markdown" style="width: calc(100% + 4.5cm); max-width: none" class="docs-article card px-6 py-8 sm:px-10 lg:px-12">
          <div v-html="renderedHtml" />
        </article>
        <div v-else class="card p-10 text-center text-gray-500 dark:text-dark-300">文档不存在</div>
      </section>

      <aside v-if="tocItems.length" class="hidden xl:block">
        <nav class="sticky top-6 border-l border-gray-200 pl-5 dark:border-dark-700" aria-label="页面目录">
          <p class="mb-3 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">本页目录</p>
          <a v-for="item in tocItems" :key="item.id" :href="`#${item.id}`" class="mb-2 block text-sm text-gray-500 transition-colors hover:text-primary-600 dark:text-dark-400 dark:hover:text-primary-300" :class="item.level > 2 ? 'pl-3' : ''">
            {{ item.text }}
          </a>
        </nav>
      </aside>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

interface DocItem { slug: string; title: string }
interface DocGroup { title: string; items: DocItem[] }
interface TocItem { id: string; text: string; level: number }

const groups: DocGroup[] = [
  { title: '入门', items: [{ slug: 'introduction', title: '简介' }, { slug: 'faq', title: '常见问题' }, { slug: 'dragon-code-guide', title: 'Dragon Code 使用指南' }] },
  { title: '快速开始', items: [{ slug: 'claude-code-quickstart', title: 'Claude Code 快速开始指南' }, { slug: 'codex-quickstart', title: 'Codex 快速开始指南' }, { slug: 'openclaw-quickstart', title: 'OpenClaw 快速开始指南' }] },
]
const allItems = groups.flatMap((group) => group.items)
const modules = import.meta.glob('@/content/docs/*.md', { query: '?raw', import: 'default' }) as Record<string, () => Promise<string>>
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const markdown = ref('')
const renderedHtml = ref('')
const tocItems = ref<TocItem[]>([])
const loading = ref(true)
const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => settings.value?.site_name || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(settings.value?.site_logo || '', { allowRelative: true, allowDataUrl: true }))
const currentSlug = computed(() => String(route.params.slug || 'introduction'))

marked.setOptions({ breaks: true, gfm: true })

function slugify(text: string, index: number) {
  const value = text.toLowerCase().replace(/[^\w一-鿿]+/g, '-').replace(/^-+|-+$/g, '')
  return value ? `${value}-${index}` : `heading-${index}`
}

async function loadDocument() {
  loading.value = true
  const loader = modules[`/src/content/docs/${currentSlug.value}.md`] || modules[`../content/docs/${currentSlug.value}.md`]
  markdown.value = loader ? await loader() : ''
  const headings: TocItem[] = []
  let headingIndex = 0
  const renderer = new marked.Renderer()
  renderer.heading = ({ text, depth }) => {
    const id = slugify(text, headingIndex++)
    headings.push({ id, text, level: depth })
    return `<h${depth} id="${id}">${text}</h${depth}>`
  }
  renderedHtml.value = DOMPurify.sanitize(marked.parse(markdown.value, { renderer }) as string)
  tocItems.value = headings.filter((item) => item.level >= 2 && item.level <= 3)
  loading.value = false
  await nextTick()
}

function navigateFromSelect(event: Event) {
  const slug = (event.target as HTMLSelectElement).value
  router.push({ name: 'DocsPage', params: { slug } })
}

watch(currentSlug, loadDocument, { immediate: true })
</script>

<style scoped>
.docs-article { line-height: 1.8; }
.docs-article :deep(h1) { @apply mb-6 border-b border-gray-200 pb-5 text-3xl font-bold dark:border-dark-700; }
.docs-article :deep(h2) { @apply mb-3 mt-10 text-2xl font-bold; }
.docs-article :deep(h3) { @apply mb-2 mt-8 text-xl font-semibold; }
.docs-article :deep(p) { @apply mb-4 text-gray-700 dark:text-dark-200; }
.docs-article :deep(ul), .docs-article :deep(ol) { @apply mb-5 pl-6 text-gray-700 dark:text-dark-200; }
.docs-article :deep(li) { @apply mb-2; }
.docs-article :deep(a) { @apply text-primary-600 underline underline-offset-4 dark:text-primary-300; }
.docs-article :deep(pre) { @apply my-5 overflow-x-auto rounded-xl bg-gray-950 p-4 text-sm text-gray-100; }
.docs-article :deep(code) { @apply rounded bg-gray-100 px-1.5 py-0.5 text-[0.9em] dark:bg-dark-800; }
.docs-article :deep(pre code) { @apply bg-transparent p-0; }
</style>
