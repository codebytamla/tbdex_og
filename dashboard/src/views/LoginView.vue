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
      : authStore.defaultDashboardRoute()
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
}

.login-card {
  width: min(100%, 27rem);
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
</style>
