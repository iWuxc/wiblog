// JavaScript Document

function iniParam() {
    var form = layui.form,laypage = layui.laypage,layedit = layui.layedit;
	
 	layer.photos({
		photos: '#details-content',
		anim: 5 //0-6的选择，指定弹出图片动画类型，默认随机（请注意，3.0之前的版本用shift参数）
	});   
	
	
	//评论和留言的编辑器
	for(var i=1;i<9;i++){
		layedit.build('demo-'+i.toString(), {
			height: 150,
			tool: ['face', '|', 'link'],
		});
	}
	

	$(".btn-reply").click(function(){
		 if ($(this).text() == '回复') {
       		$(this).parent().next().removeClass("layui-hide");
        	$('.btn-reply').text('回复');
		    $(this).text('收起');
		}
		else {
			$(this).parent().next().addClass("layui-hide");
			$(this).text('回复');
		}
	});
	
	laypage.render({
		elem: 'page',
		count: 10, //数据总数通过服务端得到
		limit: 5, //每页显示的条数。laypage将会借助 count 和 limit 计算出分页数。
		curr: 1,
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

  
	
	
	//我用的百度编辑器，按照你们自己需求改
	CodeHighlighting(); //代码高亮
    function CodeHighlighting() {
        //添加code标签
        var allPre = document.getElementsByTagName("pre");
        for (i = 0; i < allPre.length; i++) {
            var onePre = document.getElementsByTagName("pre")[i];
            var myCode = document.getElementsByTagName("pre")[i].innerHTML;
            onePre.innerHTML = '<div class="pre-title">Code</div><code class="' + onePre.className.substring((onePre.className.indexOf(":") + 1), onePre.className.indexOf(";")) + '">' + myCode + '</code>';
        }
        //添加行号
        $("code").each(function () {
            $(this).html("<ol><li>" + $(this).html().replace(/\n/g, "\n</li><li>") + "\n</li></ol>");
        });
		hljs.initHighlighting(); //对页面上的所有块应用突出显示
        //hljs.initHighlightingOnLoad(); //页面加载时执行代码高亮
    }
}

