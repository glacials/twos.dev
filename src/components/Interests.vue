<template>
  <b-container>
    <h2 class="title">Interests</h2>
    <b-row cols="1" cols-sm="2" cols-xl="4">
      <b-col>
        <b-card header="Industry" no-body>
          <b-list-group flush>
            <b-list-group-item>Climate crisis</b-list-group-item>
            <b-list-group-item>Gaming</b-list-group-item>
            <b-list-group-item>Open source</b-list-group-item>
          </b-list-group>
        </b-card>
      </b-col>
      <b-col>
        <b-card header="Technology" no-body>
          <b-list-group flush>
            <b-list-group-item>Ruby on Rails</b-list-group-item>
            <b-list-group-item>Flutter</b-list-group-item>
            <b-list-group-item>Vue</b-list-group-item>
          </b-list-group>
        </b-card>
      </b-col>
      <b-col>
        <b-card header="Recreation" no-body>
          <b-list-group flush>
            <b-list-group-item>Beer</b-list-group-item>
            <b-list-group-item>Whiskey</b-list-group-item>
            <b-list-group-item>Speedrunning</b-list-group-item>
          </b-list-group>
        </b-card>
      </b-col>
      <b-col>
        <b-card no-body>
          <template v-slot:header>
            Career
            <span class="text-muted">(non-engineering)</span>
          </template>
          <b-list-group flush>
            <b-list-group-item>Wearing many hats</b-list-group-item>
            <b-list-group-item>UX design</b-list-group-item>
            <b-list-group-item>Product management</b-list-group-item>
          </b-list-group>
        </b-card>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
export default {
  name: 'Interests',
}
</script>

<style scoped>
.col {
  margin-bottom: 2em;
}
</style>
