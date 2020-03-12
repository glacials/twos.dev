<template>
  <b-media no-body>
    <b-media-aside class="mr-0 mr-md-3">
      <b-img :src="logo" class="career-logo d-none d-md-flex" :alt="`${name} logo`" />
    </b-media-aside>
    <b-media-body>
      <h5 class="mb-0 justify-content-between d-flex">
        <span>
          <b-link :href="url" class="text-dark">{{ name }}</b-link>
          <small class="text-muted text-nowrap">{{ dates }}</small>
          <small class="text-secondary d-block">{{ role }}</small>
        </span>
        <b-img :src="logo" class="career-logo d-flex d-md-none" :alt="`${name} logo`" style="height: 2em; width: 2em" />
      </h5>
      <slot name="default" />
    </b-media-body>
  </b-media>
</template>

<script>
export default {
  name: 'CareerEntry',
  props: ['logo', 'dates', 'role', 'name', 'url'],
}
</script>

<style scoped>
.career-logo {
  width: 5em;
}

.media {
  margin-bottom: 3em;
}
</style>
