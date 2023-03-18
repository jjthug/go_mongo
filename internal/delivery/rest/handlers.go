package rest

import (
	"context"
	"net/http"

	models "go_mongo/pkg"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateUser(ctx context.Context, name string, tags []string) (ID string, err error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUsersByTags(ctx context.Context, tags []string) ([]models.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	UpdateUser(ctx context.Context, user models.User) (bool, error)
}

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s: s}
}

func (h Handler) Router() *gin.Engine {
	router := gin.Default()
	router.POST("/", h.CreateUser)

	// called as localhost:3000/getOne/{id}
	router.GET("/:userId", h.ReadOneUser)
	router.GET("/getUserFromTags", h.GetUsersFromTags)

	// called as localhost:3000/update/{id}
	router.PUT("/:userId", h.UpdateUser)

	// called as localhost:3000/delete/{id}
	router.DELETE("/:userId", h.DeleteUser)

	return router
}

func (h Handler) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, FromErr(err))
		return
	}
	id, err := h.s.CreateUser(ctx, user.Name, user.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FromErr(err))
		return
	}

	ctx.JSON(http.StatusCreated, NewSuccessMessage(id))
}

func (h Handler) ReadOneUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	user, err := h.s.GetUserByID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FromErr(err))
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, FromString("user not found"))
		return
	}
	ctx.JSON(http.StatusOK, NewSuccessMessage(user))
}

type TagsRequest struct {
	QueryTags []string `json:"queryTags" binding:"required"`
}

func (h Handler) GetUsersFromTags(ctx *gin.Context) {
	// Bind the JSON payload to a TagsRequest struct
	var tagsRequest TagsRequest
	if err := ctx.ShouldBindJSON(&tagsRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, FromErr(err))
		return
	}

	// Extract the queryTags parameter from the struct
	tags := tagsRequest.QueryTags
	users, err := h.s.GetUsersByTags(ctx, tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FromErr(err))
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, FromString("users not found"))
		return
	}
	ctx.JSON(http.StatusOK, NewSuccessMessage(users))

}

func (h Handler) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	deleted, err := h.s.DeleteUser(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FromErr(err))
		return
	}
	if !deleted {
		ctx.JSON(http.StatusNotFound, FromString("No data to delete"))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

func (h Handler) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var user struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	updated, err := h.s.UpdateUser(ctx, models.User{
		ID:   userId,
		Name: user.Name,
		Tags: user.Tags,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, FromErr(err))
		return
	}
	if !updated {
		ctx.JSON(http.StatusNotFound, FromString("No data to update"))
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
