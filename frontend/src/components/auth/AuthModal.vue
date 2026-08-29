<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="modal-overlay login-modal-overlay"
        id="auth-modal"
        role="dialog"
        aria-modal="true"
        aria-labelledby="auth-modal-title"
        @click.self="close"
      >
        <div ref="dialogRef" class="modal-content login-modal-content" @click.stop>
          <h2 id="auth-modal-title" class="sr-only">{{ t('home.login') }}</h2>
          <button
            type="button"
            class="login-modal-close"
            :aria-label="t('common.close')"
            @click="close"
          >
            <Icon name="x" size="md" aria-hidden="true" />
          </button>
          <component
            :is="activeComponent"
            :key="activeView"
            modal
            :options="activeOptions"
            @success="handleSuccess"
          />
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import LoginFormModal from './LoginFormModal.vue'
import RegisterFormModal from './RegisterFormModal.vue'
import ForgotPasswordModal from './ForgotPasswordModal.vue'
import ResetPasswordModal from './ResetPasswordModal.vue'
import EmailVerifyModal from './EmailVerifyModal.vue'
import { AUTH_MODAL_EVENT, type AuthModalOptions, type AuthModalView } from '@/utils/loginModal'

const props = withDefaults(defineProps<{
  show: boolean
  initialView?: AuthModalView
  options?: AuthModalOptions
}>(), {
  initialView: 'login',
  options: () => ({})
})
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n()
const activeView = ref<AuthModalView>(props.initialView)
const activeOptions = ref<AuthModalOptions>({ ...props.options })
const dialogRef = ref<HTMLElement | null>(null)
let previousActiveElement: HTMLElement | null = null

const componentMap = {
  login: LoginFormModal,
  register: RegisterFormModal,
  'forgot-password': ForgotPasswordModal,
  'reset-password': ResetPasswordModal,
  'email-verify': EmailVerifyModal,
}
const activeComponent = computed(() => componentMap[activeView.value])

function close() {
  emit('close')
}

function handleSuccess() {
  close()
}

function handleAuthRequest(event: Event) {
  const detail = (event as CustomEvent<{ view?: AuthModalView } & AuthModalOptions>).detail || {}
  if (detail.view) {
    activeView.value = detail.view
    activeOptions.value = { redirect: detail.redirect, token: detail.token, email: detail.email, resetMethod: detail.resetMethod, query: detail.query }
  }
}

function handleEscape(event: KeyboardEvent) {
  if (event.key === 'Escape') close()
}

watch(() => props.initialView, (view) => {
  activeView.value = view
})

watch(() => props.show, async (isOpen) => {
  if (isOpen) {
    activeView.value = props.initialView
    activeOptions.value = { ...props.options }
    previousActiveElement = document.activeElement as HTMLElement
    document.body.classList.add('modal-open')
    await nextTick()
    dialogRef.value?.querySelector<HTMLElement>('input, button, [href]')?.focus()
  } else {
    document.body.classList.remove('modal-open')
    previousActiveElement?.focus()
    previousActiveElement = null
  }
}, { immediate: true })

onMounted(() => {
  window.addEventListener('keydown', handleEscape)
  window.addEventListener(AUTH_MODAL_EVENT, handleAuthRequest)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape)
  window.removeEventListener(AUTH_MODAL_EVENT, handleAuthRequest)
  document.body.classList.remove('modal-open')
})
</script>
