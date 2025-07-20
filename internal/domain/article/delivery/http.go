package delivery

import (
	"net/http"
	"news-service/internal/domain/article/usecase"
	"news-service/package/response"
	"news-service/package/structs"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ArticleHTTP struct{
	uc usecase.IArticle
}

// GetArticles godoc
// @Summary      Get all articles
// @Description  Get articles with optional search keyword
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        keyword  query     string  false  "Keyword for search"
// @Param        page     query     int     false  "Page number"
// @Success      200      {object}  structs.Response
// @Router       /api/v1/articles [get]
func (h ArticleHTTP) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	keyword := c.QueryParam("keyword")
	orderBy := c.QueryParam("order_by")

	req := structs.RequestSearchArticle{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
		OrderBy: orderBy,
	}

	resp, err := h.uc.GetAll(ctx, req); 
	if err != nil {
		return response.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
	}
	return response.JSONSuccess(c, resp, "success get all Article")
}


// CreateArticle godoc
// @Summary      Create new article
// @Description  Create article with title and body
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        request body structs.RequestCreateArticle true "Article to create"
// @Success      201  {object}  structs.Response
// @Router /api/v1/articles [post]
func (h ArticleHTTP) Create(c echo.Context) error {
	ctx := c.Request().Context()
	req := &structs.RequestCreateArticle{}
	if err := c.Bind(req); err != nil {
		return response.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
	}
	if err := h.uc.Create(ctx, req); err != nil {
		return err
	}
	return response.JSONSuccess(c, nil, "success create Article")
}

func NewArticleHTTP(r *echo.Group, uc usecase.IArticle)  {
	u := ArticleHTTP{uc: uc}
	r.GET("", u.GetAll).Name = "article.get-all"
	r.POST("", u.Create).Name = "article.create"
}