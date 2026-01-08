package service

import (
	"context"
	"fmt"
	"time"

	"github.com/arsenh/DevLore/internal/logger"
	"github.com/arsenh/DevLore/internal/model"
	"github.com/arsenh/DevLore/internal/repository"
	"github.com/arsenh/DevLore/internal/views"
)

type ArticleService struct {
	articleRepository model.ArticleRepository
	userRepository    model.UserRepository
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		articleRepository: repository.NewDummyArticleRepository(),
		userRepository:    repository.NewDummyUserRepository(),
	}
}

func (s *ArticleService) GetDashboardData(ctx context.Context) (*views.DashboardView, error) {
	articles, err := s.articleRepository.List(ctx)
	if err != nil {
		return nil, err
	}

	//TODO: update this to get current login user FullName
	user, _ := s.userRepository.FindByID(ctx, 11)

	dashboardView := views.DashboardView{
		BaseView: views.BaseView{
			UserName: user.FullName,
		},
		Articles: make([]views.ArticleShortItem, 0),
	}

	for _, article := range articles {
		articleItem := views.ArticleShortItem{
			ID:        article.ID,
			Title:     article.Title,
			UpdatedAt: article.UpdatedAt,
		}

		dashboardView.Articles = append(dashboardView.Articles, articleItem)
	}

	return &dashboardView, nil
}

func (s *ArticleService) GetArticleById(ctx context.Context, id int) (*views.ArticleView, error) {
	article, err := s.articleRepository.FindByID(ctx, id)
	if err != nil {
		logger.L().Warnf("the article id = %d not found", id)
		return nil, err
	}

	//TODO: update this to get current login user FullName
	user, _ := s.userRepository.FindByID(ctx, 11)

	view := &views.ArticleView{
		BaseView: views.BaseView{
			UserName: user.FullName,
		},
		Article: views.ArticleFullViewItem{
			ID:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		},
	}

	return view, nil
}

func (s *ArticleService) SaveArticle(ctx context.Context, title string, content string) (int, error) {
	//TODO: this function my be receive also current UserID, but for now it will be hardcoded.
	userId := 777 //TODO: remove this value, its for testing in dummy repository.
	id := 777     // TODO: this is article id, which needs to be generated in this place.
	if title == "" || content == "" {
		return -1, fmt.Errorf("the title and content of article must be not empty")
	}

	newArticle := &model.Article{
		ID:        id,
		Title:     title,
		Content:   content,
		UserId:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.articleRepository.Create(ctx, newArticle); err != nil {
		return -1, logger.LogAndErr("failed to store new article in database")
	}
	return id, nil
}

func (s *ArticleService) DeleteArticleById(ctx context.Context, id int) error {
	//TODO: Need to check userId, only article owner can delete!
	if err := s.articleRepository.Delete(ctx, id); err != nil {
		return logger.LogAndErr("failed to delete article from database")
	}
	return nil
}

func (s *ArticleService) EditArticleById(ctx context.Context, id int, title string, content string) (*views.ArticleView, error) {
	//TODO: Need to check userId, only article owner can edit!
	article, err := s.articleRepository.FindByID(ctx, id)
	if err != nil {
		return nil, logger.LogAndErr("failed to get article by id = %d", id)
	}

	newArticle := model.Article{
		ID:        article.ID,
		Title:     title,
		Content:   content,
		UserId:    article.UserId,
		CreatedAt: article.CreatedAt,
		UpdatedAt: time.Now(),
	}

	//TODO: optimize to have direct update function in repository.
	if err := s.DeleteArticleById(ctx, id); err != nil {
		return nil, logger.LogAndErr("failed to delete article by id = %d", id)
	}

	if err := s.articleRepository.Create(ctx, &newArticle); err != nil {
		return nil, logger.LogAndErr("failed to add updated article by id = %d", id)
	}

	//TODO: update this to get current login user FullName
	user, _ := s.userRepository.FindByID(ctx, 11)

	view := &views.ArticleView{
		BaseView: views.BaseView{
			UserName: user.FullName,
		},
		Article: views.ArticleFullViewItem{
			ID:        newArticle.ID,
			Title:     newArticle.Title,
			Content:   newArticle.Content,
			CreatedAt: newArticle.CreatedAt,
			UpdatedAt: newArticle.UpdatedAt,
		},
	}

	return view, nil
}

func (s *ArticleService) SearchArticlesByQuery(ctx context.Context, query string) (*views.SearchView, error) {

	//TODO: update this to get current login user FullName
	user, _ := s.userRepository.FindByID(ctx, 11)

	view := &views.SearchView{
		DashboardView: views.DashboardView{
			BaseView: views.BaseView{
				UserName: user.FullName,
			},
		},
		Query: query,
	}

	if query == "" {
		return view, nil
	}

	articles, err := s.articleRepository.Search(ctx, query)
	if err != nil {
		return nil, logger.LogAndErr("failed to search articles with query = %s", query)
	}

	for _, article := range articles {
		articleShortView := views.ArticleShortItem{
			ID:        article.ID,
			Title:     article.Title,
			UpdatedAt: article.UpdatedAt,
		}

		view.Articles = append(view.Articles, articleShortView)
	}

	return view, nil
}
