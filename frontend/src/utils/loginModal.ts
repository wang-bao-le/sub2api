export type AuthModalView = 'login' | 'register' | 'forgot-password' | 'reset-password' | 'email-verify'

export interface AuthModalOptions {
  redirect?: string
  token?: string
  email?: string
  query?: Record<string, string>
}

export const AUTH_MODAL_EVENT = 'sub2api:open-auth-modal'
export const LOGIN_MODAL_EVENT = AUTH_MODAL_EVENT

let pendingAuthModalRequest: { view: AuthModalView; options: AuthModalOptions } | null = null

export function requestAuthModal(view: AuthModalView = 'login', options: AuthModalOptions = {}) {
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
