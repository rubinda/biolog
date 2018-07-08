<template>
<div>
  <nav class="inner-nav">
    <b-nav>
      <b-nav-item active>BioLog</b-nav-item>
      <b-nav-item class="ml-auto" v-on:click="signOut()">Odjava</b-nav-item>
    </b-nav>
  </nav>
  <div class="google-map" id="map-slo"></div>
  <div class="new-observe-btn">
    <i class="fas fa-4x fa-plus-circle"></i>
  </div>
</div>
</template>

<script>
/* global google */
/* global gapi */

export default {
  name: 'Map',

  mounted() {
    this.initMap();
  },

  methods: {
    initMap() {
      // The location of Uluru
      const geos = { lat: 46.119944, lng: 14.815333 };
      // The map, centered at Uluru
      const map = new google.maps.Map(
        document.getElementById('map-slo'),
        {
          minZoom: 9,
          zoom: 9,
          center: geos,
          streetViewControl: false,
          fullscreenControl: false,
          mapTypeControl: false,
        });
      const marker = new google.maps.Marker({});
      map.addListener('click', (event) => {
        marker.setMap(map);
        marker.setPosition(event.latLng);
      });
    },

    signOut() {
      // Remove the token from local storage and redirect the user to the landing page
      localStorage.removeItem('userToken');
      this.$router.push('/');

      // Odjavi ga tudi iz Google racuna
      const auth2 = gapi.auth2.getAuthInstance();
      auth2.signOut().then(() => {
        console.log('User signed out.');
      });
    },
  },
};
</script>

<style scoped>
  .google-map {
    height: calc(100vh - 40px);
  }

  .inner-nav {
    height: 40px;
  }

  .new-observe-btn {
    border: 1px;
    border-radius: 100%;
    height: 80px;
    width: 80px;
    position: absolute;
    bottom: 50px;
    right: 70px;
    vertical-align: middle;
  }

  .new-observe-btn i {
    color: rgb(38, 117, 42);
    line-height: 80px;
  }
</style>
