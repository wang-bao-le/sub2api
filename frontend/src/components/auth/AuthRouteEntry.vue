<template><div class="min-h-screen bg-gray-50 dark:bg-dark-950" aria-hidden="true"></div></template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { requestAuthModal, type AuthModalView } from '@/utils/loginModal'

const route = useRoute()
const router = useRouter()
const views: Record<string, AuthModalView> = {
  '/login': 'login',
  '/register': 'register',
  '/forgot-password': 'forgot-password',
  '/reset-password': 'reset-password',
  '/email-verify': 'email-verify',
}

onMounted(async () => {
  const query = Object.fromEntries(
    Object.entries(route.query).filter((entry): entry is [string, string] => typeof entry[1] === 'string')
  )
  const options = {
    redirect: query.redirect,
    token: query.token,
    email: query.email,
    query,
  }
  const view = views[route.path] || 'login'
  await router.replace('/home')
  requestAuthModal(view, options)
})
</script>
