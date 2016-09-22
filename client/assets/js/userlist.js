(function(factory){

    self['UserList'] = factory(createDom);

})(function(createDom){

function UserList(dom) {
    
    var me = this, unique = {};
    this.join = function(data){
        if (unique[data.icon]) {
            unique[data.icon][0]++;
            return
        }

        var icon = createDom(me.buildIcon(data.name, data.icon));
        unique[data.icon] = [1, icon];
        dom.appendChild(icon);
    };

    this.leave = function(data){
        if (unique[data.icon]) {
            unique[data.icon][0]--;
            if (unique[data.icon][0] <= 0) {
                dom.removeChild(unique[data.icon][1]);
                delete(unique[data.icon]);
            }
        }
    };
}

UserList.prototype = {
    format: '<span class="icon" title="{name}" style="background-image: url(http://s.gravatar.com/avatar/{icon}?s=30);"></span>',
    buildIcon: function(name, icon){
        return this.format.replace(/\{name\}/g, name).
                replace(/\{icon\}/g, icon);
    },
};

return UserList;
});
