package scheduler

import (
	"log"
	"social-app/services/news-service/service"
	"time"
)

type NewsScheduler struct {
	newsService service.NewsService
	ticker      *time.Ticker
	done        chan bool
}

func NewNewsScheduler(newsService service.NewsService) *NewsScheduler {
	return &NewsScheduler{
		newsService: newsService,
		done:        make(chan bool),
	}
}

func (s *NewsScheduler) Start() {
	log.Println("启动新闻定时任务...")
	
	// 立即执行一次
	go func() {
		log.Println("开始获取新闻...")
		if err := s.newsService.FetchAndSaveNews(); err != nil {
			log.Printf("获取新闻失败: %v", err)
		} else {
			log.Println("新闻获取完成")
		}
	}()

	// 每小时执行一次
	s.ticker = time.NewTicker(1 * time.Hour)
	
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Println("定时获取新闻...")
				if err := s.newsService.FetchAndSaveNews(); err != nil {
					log.Printf("定时获取新闻失败: %v", err)
				} else {
					log.Println("定时新闻获取完成")
				}
			case <-s.done:
				log.Println("停止新闻定时任务")
				return
			}
		}
	}()
}

func (s *NewsScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}