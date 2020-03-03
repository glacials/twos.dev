<template>
  <b-jumbotron header-level="4">
    <template v-slot:header>
      <img src="../assets/me.jpg" />
      Ben "Glacials" Carlsson
    </template>
    <template v-slot:lead>
      <div class="mb-4">
        I'm a senior backend and full-stack engineer in Seattle. I was the 91st employee at Twitch. I've worked on
        everything from social networks to authentication-as-a-service.
      </div>
    </template>
    <small>
      <CallToAction />
      <h6>
        <b-link href="https://twitter.com/glcls" class="social" v-b-tooltip.nofade title="glcls">
          <font-awesome-icon :icon="['fab', 'twitter']" />
          <span class="d-none d-print-inline">@glcls</span>
        </b-link>
        <b-link href="https://www.twitch.tv/glacials" class="social" v-b-tooltip.nofade title="glacials">
          <font-awesome-icon :icon="['fab', 'twitch']" />
          <span class="d-none d-print-inline">Glacials</span>
        </b-link>
        <b-link href="https://github.com/glacials" class="social" v-b-tooltip.nofade title="glacials">
          <font-awesome-icon :icon="['fab', 'github']" />
          <span class="d-none d-print-inline">glacials</span>
        </b-link>
        <b-link href="https://linkedin.com/in/ben" class="social" v-b-tooltip.nofade title="ben">
          <font-awesome-icon :icon="['fab', 'linkedin']" />
          <span class="d-none d-print-inline">ben</span>
        </b-link>
      </h6>
    </small>
    <small class="text-muted">
      <font-awesome-icon icon="comment-dots" />
      I speak English and a little Chinese.<br />
      <font-awesome-icon icon="comment-dots" />
      我会说英文，也会说一点中文。
    </small>
  </b-jumbotron>
</template>

<script>
import CallToAction from './CallToAction.vue'

export default {
  components: {
    CallToAction,
  },
  name: 'Summary',
}
</script>

<style scoped>
.social {
  margin-right: .25em;
}

.btn {
  margin: 1em .25em;
}

img {
  display: block;
  height: 3em;
  margin: auto;
  margin-bottom: 1em;
  width: 3em;
}

.jumbotron {
  text-align: center;
}

.jumbotron > * {
  max-width: 40rem;
  margin: auto;
  text-align: center;
}
</style>
