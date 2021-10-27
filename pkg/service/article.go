package service

import (
	"github.com/gin-gonic/gin"
	"sync"
	"wiblog/pkg/cache"
	"wiblog/pkg/cache/store"
	"wiblog/pkg/conf"
	"wiblog/pkg/model"
)

func ArticleList(c *gin.Context, page, limit, serieid int, keyword string) ([]*model.Article, int, error) {
	search := store.SearchArticles{
		Page:  page,
		Limit: limit,
		Fields: map[string]interface{}{
			store.SearchArticleDraft: false,
			store.SearchArticleTitle: keyword,
			store.SearchArticleSerieID: serieid,
		},
	}
	infos := make([]*model.Article, 0)
	articles, count, err := cache.Wi.Store.LoadArticleList(c, search)
	if err != nil {
		return nil, count, err
	}
	ids := make([]int, 0)
	for _, article := range articles {
		ids = append(ids, article.ID)
	}

	wg := sync.WaitGroup{}

	articleList := model.ArticleList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[int]*model.Article, len(articles)),
	}

	finished := make(chan bool, 1)
	//errChan := make(chan error, 1)

	for _, a := range articles {
		wg.Add(1)
		go func(a *model.Article) {
			defer wg.Done()
			articleList.Lock.Lock()
			defer articleList.Lock.Unlock()
			articleList.IdMap[a.ID] = &model.Article{
				ID:            a.ID,
				Author:        a.Author,
				Slug:          a.Slug,
				Title:         a.Title,
				Count:         a.Count,
				Content:       a.Content,
				Tags:          a.Tags,
				IsDraft:       a.IsDraft,
				IsHot:         a.IsHot,
				Cover:         a.Cover,
				CreatedFormat: dateFormat(a.CreatedAt, "2006-01-02 03:04"), //01/02 03:04:05PM 06 -0700
				ArticleUrl:    "http://" + conf.Conf.WiBlogApp.Host + "/post/" + a.Slug + ".html",
			}
		}(a)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	}

	for _, id := range ids {
		infos = append(infos, articleList.IdMap[id])
	}
	return infos, count, err
}
