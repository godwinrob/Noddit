import Vue from 'vue'
import VueDarkMode from '@growthbunker/vuedarkmode'
import App from './App.vue'
import router from './router'
import BootstrapVue from 'bootstrap-vue'
import Icon from 'vue-awesome/icons'

Vue.use(VueDarkMode, {
  theme: "dark",
  component: [
    "alert", "avatar", "badge", "button", "divider", "heading", "icon",  "progress-bar",  "spinner",
    "checkbox", "file", "input", "label", "message", "radios", "select", "tabs", "textarea", "toggle"
  ]
});
Vue.use(BootstrapVue);
Vue.component('icon', Icon);

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'



Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
