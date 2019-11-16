jQuery(function ($) {
    var ansi_up = new AnsiUp;
    var terminalBox = $("#terminal-box");
    var promptInput = $('#prompt-input');

    var getUrlParameter = function getUrlParameter(sParam) {
        var sPageURL = window.location.search.substring(1),
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

    var showMessage = function (msg) {
        var atBottom = (terminalBox.scrollTop() + 40 >= terminalBox[0].scrollHeight - terminalBox.height());
        terminalBox.append(msg);
        // If we were scrolled to the bottom before this call, remain there.
        if (atBottom) {
            terminalBox.scrollTop(terminalBox[0].scrollHeight - terminalBox.height())
        }
    };

    var escapeHtml = function (text) {
        return text
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    };

    // Check for websocket availability.
    if (!("WebSocket" in window)) {
        showMessage("Requires a browser with WebSockets support.");
    } else {
        var wsc = new Tomato.WebSocketClient();
        var url = Tomato.websocketUrl;
        var connectionIcon = $(".prompt .connection i");
        wsc.open(url);
        wsc.onopen = function () {
            showMessage('<p class="text-success">Connected.</p>');
            connectionIcon.removeClass("red").addClass("green").attr("title", "Connected");
        };
        wsc.onmessage = function (msg) {
            // console.log(msg);
            if (msg.data === "") {
                return;
            }
            // @TODO: mxp & gmcp
            try {
                var resp = JSON.parse(msg.data);
                if (resp.event === "text") {
                    msg = ansi_up.ansi_to_html(resp.content);
                    showMessage(msg);
                } else if (resp.event === "session") {
                    try {
                        /** @var {{sid:string, token:string}} session */
                        var session = JSON.parse(resp.content);
                        if (!session.sid || !session.token) {
                            showMessage("Invalid websocket session");
                            wsc.autoReconnectInterval = 0;
                            wsc.close(1000);
                            return;
                        }
                        Tomato.cookie.set("sessionid", session.sid);
                        Tomato.cookie.set("token", session.token);
                        var query_sid = getUrlParameter('sid');
                        if (query_sid && session.sid !== query_sid) {
                            window.history.pushState('', '', location.href.replace('sid=' + query_sid, 'sid=' + session.sid));
                        }
                    } catch (e) {
                        showMessage("Invalid websocket response");
                        wsc.autoReconnectInterval = 0;
                        wsc.close(1000);
                    }
                } else if (resp.event === "ping") {
                    console.log("ping...");
                } else if (resp.event === "mssp") {
                    console.log("mssp:", resp.content);
                } else if (resp.event === "atcp") {
                    console.log("atcp");
                } else if (resp.event === "mxp") {
                    console.log("mxp");
                }
            } catch (e) {
                console.log(e);
            }
        };
        wsc.onclose = function () {
            showMessage('<p class="text-danger">Disconnected.</p>');
            connectionIcon.removeClass("green").addClass("red").attr("title", "Disconnected");
        };
        wsc.onerror = function () {
            showMessage('<p class="text-danger">Error connecting.</p>');
        };
        wsc.onreconnect = function (interval) {
            interval = Math.floor(interval / 1000);
            showMessage('<p class="text-warning">Reconnecting in ' + interval + ' seconds.</p>');
        };

        connectionIcon.on("click", function () {
            if ($(this).hasClass("red")) {
                wsc.open(url);
            } else if ($(this).hasClass("green")) {
                wsc.close(1000);
            }
        });

        terminalBox.on("click", function () {
            promptInput.trigger("focus");
        });

        promptInput.inputHistory({
            size: 50,
            ignoreEmpty: true
        });

        promptInput.on("keyup", function (e) {
            if (e.key === "Enter") {
                if (!e.shiftKey) {  
                    var cmd = {
                        "type": "cmd",
                        "content": this.value + "\n"
                    };
                    try {
                        showMessage(escapeHtml(cmd.content));
                        wsc.send(JSON.stringify(cmd));
                    } catch (exception) {
                        showMessage("<p class=\"text-warning\">Couldn't send message.</p>");
                    }
                    this.value = "";
                    e.stopPropagation();
                }
            }
        });
    }
});
