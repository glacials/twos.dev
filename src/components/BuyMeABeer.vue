<template>
  <div class="mb-4">
    <b-button href="mailto:qhiiyr@gmail.com?subject=Let me buy you a beer" variant="primary" class="d-print-none">
      Buy me a beer
    </b-button>
    <b-button href="mailto:qhiiyr@gmail.com">
      Email me
      <span class="d-none d-print-inline">
        &mdash; qhiiyr@gmail.com
      </span>
    </b-button>
    <b-button variant="link" class="text-muted" v-b-popover.hover.top="'(925) 979-5797'">
      <font-awesome-icon :icon="['fas', 'mobile-alt']" />
      <span class="d-none d-print-inline">
        (925) 979-5797
      </span>
    </b-button>
  </div>
</template>

<script>
export default {
  name: 'BuyMeABeer',
}
</script>

<style scoped>
</style>
