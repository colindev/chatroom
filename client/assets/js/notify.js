(function(factory){

    self['notify'] = factory(Payload);

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

return function notify(title, data){
    if (!available) return;
    if (Notification.permission !== "granted") Notification.requestPermission();

    var pld = new Payload(data),
        n = new Notification(title, {
            icon: "http://s.gravatar.com/avatar/"+pld.icon()+"?s=80",
            body: pld.msg().replace(/</g, '&lt;').replace(/>/g, '&gt;')
        });
};

});
