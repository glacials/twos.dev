<template>
  <div id="app">
    <Summary />
    <Jobs />
    <Projects />
    <Interests />
    <Education />
    <FollowUp class="d-print-none" />
    <Footer class="d-print-none" />
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

<style lang="scss">
@import 'assets/custom-vars.scss';

@import '~bootstrap/scss/bootstrap.scss';
@import '~bootstrap-vue/src/index.scss';


@font-face {
  font-family: 'Barlow';
  src: url('./assets/fonts/Barlow-Regular.ttf');
};

@font-face {
  font-family: 'Pattaya';
  src: url('./assets/fonts/Pattaya-Regular.ttf');
};

h1, h2, h3 {
  font-family: 'Pattaya';
}

.title {
  margin-top: 2em;
}
</style>
