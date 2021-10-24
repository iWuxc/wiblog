// JavaScript Document

function iniParam() {
    $('.nav-items a').click(function () {
        layer.msg("Hello World！");
        return false; 
    });

    //页面效果
    setTimeout(function () {
        $('#menu-cb').click();
        setTimeout(function () {  $('#menu-cb').click();}, 1500);
    }, 1000);

}