package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"nebula/pkg/s3client"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// S3ResourcesController handles serving content from S3
type S3ResourcesController struct {
	s3Client *s3client.S3Client
}

// NewS3ResourcesController creates a new S3ResourcesController
func NewS3ResourcesController() (*S3ResourcesController, error) {
	client, err := s3client.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &S3ResourcesController{
		s3Client: client,
	}, nil
}

// GetArticle retrieves an article from S3
func (c *S3ResourcesController) GetArticle() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		articleID := ctx.Param("id")
		if articleID == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Article ID is required",
			})
		}

		// Get article content from S3
		key := fmt.Sprintf("articles/%s.md", articleID)
		content, err := c.s3Client.GetObject(key)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Article not found",
			})
		}

		// Convert markdown to HTML with security measures
		unsafe := blackfriday.Run(content)
		html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

		// Get authentication status
		authenticated, _ := ctx.Get("authenticated").(bool)

		// Render article
		return ctx.Render(http.StatusOK, "article.html", map[string]interface{}{
			"Title":         "Article: " + articleID,
			"Content":       template.HTML(html),
			"ID":            articleID,
			"ActivePage":    "news",
			"Authenticated": authenticated,
		})
	}
}

// ListArticles lists all articles from S3
func (c *S3ResourcesController) ListArticles() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// List articles from S3
		keys, err := c.s3Client.ListObjects("articles/")
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to list articles",
			})
		}

		// Process article keys to get titles and IDs
		var articles []map[string]string
		for _, key := range keys {
			// Skip directories or non-markdown files
			if !strings.HasSuffix(key, ".md") {
				continue
			}

			// Extract ID from key (articles/my-article.md -> my-article)
			filename := filepath.Base(key)
			id := strings.TrimSuffix(filename, filepath.Ext(filename))

			// Get article content to extract title
			content, err := c.s3Client.GetObject(key)
			if err != nil {
				continue
			}

			// Extract title from first line (assuming # Title format)
			lines := strings.SplitN(string(content), "\n", 2)
			title := id
			if len(lines) > 0 {
				title = strings.TrimPrefix(lines[0], "# ")
				title = strings.TrimSpace(title)
			}

			articles = append(articles, map[string]string{
				"id":    id,
				"title": title,
				"url":   "/news/" + id,
			})
		}

		// Get authentication status
		authenticated, _ := ctx.Get("authenticated").(bool)

		// Render articles list
		return ctx.Render(http.StatusOK, "news.html", map[string]interface{}{
			"Title":         "News",
			"Articles":      articles,
			"ActivePage":    "news",
			"Authenticated": authenticated,
		})
	}
}

// GetImage serves an image from S3
func (c *S3ResourcesController) GetImage() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		imagePath := ctx.Param("*")
		if imagePath == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Image path is required",
			})
		}

		// Get image from S3
		key := fmt.Sprintf("images/%s", imagePath)
		imageData, err := c.s3Client.GetObject(key)
		if err != nil {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": "Image not found",
			})
		}

		// Determine content type based on file extension
		contentType := "image/jpeg" // Default
		ext := strings.ToLower(filepath.Ext(imagePath))
		switch ext {
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".svg":
			contentType = "image/svg+xml"
		case ".webp":
			contentType = "image/webp"
		}

		return ctx.Blob(http.StatusOK, contentType, imageData)
	}
}

// UploadImage handles image uploads to S3
func (c *S3ResourcesController) UploadImage() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Check if user is authenticated
		authenticated, _ := ctx.Get("authenticated").(bool)
		if !authenticated {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authentication required",
			})
		}

		// Get file from form
		file, err := ctx.FormFile("image")
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Image file is required",
			})
		}

		// Open the file
		src, err := file.Open()
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to open file",
			})
		}
		defer src.Close()

		// Read file content
		fileData := make([]byte, file.Size)
		_, err = src.Read(fileData)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to read file",
			})
		}

		// Determine content type
		contentType := file.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "image/jpeg" // Default
		}

		// Generate a unique filename
		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		key := fmt.Sprintf("images/%s", filename)

		// Upload to S3
		err = c.s3Client.UploadObject(key, fileData, contentType)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to upload image",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]string{
			"url": "/images/" + filename,
		})
	}
}
