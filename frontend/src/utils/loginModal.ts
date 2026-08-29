export const LOGIN_MODAL_EVENT = 'sub2api:open-login-modal'

export function requestLoginModal() {
  window.dispatchEvent(new CustomEvent(LOGIN_MODAL_EVENT))
}
