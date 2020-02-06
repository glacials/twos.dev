<template>
  <b-container>
    <h3 class="title">Projects</h3>
    <p>I've had many, many solo side projects over the years. Here's a small selection.</p>
    <b-row cols="1" cols-md="2" cols-xl="3">
      <b-col>
        <b-card img-src="../assets/licenseland.png" title="Licenseland">
          <b-card-sub-title sub-title-tag="small">Current</b-card-sub-title>
          <b-card-text>
            Dead-simple dual-licensing for open source projects.
          </b-card-text>
          <b-badge pill>Vue</b-badge>
          <b-badge pill>Firebase</b-badge>
          <b-badge pill>Stripe Connect</b-badge>
          <template v-slot:footer>
            <b-button-group>
              <b-button href="https://license.land" variant="outline-secondary" size="sm">
                Website
              </b-button>
            </b-button-group>
          </template>
        </b-card>
      </b-col>
      <b-col>
        <b-card img-src="../assets/whatsinstandard.svg" title="What's in Standard?">
          <b-card-sub-title sub-title-tag="small">Current</b-card-sub-title>
          <b-card-text>
            A reference tool for Magic: The Gathering players.
          </b-card-text>
          <b-badge pill variant="info">20k MAU</b-badge>
          <b-badge pill>Flutter</b-badge>
          <b-badge pill>Progressive Web App</b-badge>
          <template v-slot:footer>
            <b-button-group>
              <b-button href="https://whatsinstandard.com/" variant="outline-secondary" size="sm">
                Website
              </b-button>
              <b-button href="https://play.google.com/store/apps/details?id=com.whatsinstandard" variant="outline-secondary" size="sm">
                Android app
              </b-button>
              <b-button href="https://github.com/glacials/whatsinstandard" variant="outline-secondary" size="sm">
                Source
              </b-button>
            </b-button-group>
          </template>
        </b-card>
      </b-col>
      <b-col>
        <b-card img-src="../assets/housebot.svg" title="Housebot">
          <b-card-sub-title sub-title-tag="small">Archived</b-card-sub-title>
          <b-card-text>
            A Twitch bot that lets chat participants vote to control your lights.
          </b-card-text>
          <b-badge pill>CoffeeScript</b-badge>
          <b-badge pill>Not a great idea</b-badge>
          <b-badge pill>Still made it</b-badge>
          <template v-slot:footer>
            <b-button-group>
              <b-button href="https://github.com/glacials/housebot" variant="outline-secondary" size="sm">
                Source
              </b-button>
            </b-button-group>
          </template>
        </b-card>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  name: 'Projects',
}
</script>

<style scoped>
.badge {
  margin-right: .25em;
}

.col {
  margin-bottom: 2em;
}
</style>
