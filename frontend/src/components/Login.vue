<template>
  <div class="full-viewport">
    <b-container class="h-100">
      <b-row class="h-100" align-h="center">
        <b-col class="col-lg-6 my-auto">
            <b-card-group deck>
                 <b-card title="Prijava"
                        bg-variant="dark"
                        text-variant="white"
                        footer-tag="footer"
                        class="login-card">
                    <em slot="footer">
                        <router-link to="/" class="footer-link">
                            <a>Nazaj</a>
                        </router-link>
                    </em>
                    <p class="card-text"></p>
                    <b-container>
                      <b-row align-h="center" align-v="center">
                        <b-col>
                            <div id="google-signin2"></div>
                        </b-col>
                        <b-col>
                          <img src="../assets/owl.png" alt="Sovica" width="100px">
                        </b-col>
                      </b-row>
                    </b-container>
                </b-card>
            </b-card-group>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
/* global gapi */

export default {
  name: 'Login',

  mounted() {
    if (!gapi.auth2) {
      gapi.load('auth2', () => {
        gapi.auth2.init({
          client_id: '136952852051-0sojmp8svvt6m760j4fotcuaifohng72.apps.googleusercontent.com',
          cookiepolicy: 'single_host_origin',
        });
      });
    }

    this.renderGoogleButton();
  },

  methods: {
    onSignIn(googleUser) {
      const tok = { token: googleUser.getAuthResponse().id_token };
      this.$axios.post('/login/google', tok)
        .then((response) => {
          if (response.status === 200) {
            // Login successful, store the acquired JWT in localStorage
            localStorage.setItem('userToken', response.data.token);
            // Redirect the user to the map component
            this.$router.push('map');
          } else {
            console.warn('Napaka pri prijavi');
            // TODO: inform the user about the unsuccesful login on a more
            // user friendly way
          }
        })
        .catch((e) => {
          console.log(e);
          // FIXME: do not use console in production
        });
    },

    renderGoogleButton() {
      gapi.signin2.render('google-signin2', {
        scope: 'profile email',
        width: 200,
        height: 30,
        longtitle: true,
        theme: 'dark',
        onsuccess: this.onSignIn,
      });
    },
  },
};
</script>

<style scoped>

.btn-sign-in {
  background: #fff;
  font: 16px/22px Roboto;
  padding: 4px 8px;
  border: 1px solid #ccc;
  display: inline-block;
  cursor: pointer;
}

.full-viewport {
  height: 100vh;
  min-height: 768px;
  background-color: #4c945b;
  background-image: url('../assets/iron-grip.png')
}

.footer-link {
  font-style: normal;
  color: #fff;
  font-size: 10pt;
  font-weight: 700;
  text-transform: uppercase;
}

.footer-link:hover,
.footer-link:focus {
  color: #4c945b;
  text-decoration: none;
}

.login-card {
  height: 300px;
}
</style>
