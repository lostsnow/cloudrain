jQuery(function ($) {
    var ansi_up = new AnsiUp;
    var terminalBox = $("#terminal-box");
    var promptInput = $('#prompt-input');

    var showMessage = function (msg) {
        var atBottom = (terminalBox.scrollTop() + 40 >= terminalBox[0].scrollHeight - terminalBox.height());
        terminalBox.append(msg);
        // If we were scrolled to the bottom before this call, remain there.
        if (atBottom) {
            terminalBox.scrollTop(terminalBox[0].scrollHeight - terminalBox.height())
        }
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
            if (msg.data === "") {
                return;
            }
            // @TODO: mxp & gmcp
            try {
                var resp = JSON.parse(msg.data);
                switch (resp.event) {
                    case "text":
                        msg = ansi_up.ansi_to_html(resp.content);
                        showMessage(msg);
                        break;
                    case "session":
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
                            var query_sid = Tomato.GetUrlParameter('sid');
                            if (query_sid && session.sid !== query_sid) {
                                window.history.pushState('', '', location.href.replace('sid=' + query_sid, 'sid=' + session.sid));
                            }
                        } catch (e) {
                            showMessage("Invalid websocket response");
                            wsc.autoReconnectInterval = 0;
                            wsc.close(1000);
                        }
                        break;
                    case "ping":
                        console.log("ping...");
                        break;
                    case "mssp":
                        console.log("mssp:", resp.content);
                        break;
                    case "atcp":
                        console.log("atcp");
                        break;
                    case "mxp":
                        console.log("mxp");
                        break;
                    default:
                        break;
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

        terminalBox.on("mouseup", function () {
            var selection;
            if (window.getSelection) {
                selection = window.getSelection();
            } else if (document.selection) {
                selection = document.selection.createRange();
            }

            if (selection.toString() !== "") {
                Tomato.CopySelect();
            } else {
                promptInput.trigger("focus");
            }
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
                        showMessage(Tomato.EscapeHtml(cmd.content));
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
