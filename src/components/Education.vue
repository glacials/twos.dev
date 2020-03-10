<template>
  <b-container>
    <h2 class="title">Education</h2>
    <b-media no-body>
      <b-media-aside class="mr-0 mr-md-3">
        <img src="../assets/images/chico.webp" class="school-logo d-none d-md-flex" alt="California State University, Chico logo" />
      </b-media-aside>
      <b-media-body>
        <h5 class="mb-0 justify-content-between d-flex">
          <span>
            <b-link href="https://csuchico.edu" class="text-dark">California State University, Chico</b-link>
            <small class="text-muted text-nowrap">2010 - 2013</small>
            <small class="text-secondary d-block">B.Sc. Computer Science</small>
          </span>
          <img src="../assets/images/chico.webp" class="company-logo d-flex d-md-none" alt="California State University, Chico logo logo" style="height: 2em; width: 2em" />
        </h5>
        <ul>
          <li>President, Chico State Linux User Group</li>
          <li>Member, computer science honors society UPE</li>
        </ul>
      </b-media-body>
    </b-media>
  </b-container>
</template>

<script>
export default {
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
