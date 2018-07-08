import Vue from 'vue';
import Router from 'vue-router';
import HelloWorld from '@/components/HelloWorld';
import Home from '@/components/Home';
import Login from '@/components/Login';
import Map from '@/components/Map';

Vue.use(Router);

export default new Router({
  // History mode requires properly configured backend servers,
  // see: https://router.vuejs.org/guide/essentials/history-mode.html
  mode: 'history',
  routes: [
    {
      path: '/hello',
      name: 'HelloWorld',
      component: HelloWorld,
    },
    {
      path: '/',
      name: 'Home',
      component: Home,
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
    },
    {
      path: '/map',
      name: 'Map',
      component: Map,
    },
  ],
});
