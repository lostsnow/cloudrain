<template>
  <div class="root">
    <div class="inner-container">
      <div class="ping">Ping: {{ pingTime }}ms</div>
      <div class="version">
        {{ mudName }}
        <span v-if="mudVesion && mudVesion.length > 0">
          ({{ mudVesion }})
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "StatusBar",
  computed: {
    ...mapState(["pingTime", "mudName", "mudVesion"]),
  },
};
</script>

<style scoped lang="scss">
@import "~@/styles/common.module";

.root {
  display: flex;
  height: 100%;
  align-items: center;
  font-size: 12px;
  box-sizing: border-box;

  .inner-container {
    display: flex;
    width: 100%;
    padding: 0 5px;

    .ping {
      flex-shrink: 1;
    }

    .version {
      flex-grow: 1;
      text-align: right;
    }
  }
}
</style>
