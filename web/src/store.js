import { createStore } from 'vuex';
import { Room } from './models';
import { ParseGMCP } from './gmcp';
import app from "./main";

export const store = createStore({
  state() {
    return {
      isProduction: process.env.NODE_ENV === "production",
      isConnected: false,
      reconnectError: false,
      mudName: 'CloudRain',
      mudVesion: 'v0.1.0',
      pingTime: 0,
      settings: { lines: 100 },
      gameText: "",
      allowGlobalHotkeys: true,
      forceInputFocus: { forced: false, text: '' },
      commandDictionary: [],
      commandHistory: [],
      minimapData: { name: '', rooms: [] },
      characterLocation: { x: 0, y: 0, z: 0 },
      areaTitle: '',
      roomTitle: '',
      roomObjects: [],
      autoLoginToken: "",
      playerInfo: { uuid: '', name: '' },
    }
  },
  mutations: {
    SOCKET_ONOPEN(state) {
      state.isConnected = true;
    },
    SOCKET_ONCLOSE(state) {
      if (state.isConnected) {
        state.isConnected = false;
        state.gameText = app.app.config.globalProperties.$t('socket.closed');
      } else {
        state.gameText = app.app.config.globalProperties.$t('socket.not-established');
      }
    },
    SOCKET_ONERROR(state, event) {
      console.error(state, event);
    },
    SOCKET_ONMESSAGE(state, message) {
      try {
        switch (message.event) {
          case "text":
            this.dispatch("showText", message.content);
            break;
          case "ping":
            console.log("ping...");
            break;
          case "gmcp":
            ParseGMCP(message.content);
            break;
          default:
            break;
        }
      } catch (e) {
        console.log(e);
      }
    },
    SOCKET_RECONNECT(state, count) {
      console.info("reconnect...", state, count);
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.reconnectError = true;
    },

    SET_ALLOW_GLOBAL_HOTKEYS: (state, allow) => {
      state.allowGlobalHotkeys = allow;
    },
    SET_FORCE_INPUT_FOCUS: (state, data) => {
      state.forceInputFocus = data;
    },

    SET_COMMAND_DICTIONARY: (state, dictionary) => {
      state.commandDictionary = dictionary;
    },
    APPEND_COMMAND_HISTORY: (state, command) => {
      state.commandHistory.push(command);
    },
    ADD_GAME_TEXT: (state, text) => {
      state.gameText = text;
    },

    SET_MINIMAP_DATA: (state, minimapData) => {
      state.minimapData = {
        name: minimapData.name,
        rooms: [],
      };
      minimapData.rooms.forEach(r => {
        state.minimapData.rooms.push(new Room(r));
      });
    },
    SET_CHARACTER_LOCATION: (state, loc) => {
      state.characterLocation = loc;
    },
    SET_AREA_TITLE: (state, title) => {
      state.areaTitle = title;
    },
    SET_ROOM_TITLE: (state, title) => {
      state.roomTitle = title;
    },
    SET_ROOM_OBJECTS: (state, objects) => {
      state.roomObjects = objects;
    },

    SET_PLAYER_INFO: (state, playerInfo) => {
      state.playerInfo = playerInfo;
    },

    SET_SETTINGS: (state, settings) => {
      state.settings = settings;
    },
  },
  actions: {
    setAllowGlobalHotkeys: ({ commit }, payload) => {
      commit('SET_ALLOW_GLOBAL_HOTKEYS', payload);
    },
    setForceInputFocus: ({ commit }, payload) => {
      commit('SET_FORCE_INPUT_FOCUS', payload);
    },

    sendCommand: ({ state, commit }, payload) => {
      if (!state.isConnected) {
        return;
      }

      commit('APPEND_COMMAND_HISTORY', payload.command);

      let echoCmd = payload.command;
      if (typeof payload.hidden !== 'boolean' || !payload.hidden) {
        commit('ADD_GAME_TEXT', `${echoCmd}\r\n`);
      }

      app.app.config.globalProperties.$socket.sendObj({
        type: "cmd",
        content: payload.command
      });
    },

    showText: ({ commit }, payload) => {
      commit('ADD_GAME_TEXT', payload);
    },
    setMapData: ({ commit }, payload) => {
      commit('SET_MINIMAP_DATA', payload.data);
    },
    setCharacterLocation: ({ commit }, payload) => {
      commit('SET_CHARACTER_LOCATION', payload.data);
    },
    setAreaTitle: ({ commit }, payload) => {
      commit('SET_AREA_TITLE', payload.data);
    },
    setRoomTitle: ({ commit }, payload) => {
      commit('SET_ROOM_TITLE', payload.data);
    },
    setRoomObjects: ({ commit }, payload) => {
      commit('SET_ROOM_OBJECTS', payload.data);
    },

    setPlayerInfo: ({ commit }, payload) => {
      commit('SET_PLAYER_INFO', payload.data);
    },

    setSettings: ({ commit }, payload) => {
      commit('SET_SETTINGS', payload.data);
    },
  },
})
