package admin

func articleHtml() string {
	str := `<div class="single-post">
				<div class="inner-post">
					<div class="post-img">
					<a href="%s">
					<img src="%s" alt=""/></a></div>
					<div class="post-info">
						<div class="post-title">
							<h3>
								<a href="%s">
									%s
								</a>
							</h3>
						</div>
						<div class="post-content"><p>%s</p>
						</div>
						<div class="blog-meta">
							<ul>
								<li>
									<i class="layui-icon layui-icon-note"></i>
										<a href="%s">算法</a>
								</li>
								<li><i class="layui-icon layui-icon-log"></i>%s</li>
								<li><i class="layui-icon layui-extend-wiappfangwenliang"></i>浏览（<a
										href="%s">64</a>）
								</li>
							</ul>
							<div class="post-readmore">
								<a href="%s">阅读更多</a>
							</div>
						</div>
					</div>
				</div>
				<div class="post-date one"><span>%d</span></div>
			</div>`
	return str
}
