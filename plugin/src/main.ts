import { createApp } from 'vue'
import 'highlight.js/styles/base16/classic-dark.css'

import App from './App.vue'

import router from './router'

import { createPinia } from 'pinia'

const pinia = createPinia();
pinia.use(({ store }) => {
    store.$router = router;
})

createApp(App)
    .use(router)
    .use(pinia)
    .mount('#app')
