<template>
  <div class="root">
    <div class="banner">{{ $t("map.room.target") }}</div>
    <div class="targets-list" @click="handleClick">
      <Target
        v-for="obj in sortedRoomObjects"
        :key="obj.uuid"
        :uuid="obj.uuid"
        :name="obj.name"
        :pictureKey="obj.picture"
        :objectType="obj.type"
        :title="obj.title"
        :color="obj.color"
        :visible="obj.visible"
      />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import Target from "@/components/Target";

export default {
  name: "RoomTargets",
  components: {
    Target,
  },
  computed: {
    ...mapState(["roomObjects"]),
    sortedRoomObjects: function () {
      return this.roomObjects.slice().sort((a, b) => {
        if (b.sort > a.sort) {
          return 1;
        } else if (a.sort > b.sort) {
          return -1;
        } else {
          return 0;
        }
      });
    },
  },
  mounted() {
    // TODO: for test
    this.$store.dispatch("setRoomObjects", {
      data: [
        {
          name: "ZhangTian Ping",
          title: "掌天瓶",
        },
      ],
    });
  },
};
</script>

<style scoped lang="scss">
@import "@/styles/common";

.root {
  height: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  @include defaultBorderImage;
}

.banner {
  text-align: center;
  font-size: 1.2em;
  font-weight: 500;
  margin-bottom: 3px;
}

.targets-list {
  overflow-y: scroll;
  overflow-x: hidden;
  flex-grow: 1;
  scrollbar-width: thin;
}
</style>
