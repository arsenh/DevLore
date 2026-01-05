package service

import (
	"context"

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
		logger.L().Infof("the article id = %d not found", id)
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
