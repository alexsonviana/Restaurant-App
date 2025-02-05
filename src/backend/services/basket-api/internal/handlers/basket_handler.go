package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jurabek/basket.api/internal/models"
	"github.com/jurabek/basket.api/internal/repositories"
	"github.com/pkg/errors"
)

type GetCreateDeleter interface {
	Get(ctx context.Context, customerID string) (*models.CustomerBasket, error)
	Update(ctx context.Context, item *models.CustomerBasket) error
	Delete(ctx context.Context, id string) error
}

// BasketHandler is router initializer for http
type BasketHandler struct {
	BasketRepository GetCreateDeleter
}

// NewBasketHandler creates new instance of BasketController with BasketRepository
func NewBasketHandler(r GetCreateDeleter) *BasketHandler {
	return &BasketHandler{BasketRepository: r}
}

// Create go doc
//
//	@Summary		Add a CustomerBasket
//	@Description	add by json new CustomerBasket
//	@Tags			CustomerBasket
//	@Accept			json
//	@Produce		json
//	@Param			CustomerBasket	body		models.CustomerBasket	true	"Add CustomerBasket"
//	@Success		200				{object}	models.CustomerBasket
//	@Failure		400				{object}	models.HTTPError
//	@Router			/items [post]
func (bc *BasketHandler) Create(c *gin.Context) {
	var entity models.CustomerBasket
	if err := c.BindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, models.NewHTTPError(http.StatusNotFound, err))
		return
	}

	err := bc.BasketRepository.Update(c.Request.Context(), &entity)

	if err != nil {
		httpError := models.NewHTTPError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, httpError)
		return
	}

	result, err := bc.BasketRepository.Get(c.Request.Context(), entity.CustomerID.String())
	if err != nil {
		httpError := models.NewHTTPError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, httpError)
		return
	}

	c.JSON(http.StatusOK, result)
}

// Get go doc
//
//	@Summary		Gets a CustomerBasket
//	@Description	Get CustomerBasket by ID
//	@Tags			CustomerBasket
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"CustomerBasket ID"
//	@Success		200	{object}	models.CustomerBasket
//	@Failure		400	{object}	models.HTTPError
//	@Router			/items/{id} [get]
func (bc *BasketHandler) Get(c *gin.Context) {
	id := c.Param("id")

	result, err := bc.BasketRepository.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrCartNotFound) {
			c.JSON(http.StatusNotFound, models.NewHTTPError(http.StatusNotFound, errors.Wrap(err, "itemID: "+id)))
			return
		}
		httpError := models.NewHTTPError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, httpError)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Delete go doc
//
//	@Summary		Deletes a CustomerBasket
//	@Description	Deletes CustomerBasket by ID
//	@Tags			CustomerBasket
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"CustomerBasket ID"
//	@Success		200	""
//	@Failure		400	{object}	models.HTTPError
//	@Router			/items/{id} [delete]
func (bc *BasketHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := bc.BasketRepository.Delete(c.Request.Context(), id)

	if err != nil {
		httpError := models.NewHTTPError(http.StatusBadRequest, err)
		c.JSON(http.StatusBadRequest, httpError)
		return
	}
	c.Status(http.StatusOK)
}
