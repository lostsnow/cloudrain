import jQuery from 'jquery';
import 'bootstrap/dist/css/bootstrap.css';
import '../css/app.css';
import AnsiUp from '../js/ansi_up.min.js';
import '../js/jquery.inputHistory.js';
require('expose-loader?Tomato!../js/tomato.js');

document.addEventListener("DOMContentLoaded", function () {
    let ansi_up = new AnsiUp;
    let terminalBox = $("#terminal-box");
    let promptInput = $('#prompt-input');
    let showMessage = function (msg) {
        let atBottom = (terminalBox.scrollTop() + 40 >= terminalBox[0].scrollHeight - terminalBox.height());
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
        let wsc = new Tomato.WebSocketClient();
        let url = process.env.WEBSOCKET_URL;
        let connectionIcon = $(".prompt .connection i");
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
                let resp = JSON.parse(msg.data);
                switch (resp.event) {
                    case "text":
                        msg = ansi_up.ansi_to_html(resp.content);
                        showMessage(msg);
                        break;
                    case "session":
                        try {
                            /** @var {{sid:string, token:string}} session */
                            let session = JSON.parse(resp.content);
                            if (!session.sid || !session.token) {
                                showMessage("Invalid websocket session");
                                wsc.autoReconnectInterval = 0;
                                wsc.close(1000);
                                return;
                            }
                            Tomato.cookie.set("sessionid", session.sid);
                            Tomato.cookie.set("token", session.token);
                            let query_sid = Tomato.GetUrlParameter('sid');
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
                    case "gmcp":
                        console.log("gmcp:", resp.content);
                        try {
                            /** @var {{event:{id:string}}} gmcp */
                            let gmcp = JSON.parse(resp.content);
                            console.log(gmcp.event.id);
                            if (gmcp.event && (gmcp.event.id === "login" || gmcp.event.id === "reconnect")) {
                                let cmd = {
                                    "type": "gmcp",
                                    "content": "request room.info"
                                };
                                wsc.send(JSON.stringify(cmd));
                                return;
                            }
                        } catch (e) {}
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
                Tomato.cookie.delete("sessionid");
            }
            promptInput.trigger("focus");
        });

        terminalBox.on("mouseup", function () {
            let selection;
            if (window.getSelection) {
                selection = window.getSelection();
            } else if (window.document.selection) {
                selection = window.document.selection.createRange();
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
                    let type = "cmd";
                    let content = this.value;
                    let k = this.value.substr(0, this.value.indexOf(" "));
                    if (k === "#gmcp") {
                        type = "gmcp";
                        content = this.value.substr(this.value.indexOf(" ") + 1).trim();
                    }
                    let cmd = {
                        "type": type,
                        "content": content
                    };
                    try {
                        showMessage(Tomato.EscapeHtml(cmd.content) + "\n");
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
