package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/service"
)

type AdminPostHandler struct {
	postSvc service.PostService
}

func NewPostHandler(postSvc service.PostService) *AdminPostHandler {
	return &AdminPostHandler{postSvc: postSvc}
}

func (h *AdminPostHandler) List(c *gin.Context) {
	posts, err := h.postSvc.GetAll(false)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_posts.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "admin_posts.html", gin.H{
		"title": "Posts",
		"posts": posts,
	})
}

func (h *AdminPostHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_post_form.html", gin.H{
		"title": "New Post",
	})
}

func (h *AdminPostHandler) Create(c *gin.Context) {
	coverImage, err := saveUploadedImage(c, "cover_file", "cover_image", "")
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_post_form.html", gin.H{
			"title": "New Post",
			"error": err.Error(),
		})
		return
	}

	input := service.CreatePostInput{
		Title:      c.PostForm("title"),
		Content:    c.PostForm("content"),
		Excerpt:    c.PostForm("excerpt"),
		CoverImage: coverImage,
		Published:  c.PostForm("published") == "on",
	}

	if input.Title == "" || input.Content == "" {
		c.HTML(http.StatusBadRequest, "admin_post_form.html", gin.H{
			"title": "New Post",
			"error": "Title and content are required",
			"input": input,
		})
		return
	}

	_, err = h.postSvc.Create(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_post_form.html", gin.H{
			"title": "New Post",
			"error": err.Error(),
			"input": input,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminPostHandler) Edit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}

	post, err := h.postSvc.GetByID(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}

	c.HTML(http.StatusOK, "admin_post_form.html", gin.H{
		"title": "Edit Post",
		"post":  post,
	})
}

func (h *AdminPostHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}

	post, err := h.postSvc.GetByID(id)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}

	coverImage, err := saveUploadedImage(c, "cover_file", "cover_image", post.CoverImage)
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_post_form.html", gin.H{
			"title": "Edit Post",
			"error": err.Error(),
			"post":  post,
		})
		return
	}

	input := service.CreatePostInput{
		Title:      c.PostForm("title"),
		Content:    c.PostForm("content"),
		Excerpt:    c.PostForm("excerpt"),
		CoverImage: coverImage,
		Published:  c.PostForm("published") == "on",
	}

	if _, err := h.postSvc.Update(id, input); err != nil {
		c.HTML(http.StatusInternalServerError, "admin_post_form.html", gin.H{
			"title": "Edit Post",
			"error": err.Error(),
			"post":  post,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminPostHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/posts")
		return
	}

	h.postSvc.Delete(id)
	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminPostHandler) TogglePublish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.postSvc.TogglePublish(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/posts")
}
