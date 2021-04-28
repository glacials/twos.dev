<template>
  <b-container>
    <h2 class="title">Education</h2>
    <CareerEntry
      name="California State University, Chico"
      url="https://csuchico.edu"
      logo="../assets/images/chico.webp"
      dates="2010 - 2013"
      role="B.Sc. Computer Science"
    >
      <ul>
        <li>President, Chico State Linux User Group</li>
        <li>Member, computer science honors society UPE</li>
      </ul>
    </CareerEntry>
  </b-container>
</template>

<script>
import CareerEntry from './CareerEntry.vue'

export default {
  components: {
    CareerEntry,
  },
  name: 'Education',
}
</script>

<style scoped>
.school-logo {
  width: 5em;
}

ul {
  list-style-type: none;
  margin: 0;
  padding: 0;
}
</style>
