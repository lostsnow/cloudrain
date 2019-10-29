// Tomato object setup
if (!Tomato) var Tomato = {};

Tomato.websocket = function (options) {
    var ansi_up = new AnsiUp;

    function showMessage(msg) {
        var historyBox = options.historyBox;
        var atBottom = (historyBox.scrollTop() + 40 >= historyBox[0].scrollHeight - historyBox.height());
        historyBox.append(msg);
        // If we were scrolled to the bottom before this call, remain there.
        if (atBottom) {
            historyBox.scrollTop(historyBox[0].scrollHeight - historyBox.height())
        }
    }

    function escapeHtml(text) {
        return text
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    // Register listeners.
    jQuery(function () {
        // Check for websocket availability.
        if (!("WebSocket" in window)) {
            showMessage("Requires a browser with WebSockets support.");
            return;
        }

        try {
            var socket = new WebSocket(options.address);

            socket.onopen = function () {
                showMessage('<p class="text-success">Connected.</p>');
            };

            socket.onmessage = function (msg) {
                msg = ansi_up.ansi_to_html(msg.data);

                showMessage(msg);
            };

            socket.onclose = function () {
                showMessage('<p class="text-danger">Disconnected.</p>');
            }

        } catch (exception) {
            showMessage('<p class="text-danger">Error connecting.</p>');
        }

        options.textBox.keyup(function (e) {
            if (e.keyCode === 13) {
                if (!e.shiftKey) {
                    var content = this.value + "\n";
                    try {
                        showMessage(escapeHtml(content));
                        socket.send(content);
                    } catch (exception) {
                        showMessage('<p class="text-warning">Couldn\'t send message.</p>');
                    }
                    this.value = "";
                    e.stopPropagation();
                }
            }
        });
    });
};
