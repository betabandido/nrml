package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"nrml/products"
)

type ProductHandler struct {
	productGetter products.ProductGetter
}

func NewProductHandler(
	productGetter products.ProductGetter,
) *ProductHandler {
	return &ProductHandler{
		productGetter: productGetter,
	}
}

func (h *ProductHandler) RegisterRoutes(
	router *gin.RouterGroup,
	handlers ...gin.HandlerFunc,
) {
	router.GET("/productByProductKey/:tenant/:locale/:productKey",
		append(handlers,
			func(ctx *gin.Context) {
				h.GetProductByProductKey(ctx)
			},
		)...,
	)
}

type productKeyRequest struct {
	Tenant     string `uri:"tenant" binding:"required"`
	Locale     string `uri:"locale" binding:"required"`
	ProductKey string `uri:"productKey" binding:"required"`
}

func (h *ProductHandler) GetProductByProductKey(ctx *gin.Context) {
	var request productKeyRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	product, err := h.productGetter.GetProductByProductKey(
		ctx,
		request.Tenant,
		request.Locale,
		request.ProductKey,
	)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, product)
}
