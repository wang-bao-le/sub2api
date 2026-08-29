<template>
  <AuthLayout compact :embedded="modal">
    <div class="space-y-5">
      <!-- Title -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.resetPasswordTitle') }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ t('auth.resetPasswordHint') }}
        </p>
      </div>

      <!-- Invalid Link State -->
      <div v-if="isInvalidLink" class="space-y-6">
        <div class="rounded-xl border border-amber-200 bg-amber-50 p-6 dark:border-amber-800/50 dark:bg-amber-900/20">
          <div class="flex flex-col items-center gap-4 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-amber-100 dark:bg-amber-800/50">
              <Icon name="exclamationCircle" size="lg" class="text-amber-600 dark:text-amber-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-amber-800 dark:text-amber-200">
                {{ t('auth.invalidResetLink') }}
              </h3>
              <p class="mt-2 text-sm text-amber-700 dark:text-amber-300">
                {{ t('auth.invalidResetLinkHint') }}
              </p>
            </div>
          </div>
        </div>

        <div class="text-center">
          <a
            href="/forgot-password"
            @click="handleAuthLink($event, 'forgot-password')"
            class="inline-flex items-center gap-2 font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
          >
            {{ t('auth.requestNewResetLink') }}
          </a>
        </div>
      </div>

      <!-- Success State -->
      <div v-else-if="isSuccess" class="space-y-6">
        <div class="rounded-xl border border-green-200 bg-green-50 p-6 dark:border-green-800/50 dark:bg-green-900/20">
          <div class="flex flex-col items-center gap-4 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-800/50">
              <Icon name="checkCircle" size="lg" class="text-green-600 dark:text-green-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-green-800 dark:text-green-200">
                {{ t('auth.passwordResetSuccess') }}
              </h3>
              <p class="mt-2 text-sm text-green-700 dark:text-green-300">
                {{ t('auth.passwordResetSuccessHint') }}
              </p>
            </div>
          </div>
        </div>

        <div class="text-center">
          <a
            href="/login"
            @click="handleAuthLink($event, 'login')"
            class="btn btn-primary inline-flex items-center gap-2"
          >
            <Icon name="login" size="md" />
            {{ t('auth.signIn') }}
          </a>
        </div>
      </div>

      <!-- Form State -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Email (readonly) -->
        <div>
          <label for="email" class="input-label">
            {{ t('auth.emailLabel') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="mail" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="email"
              :value="email"
              type="email"
              readonly
              disabled
              class="input pl-11 bg-gray-50 dark:bg-dark-700"
            />
          </div>
        </div>

        <!-- New Password Input -->
        <div v-if="isCodeMode">
          <label for="verifyCode" class="input-label">{{ t('auth.verificationCode') }}</label>
          <input
            id="verifyCode"
            v-model="formData.verifyCode"
            type="text"
            inputmode="numeric"
            maxlength="6"
            autocomplete="one-time-code"
            :disabled="isLoading"
            class="input"
            :class="{ 'input-error': errors.verifyCode }"
            :placeholder="t('auth.verificationCodePlaceholder')"
            :aria-invalid="Boolean(errors.verifyCode)"
            :aria-describedby="errors.verifyCode ? 'reset-verify-code-error' : undefined"
            @input="errors.verifyCode = ''; errorMessage = ''"
          />
          <p v-if="errors.verifyCode" id="reset-verify-code-error" class="input-error-text" role="alert">{{ errors.verifyCode }}</p>
          <button
            type="button"
            class="mt-2 text-sm font-medium text-gray-700 dark:text-gray-300"
            :disabled="resendCooldown > 0 || isResending"
            @click="resendCode"
          >
            {{ resendCooldown > 0 ? t('auth.resendCountdown', { countdown: resendCooldown }) : t('auth.resendCode') }}
          </button>
        </div>

        <div>
          <label for="password" class="input-label">
            {{ t('auth.newPassword') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="password"
              v-model="formData.password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="new-password"
              :disabled="isLoading"
              class="input pl-11 pr-11"
              :class="{ 'input-error': errors.password }"
              :placeholder="t('auth.newPasswordPlaceholder')"
              :aria-invalid="Boolean(errors.password)"
              :aria-describedby="errors.password ? 'reset-password-error' : undefined"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              :aria-label="showPassword ? t('auth.hidePassword') : t('auth.showPassword')"
              :aria-pressed="showPassword"
              aria-controls="password"
              class="absolute inset-y-0 right-0 flex min-w-11 items-center justify-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-black/40 dark:hover:text-dark-300 dark:focus-visible:ring-white/50"
            >
              <Icon v-if="showPassword" name="eyeOff" size="md" aria-hidden="true" />
              <Icon v-else name="eye" size="md" aria-hidden="true" />
            </button>
          </div>
          <p v-if="errors.password" id="reset-password-error" class="input-error-text" role="alert">
            {{ errors.password }}
          </p>
        </div>

        <!-- Confirm Password Input -->
        <div>
          <label for="confirmPassword" class="input-label">
            {{ t('auth.confirmPassword') }}
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="confirmPassword"
              v-model="formData.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              autocomplete="new-password"
              :disabled="isLoading"
              class="input pl-11 pr-11"
              :class="{ 'input-error': errors.confirmPassword }"
              :placeholder="t('auth.confirmPasswordPlaceholder')"
              :aria-invalid="Boolean(errors.confirmPassword)"
              :aria-describedby="errors.confirmPassword ? 'reset-confirm-password-error' : undefined"
            />
            <button
              type="button"
              @click="showConfirmPassword = !showConfirmPassword"
              :aria-label="showConfirmPassword ? t('auth.hidePassword') : t('auth.showPassword')"
              :aria-pressed="showConfirmPassword"
              aria-controls="confirmPassword"
              class="absolute inset-y-0 right-0 flex min-w-11 items-center justify-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-black/40 dark:hover:text-dark-300 dark:focus-visible:ring-white/50"
            >
              <Icon v-if="showConfirmPassword" name="eyeOff" size="md" aria-hidden="true" />
              <Icon v-else name="eye" size="md" aria-hidden="true" />
            </button>
          </div>
          <p v-if="errors.confirmPassword" id="reset-confirm-password-error" class="input-error-text" role="alert">
            {{ errors.confirmPassword }}
          </p>
        </div>

        <div v-if="errorMessage" class="flex items-start gap-3 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-900/50 dark:bg-red-950/30 dark:text-red-300" role="alert" aria-live="polite">
          <Icon name="exclamationCircle" size="md" class="mt-0.5 shrink-0" aria-hidden="true" />
          <p>{{ errorMessage }}</p>
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading"
          class="btn btn-primary w-full"
        >
          <svg
            v-if="isLoading"
            class="-ml-1 mr-2 h-4 w-4 animate-spin text-white"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <Icon v-else name="checkCircle" size="md" class="mr-2" />
          {{ isLoading ? t('auth.resettingPassword') : t('auth.resetPassword') }}
        </button>
      </form>
    </div>

    <!-- Footer -->
    <template #footer>
      <p class="text-gray-500 dark:text-dark-400">
        {{ t('auth.rememberedPassword') }}
        <a
          href="/login"
          @click="handleAuthLink($event, 'login')"
          class="font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
        >
          {{ t('auth.signIn') }}
        </a>
      </p>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores'
import { resetPassword } from '@/api/auth'
import { isValidPassword } from '@/utils/passwordPolicy'
import { requestAuthModal, type AuthModalOptions, type AuthModalView } from '@/utils/loginModal'
import { extractI18nErrorMessage } from '@/utils/apiError'

const { t } = useI18n()
const { modal = false, options } = defineProps<{ modal?: boolean; options?: AuthModalOptions }>()

function handleAuthLink(event: MouseEvent, view: AuthModalView): void {
  if (!modal) return
  event.preventDefault()
  requestAuthModal(view)
}

// ==================== Router & Stores ====================

const route = useRoute()
const appStore = useAppStore()

// ==================== State ====================

const isLoading = ref<boolean>(false)
const isResending = ref<boolean>(false)
const resendCooldown = ref<number>(0)
const isSuccess = ref<boolean>(false)
const errorMessage = ref<string>('')
const showPassword = ref<boolean>(false)
const showConfirmPassword = ref<boolean>(false)

// URL parameters
const email = ref<string>('')
const token = ref<string>('')
const isCodeMode = computed(() => (options?.resetMethod || (!token.value ? 'code' : 'token')) === 'code')

const formData = reactive({
  verifyCode: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  verifyCode: '',
  password: '',
  confirmPassword: ''
})

// Token links require a token; code-based resets only require the email.
const isInvalidLink = computed(() => !email.value || (!isCodeMode.value && !token.value))

// ==================== Lifecycle ====================

onMounted(() => {
  // Get email and token from URL query parameters
  email.value = options?.email || (route.query.email as string) || ''
  token.value = options?.token || (route.query.token as string) || ''

})

// ==================== Validation ====================

function validateForm(): boolean {
  errors.verifyCode = ''
  errors.password = ''
  errors.confirmPassword = ''

  let isValid = true

  if (isCodeMode.value && !/^\d{6}$/.test(formData.verifyCode.trim())) {
    errors.verifyCode = t('auth.invalidCode')
    isValid = false
  }

  // Password validation
  if (!formData.password) {
    errors.password = t('auth.passwordRequired')
    isValid = false
  } else if (!isValidPassword(formData.password)) {
    errors.password = t('auth.passwordMinLength')
    isValid = false
  }

  // Confirm password validation
  if (!formData.confirmPassword) {
    errors.confirmPassword = t('auth.confirmPasswordRequired')
    isValid = false
  } else if (formData.password !== formData.confirmPassword) {
    errors.confirmPassword = t('auth.passwordsDoNotMatch')
    isValid = false
  }

  return isValid
}

// ==================== Form Handlers ====================

async function handleSubmit(): Promise<void> {
  errorMessage.value = ''

  if (!validateForm()) {
    return
  }

  isLoading.value = true

  try {
    await resetPassword({
      email: email.value,
      ...(isCodeMode.value ? { verify_code: formData.verifyCode.trim() } : { token: token.value }),
      new_password: formData.password
    })

    isSuccess.value = true
    appStore.showSuccess(t('auth.passwordResetSuccess'))
  } catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { detail?: string; code?: string } } }

    // Check for invalid/expired token error
    if (err.response?.data?.code === 'INVALID_RESET_TOKEN') {
      errorMessage.value = t('auth.invalidOrExpiredToken')
    } else if (err.response?.data?.code === 'INVALID_VERIFY_CODE' || err.response?.data?.code === 'VERIFY_CODE_MAX_ATTEMPTS') {
      errors.verifyCode = t('auth.invalidCode')
      errorMessage.value = t('auth.invalidCode')
    } else {
      errorMessage.value = extractI18nErrorMessage(error, t, 'auth.errors', t('auth.resetPasswordFailed'))
    }

  } finally {
    isLoading.value = false
  }
}

async function resendCode(): Promise<void> {
  if (!isCodeMode.value || resendCooldown.value > 0 || isResending.value) return
  requestAuthModal('forgot-password', { email: email.value, resetMethod: 'code' })
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
