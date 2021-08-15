<template>
  <div id="app">
    <div class="container-wrapper">
      <div class="container-left">
        <div class="container-maintext">
          <MainText :windowHeight="windowHeight" />
        </div>
        <div class="container-input">
          <InputBox />
        </div>
        <div class="container-bars">
          <Vitals />
        </div>
      </div>
      <div class="container-right" :style="{ display: rightSidebar }">
        <div class="container-minimap">
          <Minimap />
        </div>
        <div class="container-targets">
          <RoomTargets />
        </div>
      </div>
    </div>
    <div class="status-bar-container">
      <StatusBar />
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import MainText from "@/components/MainText";
import InputBox from "@/components/InputBox";
import Vitals from "@/components/Vitals";
import Minimap from "@/components/Minimap";
import RoomTargets from "@/components/RoomTargets";
import StatusBar from "@/components/StatusBar";

export default {
  name: "App",
  components: {
    MainText,
    InputBox,
    Vitals,
    Minimap,
    RoomTargets,
    StatusBar,
  },
  data: () => {
    return {
      windowHeight: 0,
      windowWidth: 0,
      rightSidebar: "flex",
    };
  },
  computed: {
    ...mapState(["allowGlobalHotkeys"]),
  },
  methods: {
    onWindowResize() {
      this.windowHeight = window.innerHeight;
      this.windowWidth = window.innerWidth;

      if (this.windowWidth < 784) {
        this.showRightSidebar = false;
        this.rightSidebar = "none";
      } else {
        this.showRightSidebar = true;
        this.rightSidebar = "flex";
      }
    },

    onKeyUp(event) {
      if (!this.allowGlobalHotkeys) {
        return;
      }

      let moveCommand = "";

      switch (event.key.toLowerCase()) {
        case "w":
          moveCommand = "go north";
          break;
        case "a":
          moveCommand = "go west";
          break;
        case "s":
          moveCommand = "go south";
          break;
        case "d":
          moveCommand = "go east";
          break;
        case "q":
          moveCommand = "go up";
          break;
        case "e":
          moveCommand = "go down";
          break;
        case "enter":
          this.$store.dispatch("setForceInputFocus", { forced: true });
          break;
      }

      if (moveCommand.length > 0) {
        this.$store.dispatch("sendCommand", {
          command: moveCommand,
          hidden: true,
        });
      }
    },
  },
  mounted() {
    this.onWindowResize();

    window.addEventListener("resize", this.onWindowResize);

    window.addEventListener("keyup", this.onKeyUp);
  },
  unmounted() {
    window.removeEventListener("resize", this.onWindowResize);
  },
};
</script>

<style lang="scss">
@import "~@/styles/common";
$backgroundNormal: #111;
$backgroundLight: #1b1b1b;
$sidebarWidth: 250px;

html,
body {
  padding: 0;
  margin: 0;
  height: 100%;
  background-color: $bg-color;
  user-select: none;
}

::-webkit-scrollbar {
  width: 3px;
  height: 3px;
}
::-webkit-scrollbar-button {
  background-color: #666;
}
::-webkit-scrollbar-track {
  background-color: #646464;
}
::-webkit-scrollbar-track-piece {
  background-color: #111;
}
::-webkit-scrollbar-thumb {
  height: 50px;
  background-color: #333;
  border-radius: 0px;
}
::-webkit-scrollbar-corner {
  background-color: #646464;
}
::-webkit-resizer {
  background-color: #666;
}

#app {
  font-family: "Montserrat", sans-serif;
  font-size: 14px;
  margin: 0;
  padding: 0;
  color: $defaultTextColor;
  display: flex;
  flex-direction: column;
  height: 100%;

  .container-wrapper {
    flex-grow: 1;
    display: flex;

    .container-left {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      position: relative;
      padding: 4px 2px;

      .container-maintext {
        flex-grow: 1;
        margin-bottom: 2px;
      }

      .container-input {
        flex-shrink: 1;
        margin-top: 2px;
        margin-bottom: 2px;
      }

      .container-bars {
        margin-top: 2px;
      }
    }

    .container-right {
      flex-basis: $sidebarWidth;
      min-width: $sidebarWidth;
      background-color: $bg-color-light;
      display: flex;
      flex-direction: column;
      padding: 4px;
      border-left: solid 4px $bg-color-dark;

      .container-minimap {
        flex-basis: 250px;
        margin-bottom: 2px;
      }

      .container-targets {
        flex-grow: 1;
        flex-basis: 100px; /* This can be any number; forces div to respect flex box height. */
        min-height: 100px;
        margin-top: 2px;
      }
    }
  }

  .status-bar-container {
    flex-basis: 30px;
    position: relative;
    background-color: $bg-color;
    padding: 2px;
    margin-top: 2px;
    background-image: url(../public/gfx/status-bg-01.png);
  }
}
</style>