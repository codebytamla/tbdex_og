<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Message from 'primevue/message'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
})

async function submit() {
  const success = await authStore.login(form.email, form.password)
  if (!success) return

  const redirect =
    typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/')
      ? route.query.redirect
      : authStore.defaultOverviewRoute()
  await router.push(redirect)
}
</script>

<template>
  <section class="login-page">
    <Card class="login-card">
      <template #title>Sign in</template>
      <template #subtitle>Use your account to access BorderDesk</template>
      <template #content>
        <form class="form" @submit.prevent="submit">
          <label class="field">
            <span>Email</span>
            <InputText v-model="form.email" type="email" required autocomplete="email" />
          </label>

          <label class="field">
            <span>Password</span>
            <Password
              v-model="form.password"
              required
              :feedback="false"
              toggleMask
              fluid
              autocomplete="current-password"
            />
          </label>

          <Message v-if="authStore.errorMessage" severity="error" size="small" variant="simple">
            {{ authStore.errorMessage }}
          </Message>

          <Button
            type="submit"
            label="Sign in"
            :loading="authStore.loading"
            :disabled="!form.email || !form.password"
          />
        </form>
      </template>
    </Card>
  </section>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 1.5rem;
  background: #f8f9fb;
}

.login-card {
  width: min(100%, 26rem);
  border: 1px solid #e5e7eb;
  box-shadow: var(--shadow-lg, 0 12px 24px -4px rgb(0 0 0 / 0.08));
  border-radius: var(--radius-xl, 18px) !important;
}

:deep(.login-card .p-card-title) {
  font-size: 1.375rem;
  font-weight: 800;
  letter-spacing: -0.025em;
  color: #111827;
}

:deep(.login-card .p-card-subtitle) {
  color: #6b7280;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

:deep(.login-card .p-card-content) {
  padding: 1.5rem;
}

.form {
  display: grid;
  gap: 1rem;
}

.field {
  display: grid;
  gap: 0.375rem;
}

.field span {
  font-size: 0.825rem;
  color: #374151;
  font-weight: 600;
  letter-spacing: -0.005em;
}

:deep(.form .p-button) {
  margin-top: 0.25rem;
  height: 2.625rem;
  font-size: 0.935rem;
  font-weight: 700;
  background: var(--p-primary-600);
  border-color: var(--p-primary-600);
  border-radius: var(--radius-md, 10px);
}

:deep(.form .p-button:hover:not(:disabled)) {
  background: var(--p-primary-700);
  border-color: var(--p-primary-700);
}
</style>
