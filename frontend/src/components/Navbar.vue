<template>
<b-navbar toggleable="md" id="main-nav" class="fixed-top" ref="landingNav"
  v-bind:class="[navbarShrink ? shrinkClass : '']">
  <div class="container">
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
    <b-navbar-brand id="brand-name" href="#">BioLog</b-navbar-brand>
    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav class="ml-auto">
        <router-link tag="li" to="/login">
          <a class="landing-link">Prijava</a>
        </router-link>
      </b-navbar-nav>
    </b-collapse>
  </div>
</b-navbar>
</template>

<script>
import axios from 'axios';


export default {
  name: 'BiologNavbar',

  mounted() {
    window.addEventListener('scroll', this.updateScroll);
  },

  data() {
    return {
      navbarShrink: false,
      shrinkClass: 'navbar-shrink',
    };
  },

  methods: {
    updateScroll() {
      // Check if the navbar should get a solid background
      if (window.scrollY > 250) {
        this.navbarShrink = true;
      } else {
        this.navbarShrink = false;
      }
    },

    redirectLogin() {
      axios.get('https://127.0.0.1:4000/api/v1/login/google')
        .then((response) => {
          console.log(response);
        });
    },
  },

  destroyed() {
    window.removeEventListener('scroll', this.updateScroll);
  },
};
</script>

<style>

  #brand-name {
    font-size: 20pt;
    font-weight: 600;
    color: rgba(255, 255, 255, .7);
  }

  #brand-name:hover,
  #brand-name:focus {
    text-decoration: none;
    color: #fff;
  }

  .nav-wrapper {
    position: absolute;
    width: 80%;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    margin: auto;
  }

  #main-nav {
    transition: all .35s;
  }

  .navbar-shrink #brand-name {
    color: rgb(130, 130, 130);
  }

  .navbar-shrink #brand-name:hover,
  .navbar-shrink #brand-name:focus
  {
    color: rgb(142, 230, 154);
  }

  .navbar-shrink a {
    color: rgb(130, 130, 130);
  }

  .navbar-shrink a:focus,
  .navbar-shrink a:hover,
  .navbar-shrink a:active {
    color: rgb(142, 230, 154);
  }

  .navbar-shrink {
    background-color: #fff;
    border-bottom: 1px solid rgba(33,37,41,.1);
  }

  .landing-link {
    font-size: .9rem;
    color: rgba(255,255,255,.7);
    text-transform: uppercase;
    text-decoration: none;
    letter-spacing: 2px;
    font-weight: 700;
    text-decoration-style: solid;
    padding: .5rem 1rem;
  }

  .landing-link:active,
  .landing-link:focus,
  .landing-link:hover {
    color: #fff;
    text-decoration: none;
  }
</style>
