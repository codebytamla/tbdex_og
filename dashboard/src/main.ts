import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'
import Aura from '@primeuix/themes/aura'
import { definePreset, palette } from '@primeuix/themes'

import App from './App.vue'
import router from './router'

const app = createApp(App)
const primaryPalette = palette('{teal}') as Record<string, string>
const lightSurfacePalette = palette('{slate}') as Record<string, string>
const darkSurfacePalette: Record<string, string> = {
  0: lightSurfacePalette['950'] || '#020617',
  50: lightSurfacePalette['900'] || '#0f172a',
  100: lightSurfacePalette['800'] || '#1e293b',
  200: lightSurfacePalette['700'] || '#334155',
  300: lightSurfacePalette['600'] || '#475569',
  400: lightSurfacePalette['500'] || '#64748b',
  500: lightSurfacePalette['400'] || '#94a3b8',
  600: lightSurfacePalette['300'] || '#cbd5e1',
  700: lightSurfacePalette['200'] || '#e2e8f0',
  800: lightSurfacePalette['100'] || '#f1f5f9',
  900: lightSurfacePalette['50'] || '#f8fafc',
  950: '#ffffff',
}

const MonoPreset = definePreset(Aura, {
  semantic: {
    primary: primaryPalette,
    colorScheme: {
      light: {
        surface: {
          0: '#ffffff',
          ...lightSurfacePalette,
        },
      },
      dark: {
        surface: darkSurfacePalette,
      },
    },
    focusRing: {
      width: '1px',
      style: 'solid',
      color: '{primary.400}',
      offset: '2px',
    },
  },
  components: {
    card: {
      colorScheme: {
        dark: {
          root: {
            background: '{surface.100}',
            color: '{surface.900}',
          },
          subtitle: {
            color: '{surface.600}',
          },
        },
      },
    },
    panelmenu: {
      root: {
        gap: '0.125rem',
      },
      panel: {
        background: 'transparent',
        borderWidth: '0',
        borderColor: 'transparent',
        color: '{surface.700}',
        padding: '0',
        borderRadius: '0.7rem',
        first: {
          borderWidth: '0',
          topBorderRadius: '0.7rem',
        },
        last: {
          borderWidth: '0',
          bottomBorderRadius: '0.7rem',
        },
      },
      item: {
        color: '{surface.600}',
        focusBackground: 'color-mix(in srgb, {primary.color} 8%, {surface.0})',
        focusColor: '{primary.600}',
        gap: '0.5rem',
        padding: '0.48rem 0.62rem',
        borderRadius: '0.7rem',
        icon: {
          color: '{surface.400}',
          focusColor: '{primary.600}',
        },
      },
      submenu: {
        indent: '0.8rem',
      },
      submenuIcon: {
        color: '{surface.400}',
        focusColor: '{primary.500}',
      },
    },
  },
})

app.use(PrimeVue, {
  theme: {
    preset: MonoPreset,
    options: {
      darkModeSelector: false,
    },
  },
})
app.use(ConfirmationService)
app.use(ToastService)
app.use(createPinia())
app.use(router)

app.mount('#app')
