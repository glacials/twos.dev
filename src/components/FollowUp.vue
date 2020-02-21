<template>
  <b-container>
    <h2 class="title">Follow up</h2>
    <p>If you're in the Seattle area, you can win my heart by buying me a beer.</p>
    <BuyMeABeer />
  </b-container>
</template>

<script>
import BuyMeABeer from './BuyMeABeer.vue'

export default {
  components: {
    BuyMeABeer,
  },
  name: 'FollowUp',
}
</script>

<style scoped>
</style>
