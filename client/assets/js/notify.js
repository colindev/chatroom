(function(factory){

    self['notifier'] = factory(Payload);

})(function(Payload){

var available = true;
document.addEventListener('DOMContentLoaded', function(e){
    if (!Notification) {
        console.warn('Notification not available');
        available = false;
        return
    }

    if (Notification.permission !== "granted") Notification.requestPermission();
        
});

function Notifier() {

}

Notifier.prototype = {
    isQuite: false,
    quite: function(bool){
        this.isQuite = !!bool;
    },
    notify: function(title, data){
        if (!available) return;
        if (Notification.permission !== "granted") Notification.requestPermission();
        if (this.isQuite) return;

        var pld = new Payload(data),
            n = new Notification(title, {
                icon: "http://s.gravatar.com/avatar/"+pld.icon()+"?s=80",
                body: pld.msg().replace(/</g, '&lt;').replace(/>/g, '&gt;')
            });
    }
};

return new Notifier();

});
