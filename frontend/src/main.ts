import '@/assets/styles.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAccountsStore } from './stores'

const app = createApp(App)

app.use(createPinia())
app.use(router)

useAccountsStore().loadAccounts()

app.mount('#app')
