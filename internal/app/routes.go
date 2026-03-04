package app

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/cyllective/olim/internal/db"
)

// GET /
//
// Returns the HTML for creating a new secret (landing page)
func getPageIndex(c echo.Context) error {
	return c.Render(http.StatusOK, "new", nil)
}

// GET /view
//
// Returns the HTML for viewing a secret
func getPageView(c echo.Context) error {
	return c.Render(http.StatusOK, "view", nil)
}

type newSecretStringRequest struct {
	Content     string `json:"content" validate:"required,base64"`
	ExpireHours int    `json:"expire_hours" validate:"required,oneof=1 8 24 48 72 120 192 720"`
}

// POST /api/string/new
//
// Creates a new secret string based on NewSecretStringRequest
func postStringNew(c echo.Context) error {
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 1<<20) // Limit to 1 MB

	var req newSecretStringRequest
	if err := c.Bind(&req); err != nil {
		log.Warning("new secret string request binding failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := c.Validate(req); err != nil {
		log.Warning("new secret string request validation failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	bytes, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		log.Error("could not decode base64 request content '%s': %s", req.Content, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "bad request"})
	}

	now := time.Now()
	secret := &db.SecretString{
		ID:        uuid.NewString(), // Can panic, should be covered by gin.Recovery()
		Value:     bytes,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(req.ExpireHours) * time.Hour),
	}

	if err := database.AddSecretString(secret); err != nil {
		log.Error("adding entry to db failed: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create db entry"})
	}

	log.Info("created new string entry %s", secret.ID)
	return c.JSON(http.StatusOK, map[string]string{"id": secret.ID})
}

type fetchSecretStringRequest struct {
	ID string `param:"id" validate:"required,uuid4"`
}

// GET /api/string/fetch/:id
//
// Gets a secret string from the database based on FetchSecretStringRequest
func getStringFetch(c echo.Context) error {
	var req fetchSecretStringRequest
	if err := c.Bind(&req); err != nil {
		log.Warning("fetch secret string request binding failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := c.Validate(req); err != nil {
		log.Warning("fetch secret string request validation failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	secret, err := database.GetAndDeleteSecretStringByID(req.ID)
	if err != nil {
		log.Warning("getting secret string %s from db failed: %s", req.ID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "secret not found or already viewed"})
	}
	log.Info("fetched and deleted string entry %s", secret.ID)

	b64secret := base64.StdEncoding.EncodeToString(secret.Value)

	return c.JSON(http.StatusOK, map[string]string{"content": b64secret})
}

type newSecretFileRequest struct {
	Name        string `json:"name" validate:"required,base64"`
	Content     string `json:"content" validate:"required,base64"`
	ExpireHours int    `json:"expire_hours" validate:"required,oneof=1 8 24 48 72"`
}

// POST /api/file/new
//
// Creates a new secret file based on NewSecretFileRequest
func postFileNew(c echo.Context) error {
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, 1<<20) // Limit to 1 MB

	var req newSecretFileRequest
	if err := c.Bind(&req); err != nil {
		log.Warning("new secret file request binding failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := c.Validate(req); err != nil {
		log.Warning("new secret file request validation failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	nameBytes, err := base64.StdEncoding.DecodeString(req.Name)
	if err != nil {
		log.Error("could not decode base64 request name '%s': %s", req.Name, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "bad request"})
	}

	contentBytes, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		log.Error("could not decode base64 request content '%s': %s", req.Content, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "bad request"})
	}

	now := time.Now()
	secret := &db.SecretFile{
		ID:        uuid.NewString(), // Can panic, should be covered by gin.Recovery()
		Name:      string(nameBytes),
		Value:     contentBytes,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(req.ExpireHours) * time.Hour),
	}

	if err := database.AddSecretFile(secret); err != nil {
		log.Error("adding entry to db failed: %s", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create db entry"})

	}

	log.Info("created new file entry %s", secret.ID)
	return c.JSON(http.StatusOK, map[string]string{"id": secret.ID})
}

type fetchSecretFileRequest struct {
	ID string `param:"id" validate:"required,uuid4"`
}

// GET /api/file/fetch/:id
//
// Gets a secret file from the database based on FetchSecretFileRequest
func getFileFetch(c echo.Context) error {
	var req fetchSecretFileRequest
	if err := c.Bind(&req); err != nil {
		log.Warning("fetch secret file request binding failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	if err := c.Validate(req); err != nil {
		log.Warning("fetch secret file request validation failed: %s", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "bad request"})
	}

	secret, err := database.GetAndDeleteSecretFileByID(req.ID)
	if err != nil {
		log.Warning("getting secret file %s from db failed: %s", req.ID, err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "secret not found or already viewed"})
	}
	log.Info("fetched and deleted file entry %s", secret.ID)

	b64name := base64.StdEncoding.EncodeToString([]byte(secret.Name))
	b64content := base64.StdEncoding.EncodeToString(secret.Value)

	return c.JSON(http.StatusOK, map[string]string{
		"name":    b64name,
		"content": b64content,
	})
}
