package service

import (
	"strings"
	"time"
	"unicode"

	"bytes"

	"github.com/google/uuid"
	"github.com/syarifmhidayatullah/portfolio/internal/model"
	"github.com/syarifmhidayatullah/portfolio/internal/repository"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type PostService interface {
	GetAll(onlyPublished bool) ([]model.Post, error)
	GetBySlug(slug string) (*model.Post, string, error) // returns post + rendered HTML
	GetByID(id uuid.UUID) (*model.Post, error)
	Create(input CreatePostInput) (*model.Post, error)
	Update(id uuid.UUID, input CreatePostInput) (*model.Post, error)
	Delete(id uuid.UUID) error
	TogglePublish(id uuid.UUID) error
	Count() (int64, error)
	Recent(limit int) ([]model.Post, error)
}

type CreatePostInput struct {
	Title      string
	Content    string
	Excerpt    string
	CoverImage string
	Tags       string // comma-separated raw input
	Published  bool
}

type postService struct {
	repo postRepo
	md   goldmark.Markdown
}

type postRepo interface {
	repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithHardWraps(), html.WithUnsafe()),
	)
	return &postService{repo: repo, md: md}
}

func (s *postService) GetAll(onlyPublished bool) ([]model.Post, error) {
	return s.repo.FindAll(onlyPublished)
}

func (s *postService) GetBySlug(slug string) (*model.Post, string, error) {
	post, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	if err := s.md.Convert([]byte(post.Content), &buf); err != nil {
		return nil, "", err
	}
	return post, buf.String(), nil
}

func (s *postService) GetByID(id uuid.UUID) (*model.Post, error) {
	return s.repo.FindByID(id)
}

func (s *postService) Create(input CreatePostInput) (*model.Post, error) {
	slug := slugify(input.Title)
	excerpt := input.Excerpt
	if excerpt == "" {
		excerpt = truncate(stripMarkdown(input.Content), 160)
	}

	post := &model.Post{
		Title:      input.Title,
		Slug:       slug,
		Excerpt:    excerpt,
		Content:    input.Content,
		CoverImage: input.CoverImage,
		Tags:       parseTags(input.Tags),
		Published:  input.Published,
	}

	if input.Published {
		now := time.Now()
		post.PublishedAt = &now
	}

	return post, s.repo.Create(post)
}

func (s *postService) Update(id uuid.UUID, input CreatePostInput) (*model.Post, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	post.Title = input.Title
	post.Content = input.Content
	post.CoverImage = input.CoverImage
	post.Tags = parseTags(input.Tags)

	excerpt := input.Excerpt
	if excerpt == "" {
		excerpt = truncate(stripMarkdown(input.Content), 160)
	}
	post.Excerpt = excerpt

	if input.Published && !post.Published {
		now := time.Now()
		post.PublishedAt = &now
	}
	post.Published = input.Published

	return post, s.repo.Update(post)
}

func (s *postService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *postService) TogglePublish(id uuid.UUID) error {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	post.Published = !post.Published
	if post.Published && post.PublishedAt == nil {
		now := time.Now()
		post.PublishedAt = &now
	}
	return s.repo.Update(post)
}

func (s *postService) Count() (int64, error) {
	return s.repo.Count()
}

func (s *postService) Recent(limit int) ([]model.Post, error) {
	return s.repo.Recent(limit)
}

// slugify converts a title to a URL-friendly slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r) || r == '-':
			b.WriteRune('-')
		}
	}
	slug := strings.Trim(b.String(), "-")
	// Add timestamp suffix to ensure uniqueness
	return slug + "-" + time.Now().Format("20060102")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func parseTags(raw string) []string {
	var result []string
	for _, t := range strings.Split(raw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}

func stripMarkdown(s string) string {
	// Very simple: remove markdown symbols
	replacer := strings.NewReplacer(
		"#", "", "*", "", "_", "", "`", "", "[", "", "]", "",
		"(", "", ")", "", "!", "", ">", "",
	)
	return strings.TrimSpace(replacer.Replace(s))
}
