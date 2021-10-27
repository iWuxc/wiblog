// JavaScript Document
function iniParam() {
    var e = layui.laypage;

    $("#keyword").focus(function () {
        $(this).parent().addClass("search-border")
    }).blur(function () {
        $(this).parent().removeClass("search-border")
    }).keydown(function (i) {
        if (i.which == 13) {
            h($(this).val());
            return false
        }
    });

    var a = $("#ArticleListCount").val();
    var b = $("#CategoriesId").val();
    var c = 1;
    var f = 5;

    d(c, f, null, null);

    // g(c, a, f, null, b);

    function d(j, l, k, i)   {
        $.httpAsyncPost("/web/article/list", {
            page: j,
            limit: l,
            keyword: k,
            serieid: i
        }, function (n) {
            if (n.code === 0) {
                if (!k) {
                    $("#single-list").html(n.data.html)
                } else {
                    var m = '<div class="card erach-tip">';
                    m += "<p>";
                    m += "<span>" + k + "</span> 为您找到 <strong>" + n.data.count + "</strong> 个相关结果";
                    m += "</p>";
                    m += "</div>";
                    $("#single-list").html(m + n.data.html)
                }
                g(j, n.data.count, l, k, i)
            } else {
                $.layerMsg(n.message, n.state)
            }
            setTimeout(function () {
                $.loading(false)
            }, 500)
        })
    }

    function g(j, m, l, k, i) {
        e.render({
            elem: "page",
            count: m,
            limit: l,
            curr: j,
            first: "首页",
            last: "尾页",
            layout: ["prev", "page", "next", "skip"],
            jump: function (o, n) {
                if (!n) {
                    $.loading(true);
                    d(o.curr, o.limit, k, i);
                    $("body,html").animate({
                        scrollTop: $("#am").offset().top
                    }, 500)
                }
            }
        })
    }

    $("#search").click(function () {
        var i = $("#keyword").val();
        h(i)
    });

    $(".lmNav li").click(function () {
        var serieid = $(this).find(".layui-col-xs10").data("id");  //文本内容
        $.loading(true)
        d(c, f, null, serieid)
    })


    function h(i) {
        if (!i) {
            layer.tips("关键字都没输入想搜啥呢...", "#keyword", {
                tips: [1, "#659FFD"]
            });
            $("#keyWord")[0].focus();
            return
        }
        $.loading(true);
        d(c, f, i, null)
    }
}