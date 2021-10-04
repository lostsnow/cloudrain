<template>
  <div>
    <div
      class="targets-container"
      ref="container"
      :style="{
        borderColor: color ? `#${color}` : '',
        borderStyle: visible ? 'solid' : 'dashed',
        opacity: visible ? 1 : 0.7,
      }"
    >
      <div class="picture">
        <div
          class="picture-container"
          :style="{ backgroundImage: getBackgroundUrl() }"
        ></div>
      </div>
      <div class="name">
        <div class="name-container">
          <div>{{ title }}</div>
          <div class="alt">{{ name }}</div>
        </div>
        <div
          class="you"
          :class="{ selected: uuid === objectTargetUUID }"
          v-if="uuid === playerInfo.id"
        >
          {{ $t('char.you') }}
        </div>
      </div>
      <div
        class="overlay"
      ></div>
    </div>
  </div>
</template>

<script>
const OBJECT_TYPE_CHARACTER = 0;
const OBJECT_TYPE_MOB = 1;
const OBJECT_TYPE_ITEM = 2;

import { mapState } from "vuex";

export default {
  name: "Target",
  props: [
    "uuid",
    "name",
    "objectType",
    "pictureKey",
    "title",
    "color",
    "visible",
  ],
  computed: {
    ...mapState(["playerInfo", "objectTargetUUID"]),
  },
  watch: {
    objectTargetUUID: function (target) {
      if (this.uuid === target) {
        this.$refs["container"].classList.add("selected");
      } else {
        this.$refs["container"].classList.remove("selected");
      }
    },
  },
  mounted() {
    switch (this.objectType) {
      case OBJECT_TYPE_CHARACTER:
        this.$refs["container"].classList.add("is-character");
        break;
      case OBJECT_TYPE_MOB:
        this.$refs["container"].classList.add("is-mob");
        break;
      case OBJECT_TYPE_ITEM:
        this.$refs["container"].classList.add("is-item");
        if (this.color.length > 0) {
          this.$refs["container"].style.borderColor = this.color;
        }
        break;
    }

    if (this.uuid === this.objectTargetUUID) {
      this.$refs["container"].classList.add("selected");
    }
  },
  methods: {
    getBackgroundUrl() {
      if (!this.pictureKey) {
        // @TODO: default icon
        return ``;
      }
      return `url(/oi/${this.pictureKey})`;
    },
  },
};
</script>


<style lang="scss" scoped>
.targets-container {
  transition: all 0.1s ease-in-out;
  transform: scale(1);
  display: flex;
  border-width: 0 !important;

  &.can-drop-item {
    transform: scale(1.1) !important;
  }

  &.selected {
    border: 1px solid #ffeb3b !important;
    background-color: #231f00;
  }

  &.mouse-down {
    transform: scale(1.01) !important;
  }

  &.is-character {
    border: 1px solid #353535;
  }

  &.is-mob {
    border: 1px solid #673604;

    .name {
      color: #d48a3e;
    }
  }

  &.is-item {
    border: 1px solid #fff;

    .name {
      color: #fff;
    }
  }

  &:hover {
    cursor: pointer;
  }

  .picture {
    flex-basis: 50px;

    .picture-container {
      height: 50px;
      box-shadow: inset 0px 0px 5px 0px #3a3a3a;
      background-size: contain;
    }
  }

  .name {
    flex-grow: 1;
    display: flex;
    align-items: center;
    margin-left: 10px;

    .name-container {
      font-weight: 600;

      .alt {
        font-weight: 400;
        font-size: 12px;
      }
    }

    .you {
      position: absolute;
      right: -1px;
      top: -1px;
      background-color: #353535;
      padding: 2px 5px;
      border: 1px solid #353535;
      text-transform: uppercase;
      font-size: 12px;
      transition: all 0.1s ease-in-out;
      border-bottom-left-radius: 3px;

      &.selected {
        background-color: #eedb38;
        border: 1px solid #eedb38;
        color: #000;
      }
    }
  }

  .overlay {
    position: absolute;
    top: 0px;
    left: 0px;
    height: 100%;
    width: 100%;
    z-index: 999;
  }
}
</style>