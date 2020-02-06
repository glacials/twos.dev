import Vue from 'vue'
Vue.use(IconsPlugin)

// Font Awesome
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
Vue.component('font-awesome-icon', FontAwesomeIcon)
import { library } from '@fortawesome/fontawesome-svg-core'

// Font Awesome - specific icon imports
import { faTwitter, faTwitch, faGithub, faLinkedin, faOsi } from '@fortawesome/free-brands-svg-icons'
library.add(faTwitter,faTwitch, faGithub, faLinkedin, faOsi)
import { faMobileAlt, faCommentDots } from '@fortawesome/free-solid-svg-icons'
library.add(faMobileAlt, faCommentDots)

// Bootstrap
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
Vue.use(BootstrapVue)
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import App from './App.vue'

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')
