package admin

func articleHtml() string {
	str := `<section class='article-item zoomIn article'>
				<div class='fc-flag'>%s</div>
				<h5 class='title'>
					<span class='fc-blue'>【原创】</span>
					<a href='%s'>%s</a>
				</h5>
				<div class='time'>
					<span class='day'>%s</span>
					<span class='month fs-18'>%s<span class='fs-14'>月</span></span>
					<span class='year fs-18 ml10'>%s</span>
				</div>
				<div class='content'>
					<a href='%s' class='cover img-light'>
						<img src=%s />
					</a>
					%s
				</div>
				<div class='read-more'>
					<a href='%s' class='fc-black f-fwb'>继续阅读</a>
				</div>
				<aside class='f-oh footer'>
					<div class='f-fl tags'>
						<span class='fa fa-tags fs-16'></span>
						<a class='tag'>%s</a>
					</div>
					<div class='f-fr'>
						<span class='read'>
							<i class='fa fa-eye fs-16'></i>
							<i class='num'>20123</i>
						</span>
						<span class='ml20'>
							<i class='fa fa-comments fs-16'></i>
							<a href = 'javascript:void(0)' class='num fc-grey'>10</a>
						</span>
					</div>
				</aside>
			</section>`
	return str
}
