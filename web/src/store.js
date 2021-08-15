import { createStore } from 'vuex';
import { Room } from './models';

export const store = createStore({
  state() {
    return {
      isProduction: process.env.NODE_ENV === "production",
      mudName: 'CloudRain',
      mudVesion: 'v0.1.0',
      pingTime: 0,
      isConnected: false,
      settings: { lines: 100 },
      gameText: [],
      allowGlobalHotkeys: true,
      forceInputFocus: { forced: false, text: '' },
      commandDictionary: [],
      commandHistory: [],
      minimapData: { name: '', rooms: [] },
      characterLocation: { x: 0, y: 0, z: 0 },
      areaTitle: '',
      roomTitle: '',
      roomObjects: [],
      playerInfo: { uuid: '', name: '' },
    }
  },
  mutations: {
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
      state.gameText.push({
        id: state.gameText.length,
        html: text
          .replace(/\n/g, "<br>")
          .replace(/\[b\]/g, "<span style='font-weight:600'>")
          .replace(/\[\/b\]/g, "</span>")
          .replace(/\[cmd=([^\]]*)\]/g, "<a href='#' class='inline-command' onclick='window.Armeria.$store.dispatch(\"sendCommand\", {command:\"$1\"})'>")
          .replace(/\[\/cmd\]/g, "</a>")
      });
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
        commit('ADD_GAME_TEXT', `<div class="inline-loopback">${echoCmd}</div>`);
      }
    },

    showText: ({ commit }, payload) => {
      commit('ADD_GAME_TEXT', payload.data);
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
