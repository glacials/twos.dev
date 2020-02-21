<template>
  <b-col>
    <b-card no-body>
      <b-card-img :src="logo" :alt="`${logo} logo`" class="d-none d-md-block d-print-none" top />
      <b-card-body>
        <div class="justify-content-between d-flex">
          <div>
            <b-card-title class="justify-content-between d-flex mb-0">
              {{ name }}
            </b-card-title>
            <b-card-sub-title sub-title-tag="small">
              {{ status }}
            </b-card-sub-title>
          </div>
          <b-img :src="logo" class="d-md-none small-logo" />
        </div>
        <b-card-text>
          <slot name="description" />
        </b-card-text>
        <slot name="tags" />
      </b-card-body>
      <b-card-footer>
        <b-button-group>
          <slot name="buttons" />
        </b-button-group>
      </b-card-footer>
    </b-card>
  </b-col>
</template>

<script>
export default {
  name: 'Project',
  props: ['name', 'logo', 'status'], // 'logo' can only work because of the 'Project': 'logo' line in vue.config.js
}
</script>

<style scoped>
.badge {
  margin-right: .25em;
}

.col {
  margin-bottom: 2em;
}

.small-logo {
  height: 3em;
  width: 3em;
}
</style>
