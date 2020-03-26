<template>
  <div class="mb-4">
    <div class="mb-1">
      <b-button href="https://calendly.com/glacials" class="d-print-none" variant="primary">
        Schedule a call
      </b-button>
      <b-button href="mailto:qhiiyr@gmail.com">
        Email me
        <span class="d-none d-print-inline">
          &mdash; qhiiyr@gmail.com
        </span>
      </b-button>
    </div>
    <div>
      <b-button href="../ben-carlsson-resume-2020-03.pdf" variant="link">
        <font-awesome-icon :icon="['fas', 'download']" />
        Download my resume
      </b-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CallToAction',
}
</script>

<style>
</style>
