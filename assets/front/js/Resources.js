// JavaScript Document

function iniParam() {
    var element = layui.element;
    element.render('tab');

    //Hash地址的定位 页面刷新后切换回原tab界面
    var layid = location.hash.replace(/^#tabIndex=/, '');
    if (!layid)
        element.tabChange('tabResource', $('#tabResource > ul > li').eq(0).attr('lay-id'));
    else
        element.tabChange('tabResource', layid);

    element.on('tab(tabResource)', function (elem) {
        location.hash = 'tabIndex=' + $(this).attr('lay-id');
    });


}