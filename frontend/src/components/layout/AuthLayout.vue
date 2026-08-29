<template>
  <div
    :class="embedded
      ? 'relative w-full overflow-hidden bg-transparent'
      : 'relative grid min-h-[100dvh] overflow-hidden bg-gray-50 dark:bg-dark-950 lg:grid-cols-[42%_58%]'"
  >
    <!-- Brand panel -->
    <aside
      :class="embedded
        ? 'relative flex min-h-[7.5rem] items-end overflow-hidden bg-gradient-to-br from-primary-600 to-primary-500 px-6 py-5 text-white sm:px-8'
        : 'relative flex min-h-[9.5rem] items-end overflow-hidden bg-gray-950 px-6 py-7 text-white dark:bg-black sm:px-10 lg:min-h-[100dvh] lg:items-center lg:px-14 lg:py-12'"
    >
      <div class="relative z-10 w-full max-w-lg">
        <template v-if="settingsLoaded">
          <div :class="embedded ? 'flex items-center gap-3' : 'flex items-center gap-3 lg:block'">
            <div
              :class="embedded
                ? 'inline-flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white/15 ring-1 ring-white/25'
                : 'inline-flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/15 lg:mb-8 lg:h-16 lg:w-16 lg:rounded-2xl'"
            >
              <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
            </div>
            <div>
              <h1 :class="embedded ? 'text-xl font-semibold tracking-tight text-white' : 'text-xl font-semibold tracking-tight text-white lg:text-5xl lg:font-bold'">
                {{ siteName }}
              </h1>
              <p :class="embedded ? 'mt-1 text-xs text-white/75' : 'mt-1 text-xs text-gray-400 lg:mt-3 lg:max-w-sm lg:text-base lg:leading-7'">
                {{ siteSubtitle }}
              </p>
            </div>
          </div>
          <div v-if="!embedded" class="mt-6 hidden h-px w-24 bg-white/30 lg:block"></div>
        </template>
      </div>
    </aside>

    <!-- Authentication surface -->
    <main :class="embedded ? 'relative flex min-w-0 flex-col items-center' : 'relative flex min-w-0 flex-col items-center justify-center px-4 py-8 sm:px-8 lg:px-12 lg:py-10'">
      <div class="relative z-10 w-full max-w-md">
        <div
          class="w-full rounded-2xl border border-gray-200 bg-white shadow-lg dark:border-dark-700 dark:bg-dark-800"
          :class="embedded || compact ? 'p-5' : 'p-8'"
        >
          <slot />
          <div class="mt-4 text-center text-sm">
            <slot name="footer" />
          </div>
        </div>

        <div v-if="!embedded" class="mt-5 text-center text-xs text-gray-400 dark:text-dark-500">
          &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'

const props = withDefaults(defineProps<{ compact?: boolean; embedded?: boolean }>(), {
  compact: false
})

const embedded = computed(() => props.embedded === true)

const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform')
const settingsLoaded = computed(() => appStore.publicSettingsLoaded)

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>
