<template>
  <b-container>
    <h2 class="title">Follow up</h2>
    <p>If you're in the Seattle area, you can win my heart by buying me a beer.</p>
    <CallToAction />
  </b-container>
</template>

<script>
import CallToAction from './CallToAction.vue'

export default {
  components: {
    CallToAction,
  },
  name: 'FollowUp',
}
</script>

<style scoped>
</style>
