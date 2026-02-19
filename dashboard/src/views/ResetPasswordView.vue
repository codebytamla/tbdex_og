<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Message from 'primevue/message'
import Password from 'primevue/password'

const route = useRoute()
const router = useRouter()

const form = reactive({
  token: typeof route.query.token === 'string' ? route.query.token : '',
  new_password: '',
  confirm_password: '',
})

const loading = ref(false)
const errorMessage = ref('')
const success = ref(false)
const API_V1_PREFIX = '/api/v1'

const canSubmit = computed(
  () =>
    !!form.token.trim() &&
    !!form.new_password &&
    !!form.confirm_password &&
    form.new_password === form.confirm_password &&
    !loading.value,
)

function apiUrl(path: string) {
  const baseUrl = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/+$/, '')
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const versionedPath =
    normalizedPath === API_V1_PREFIX || normalizedPath.startsWith(`${API_V1_PREFIX}/`)
      ? normalizedPath
      : `${API_V1_PREFIX}${normalizedPath}`
  if (baseUrl) return `${baseUrl}${versionedPath}`
  return versionedPath
}

function parseApiError(payload: unknown) {
  if (typeof payload === 'string' && payload.trim()) return payload.trim()
  if (payload && typeof payload === 'object') {
    const record = payload as Record<string, unknown>
    for (const key of ['error', 'message', 'detail']) {
      const value = record[key]
      if (typeof value === 'string' && value.trim()) return value.trim()
    }
  }
  return 'Unable to reset password.'
}

async function submit() {
  if (!canSubmit.value) return
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(apiUrl('/auth/reset-password'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        token: form.token.trim(),
        new_password: form.new_password,
      }),
    })

    const payload = (await response.json().catch(() => ({}))) as unknown
    if (!response.ok) throw new Error(parseApiError(payload))
    success.value = true
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Unable to reset password.'
  } finally {
    loading.value = false
  }
}

async function goToLogin() {
  await router.push('/auth/login')
}
</script>

<template>
  <section class="reset-page">
    <Card class="reset-card">
      <template #title>Reset password</template>
      <template #subtitle>Use your reset token to set a new password</template>
      <template #content>
        <div v-if="success" class="success-wrap">
          <Message severity="success" size="small" variant="simple">
            Password reset complete. You can now sign in with your new password.
          </Message>
          <Button label="Go to sign in" @click="goToLogin" />
        </div>

        <form v-else class="form" @submit.prevent="submit">
          <label class="field">
            <span>Reset token</span>
            <InputText v-model="form.token" required autocomplete="off" />
          </label>

          <label class="field">
            <span>New password</span>
            <Password
              v-model="form.new_password"
              required
              :feedback="false"
              toggleMask
              fluid
              autocomplete="new-password"
            />
          </label>

          <label class="field">
            <span>Confirm new password</span>
            <Password
              v-model="form.confirm_password"
              required
              :feedback="false"
              toggleMask
              fluid
              autocomplete="new-password"
            />
          </label>

          <Message
            v-if="form.confirm_password && form.new_password !== form.confirm_password"
            severity="warn"
            size="small"
            variant="simple"
          >
            Password confirmation does not match.
          </Message>

          <Message v-if="errorMessage" severity="error" size="small" variant="simple">
            {{ errorMessage }}
          </Message>

          <div class="actions">
            <Button type="button" text label="Back to sign in" @click="goToLogin" />
            <Button type="submit" label="Reset password" :loading="loading" :disabled="!canSubmit" />
          </div>
        </form>
      </template>
    </Card>
  </section>
</template>

<style scoped>
.reset-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1.5rem;
}

.reset-card {
  width: min(100%, 30rem);
  border: 1px solid var(--p-surface-200);
  box-shadow: none;
}

.form {
  display: grid;
  gap: 0.85rem;
}

.field {
  display: grid;
  gap: 0.38rem;
}

.field span {
  font-size: 0.83rem;
  color: var(--p-surface-700);
  font-weight: 600;
}

.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.success-wrap {
  display: grid;
  gap: 0.8rem;
}
</style>
