// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Home from '@/components/Home';
import 'bootstrap-vue/dist/bootstrap-vue.css';
import 'bootstrap/dist/css/bootstrap.css';
import BootstrapVue from 'bootstrap-vue/dist/bootstrap-vue.esm';
import Vue from 'vue';
import axios from 'axios';
import App from './App';
import router from './router';

const axiosConf = {
  baseURL: 'https://localhost:4000/api/v1',
  timeout: 30000,
};


Vue.use(BootstrapVue);
Vue.config.productionTip = false;
Vue.prototype.$axios = axios.create(axiosConf);

// Custom components
Vue.component('BiologNavbar', require('./components/Navbar.vue').default);

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App, Home },
  template: '<App/>',
});
