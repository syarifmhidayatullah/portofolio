package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type BlogHandler struct {
	postSvc service.PostService
}

func NewBlogHandler(postSvc service.PostService) *BlogHandler {
	return &BlogHandler{postSvc: postSvc}
}

func (h *BlogHandler) List(c *gin.Context) {
	posts, err := h.postSvc.GetAll(true)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "blog_list.html", gin.H{
		"title": "Blog",
		"posts": posts,
	})
}

func (h *BlogHandler) Detail(c *gin.Context) {
	slug := c.Param("slug")
	post, renderedHTML, err := h.postSvc.GetBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": "Post not found"})
		return
	}

	c.HTML(http.StatusOK, "blog_detail.html", gin.H{
		"title":       post.Title,
		"post":        post,
		"renderedHTML": renderedHTML,
	})
}
