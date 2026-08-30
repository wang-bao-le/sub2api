<script setup lang="ts">
import { RouterView, useRouter, useRoute } from 'vue-router'
import { onMounted, onBeforeUnmount, watch, ref } from 'vue'
import Toast from '@/components/common/Toast.vue'
import NavigationProgress from '@/components/common/NavigationProgress.vue'
import AdminComplianceDialog from '@/components/admin/AdminComplianceDialog.vue'
import { resolveRouteDocumentTitle } from '@/router/title'
import AnnouncementPopup from '@/components/common/AnnouncementPopup.vue'
import { useAppStore, useAuthStore, useSubscriptionStore, useAnnouncementStore, useAdminComplianceStore, useAdminSettingsStore } from '@/stores'
import { getSetupStatus } from '@/api/setup'
import AuthModal from '@/components/auth/AuthModal.vue'
import { AUTH_MODAL_EVENT, consumeAuthModalRequest, type AuthModalOptions, type AuthModalView } from '@/utils/loginModal'

const showAuthModal = ref(false)
const authModalView = ref<AuthModalView>('login')
const authModalOptions = ref<AuthModalOptions>({})
const router = useRouter()
const route = useRoute()

function openAuthModal(event: Event) {
  const detail = (event as CustomEvent<{ view?: AuthModalView } & AuthModalOptions>).detail || {}
  if (authStore.isAuthenticated && (detail.view || 'login') === 'login') return
  authModalView.value = detail.view || 'login'
  if (detail.redirect !== undefined || detail.token !== undefined || detail.email !== undefined || detail.query !== undefined || !showAuthModal.value) {
    authModalOptions.value = {
      redirect: detail.redirect,
      token: detail.token,
      email: detail.email,
      query: detail.query,
    }
  }
  showAuthModal.value = true
}

function closeAuthModal() {
  showAuthModal.value = false
  authModalOptions.value = {}
}

watch(() => route.fullPath, () => {
  if (showAuthModal.value) closeAuthModal()
})

const appStore = useAppStore()
const authStore = useAuthStore()
const subscriptionStore = useSubscriptionStore()
const announcementStore = useAnnouncementStore()
const adminComplianceStore = useAdminComplianceStore()
const adminSettingsStore = useAdminSettingsStore()

function updateDocumentTitle() {
  const customMenuItems = [
    ...(appStore.cachedPublicSettings?.custom_menu_items ?? []),
    ...(authStore.isAdmin ? adminSettingsStore.customMenuItems : []),
  ]
  document.title = resolveRouteDocumentTitle(route, appStore.siteName, customMenuItems)
}

watch(
  [
    () => route.fullPath,
    () => route.meta.title,
    () => route.meta.titleKey,
    () => appStore.siteName,
    () => appStore.cachedPublicSettings?.custom_menu_items,
    () => authStore.isAdmin,
    () => adminSettingsStore.customMenuItems,
  ],
  updateDocumentTitle,
  { deep: true }
)

// Watch for authentication state and manage subscription data + announcements
function onVisibilityChange() {
  if (document.visibilityState === 'visible' && authStore.isAuthenticated) {
    announcementStore.fetchAnnouncements()
  }
}

function onAdminComplianceRequired(event: Event) {
  const detail = (event as CustomEvent<Record<string, string>>).detail || {}
  adminComplianceStore.requireAcknowledgement(detail)
}

watch(
  () => authStore.isAuthenticated,
  (isAuthenticated, oldValue) => {
    if (isAuthenticated) {
      if (authStore.isAdmin) {
        adminComplianceStore.fetchStatus().catch((error) => {
          console.error('Failed to fetch admin compliance status:', error)
        })
      }

      // User logged in: preload subscriptions and start polling
      subscriptionStore.fetchActiveSubscriptions().catch((error) => {
        console.error('Failed to preload subscriptions:', error)
      })
      subscriptionStore.startPolling()

      // Announcements: new login vs page refresh restore
      if (oldValue === false) {
        // New login: delay 3s then force fetch
        setTimeout(() => announcementStore.fetchAnnouncements(true), 3000)
      } else {
        // Page refresh restore (oldValue was undefined)
        announcementStore.fetchAnnouncements()
      }

      // Register visibility change listener
      document.addEventListener('visibilitychange', onVisibilityChange)
    } else {
      // User logged out: clear data and stop polling
      subscriptionStore.clear()
      announcementStore.reset()
      adminComplianceStore.reset()
      document.removeEventListener('visibilitychange', onVisibilityChange)
    }
  },
  { immediate: true }
)

// Route change trigger (throttled by store)
router.afterEach(() => {
  if (authStore.isAuthenticated) {
    announcementStore.fetchAnnouncements()
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
  window.removeEventListener('admin-compliance-required', onAdminComplianceRequired)
  window.removeEventListener(AUTH_MODAL_EVENT, openAuthModal)
})

onMounted(async () => {
  window.addEventListener('admin-compliance-required', onAdminComplianceRequired)
  window.addEventListener(AUTH_MODAL_EVENT, openAuthModal)
  const pendingRequest = consumeAuthModalRequest()
  if (pendingRequest) {
    openAuthModal(new CustomEvent(AUTH_MODAL_EVENT, {
      detail: { view: pendingRequest.view, ...pendingRequest.options },
    }))
  }

  // Check if setup is needed
  try {
    const status = await getSetupStatus()
    if (status.needs_setup && route.path !== '/setup') {
      router.replace('/setup')
      return
    }
  } catch {
    // If setup endpoint fails, assume normal mode and continue
  }

  // Load public settings into appStore (will be cached for other components)
  await appStore.fetchPublicSettings()

  // Re-resolve document title now that site settings are available
  updateDocumentTitle()
})
</script>

<template>
  <NavigationProgress />
  <RouterView />
  <Toast />
  <AnnouncementPopup />
  <AdminComplianceDialog />
  <AuthModal
    :show="showAuthModal"
    :initial-view="authModalView"
    :options="authModalOptions"
    @close="closeAuthModal"
  />
</template>
