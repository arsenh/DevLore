package service

import (
	"context"

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
		UserName: user.FullName,
		Articles: make([]views.ArticleItem, 0),
	}

	for _, article := range articles {
		articleItem := views.ArticleItem{
			ID:        article.ID,
			Title:     article.Title,
			UpdatedAt: article.UpdatedAt,
		}

		dashboardView.Articles = append(dashboardView.Articles, articleItem)
	}

	return &dashboardView, nil
}
