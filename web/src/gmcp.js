import { store } from './store'

const CorePing = 'Core.Ping';
const CoreInfo = 'Core.Info';
const CharRegister = 'Char.Register';
const CharLogin = 'Char.Login';

export function ParseGMCP(msg) {
  try {
    if (!msg) {
      return;
    }
    let gmcpInfo = JSON.parse(msg);
    let key = gmcpInfo["key"];
    let value = gmcpInfo["value"];
    if (!key) {
      return;
    }

    switch (key) {
      case CorePing:
        break;
      case CoreInfo:
        if (value.NAME) {
          store.state.mudName = value.NAME;
        }
        break;
      case CharRegister:
      case CharLogin:
        if (typeof value.code === "undefined") {
          return;
        }

        if (value.code == 0) {
          store.commit("SET_LOGIN_TOKEN", {id: value.id, token: value.token});
          return;
        }
        switch (value.err) {
          case "ERR_REGISTER":
              // @TODO: register failed
            break;
          case "ERR_LOGIN_PASS":
            // @TODO: login by pass failed
            break;
          case "ERR_LOGIN_TOKEN":
              // @TODO: redirect login by pass
            break;
        }
        break;
      default:
        console.log("gmcp:", gmcpInfo);
        break;
    }
  } catch (e) {
    //
  }
}

export function SendGMCPString(key, payload) {
  store.dispatch("sendCommand", {
    type: "gmcp",
    command: `${key} ${payload}`,
  });
}

export function SendGMCP(key, payload) {
  let p = JSON.stringify(payload);
  SendGMCPString(key, p);
}
