// Tomato object setup
if (!Tomato) var Tomato = {};

// @see: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
Tomato.WebSocketClient = function (options = {}) {
    this.number = 0;    // Message number
    this.autoReconnectRetry = 0;

    if (!options.autoReconnectMaxRetry) {
        options.autoReconnectMaxRetry = 10;
    }
    this.autoReconnectMaxRetry = options.autoReconnectMaxRetry;
    if (!options.autoReconnectInterval) {
        options.autoReconnectInterval = 15 * 1000;
    }
    this.autoReconnectInterval = options.autoReconnectInterval;  // ms

    this.options = options;
};
Tomato.WebSocketClient.prototype.open = function (url) {
    var that = this;
    this.url = url;
    this.instance = new WebSocket(this.url);
    this.instance.onopen = function () {
        that.onopen();
        if (that.reconnectTimeout) {
            that.autoReconnectRetry = 0;
            that.autoReconnectInterval = that.options.autoReconnectInterval;
            clearTimeout(that.reconnectTimeout);
        }
    };
    this.instance.onmessage = function (data, flags) {
        that.number++;
        that.onmessage(data, flags, this.number);
    };
    this.instance.onclose = function (e) {
        that.onclose(e);

        // @see: https://tools.ietf.org/html/rfc6455#section-7.4.1
        switch (e.code) {
            case 1000:  // CLOSE_NORMAL
                console.log("WebSocket: closed");
                break;
            default:    // Abnormal closure
                that.reconnect(e);
                break;
        }
    };
    this.instance.onerror = function (e) {
        switch (e.code) {
            case 'ECONNREFUSED':
                that.reconnect(e);
                break;
            default:
                that.onerror(e);
                break;
        }
    };
};
Tomato.WebSocketClient.prototype.send = function (data, option) {
    try {
        this.instance.send(data, option);
    } catch (e) {
        this.instance.emit('error', e);
    }
};
Tomato.WebSocketClient.prototype.close = function (code, reason) {
    try {
        this.instance.close(code, reason);
    } catch (e) {
        this.instance.emit('error', e);
    }
};
Tomato.WebSocketClient.prototype.reconnect = function (e) {
    if (this.autoReconnectInterval <= 0) {
        return;
    }
    var reconnectInterval = this.autoReconnectInterval * (1 + Math.log(this.autoReconnectRetry + 1));
    console.log(`WebSocketClient: retry in ${reconnectInterval}ms`, e);
    this.instance.removeEventListener("open", this.instance.onopen);
    this.instance.removeEventListener("message", this.instance.onmessage);
    this.instance.removeEventListener("close", this.instance.onclose);
    this.instance.removeEventListener("error", this.instance.onerror);
    var that = this;
    if (this.reconnectTimeout) {
        clearTimeout(this.reconnectTimeout);
    }

    that.onreconnect(reconnectInterval);
    this.reconnectTimeout = setTimeout(function () {
        console.log("WebSocketClient: reconnecting...");
        that.autoReconnectRetry++;
        that.open(that.url);
    }, reconnectInterval);
};
Tomato.WebSocketClient.prototype.onopen = function (e) {
    console.log("WebSocketClient: open", arguments);
};
Tomato.WebSocketClient.prototype.onmessage = function (data, flags, number) {
    console.log("WebSocketClient: message", arguments);
};
Tomato.WebSocketClient.prototype.onerror = function (e) {
    console.log("WebSocketClient: error", arguments);
};
Tomato.WebSocketClient.prototype.onclose = function (e) {
    console.log("WebSocketClient: closed", arguments);
};
Tomato.WebSocketClient.prototype.onreconnect = function (e) {
    console.log("WebSocketClient: reconnect", arguments);
};