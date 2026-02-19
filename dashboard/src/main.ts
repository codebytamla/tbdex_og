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
const lightSurfacePalette = palette('{zinc}') as Record<string, string>
const darkSurfacePalette: Record<string, string> = {
  0: lightSurfacePalette['950'] || '#09090b',
  50: lightSurfacePalette['900'] || '#18181b',
  100: lightSurfacePalette['800'] || '#27272a',
  200: lightSurfacePalette['700'] || '#3f3f46',
  300: lightSurfacePalette['600'] || '#52525b',
  400: lightSurfacePalette['500'] || '#71717a',
  500: lightSurfacePalette['400'] || '#a1a1aa',
  600: lightSurfacePalette['300'] || '#d4d4d8',
  700: lightSurfacePalette['200'] || '#e4e4e7',
  800: lightSurfacePalette['100'] || '#f4f4f5',
  900: lightSurfacePalette['50'] || '#fafafa',
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
        color: '{surface.700}',
        focusBackground: 'color-mix(in srgb, {primary.color} 14%, {surface.0})',
        focusColor: '{primary.300}',
        gap: '0.5rem',
        padding: '0.48rem 0.62rem',
        borderRadius: '0.7rem',
        icon: {
          color: '{surface.500}',
          focusColor: '{primary.700}',
        },
      },
      submenu: {
        indent: '0.8rem',
      },
      submenuIcon: {
        color: '{surface.500}',
        focusColor: '{primary.700}',
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
