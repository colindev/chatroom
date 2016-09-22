(function(factory){
    
    // 只先調整成方便改用requirejs
    self['MessageBox'] = factory(createDom);

})(function(createDom){

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
        dom.appendChild(createDom(me.buildMsg(data.icon, data.name, data.msg)));
        me.GCMsg(dom);
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
                '<p>{msg}</p>'+
            '</div>',

    // 產生訊息模組字串
    buildMsg: function(icon, name, msg) {
        return this.format.
            replace(/\{icon\}/g, icon).
            replace(/\{name\}/g, name).
            replace(/\{msg\}/g, msg);
    }
};

return MessageBox;
});
