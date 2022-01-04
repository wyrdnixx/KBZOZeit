import Vue from 'vue'
import App from './App.vue'

import VueCookies from 'vue-cookies'
import moment from 'vue-moment'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'

Vue.config.productionTip = false


// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue)
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin)
Vue.use(VueCookies)
Vue.use(moment)


new Vue({
  render: h => h(App),
}).$mount('#app')
