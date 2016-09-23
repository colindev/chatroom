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
            return d.toLocaleTimeString();
        },
        msg: function(def){
            return this.data && this.data.msg || def || '';
        },
        isSelf: function(data){
            return this.icon() === data.profile.icon;
        }
    };

    return Payload;
});
