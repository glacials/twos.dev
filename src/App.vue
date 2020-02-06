<template>
  <div id="app">
    <Summary />
    <Jobs />
    <Projects />
    <Interests />
    <Education />
    <FollowUp />
    <Footer />
  </div>
</template>

<script>
import Education from './components/Education.vue'
import FollowUp from './components/FollowUp.vue'
import Footer from './components/Footer.vue'
import Interests from './components/Interests.vue'
import Jobs from './components/Jobs.vue'
import Projects from './components/Projects.vue'
import Summary from './components/Summary.vue'

export default {
  name: 'app',
  components: {
    Education,
    FollowUp,
    Footer,
    Interests,
    Jobs,
    Projects,
    Summary,
  }
}
</script>

<style>
.title {
  margin-top: 2em;
}
</style>
