layui.use(['jquery', 'flow'], function () {
    var $ = layui.jquery;
    var flow = layui.flow;
    article.Init($);//初始化共用js
    flow.load({
        elem: "#LAY_bloglist",
        done: function (page, next) {
            var pagecount = $(".bloglist").attr("data-pagecount"),
                pagesize = $(".bloglist").attr("data-pagesize"),
                lis = [];
            $.ajax({
                type: "POST",
                url: "/web/article/list",
                data: {
                    page: page,
                    pagesize: pagesize
                },
                success: function (res) {
                    if (res.code == 0 && res.data) {
                        let articles = res.data
                        console.log(articles)
                        // let html = "";
                        // for(let i = 0; i < articles.length; i++) {
                        //     html += `<section class='article-item zoomIn article'>
                        //                 <div class='fc-flag'>置顶</div>
                        //                 <h5 class='title'>
                        //                     <span class='fc-blue'>【分类】</span>
                        //                     <a href='{11}'>${articles[i].Title}</a>
                        //                 </h5>
                        //                 <div class='time'>
                        //                     <span class='day'>28</span>
                        //                     <span class='month fs-18'>10<span class='fs-14'>月</span></span>
                        //                     <span class='year fs-18 ml10'>2021</span>
                        //                 </div>
                        //                 <div class='content'>
                        //                     <a href='${articles[i].Title}' class='cover img-light'>
                        //                         <img src='{{.Domain}}' >
                        //                     </a>
                        //                     {7}
                        //                 </div>
                        //                 <div class='read-more'>
                        //                     <a href='{13}' class='fc-black f-fwb'>继续阅读</a>
                        //                 </div>
                        //                 <aside class='f-oh footer'>
                        //                     <div class='f-fl tags'>
                        //                         <span class='fa fa-tags fs-16'></span>
                        //                         <a class='tag'>{8}</a>
                        //                     </div>
                        //                     <div class='f-fr'>
                        //                         <span class='read'>
                        //                             <i class='fa fa-eye fs-16'></i>
                        //                             <i class='num'>{9}</i>
                        //                         </span>
                        //                         <span class='ml20'>
                        //                             <i class='fa fa-comments fs-16'></i>
                        //                             <a href = 'javascript:void(0)' class='num fc-grey'>{10}</a>
                        //                         </span>
                        //                     </div>
                        //                 </aside>
                        //             </section>`;
                        // }
                        // lis.push(html);
                        next(lis.join(""), page < pagecount);
                    }

                }
            })
        }
    })
});
var article = {};
article.Init = function ($) {
    //var $ = layui.jquery,
    var slider = 0;
    blogtype();
    //类别导航开关点击事件
    $('.category-toggle').click(function (e) {
        e.stopPropagation();    //阻止事件冒泡
        categroyIn();
    });
    //类别导航点击事件，用来关闭类别导航
    $('.article-category').click(function () {
        categoryOut();
    });
    //遮罩点击事件
    $('.blog-mask').click(function () {
        categoryOut();
    });
    $('.f-qq').on('click', function () {
        window.open('http://connect.qq.com/widget/shareqq/index.html?url=' + $(this).attr("href") + '&sharesource=qzone&title=' + $(this).attr("title") + '&pics=' + $(this).attr("cover") + '&summary=' + $(this).attr("desc") + '&desc=你的分享简述' + $(this).attr("desc"));
    });
    $("body").delegate(".fa-times", "click", function () {
        $(".search-result").hide().empty(); $("#searchtxt").val("");
        $(".search-icon i").removeClass("fa-times").addClass("fa-search");
    });
    //显示类别导航
    function categroyIn() {
        $('.category-toggle').addClass('layui-hide');
        $('.article-category').unbind('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend');
        $('.blog-mask').unbind('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend');
        $('.blog-mask').removeClass('maskOut').addClass('maskIn');
        $('.blog-mask').removeClass('layui-hide').addClass('layui-show');
        $('.article-category').removeClass('categoryOut').addClass('categoryIn');
        $('.article-category').addClass('layui-show');
    }
    //隐藏类别导航
    function categoryOut() {
        $('.blog-mask').on('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend', function () {
            $('.blog-mask').addClass('layui-hide');
        });
        $('.article-category').on('webkitAnimationEnd mozAnimationEnd MSAnimationEnd oanimationend animationend', function () {
            $('.article-category').removeClass('layui-show');
            $('.category-toggle').removeClass('layui-hide');
        });
        $('.blog-mask').removeClass('maskIn').addClass('maskOut').removeClass('layui-show');
        $('.article-category').removeClass('categoryIn').addClass('categoryOut');
    }
    function blogtype() {
        var i = $("#blogtypeid").val();
        i != 0 && (t = parseInt(i) * 40, $(".slider").css({
            top: t + "px"
        }));
        $('#category li').hover(function () {
            $(this).addClass('current');
            var num = $(this).attr('data-index');
            $('.slider').css({ 'top': ((parseInt(num) - 1) * 40) + 'px' });
        }, function () {
            $(this).removeClass('current');
            $('.slider').css({ 'top': slider });
        });
        $(window).scroll(function (event) {
            var winPos = $(window).scrollTop();
            if (winPos > 750)
                $('#categoryandsearch').addClass('fixed');
            else
                $('#categoryandsearch').removeClass('fixed');
        });
    };
    $("#searchtxt").on("keyup", function () {
        setTimeout(function () {
            "" == $("#searchtxt").val().trim() ? $(".search-result").empty().hide() : $.ajax({
                type: "post",
                url: "/Article/SearchResult",
                data: {
                    context: $("#searchtxt").val().trim()
                },
                dataType: "json",
                success: function (a) {
                    "[]" != a ? ($(".search-result").show().empty(), $.each(a,
                        function (t, i) {
                            $(".search-result").append('<li class="child"><a href="/Article/Detail/' + i.Id + '" style="display:block" target="_blank">' + i.Title.toLowerCase().replace($("#searchtxt").val().trim().toLowerCase(), '<b style="color:#6bc30d">' + $("#searchtxt").val().trim() + "<\/b>") + "<\/a><\/li>")
                        })) : $(".search-result").hide().empty()
                },
                complete: function () {
                    "" != $("#searchtxt").val().trim() && $(".search-icon i").removeClass("fa-search").addClass("fa-times")
                }
            })
        },
            500)
    });
    $("body").delegate(".fa-times", "click", function () {
        $(".search-result").hide().empty();
        $("#searchtxt").val("");
        $(".search-icon i").removeClass("fa-times").addClass("fa-search")
    });
};

