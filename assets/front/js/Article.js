// JavaScript Document

function iniParam() {
    var laypage = layui.laypage;

    //页面效果
    $("#keyWord").focus(function () {
        $(this).parent().addClass("search-border");
    }).blur(function () {
        $(this).parent().removeClass("search-border");
    }).keydown(function (e) { 
        if (e.which == 13) { //监听回车事件
            search($(this).val());
            return false;
        }
    });

    //搜索
    $('#search').click(function () {
        var value = $("#keyWord").val();
        search(value);
    });

    function search(value) {
        if (!value) {
            layer.tips('关键字都没输入想搜啥呢...', '#keyWord', { tips: [1, '#659FFD'] });
            $("#keyWord")[0].focus(); //使文本框获得焦点
            return;
        }
      	layer.msg("没想到你居然搜这种东西："+value+"！！！");
    }
	
	laypage.render({
		elem: 'page',
		count: 50, //数据总数通过服务端得到
		limit: 5, //每页显示的条数。laypage将会借助 count 和 limit 计算出分页数。
		curr: 1, //起始页。一般用于刷新类型的跳页以及HASH跳页。
		first: '首页',
		last: '尾页',
		layout: ['prev', 'page', 'next', 'skip'],
		//theme: "page",
		jump: function (obj, first) {
			if (!first) { //首次不执行
				layer.msg("第"+obj.curr+"页");

			}
		}
	});
  

}
