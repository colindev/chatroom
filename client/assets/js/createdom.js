(function(factory){
    
    self['createDom'] = factory();

})(function(){

function createDom(html) {
    if (!arguments.callee.div) {
        arguments.callee.div = document.createElement('div');
    }
    var div = arguments.callee.div;

    div.innerHTML = html;

    return div.removeChild(div.firstChild);
}


return createDom;
});
