<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <router-link
            to="/docs"
            class="flex h-10 shrink-0 items-center justify-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
            <span class="hidden sm:inline">{{ t('home.docs') }}</span>
          </router-link>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <component
            :is="isAuthenticated ? 'router-link' : 'button'"
            :to="isAuthenticated ? dashboardPath : undefined"
            type="button"
            @click="handleLoginClick"
            data-auth-entry
            :aria-haspopup="isAuthenticated ? undefined : 'dialog'"
            aria-controls="auth-modal"
            class="home-auth-entry inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </component>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <component
          :is="isAuthenticated ? 'router-link' : 'button'"
          :to="isAuthenticated ? dashboardPath : undefined"
          type="button"
          @click="handleLoginClick"
          data-auth-entry
          :aria-haspopup="isAuthenticated ? undefined : 'dialog'"
          aria-controls="auth-modal"
          class="home-auth-entry mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </component>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-hidden bg-white dark:bg-gradient-to-br dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
  >

    <!-- Header -->
    <header class="relative z-20 border-b border-gray-200/70 px-4 py-4 dark:border-dark-800/70 sm:px-6">
      <nav class="mx-auto grid max-w-6xl grid-cols-[minmax(0,1fr)_auto] items-center gap-4 md:grid-cols-[minmax(0,1fr)_auto_minmax(0,1fr)]">
        <!-- Logo -->
        <router-link to="/" class="flex min-w-0 items-center gap-3" :aria-label="t('home.navHome')">
          <div class="h-10 w-10 overflow-hidden rounded-xl">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="truncate text-xl font-semibold tracking-tight text-gray-950 dark:text-white sm:text-2xl">{{ siteName }}</span>
        </router-link>

        <!-- Primary navigation -->
        <div class="hidden items-center justify-center gap-9 md:flex" :aria-label="t('home.primaryNavigation')">
          <router-link
            to="/"
            class="rounded-md px-1 py-2 text-base font-semibold text-gray-950 transition-colors hover:text-primary-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 dark:text-white dark:hover:text-primary-300"
            aria-current="page"
          >
            {{ t('home.navHome') }}
          </router-link>
          <router-link
            to="/docs"
            class="rounded-md px-1 py-2 text-base font-medium text-gray-600 transition-colors hover:text-gray-950 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 dark:text-dark-300 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </router-link>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center justify-end gap-2 sm:gap-3">
          <router-link
            to="/docs"
            class="rounded-md px-2 py-2 text-sm font-medium text-gray-600 transition-colors hover:text-gray-950 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 dark:text-dark-300 dark:hover:text-white md:hidden"
          >
            {{ t('home.docs') }}
          </router-link>
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Announcement Bell -->
          <AnnouncementBell v-if="isAuthenticated" />

          <!-- Model Plaza Link -->
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex items-center gap-1.5 rounded-lg p-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="group inline-flex items-center justify-center rounded-full border border-gray-900 bg-gray-50 px-5 py-1.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-900 hover:text-white dark:border-gray-200 dark:bg-dark-900 dark:text-gray-100 dark:hover:bg-gray-100 dark:hover:text-gray-900"
          >
            <span class="font-semibold text-gray-900 transition-colors group-hover:text-white dark:text-gray-100 dark:group-hover:text-gray-900">{{ t('home.dashboard') }}</span>
          </router-link>
          <button
            v-else
            type="button"
            @click="handleLoginClick"
            data-auth-entry
            aria-haspopup="dialog"
            aria-controls="auth-modal"
            class="home-auth-entry group inline-flex items-center justify-center rounded-full border border-gray-900 bg-gray-50 px-5 py-1.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-900 hover:text-white dark:border-gray-200 dark:bg-dark-900 dark:text-gray-100 dark:hover:bg-gray-100 dark:hover:text-gray-900"
          >
            <span class="font-semibold text-gray-900 transition-colors group-hover:text-white dark:text-gray-100 dark:group-hover:text-gray-900">
              {{ t('home.login') }}
            </span>
          </button>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content -->
          <div class="flex-1 pt-12 text-center">
            <div class="mb-8 flex flex-wrap items-center justify-center gap-4 text-black">
              <span class="rounded-full border-2 border-black px-6 py-2 text-2xl font-medium md:text-3xl">
                {{ t('home.heroLabel') }}
              </span>
              <span class="hero-title-animated text-5xl font-bold tracking-tight md:text-7xl">
                {{ t('home.heroTitle') }}
              </span>
              <p class="mt-4 w-full text-center text-base font-normal leading-relaxed text-gray-700 md:text-xl">
                <span class="block">{{ t('home.heroDescriptionLine1') }}</span>
                <span class="block">{{ t('home.heroDescriptionLine2') }}</span>
              </p>
            </div>

            <!-- CTA Button -->
            <div>
          <component
            :is="isAuthenticated ? 'router-link' : 'button'"
            :to="isAuthenticated ? dashboardPath : undefined"
            type="button"
            @click="handleLoginClick"
                data-auth-entry
                :aria-haspopup="isAuthenticated ? undefined : 'dialog'"
                aria-controls="auth-modal"
                class="home-auth-entry btn bg-black px-8 py-3 text-base text-white shadow-lg shadow-black/20 hover:bg-gray-800 hover:shadow-black/30"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
          </component>
            </div>
          </div>

        </div>

        <!-- Feature Tags - Centered -->
        <div class="mb-12 flex flex-wrap items-center justify-center gap-4 md:gap-6">
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="swap" size="sm" class="text-black" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.subscriptionToApi')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="shield" size="sm" class="text-black" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.stickySession')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="chart" size="sm" class="text-black" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.realtimeBilling')
            }}</span>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="providers-marquee mb-16" aria-label="Supported AI models">
          <div class="providers-track">
            <div
              v-for="group in 3"
              :key="group"
              class="providers-group"
              :aria-hidden="group > 1"
            >
              <div
                v-for="provider in supportedProviders"
                :key="provider.id"
                class="provider-item flex shrink-0 items-center gap-3 rounded-xl border border-gray-200/70 bg-white/70 px-5 py-3 backdrop-blur-sm dark:border-dark-700/70 dark:bg-dark-800/70"
              >
                <img :src="provider.logo" :alt="group === 1 ? `${provider.name} logo` : ''" class="h-8 w-8 object-contain" loading="lazy" />
                <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ provider.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import AnnouncementBell from '@/components/common/AnnouncementBell.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { requestLoginModal } from '@/utils/loginModal'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

const supportedProviders = [
  { id: 'openai', name: 'OpenAI', logo: '/provider-logos/openai-custom.jpg' },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    logo: '/provider-logos/deepseek-custom.jpg',
  },
  {
    id: 'qwen',
    name: '通义千问',
    logo: '/provider-logos/qwen-custom.png',
  },
  { id: 'kimi', name: 'Kimi', logo: '/provider-logos/kimi.ico' },
]
// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

function handleLoginClick(event: MouseEvent) {
  if (!isAuthenticated.value) {
    event.preventDefault()
    if (event.currentTarget instanceof HTMLElement) {
      event.currentTarget.blur()
    }
    requestLoginModal()
  }
}
// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.home-auth-entry {
  @apply focus:outline-none focus-visible:ring-2 focus-visible:ring-black/60 focus-visible:ring-offset-2 dark:focus-visible:ring-white/70 dark:focus-visible:ring-offset-dark-950;
}

.hero-title-animated {
  color: transparent;
  background: linear-gradient(
    100deg,
    #111827 0%,
    #111827 38%,
    #4f8cff 50%,
    #111827 62%,
    #111827 100%
  );
  background-size: 240% 100%;
  background-position: 120% 0;
  background-clip: text;
  -webkit-background-clip: text;
  animation: hero-title-highlight 8s ease-in-out infinite;
}

.providers-marquee {
  position: relative;
  overflow: hidden;
  mask-image: linear-gradient(90deg, transparent, #000 12%, #000 88%, transparent);
  -webkit-mask-image: linear-gradient(90deg, transparent, #000 12%, #000 88%, transparent);
}

.providers-track {
  display: flex;
  width: max-content;
  animation: providers-marquee 38s linear infinite;
}

.providers-group {
  display: flex;
  gap: 3.5rem;
  padding-right: 3.5rem;
}

.providers-marquee:hover .providers-track {
  animation-play-state: paused;
}

@keyframes providers-marquee {
  from {
    transform: translateX(0);
  }

  to {
    transform: translateX(-33.333333%);
  }
}

:global(.dark) .hero-title-animated {
  background-image: linear-gradient(
    100deg,
    #f8fafc 0%,
    #f8fafc 38%,
    #70a7ff 50%,
    #f8fafc 62%,
    #f8fafc 100%
  );
}

@keyframes hero-title-highlight {
  from {
    background-position: 120% 0;
  }

  62.5% {
    background-position: -20% 0;
  }

  to {
    background-position: -20% 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .hero-title-animated {
    color: #111827;
    background: none;
    animation: none;
  }

  :global(.dark) .hero-title-animated {
    color: #f8fafc;
  }

  .providers-track {
    animation: none;
  }
}
</style>
