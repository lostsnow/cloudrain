import { store } from './store'

const CoreInfo = 'Core.Info';

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
      case CoreInfo:
        if (!value) {
          return;
        }
        store.state.mudName = value["NAME"];
        break;
      default:
        console.log("gmcp:", gmcpInfo);
        break;
    }
  } catch (e) {
    //
  }
}
