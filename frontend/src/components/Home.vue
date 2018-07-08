<template>
  <div>
    <biolog-navbar></biolog-navbar>
    <header class="home-header">
      <b-container class="h-100">
         <b-row class="h-100" align-h="center">
          <b-col class="col-lg-6 my-auto">
            <div class="mx-auto">
              <h1 class="mb-5">
                Spletna aplikacija, ki nudi podporo pri popisu vrst po Sloveniji.
              </h1>
              <b-button size="lg" variant="outline-secondary btn-outline">
                Več o aplikaciji
              </b-button>
            </div>
          </b-col>
          <b-col class="col-lg-4 my-auto">
            <img src="../assets/biolog-net-jump@3x.png" alt="BioLog maskota" width="300px">
          </b-col>
        </b-row>
      </b-container>
    </header>
    <div class="landing-section section-yellow about">
      <b-container>
        <b-row>
          <b-col>
            <h2 class="section-heading text-dark">Izvor aplikacije</h2>
            <hr class="my-4 border-dark" />
          </b-col>
        </b-row>
        <b-row align-h="center">
          <b-col class="col-lg-6" align-self="center">
            <p class="section-paragraph text-dark">
              Ideja za aplikacijo izhaja iz dogodka Dragonhack, kjer je skupina
              <i>&bdquo;Die Schweinebraten Jungs&rdquo;</i> zasnovala platformo za
              popis vrst po Sloveniji. Ekipa je nato v sodelovanju z DŠB aplikacijo
              še dodelala.
            </p>
          </b-col>
          <b-col class="col-lg-4">
            <img src="../assets/bee.png" alt="cebela" width="200px">
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div class="landing-section info">
      <b-container>
        <b-row>
          <b-col>
            <h2 class="section-heading text-dark">Od študentov, za študente</h2>
            <hr class="my-4 border-green" />
          </b-col>
        </b-row>
        <b-row aligh-h="center">
          <b-col class="col-lg-3">
            <i class="far fa-smile fa-3x my-4 text-green"></i>
            <p class="section-paragraph text-dark">
              Uporaba aplikacije je brezplačna.
            </p>
          </b-col>
          <b-col class="col-lg-3">
            <i class="fas fa-code fa-3x my-4 text-green"></i>
            <p class="section-paragraph text-dark">
              Projekt je odprtokoden, torej je vsa koda dostopna na spletu.
            </p>
          </b-col>
          <b-col class="col-lg-3">
            <i class="far fa-handshake fa-3x my-4 text-green"></i>
            <p class="section-paragraph text-dark">
              Pomoč pri razvoju je vedno dobrodošla.
            </p>
          </b-col>
          <b-col class="col-lg-3">
            <i class="far fa-comments fa-3x my-4 text-green"></i>
            <p class="section-paragraph">
              Veseli bomo tudi vaših predlogov in komentarjev.
            </p>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div class="landing-section contact">
      <b-container>
        <b-row>
          <b-col>
            <h2 class="section-heading text-dark">Kontaktirajte nas</h2>
          </b-col>
        </b-row>
        <b-row align-h="center">
          <b-col class="col-lg-6">
            <ul class="list-inline list-social">
              <li class="list-inline-item social-email" title="Epošta"
                v-b-tooltip.hover>
                <a href="#">
                  <i class="far fa-envelope"></i>
                </a>
              </li>
              <li class="list-inline-item social-facebook" title="Facebook skupina"
                v-b-tooltip.hover>
                <a href="#">
                  <i class="fab fa-facebook-f"></i>
                </a>
              </li>
              <li class="list-inline-item social-discord" title="Discord strežnik"
                v-b-tooltip.hover>
                <a href="#">
                  <i class="fab fa-discord"></i>
                </a>
              </li>
              <li class="list-inline-item social-github" title="Github repozitorij"
                v-b-tooltip.hover>
                <a href="#">
                  <i class="fab fa-github"></i>
                </a>
              </li>
            </ul>
          </b-col>
        </b-row>
      </b-container>
    </div>
    <div class="section-landing footer">
      <b-container>
        <b-row align-h="center">
          <b-col class="col-lg-4">
              Made with <i class="fas fa-heart text-red"></i> by Schweinebraten Jungs, 2018
          </b-col>
        </b-row>
      </b-container>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home',

  data() {
    return {
      navstyle: {
        backgroundColor: 'transparent',
      },
    };
  },
};
</script>

<style scoped>
  .home-header {
    background: url(../assets/bg-pattern.png),
      linear-gradient(to left, #4c945b, #d2db7d); /* #7b4397,#dc2430 */
    height: 100vh;
    min-height: 768px;
  }

  .btn-outline {
    border-radius: 300px;
    border: 1px, solid;
    color: #fff;
    border-color: #fff;
  }

  .btn-outline.active,
  .btn-outline:active,
  .btn-outline:focus,
  .btn-outline:hover {
    color: #fff;
    border-color: #d2db7d;
    background-color: #d2db7d;
  }

  .landing-section {
    min-height: 400px;
  }

  .section-yellow {
    background-color: rgb(70, 165, 75); /* FFF080 */
    background: linear-gradient(#65b84c,#30c45c);
    color: #2c3e50;
  }

  .about {
    padding: 150px 0;
  }

  .info {
    padding: 100px 0;
  }

  .section-heading {
    font-size: 20pt;
    font-weight: 700;
  }

  .section-paragraph {
    font-weight: 500;
    font-size: 16pt;
  }

  .text-faded {
    color: rgba(255,255,255,.7);
  }

  .text-green {
    color: #4c945b;
  }

  hr {
    max-width: 60px;
    border-width: 3px;
  }

  .border-green {
    border-color: #4c945b;
  }

  .contact {
    padding: 100px 0;
    background: linear-gradient(#fdcc52,#fdc539);
  }

  .contact h2 {
    margin-bottom: 50px;
  }

  .contact ul.list-social li a {
    font-size: 40px;
    height: 80px;
    width: 80px;
    display: block;
    color: #fff;
    border-radius: 100%;
  }

  .contact ul.list-social li a i {
    line-height: 80px;
  }
  .contact ul.list-social li.social-email a {
    background-color: #dd4b39;
  }

  .contact ul.list-social li.social-facebook a {
    background-color: #3b5998;
  }

  .contact ul.list-social li.social-discord a {
    background-color: #7289da;
  }

  .contact ul.list-social li.social-github a {
    background-color: #24292e;
  }

  .footer {
    padding: 10px 0;
    color: #fff;
    background-color: #222;
    font-size: 10pt;
  }

  .footer i {
    color: red;
  }
</style>
