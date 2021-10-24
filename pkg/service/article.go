package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"wiblog/pkg/cache"
	"wiblog/pkg/cache/store"
	"wiblog/pkg/model"
)

func List(c *gin.Context, page, limit int) ([]*model.Article, int, error) {
	search := store.SearchArticles{
		Page: page,
		Limit: limit,
		Fields: map[string]interface{}{
			store.SearchArticleDraft: false,
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
		Lock: new(sync.Mutex),
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
				ID: a.ID,
				Author: a.Author,
				Slug: a.Slug,
				Title: a.Title,
				Count: a.Count,
				Content: a.Content,
				Tags: a.Tags,
				IsDraft: a.IsDraft,
				IsHot: a.IsHot,
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

	for _, id := range ids{
		infos = append(infos, articleList.IdMap[id])
	}

	fmt.Println(articleList.IdMap)
	return infos, count, err
}
