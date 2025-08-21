package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"social-app/services/news-service/models"
	"social-app/services/news-service/repository"
	"strings"
	"time"
)

type NewsService interface {
	GetNews(page, limit int) ([]*models.News, error)
	GetNewsByCategory(category models.NewsCategory, page, limit int) ([]*models.News, error)
	GetNewsDetail(id string) (*models.News, error)
	GetCategories() ([]models.CategoryInfo, error)
	FetchAndSaveNews() error
	RefreshNews() error
}

type newsService struct {
	newsRepo repository.NewsRepository
	apiKey   string
}

func NewNewsService(newsRepo repository.NewsRepository) NewsService {
	apiKey := os.Getenv("NEWS_API_KEY")
	fmt.Printf("🔑 初始化新闻服务，API密钥: %s\n", maskAPIKeyStatic(apiKey))
	
	return &newsService{
		newsRepo: newsRepo,
		apiKey:   apiKey,
	}
}

// 静态掩码函数
func maskAPIKeyStatic(key string) string {
	if key == "" {
		return "未设置"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

func (s *newsService) GetNews(page, limit int) ([]*models.News, error) {
	return s.newsRepo.GetAll(page, limit)
}

func (s *newsService) GetNewsByCategory(category models.NewsCategory, page, limit int) ([]*models.News, error) {
	return s.newsRepo.GetByCategory(category, page, limit)
}

func (s *newsService) GetNewsDetail(id string) (*models.News, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	s.newsRepo.IncrementViewCount(id)

	return news, nil
}

func (s *newsService) GetCategories() ([]models.CategoryInfo, error) {
	categories := models.GetAllCategories()
	
	// 获取每个分类的新闻数量
	for i, category := range categories {
		count, err := s.newsRepo.GetCategoryCount(category.Category)
		if err != nil {
			count = 0
		}
		categories[i].Count = count
	}

	return categories, nil
}

func (s *newsService) FetchAndSaveNews() error {
	if s.apiKey == "" {
		fmt.Println("没有API密钥，创建测试数据...")
		return s.createTestNews()
	}
	
	fmt.Printf("使用API密钥获取真实新闻数据: %s\n", s.maskAPIKey(s.apiKey))

	// 定义要获取的分类和对应的关键词（使用NewsAPI支持的分类）
	categoryKeywords := map[models.NewsCategory]string{
		models.CategoryTech:          "technology",
		models.CategorySports:        "sports",
		models.CategoryBusiness:      "business",
		models.CategoryHealth:        "health",
		models.CategoryEntertainment: "entertainment",
		models.CategoryGeneral:       "general",
		// 注意：NewsAPI不直接支持politics和military分类，我们使用general然后过滤
		models.CategoryPolitics:      "general",
		models.CategoryMilitary:      "general",
	}

	for category, keyword := range categoryKeywords {
		err := s.fetchCategoryNews(category, keyword)
		if err != nil {
			fmt.Printf("获取 %s 分类新闻失败: %v\n", category, err)
			continue
		}
		
		// 避免API限制，每个分类之间间隔1秒
		time.Sleep(1 * time.Second)
	}

	// 清理30天前的旧新闻
	s.newsRepo.DeleteOldNews(30)

	return nil
}

func (s *newsService) fetchCategoryNews(category models.NewsCategory, keyword string) error {
	// 使用NewsAPI获取新闻
	var url string
	if keyword == "general" {
		url = fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&pageSize=20&apiKey=%s", s.apiKey)
	} else {
		url = fmt.Sprintf("https://newsapi.org/v2/top-headlines?category=%s&country=us&pageSize=20&apiKey=%s", keyword, s.apiKey)
	}
	
	fmt.Printf("正在获取 %s 分类新闻: %s\n", category, url[:50]+"...")
	
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiResponse models.NewsAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return err
	}

	if apiResponse.Status != "ok" {
		return fmt.Errorf("API返回错误状态: %s", apiResponse.Status)
	}

	// 保存新闻到数据库
	for _, article := range apiResponse.Articles {
		// 检查新闻是否已存在
		exists, err := s.newsRepo.CheckExists(article.URL)
		if err != nil || exists {
			continue
		}

		news := &models.News{
			Title:       article.Title,
			Description: article.Description,
			Content:     article.Content,
			ImageURL:    article.URLToImage,
			SourceURL:   article.URL,
			Source:      article.Source.Name,
			Author:      article.Author,
			Category:    category,
			Tags:        s.extractTags(article.Title + " " + article.Description),
			PublishedAt: article.PublishedAt,
			ViewCount:   0,
		}

		if err := s.newsRepo.Create(news); err != nil {
			fmt.Printf("保存新闻失败: %v\n", err)
		}
	}

	return nil
}

func (s *newsService) RefreshNews() error {
	return s.FetchAndSaveNews()
}

// 从标题和描述中提取标签
func (s *newsService) extractTags(text string) []string {
	// 简单的标签提取逻辑，可以后续优化
	words := strings.Fields(strings.ToLower(text))
	tagMap := make(map[string]bool)
	var tags []string

	// 常见的新闻关键词
	keywords := []string{
		"technology", "tech", "ai", "artificial", "intelligence",
		"sports", "football", "basketball", "soccer", "olympics",
		"politics", "government", "election", "policy", "law",
		"military", "defense", "army", "navy", "security",
		"entertainment", "movie", "music", "celebrity", "film",
		"business", "economy", "market", "finance", "company",
		"health", "medical", "hospital", "doctor", "medicine",
	}

	for _, word := range words {
		for _, keyword := range keywords {
			if strings.Contains(word, keyword) && !tagMap[keyword] {
				tagMap[keyword] = true
				tags = append(tags, keyword)
				if len(tags) >= 5 { // 最多5个标签
					break
				}
			}
		}
		if len(tags) >= 5 {
			break
		}
	}

	return tags
}

// 创建测试新闻数据
func (s *newsService) createTestNews() error {
	testNews := []*models.News{
		{
			Title:       "人工智能技术取得重大突破",
			Description: "最新的AI技术在多个领域展现出惊人的能力，为未来发展奠定基础。",
			Content:     "人工智能技术在过去一年中取得了显著进展，特别是在自然语言处理和计算机视觉领域。专家预测，这些技术将在未来几年内彻底改变我们的工作和生活方式。",
			ImageURL:    "", // 暂时移除图片，避免加载问题
			SourceURL:   "https://example.com/ai-breakthrough",
			Source:      "科技日报",
			Author:      "张三",
			Category:    models.CategoryTech,
			Tags:        []string{"AI", "技术", "突破"},
			PublishedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			Title:       "体育赛事精彩纷呈",
			Description: "本周末的体育赛事为观众带来了精彩的比赛和难忘的时刻。",
			Content:     "本周末举行的多项体育赛事吸引了全球观众的关注。从足球到篮球，从网球到游泳，各项比赛都展现了运动员们的精湛技艺和顽强拼搏精神。",
			ImageURL:    "",
			SourceURL:   "https://example.com/sports-news",
			Source:      "体育周刊",
			Author:      "李四",
			Category:    models.CategorySports,
			Tags:        []string{"体育", "比赛", "运动员"},
			PublishedAt: time.Now().Add(-4 * time.Hour),
		},
		{
			Title:       "经济形势稳中向好",
			Description: "最新经济数据显示，各项指标呈现积极态势，为未来发展注入信心。",
			Content:     "根据最新发布的经济数据，本季度GDP增长超出预期，就业率保持稳定，通胀水平控制在合理区间。专家分析认为，这些积极信号表明经济正在稳步复苏。",
			ImageURL:    "",
			SourceURL:   "https://example.com/business-news",
			Source:      "财经时报",
			Author:      "王五",
			Category:    models.CategoryBusiness,
			Tags:        []string{"经济", "GDP", "就业"},
			PublishedAt: time.Now().Add(-6 * time.Hour),
		},
		{
			Title:       "健康生活新理念",
			Description: "专家分享最新的健康生活方式，帮助人们提升生活质量。",
			Content:     "随着人们对健康的重视程度不断提高，专家们提出了一系列新的健康生活理念。这些理念涵盖了饮食、运动、睡眠和心理健康等多个方面。",
			ImageURL:    "",
			SourceURL:   "https://example.com/health-news",
			Source:      "健康周刊",
			Author:      "赵六",
			Category:    models.CategoryHealth,
			Tags:        []string{"健康", "生活方式", "养生"},
			PublishedAt: time.Now().Add(-8 * time.Hour),
		},
		{
			Title:       "科技创新推动产业升级",
			Description: "新兴科技正在重塑传统产业，为经济发展注入新动力。",
			Content:     "随着5G、物联网、大数据等技术的快速发展，传统制造业正在经历数字化转型。这些技术不仅提高了生产效率，还创造了全新的商业模式。",
			ImageURL:    "",
			SourceURL:   "https://example.com/tech-innovation",
			Source:      "创新周报",
			Author:      "陈七",
			Category:    models.CategoryTech,
			Tags:        []string{"科技", "创新", "产业升级"},
			PublishedAt: time.Now().Add(-10 * time.Hour),
		},
		{
			Title:       "国际政治局势分析",
			Description: "专家解读当前国际政治形势，分析未来发展趋势。",
			Content:     "当前国际政治格局正在发生深刻变化，多极化趋势日益明显。各国在应对全球性挑战时需要加强合作，共同维护世界和平与稳定。",
			ImageURL:    "",
			SourceURL:   "https://example.com/politics-analysis",
			Source:      "国际观察",
			Author:      "刘八",
			Category:    models.CategoryPolitics,
			Tags:        []string{"政治", "国际", "分析"},
			PublishedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			Title:       "娱乐圈新动态",
			Description: "明星们的最新动态和娱乐圈的热门话题。",
			Content:     "本周娱乐圈热闹非凡，多位明星发布了新作品，各种颁奖典礼和活动也接连举行。粉丝们对偶像的支持热情不减，社交媒体上讨论热烈。",
			ImageURL:    "",
			SourceURL:   "https://example.com/entertainment-news",
			Source:      "娱乐快报",
			Author:      "周九",
			Category:    models.CategoryEntertainment,
			Tags:        []string{"娱乐", "明星", "新作品"},
			PublishedAt: time.Now().Add(-14 * time.Hour),
		},
		{
			Title:       "军事科技发展现状",
			Description: "分析当前军事科技的发展水平和未来趋势。",
			Content:     "现代军事科技正在快速发展，人工智能、无人系统、网络安全等技术在军事领域的应用越来越广泛。这些技术的发展对国防安全具有重要意义。",
			ImageURL:    "",
			SourceURL:   "https://example.com/military-tech",
			Source:      "军事评论",
			Author:      "吴十",
			Category:    models.CategoryMilitary,
			Tags:        []string{"军事", "科技", "国防"},
			PublishedAt: time.Now().Add(-16 * time.Hour),
		},
	}

	for _, news := range testNews {
		// 检查是否已存在
		exists, err := s.newsRepo.CheckExists(news.SourceURL)
		if err != nil || exists {
			continue
		}

		if err := s.newsRepo.Create(news); err != nil {
			fmt.Printf("创建测试新闻失败: %v\n", err)
		}
	}

	return nil
}

// 掩码API密钥用于日志显示
func (s *newsService) maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}