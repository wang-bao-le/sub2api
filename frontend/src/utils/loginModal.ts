import { isNavigationFailure, NavigationFailureType, type Router } from 'vue-router'

export type AuthModalView = 'login' | 'register' | 'forgot-password' | 'reset-password' | 'email-verify'

export interface AuthModalOptions {
  redirect?: string
  token?: string
  email?: string
  resetMethod?: 'code' | 'token'
  query?: Record<string, string>
}

export const AUTH_MODAL_EVENT = 'sub2api:open-auth-modal'
export const LOGIN_MODAL_EVENT = AUTH_MODAL_EVENT
const EXPLICIT_LOGOUT_KEY = 'sub2api:explicit-logout'

let pendingAuthModalRequest: { view: AuthModalView; options: AuthModalOptions } | null = null

export function markExplicitLogout() {
  try {
    sessionStorage.setItem(EXPLICIT_LOGOUT_KEY, '1')
  } catch {
    // Ignore storage failures; the modal guard remains best-effort.
  }
}

export function isExplicitLogout() {
  try {
    return sessionStorage.getItem(EXPLICIT_LOGOUT_KEY) === '1'
  } catch {
    return false
  }
}

export function clearExplicitLogout() {
  try {
    sessionStorage.removeItem(EXPLICIT_LOGOUT_KEY)
  } catch {
    // Ignore storage failures.
  }
}

export function requestAuthModal(view: AuthModalView = 'login', options: AuthModalOptions = {}) {
  clearExplicitLogout()
  pendingAuthModalRequest = { view, options }
  window.dispatchEvent(new CustomEvent(AUTH_MODAL_EVENT, { detail: { view, ...options } }))
}

export function consumeAuthModalRequest() {
  const request = pendingAuthModalRequest
  pendingAuthModalRequest = null
  return request
}

export function requestLoginModal(options: AuthModalOptions = {}) {
  requestAuthModal('login', options)
}

export function resolveAuthRedirect(
  redirect?: string,
  routeRedirect?: unknown,
  defaultRedirect = '/dashboard'
): string {
  const candidate = typeof redirect === 'string' && redirect.trim()
    ? redirect
    : typeof routeRedirect === 'string' && routeRedirect.trim()
      ? routeRedirect
      : ''

  // The marketing page is an entry point, not a post-login destination.
  if (!candidate || candidate === '/' || candidate === '/home' || candidate.startsWith('/home?')) {
    return defaultRedirect
  }

  return candidate
}

export async function navigateAfterAuth(router: Router, target: string): Promise<boolean> {
  try {
    const failure = await router.push(target)
    return !failure || isNavigationFailure(failure, NavigationFailureType.duplicated)
  } catch {
    return false
  }
}
