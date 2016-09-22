(function(factory){
    
    // 只先調整成方便改用requirejs
    self['MessageBox'] = factory(createDom, Payload);

})(function(createDom, Payload){

// 訊息顯示區
function MessageBox(dom) {
    // 避免重複對同一個dom建構
    if (dom._message_box_lock) {
        throw new Error("duplicate handle dom");
    }
    dom._message_box_lock = true;

    // private function
    var me = this;
    this.push = function(data){
        var pld = new Payload(data);
        dom.appendChild(createDom(me.buildMsg(pld.icon(), pld.name(), pld.msg(), pld.time())));
        me.GCMsg(dom);
    };
    this.joined = function(data){
        var pld = new Payload(data);
        dom.appendChild(createDom('<div class="system-msg joined">'+pld.time()+' '+pld.name()+' joined</div>'));
    };
    this.leaved = function(data){
        var pld = new Payload(data);
        dom.appendChild(createDom('<div class="system-msg leaved">'+pld.time()+' '+pld.name()+' leaved</div>'))
    };
    
}
MessageBox.prototype = {
    // 訊息保留限制數量
    msgLimit: 100,
    GCMsg: function(dom) {
        while (dom.children.length > this.msgLimit) {
            dom.removeChild(dom.firstChild);
        }
    },
    // placeholder usage: icon / name / msg
    format: '<div class="msg" title="{name}">'+
                '<span class="name" style="background-image: url(http://s.gravatar.com/avatar/{icon}?s=30);"></span>'+
                '<span class="time">{time}</span>'+
                '<p>{msg}</p>'+
            '</div>',

    // 產生訊息模組字串
    buildMsg: function(icon, name, msg, time) {
        return this.format.
            replace(/\{icon\}/g, icon).
            replace(/\{name\}/g, name).
            replace(/\{time\}/g, time).
            replace(/\{msg\}/g, msg);
    }
};

return MessageBox;
});
