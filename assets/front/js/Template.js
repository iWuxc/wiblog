// JavaScript Document

layui.use(['element', 'laypage', 'form', 'util', 'layer', 'flow','table','layedit'], function () {
    try {
        var util = layui.util, layer = layui.layer;
	console.log("Q-Blog官网：www.qbl.link");	
        $(document).ready(function () { //DOM树加载完毕执行，不必等到页面中图片或其他外部文件都加载完毕
            //页面加载完成后，速度太快会导到loading层闪烁，影响体验，所以在此加上500毫秒延迟
            setTimeout(function () { $("#loading").hide(); }, 500);
        });

        //初始化WOW.js
        new WOW().init();

        //导航点击效果
        $('header nav > ul > li a').click(function () {
            $('header nav > ul > li').removeClass("nav-select-this").find("a").removeClass("nav-a-click");
            $(this).addClass("nav-a-click").parent().addClass("nav-select-this");
        });

        //固定块
        util.fixbar({
            css: { right: 10, bottom: 40, },
            bar1: "&#xe602e;",
            click: function (type) {
                if (type === 'bar1') {
                    layer.tab({
                        area: ['300', '290px'],
                        resize: false, 
                        shadeClose: true,                  
                        scrollbar: false,
                        anim: 4,                 
                        tab: [{
                            title: '微信',
                            content: '<img src="images/zsm.jpg" style="width:255px;" oncontextmenu="return false;" ondragstart="return false;" />',
                        },
                        {
                            title: '支付宝',
                            content: '<img src="images/zfb.jpg" style="width:255px;" oncontextmenu="return false;" ondragstart="return false;" />',
                        }],
                        success: function (layero, index) {
                            $("#" + layero[0].id + " .layui-layer-content").css("overflow", "hidden");
                            $("#" + layero[0].id + " .layui-layer-title span").css("padding", "0px");
                            layer.tips('本站收获的所有打赏费用均只用于服务器日常维护以及本站开发用途，感谢您的支持！', "#" + layero[0].id, {
                                tips: [1, '#FFB800'],
                                time: 0, 
                            });
                        },
                        end: function () {
                            layer.closeAll('tips'); 
                        }
                    });
                }
            }
        });

        //使刷新页面后，此页面导航仍然选中高亮显示，自己根据你实际情况修改
        var pathnameArr = window.location.pathname.split("/");
        var pathname = pathnameArr[pathnameArr.length - 1];
            if (pathname.indexOf(".html") < 0)
                pathname += ".html";
        if (!!pathname) {
            $('header nav > ul > li').removeClass("nav-select-this").find("a").removeClass("nav-a-click");
            $('header nav > ul > li').each(function () {
                var hrefArr = $(this).find("a").attr('href').split("/");
                var href = hrefArr[hrefArr.length - 1];
                if (pathname.toLowerCase() === href.toLowerCase()) {
                    $(this).addClass("nav-select-this").find("a").addClass("nav-a-click");
                    return false;
                }
            });
        }

        iniParam();

        //登录图标
        var anim = setInterval(function () {
            if ($(".animated-circles").hasClass("animated")) {
                $(".animated-circles").removeClass("animated");
            } else {
                $(".animated-circles").addClass('animated');
            }
        }, 3000);
        var wait = setInterval(function () {
            $(".livechat-hint").removeClass("show_hint").addClass("hide_hint");
            clearInterval(wait);
        }, 4500);
        $(".livechat-girl").hover(function () {
            clearInterval(wait);
            $(".livechat-hint").removeClass("hide_hint").addClass("show_hint");
        }, function () {
            $(".livechat-hint").removeClass("show_hint").addClass("hide_hint");
        }).click(function () {
            login();
        });

        //设置样式
        function setStyle(flag) {
            if (!flag) { //未登录
                $('.girl').attr("src", "images/a.png").css("border-radius", "0px");
                $('.livechat-girl').css({ right: "-35px", bottom: "75px" }).removeClass("red-dot");
                $('.rd-notice-content').text('嘿，来试试登录吧！');
                return;
            }
            clearInterval(anim);
			$('.girl').attr("src", "images/nan.png").css("border-radius", "50px");
			$('.rd-notice-content').text('欢迎您，渣渣辉！');
			$('.livechat-girl').css({ right: "0px", bottom: "80px" });	
        }
        //登录事件
        function login() {
			layer.msg("登录成功！");
			setStyle(true);
        }

    }
    catch (e) {
        layui.hint().error(e);
    }         
});

//百度统计
var _hmt = _hmt || [];
(function () {
    var hm = document.createElement("script");
    hm.src = "https://hm.baidu.com/hm.js?132af61e5d1e0bde3638f1ee143bfdb0";
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(hm, s);
})();

