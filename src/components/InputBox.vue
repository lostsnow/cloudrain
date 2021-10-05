<template>
  <div class="root" :class="{ active: isFocused }">
    <input
      class="input-box"
      ref="inputBox"
      type="text"
      v-model="textToSend"
      @keyup.enter="handleSendText"
      @keyup.escape="handleRemoveFocus"
      @keydown="handleKeyDown"
      @focus="handleFocus"
      @blur="handleBlur"
    />
    <div
      class="hotkey-overlay"
      v-if="!isFocused"
      @click="handleHotkeyOverlayClick"
    >
      {{ $t('ui.input-box.click-to-active') }}
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "InputBox",
  data: () => {
    return {
      textToSend: "",
      password: "",
      isFocused: false,
      lastCommandHistoryIndex: -1,
    };
  },
  computed: {
    expandedCommandDictionary: function () {
      const dict = [];
      this.commandDictionary.forEach((d) => dict.push(d));
      this.commandDictionary.forEach((cmd) => {
        if (cmd.altNames) {
          cmd.altNames.forEach((alt) => {
            let newCmd = Object.assign({}, cmd);
            newCmd.name = alt;
            newCmd.altNames = [];
            dict.push(newCmd);
          });
        }
      });
      return dict;
    },
    ...mapState(["forceInputFocus", "commandHistory", "commandDictionary"]),
  },
  mounted() {
    this.$refs["inputBox"].focus();
  },
  watch: {
    forceInputFocus: function (data) {
      if (data.forced) {
        this.$refs["inputBox"].focus();
        if (data.text) {
          this.textToSend = data.text;
        }
        this.$store.dispatch("setForceInputFocus", { forced: false });
      }
    },
  },
  methods: {
    selectAll: function () {
      this.$refs["inputBox"].select();
    },

    getLastCommand() {
      let retrieveIndex = 0;

      if (this.lastCommandHistoryIndex === -1) {
        retrieveIndex = this.commandHistory.length - 1;
        this.lastCommandHistoryIndex = retrieveIndex;
      } else if (this.lastCommandHistoryIndex > 0) {
        retrieveIndex = this.lastCommandHistoryIndex - 1;
        this.lastCommandHistoryIndex = retrieveIndex;
      }

      return this.commandHistory[retrieveIndex];
    },

    getNextCommand() {
      let retrieveIndex = this.lastCommandHistoryIndex;

      if (retrieveIndex === -1) {
        retrieveIndex = this.commandHistory.length - 1;
        this.lastCommandHistoryIndex = retrieveIndex;
      } else if (
        this.lastCommandHistoryIndex <
        this.commandHistory.length - 1
      ) {
        retrieveIndex = this.lastCommandHistoryIndex + 1;
        this.lastCommandHistoryIndex = retrieveIndex;
      }

      return this.commandHistory[retrieveIndex];
    },

    handleSendText() {
      let command = this.textToSend;

      this.$store.dispatch("sendCommand", {
        command: command,
      });

      this.textToSend = "";
      this.lastCommandHistoryIndex = -1;
    },

    handleRemoveFocus(event) {
      this.$nextTick(() => {
        event.target.blur();
      });
    },

    handleFocus() {
      this.isFocused = true;
      this.$store.dispatch("setAllowGlobalHotkeys", false);
    },

    handleBlur() {
      this.$store.dispatch("setAllowGlobalHotkeys", true);
      this.$nextTick(() => {
        this.isFocused = false;
      });
    },

    handleKeyDown(e) {
      if (e.key === "ArrowUp") {
        this.textToSend = this.getLastCommand();
        setTimeout(this.selectAll, 10);
      } else if (e.key === "ArrowDown") {
        this.textToSend = this.getNextCommand();
        setTimeout(this.selectAll, 10);
      }
    },

    handleHotkeyOverlayClick() {
      this.$refs["inputBox"].focus();
    },
  },
};
</script>

<style scoped lang="scss">
@import "@/styles/common.module";
$height: 35px;

.root {
  position: relative;
  background: $bg-color;
  padding-left: 5px;
  border: $defaultBorder;
}

.input-box {
  display: block;
  width: 100%;
  padding: 0;
  border-width: 0;
  border: 0;
  height: $height;
  color: $defaultTextColor;
  font-family: $monoFont;
  font-weight: 500;
  font-size: 13px;
}

.root.active {
  border: $defaultBorder;
}

.root.active .input-box {
  background-color: $bg-color;
}

.input-box:focus {
  margin: 0;
  padding: 0;
  border: 0;
  outline: none;
}

.hotkey-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 10;
  height: $height;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: $bg-color;
  color: $defaultTextColor;
}

.hotkey-overlay:hover {
  cursor: pointer;
  color: $defaultTextColor;
}

.command-helper-overlay {
  position: absolute;
  background: $bg-color;
  background: $bg-color;
  width: 99%;
  padding: 20px 5px 10px 5px;
  font-size: 12px;
}
</style>
