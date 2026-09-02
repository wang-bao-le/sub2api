import { describe, expect, it } from 'vitest'

import { resolveAuthRedirect } from '../loginModal'

describe('resolveAuthRedirect', () => {
  it('uses the role-specific default when no redirect is provided', () => {
    expect(resolveAuthRedirect(undefined, undefined, '/admin/dashboard')).toBe('/admin/dashboard')
    expect(resolveAuthRedirect(undefined, undefined, '/dashboard')).toBe('/dashboard')
  })

  it('uses the role-specific default for marketing-page redirects', () => {
    expect(resolveAuthRedirect('/home', undefined, '/admin/dashboard')).toBe('/admin/dashboard')
    expect(resolveAuthRedirect(undefined, '/home?login=1', '/admin/dashboard')).toBe('/admin/dashboard')
  })

  it('preserves an explicit redirect target', () => {
    expect(resolveAuthRedirect('/profile', undefined, '/admin/dashboard')).toBe('/profile')
    expect(resolveAuthRedirect(undefined, '/admin/users', '/admin/dashboard')).toBe('/admin/users')
  })
})
