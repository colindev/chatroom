(function(factory){

    self['UserList'] = factory(createDom, Payload);

})(function(createDom, Payload){

function UserList(dom) {
    
    var me = this, unique = {};
    this.join = function(data){

        var pld = new Payload(data);

        if (unique[pld.icon()]) {
            unique[pld.icon()][0]++;
            return
        }

        var icon = createDom(me.buildIcon(pld.name(), pld.icon(), pld.time()));
        unique[pld.icon()] = [1, icon];
        dom.appendChild(icon);
    };

    this.leave = function(data){

        var pld = new Payload(data);

        if (unique[pld.icon()]) {
            unique[pld.icon()][0]--;
            if (unique[pld.icon()][0] <= 0) {
                dom.removeChild(unique[pld.icon()][1]);
                delete(unique[pld.icon()]);
            }
        }
    };
}

UserList.prototype = {
    format: '<span class="icon" title="{name} / {time}" style="background-image: url(http://s.gravatar.com/avatar/{icon}?s=30);"></span>',
    buildIcon: function(name, icon, time){
        return this.format.replace(/\{name\}/g, name).
                replace(/\{time\}/g, time).
                replace(/\{icon\}/g, icon);
    },
};

return UserList;
});
