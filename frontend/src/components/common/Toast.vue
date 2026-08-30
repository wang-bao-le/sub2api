<template>
  <Teleport to="body">
    <div
      class="pointer-events-none fixed inset-x-0 top-0 z-[9999] space-y-3"
      aria-live="polite"
      aria-atomic="true"
    >
      <TransitionGroup
        enter-active-class="transition ease-out duration-300"
        enter-from-class="opacity-0 -translate-y-full"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-200"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-full"
      >
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="[
            'pointer-events-auto w-full overflow-hidden rounded-none border-x-0 border-t-0 shadow-xl shadow-black/10 backdrop-blur-xl',
            'bg-black/75 dark:bg-black/75',
            getBorderColor(toast.type)
          ]"
        >
          <div class="p-4">
            <div class="grid grid-cols-[1fr_minmax(0,1fr)_1fr] items-center gap-3">
              <!-- Icon -->
              <div class="flex-shrink-0">
                <Icon
                  :name="getToastIconName(toast.type)"
                  size="md"
                  :class="getIconColor(toast.type)"
                  aria-hidden="true"
                />
              </div>

              <!-- Content -->
              <div class="min-w-0 text-center">
                <p v-if="toast.title" class="text-sm font-semibold text-white">
                  {{ toast.title }}
                </p>
                <p
                  :class="[
                    'text-sm leading-relaxed',
                    toast.title
                      ? 'mt-1 text-gray-200'
                      : 'text-white'
                  ]"
                >
                  {{ toast.message }}
                </p>
              </div>

              <!-- Close button -->
              <button
                @click="removeToast(toast.id)"
                class="-m-1 flex-shrink-0 justify-self-end rounded p-1 text-gray-300 transition-colors hover:bg-white/10 hover:text-white"
                aria-label="Close notification"
              >
                <Icon name="x" size="sm" />
              </button>
            </div>
          </div>

        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const toasts = computed(() => appStore.toasts)

const getToastIconName = (type: string): 'checkCircle' | 'xCircle' | 'exclamationTriangle' | 'infoCircle' => {
  switch (type) {
    case 'success':
      return 'checkCircle'
    case 'error':
      return 'xCircle'
    case 'warning':
      return 'exclamationTriangle'
    case 'info':
    default:
      return 'infoCircle'
  }
}

const getIconColor = (type: string): string => {
  const colors: Record<string, string> = {
    success: 'text-green-500',
    error: 'text-red-500',
    warning: 'text-yellow-500',
    info: 'text-blue-500'
  }
  return colors[type] || colors.info
}

const getBorderColor = (type: string): string => {
  const colors: Record<string, string> = {
    success: 'border-green-500',
    error: 'border-red-500',
    warning: 'border-yellow-500',
    info: 'border-blue-500'
  }
  return colors[type] || colors.info
}

const removeToast = (id: string) => {
  appStore.hideToast(id)
}
</script>
