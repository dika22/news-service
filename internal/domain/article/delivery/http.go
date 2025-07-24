package delivery

import (
	"net/http"
	"strconv"

	"github.com/dika22/news-service/internal/domain/article/usecase"
	"github.com/dika22/news-service/package/response"
	"github.com/dika22/news-service/package/structs"
	"github.com/dika22/news-service/package/validator"
	"github.com/spf13/cast"

	"github.com/labstack/echo/v4"
)

type ArticleHTTP struct{
	uc usecase.IArticle
	v *validator.Validator
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

	req := new(structs.RequestCreateArticle)
	if err := h.v.UnmarshallJSONValidate(c, req); err != nil {
		return response.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
	}
	
	if err := h.uc.Create(ctx, req); err != nil {
		return err
	}
	return response.JSONSuccess(c, req, "success create Article")
}

// UpdateArticle godoc
// @Summary      Update new article
// @Description  Update article with title and body
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        request body structs.RequestUpdatePublishArticle true "Article to update"
// @Success      200  {object}  structs.Response
// @Router /api/v1/articles/:id [put]
func (h ArticleHTTP) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if id == "" {
		return response.JSONResponse(c, http.StatusBadRequest, "error", "id required", nil)
	}
	req := new(structs.RequestUpdatePublishArticle)
	if err := h.v.UnmarshallJSONValidate(c, req); err != nil {
		return response.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
	}
	req.ID  = cast.ToInt64(id)
	if err := h.uc.UpdatePublishArticle(ctx, req); err != nil {
		return err
	}
	return response.JSONSuccess(c, req, "success update Article")
}

// GetArticleByID godoc
// @Summary      Get article by id
// @Description  Get article with id
// @Tags         articles
// @Accept       json
// @Produce      json
// @Param        id path int true "Article ID"
// @Success      200      {object}  structs.Response
// @Router       /api/v1/articles/:id [get]
func (h ArticleHTTP) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if id == "" {
		return response.JSONResponse(c, http.StatusBadRequest, "error", "id required", nil)
	}
	
	resp, err := h.uc.GetByID(ctx, cast.ToInt64(id))
	if err != nil {
		return response.JSONResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
	}
	return response.JSONSuccess(c, resp, "success get Article by id")
}

func NewArticleHTTP(r *echo.Group, uc usecase.IArticle, v *validator.Validator)  {
	u := ArticleHTTP{uc: uc, v: v}
	r.GET("", u.GetAll).Name = "article.get-all"
	r.POST("", u.Create).Name = "article.create"
	r.PUT("/:id", u.Update).Name = "article.update"
	r.GET("/:id", u.GetByID).Name = "article.get-by-id"
}