(function(factory){

    self['Payload'] = factory();

})(function(){
    
    function Payload(data){
        this.data = data;
    }

    Payload.prototype = {
        name: function(){ return this.data && this.data.profile && this.data.profile.name;},
        icon: function(){ return this.data && this.data.profile && this.data.profile.icon;},
        time: function(){
            var d = new Date;
            d.setTime(this.data.time * 1000)
            return d.toLocaleString();
        },
        msg: function(){
            return this.data && this.data.msg;
        }
    };

    return Payload;
});
