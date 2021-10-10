<template>
  <div class="root" :style="{ height: containerHeight }">
    <div class="scrollable-container" ref="mainTextContainer" />
  </div>
</template>

<script>
require("xterm/css/xterm.css");

import { mapState } from "vuex";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { Unicode11Addon } from "xterm-addon-unicode11";
import { UnicodeV11AddonAmbiguous } from "../utils/xterm-addon/UnicodeV11AddonAmbiguous";

export default {
  name: "MainText",
  data: function () {
    return {
      lineNumber: 0,
      lastItemTooltipUUID: "",
    };
  },
  props: {
    windowWidth: Number,
    windowHeight: Number,
  },
  computed: {
    ...mapState(["gameText", "settings"]),
    containerWidth() {
      const width = this.windowWidth;
      return `${width}px`;
    },
    containerHeight() {
      const height = this.windowHeight - 37 - 45 - 2 - 35;
      return `${height}px`;
    },
  },
  mounted() {
    // @TODO:ambiguous character width
    // @see: https://github.com/xtermjs/xterm.js/issues/2668
    const term = new Terminal({
      fontFamily:
        "'Noto Sans Mono CJK SC', 'PingFang SC', 'STHeitiSC-Light', SimHei, NSimSun, monospace",
      lineHeight: 1,
    });

    if (this.$store.state.settings.useAmbiguousReplace) {
      const unicode11Addon = new Unicode11Addon();
      term.loadAddon(unicode11Addon);
      term.unicode.activeVersion = "11";
    } else {
      const unicode11Addon = new UnicodeV11AddonAmbiguous();
      term.loadAddon(unicode11Addon);
      term.unicode.activeVersion = "11-ambiguous";
    }

    const fitAddon = new FitAddon();
    this.fitAddon = fitAddon;
    term.loadAddon(fitAddon);
    term.open(this.$refs.mainTextContainer);
    fitAddon.fit();

    this.terminal = term;
    this.fit();
  },
  watch: {
    gameText: function (msg) {
      if (msg === "") {
        return;
      }
      this.terminal.write(msg);
    },
    containerWidth: function () {
      this.fit();
    },
    containerHeight: function () {
      this.fit();
    },
  },
  methods: {
    fit() {
      const term = this.terminal;
      const fitAddon = this.fitAddon;
      term.element.style.display = "none";
      setTimeout(function () {
        fitAddon.fit();
        term.element.style.display = "";
        term.refresh(0, term.rows - 1);
      }, 10);
    },
  },
};
</script>

<style scoped lang="scss">
@import "~@/styles/common.module";

.root {
  box-sizing: border-box;
  display: flex;

  /*border: $defaultBorder;*/
  @include defaultBorderImage;
}

.root .item-drag-overlay {
  position: absolute;
  z-index: 100;
  top: 0;
  bottom: 0;
  background-color: #000000b8;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 22px;
  color: #666;
  border: 2px dashed #666;
}

.root .item-drag-overlay.item-over {
  background-color: #1d1c1cb8;
  border: 2px dashed #aaa;
  color: #aaa;
}

.scrollable-container {
  flex-grow: 1;
}

.scrollable-contai .terminal {
  color: #cacaca;
  user-select: text;
  font-size: 13px;
  font-family: $monoFont;
}
</style>
