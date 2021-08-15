import { createApp } from 'vue'
import App from '@/App.vue'
import { store } from '@/store'
import VueAnimXYZ from '@animxyz/vue3'
import { setupI18n } from '@/i18n'
import enUS from '@/locales/en-US.yml'
import zhCN from '@/locales/zh-CN.yml'
import VueNativeSock from 'vue-native-websocket-vue3';

const i18n = setupI18n({
  globalInjection: true,
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en-US',
  messages: {
    'en-US': enUS,
    'zh-CN': zhCN
  }
})

const app = createApp(App)

app.use(i18n)
app.use(store)
app.use(VueAnimXYZ)
app.use(VueNativeSock, process.env.VUE_APP_WEBSOCKET_URL, { store: store, format: 'json' })

window['CloudRain'] = app.mount('#app')

export default {
  app
}
