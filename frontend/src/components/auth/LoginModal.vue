<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="modal-overlay login-modal-overlay"
        role="dialog"
        aria-modal="true"
        aria-labelledby="login-modal-title"
        @click.self="close"
      >
        <div ref="dialogRef" class="modal-content login-modal-content" @click.stop>
          <h2 id="login-modal-title" class="sr-only">{{ t('home.login') }}</h2>
          <button
            type="button"
            class="login-modal-close"
            :aria-label="t('common.close')"
            @click="close"
          >
            <Icon name="x" size="md" aria-hidden="true" />
          </button>
          <LoginView modal />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import LoginView from '@/views/auth/LoginView.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n()
const dialogRef = ref<HTMLElement | null>(null)
let previousActiveElement: HTMLElement | null = null

function close() {
  emit('close')
}

function handleEscape(event: KeyboardEvent) {
  if (event.key === 'Escape') close()
}

watch(
  () => props.show,
  async (isOpen) => {
    if (isOpen) {
      previousActiveElement = document.activeElement as HTMLElement
      document.body.classList.add('modal-open')
      await nextTick()
      dialogRef.value?.querySelector<HTMLElement>('input, button, [href]')?.focus()
    } else {
      document.body.classList.remove('modal-open')
      previousActiveElement?.focus()
      previousActiveElement = null
    }
  },
  { immediate: true }
)

onMounted(() => window.addEventListener('keydown', handleEscape))
onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape)
  document.body.classList.remove('modal-open')
})
</script>
