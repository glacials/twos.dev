<template>
  <b-container>
    <h2 class="title">Projects</h2>
    <p>I've had many, many solo side projects over the years. Here's a small selection.</p>
    <b-row cols="1" cols-md="2" cols-xl="3">
      <Project name="Licenseland" logo="../assets/images/licenseland.webp" status="Current">
        <template v-slot:tags>
          <b-badge pill>Vue</b-badge>
          <b-badge pill>Firebase</b-badge>
          <b-badge pill>Stripe Connect</b-badge>
        </template>
        <template v-slot:description>
          Dead simple dual-licensing for open source projects.
        </template>
        <template v-slot:buttons>
          <b-button href="https://license.land" variant="outline-primary" size="sm">
            Website
          </b-button>
        </template>
      </Project>
      <Project name="What's in Standard?" logo="../assets/images/whatsinstandard.svg" status="Current">
        <template v-slot:tags>
          <b-badge pill variant="info">20k MAU</b-badge>
          <b-badge pill>Flutter</b-badge>
          <b-badge pill>Progressive Web App</b-badge>
        </template>
        <template v-slot:description>
          A reference tool for Magic: The Gathering players.
        </template>
        <template v-slot:buttons>
          <b-button href="https://whatsinstandard.com/" variant="outline-primary" size="sm">
            Website
          </b-button>
          <b-button href="https://play.google.com/store/apps/details?id=com.whatsinstandard" variant="outline-primary" size="sm">
            Android app
          </b-button>
          <b-button href="https://github.com/glacials/whatsinstandard" variant="outline-primary" size="sm">
            Source
          </b-button>
        </template>
      </Project>
      <Project name="Housebot" logo="../assets/images/housebot.svg" status="Archived">
        <template v-slot:tags>
          <b-badge pill>CoffeeScript</b-badge>
          <b-badge pill>Not a great idea</b-badge>
          <b-badge pill>Made it anyway</b-badge>
        </template>
        <template v-slot:description>
          A Twitch bot that lets chat participants vote to control your lights.
        </template>
        <template v-slot:buttons>
          <b-button href="https://github.com/glacials/housebot" variant="outline-primary" size="sm">
            Source
          </b-button>
        </template>
      </Project>
    </b-row>
  </b-container>
</template>

<script>
import Project from './Project.vue'

export default {
  components: {
    Project,
  },
  name: 'Projects',
}
</script>

<style scoped>
</style>
