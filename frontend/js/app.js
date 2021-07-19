import jQuery from 'jquery';
import 'bootstrap/dist/css/bootstrap.css';
import '../css/app.css';
import AnsiUp from '../js/ansi_up.min.js';
import '../js/jquery.inputHistory.js';
import Tomato from "expose-loader?exposes=Tomato!../js/tomato.js";

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

let render = function(protocol, data){
    if (data === "") {
        return;
    }

    if (protocol === "ascii") {
        let msg = ansi_up.ansi_to_html(data);
        showMessage(msg);
        return;
    }
    // @TODO: mxp & gmcp
    try {
        let resp = JSON.parse(data);
        switch (resp.event) {
            case "text":
                let msg = ansi_up.ansi_to_html(resp.content);
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
                break;
            default:
                break;
        }
    } catch (e) {
        console.log(e);
    }
}

document.addEventListener("DOMContentLoaded", function () {
    // Check for websocket availability.
    if (!("WebSocket" in window)) {
        showMessage("Requires a browser with WebSockets support.");
    } else {
        let wsc = new Tomato.WebSocketClient();
        let url = process.env.WEBSOCKET_URL;
        let protocol = process.env.WEBSOCKET_PROTOCOL;
        let connectionIcon = $(".prompt .connection i");
        if (protocol === "cloudrain") {
            wsc.open(url);
        } else {
            wsc.open(url, protocol);
        }

        wsc.onopen = function () {
            showMessage('<p class="text-success">Connected.</p>');
            connectionIcon.removeClass("red").addClass("green").attr("title", "Connected");
        };
        wsc.onmessage = function (msg) {
            if (msg.data instanceof Blob) {
                let reader = new FileReader();

                reader.addEventListener("loadend", function() {
                    render(protocol, reader.result);
                });

                reader.readAsText(msg.data);
            } else {
                render(protocol, msg.data);
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
                if (protocol === "cloudrain") {
                    wsc.open(url);
                } else {
                    wsc.open(url, protocol);
                }
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
                    let value = this.value.trim();
                    let cmd = {
                        "type": "cmd",
                        "content": value
                    };
                    let k = value.substr(0, value.indexOf(" "));
                    if (k === "#gmcp") {
                        cmd.type = "gmcp";
                        cmd.content = value.substr(value.indexOf(" ") + 1).trim();
                    }

                    try {
                        if (cmd.type === "cmd") {
                            showMessage(Tomato.EscapeHtml(cmd.content) + "\n");
                        }
                        if (protocol === "ascii") {
                            wsc.send(cmd.content + "\n");
                        } else {
                            wsc.send(JSON.stringify(cmd));
                        }
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
