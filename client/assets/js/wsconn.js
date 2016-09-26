(function(factory){

    self['WSConn'] = factory();

})(function(){

function WSConn() {

    var me = this, conn;

    this.connect = function(url){
        if (conn) return;

        me.trigger(WSConn.CONNECT, {});
        conn = new WebSocket(url);
        conn.onopen = function(e) {
            me.trigger(WSConn.OPEN, e);
        };
        conn.onclose = function(e){
            conn = null;
            me.trigger(WSConn.CLOSE, e);
        };
        conn.onmessage = function(e){
            me.trigger(WSConn.MESSAGE, e);
        };

        return me;
    };
    this.disconnect = function(){

        if (conn) {
            me.trigger(WSConn.CLOSE, {});
            conn.onclose = null;
            conn.onopen = null;
            conn.onmessage = null;
            conn.close();
            conn = null;
        }
        return me;
    };
    this.send = function(msg){

        var event = {
            data: msg
        };
        me.trigger(me.SEND, event);
        
        if ( ! conn) {
            throw new Error('please connect before send something');
        }
        conn.send(event.data);

        return me;
    }
}

WSConn.CONNECT = 'connect';
WSConn.OPEN = 'open';
WSConn.CLOSE = 'close';
WSConn.MESSAGE = 'message';
WSConn.SEND = 'send';

WSConn.prototype = {
    _handlers: {},
    on: function(type, handler){
        if ( ! this._handlers[type]) this._handlers[type] = [];
        this._handlers[type].push(handler);

        return this;
    },
    off: function(type, handler) {
        if (this._handlers[type]) {
            var i = this._handlers[type].length;
            while (i--) {
                if (this._handlers[type][i] == handler) {
                    this._handlers[type].splice(i, 1);
                }
            } 
        }

        return this;
    },
    trigger: function(type, event) {
        if (this._handlers[type]) {
            for (var i = 0; i < this._handlers[type].length; i++) {
                this._handlers[type][i].call(this, event);
            }
        }

        return this;
    }
};

return WSConn;
});
