// Tomato object setup
const Tomato = {};

// @see: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
Tomato.WebSocketClient = function (options) {
    options = options || {};
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
Tomato.WebSocketClient.prototype.open = function (url, protocol) {
    let that = this;
    this.url = url;
    this.instance = new WebSocket(this.url, protocol);
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
Tomato.WebSocketClient.prototype.send = function (data) {
    try {
        this.instance.send(data);
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
    if (this.autoReconnectInterval <= 0 || this.autoReconnectRetry >= this.autoReconnectMaxRetry) {
        return;
    }
    let reconnectInterval = this.autoReconnectInterval * (1 + Math.log(this.autoReconnectRetry + 1));
    console.log("WebSocketClient: retry in " + reconnectInterval + "ms", e);
    this.instance.removeEventListener("open", this.instance.onopen);
    this.instance.removeEventListener("message", this.instance.onmessage);
    this.instance.removeEventListener("close", this.instance.onclose);
    this.instance.removeEventListener("error", this.instance.onerror);
    let that = this;
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

Tomato.cookie = {
    set: function (name, value, lifetime) {
        let expires = "";
        if (lifetime) {
            let date = new Date();
            date.setTime(date.getTime() + (lifetime * 1000));
            expires = "; expires=" + date.toUTCString();
        }
        document.cookie = name + "=" + (value || "") + expires + "; path=/";
    },
    get: function (name) {
        let nameEQ = name + "=";
        let ca = document.cookie.split(';');
        for (let i = 0; i < ca.length; i++) {
            let c = ca[i];
            while (c.charAt(0) === ' ') {
                c = c.substring(1, c.length);
            }
            if (c.indexOf(nameEQ) === 0) {
                return c.substring(nameEQ.length, c.length);
            }
        }
        return null;
    },
    delete: function (name) {
        document.cookie = name + '=; Max-Age=-99999999;';
    }
};

Tomato.GetUrlParameter = function (sParam) {
    let sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
};

/**
 * @return {string}
 */
Tomato.EscapeHtml = function (html) {
    return html
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
};

Tomato.CopySelect = function () {
    if (!window.document.queryCommandSupported || !window.document.queryCommandSupported("copy")) {
        return;
    }
    window.document.execCommand("copy");
};

module.exports = Tomato;
