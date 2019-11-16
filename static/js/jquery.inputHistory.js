// Modified from https://github.com/jch/jquery.inputHistory

(function () {
    (function ($) {
        var InputHistory, normalizeKeyHandler;
        InputHistory = (function () {
            function InputHistory(options) {
                this.size = options.size || 50;
                this.useLatest = options.useLatest || false;
                this.ignoreEmpty = options.ignoreEmpty || false;
                this.values = [];
                this.index = 0;
                this.noMove = false;
            }

            InputHistory.prototype.push = function (message) {
                this.moving = false;
                if (!(this.ignoreEmpty && message.length === 0)) {
                    this.values.unshift(message);
                    return this.values.splice(this.size);
                }
            };

            InputHistory.prototype.prev = function () {
                var message;
                var size = this.size;
                if (size > this.values.length) {
                    size = this.values.length;
                }
                if ((this.useLatest && !this.moving) || this.index > size - 1) {
                    this.noMove = true;
                    return this.values[this.index - 1];
                }
                message = this.values[this.index];
                this.index += 1;
                this.moving = true;
                this.noMove = false;
                return message;
            };

            InputHistory.prototype.next = function () {
                this.index -= 1;
                this.moving = true;
                this.noMove = false;
                if (this.index < 0) {
                    this.index = 0;
                    return "";
                } else {
                    return this.values[this.index];
                }
            };

            return InputHistory;

        })();
        normalizeKeyHandler = function (raw, elseHandler) {
            elseHandler || (elseHandler = function (e) {
            });
            switch (typeof raw) {
                case 'number':
                    return function (e) {
                        return e.keyCode === raw;
                    };
                case 'string':
                    return function (e) {
                        return "" + e.keyCode === raw;
                    };
                case 'function':
                    return raw;
                default:
                    return elseHandler;
            }
        };
        return $.fn.inputHistory = function (options) {
            var history,
                _this = this;
            options || (options = {});
            options.data || (options.data = 'inputHistory');
            options.store = normalizeKeyHandler(options.store, function (e) {
                return e.keyCode === 13 && !e.shiftKey;
            });
            options.prev = normalizeKeyHandler(options.prev, function (e) {
                return e.keyCode === 38 || (e.ctrlKey && e.keyCode === 80);
            });
            options.next = normalizeKeyHandler(options.next, function (e) {
                return e.keyCode === 40 || (e.ctrlKey && e.keyCode === 78);
            });
            history = this.data(options.data) || new InputHistory(options);
            this.data(options.data, history);
            this.bind('keydown', function (e) {
                if (options.store(e)) {
                    history.push(_this.val());
                }
                if (options.prev(e)) {
                    _this.val(history.prev());
                    if (!history.noMove) {
                        _this.select();
                    } else {
                        var el = _this.get()[0];
                        el.selectionStart = el.selectionEnd = el.value.length;
                    }
                    e.preventDefault();
                }
                if (options.next(e)) {
                    return _this.val(history.next()) && _this.select() && e.preventDefault();
                }
            });
            return this;
        };
    })(jQuery);

})();
