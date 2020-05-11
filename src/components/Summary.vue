<template>
  <b-jumbotron header-level="4">
    <template v-slot:header>
      <div>
        <img src="../assets/images/me.webp" alt="Photograph of Ben Carlsson" style="height: 3em" />
      </div>
      Ben Carlsson
    </template>
    <template v-slot:lead>
      <div class="mb-4">
        I'm a senior full-stack engineer in Seattle. I was an early employee at Twitch
        and I've worked on everything from social networks to authentication-as-a-service.
      </div>
    </template>
    <small>
      <CallToAction />
      <h4>
        <b-link href="https://twitter.com/glcls" class="social" v-b-tooltip.nofade title="glcls">
          <font-awesome-icon :icon="['fab', 'twitter']" />
          <span class="d-none d-print-inline">@glcls</span>
          <span class="sr-only">Twitter: @glcls</span>
        </b-link>
        <b-link href="https://www.twitch.tv/glacials" class="social" v-b-tooltip.nofade title="glacials">
          <font-awesome-icon :icon="['fab', 'twitch']" />
          <span class="d-none d-print-inline">Glacials</span>
          <span class="sr-only">Twitch: Glacials</span>
        </b-link>
        <b-link href="https://github.com/glacials" class="social" v-b-tooltip.nofade title="glacials">
          <font-awesome-icon :icon="['fab', 'github']" />
          <span class="d-none d-print-inline">glacials</span>
          <span class="sr-only">GitHub: glacials</span>
        </b-link>
        <b-link href="https://linkedin.com/in/ben" class="social" v-b-tooltip.nofade title="ben">
          <font-awesome-icon :icon="['fab', 'linkedin']" />
          <span class="d-none d-print-inline">ben</span>
          <span class="sr-only">LinkedIn: ben</span>
        </b-link>
      </h4>
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
  margin: .25em;
}

.btn {
  margin: 1em .25em;
}

img {
  margin: auto;
  margin-bottom: 1em;
}

.jumbotron {
  text-align: center;
}

.jumbotron > * {
  max-width: 40rem;
  margin: auto;
  text-align: center;
}

.jumbotron > .display-4 {
  max-width: 100%;
}
</style>
