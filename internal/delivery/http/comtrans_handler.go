package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/internal/service"
	"dubai-auto/internal/utils"
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/files"
)

type ComtransHandler struct {
	service   *service.ComtransService
	validator *auth.Validator
}

func NewComtransHandler(service *service.ComtransService, validator *auth.Validator) *ComtransHandler {
	return &ComtransHandler{service, validator}
}

// GetComtransCategories godoc
// @Summary Get commercial transport categories
// @Description Get commercial transport categories
// @Tags comtrans
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Success 200 {array} model.GetComtransCategoriesResponse
// @Failure 500 {object} model.ResultMessage
// @Router /api/v1/comtrans/categories [get]
func (h *ComtransHandler) GetComtransCategories(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	ctx := c.Context()
	return utils.FiberResponse(c, h.service.GetComtransCategories(ctx, lang))

}

// GetComtransParameters godoc
// @Summary Get commercial transport parameters
// @Description Get commercial transport parameters
// @Tags comtrans
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Success 200 {array} model.GetComtransParametersResponse
// @Failure 500 {object} model.ResultMessage
// @Router /api/v1/comtrans/categories/{category_id}/parameters [get]
func (h *ComtransHandler) GetComtransParameters(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetComtransParameters(ctx, categoryID, lang))

}

// GetComtransBrands godoc
// @Summary Get commercial transport brands
// @Description Get commercial transport brands
// @Tags comtrans
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Success 200 {array} model.GetComtransBrandsResponse
// @Failure 500 {object} model.ResultMessage
// @Router /api/v1/comtrans/categories/{category_id}/brands [get]
func (h *ComtransHandler) GetComtransBrands(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetComtransBrands(ctx, categoryID, lang))

}

// GetComtransModelsByBrandID godoc
// @Summary Get commercial transport models by brand ID
// @Description Get commercial transport models by brand ID
// @Tags comtrans
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Param brand_id path string true "Brand ID"
// @Success 200 {array} model.GetComtransModelsResponse
// @Failure 500 {object} model.ResultMessage
// @Router /api/v1/comtrans/categories/{category_id}/brands/{brand_id}/models [get]
func (h *ComtransHandler) GetComtransModelsByBrandID(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	brandID := c.Params("brand_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetComtransModelsByBrandID(ctx, categoryID, brandID, lang))
}

// GetComtrans godoc
// @Summary Get commercial transports
// @Description Get commercial transports
// @Tags comtrans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Success 200 {array} model.GetComtransResponse
// @Failure 500 {object} model.ResultMessage
// @Router /api/v1/comtrans [get]
func (h *ComtransHandler) GetComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetComtrans(ctx, lang))

}

// CreateComtrans godoc
// @Summary Create commercial transport
// @Description Create commercial transport
// @Tags comtrans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comtrans body model.CreateComtransRequest true "Commercial Transport"
// @Success 200 {object} model.SuccessWithId
// @Failure 500 {object} model.ResultMessage
// @Failure 400 {object} model.ResultMessage
// @Router /api/v1/comtrans [post]
func (h *ComtransHandler) CreateComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	var comtrans model.CreateComtransRequest

	if err := c.BodyParser(&comtrans); err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  err,
		})
	}

	userID := c.Locals("id").(int)

	createdComtrans := h.service.CreateComtrans(ctx, comtrans, userID)

	return utils.FiberResponse(c, createdComtrans)
}

// CreateComtransImages godoc
// @Summary      Upload commercial transport images
// @Description  Uploads images for a commercial transport (max 10 files)
// @Tags         comtrans
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        comtrans_id      path      int     true   "Commercial Transport ID"
// @Param        images  formData  file    true   "Commercial transport images (max 10)"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /comtrans/{comtrans_id}/images [post]
func (h *ComtransHandler) CreateComtransImages(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid commercial transport ID"),
		})

	}

	form, _ := c.MultipartForm()

	if form == nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		})

	}

	images := form.File["images"]

	if len(images) > 10 {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 10 files"),
		})

	}

	paths, status, err := files.SaveFiles(images, config.ENV.STATIC_PATH+"comtrans/"+strconv.Itoa(id), config.ENV.DEFAULT_IMAGE_WIDTHS)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: status,
			Error:  err,
		})

	}

	data := h.service.CreateComtransImages(ctx, id, paths)
	return utils.FiberResponse(c, data)
}

// CreateComtransVideos godoc
// @Summary      Upload commercial transport videos
// @Description  Uploads videos for a commercial transport (max 1 files)
// @Tags         comtrans
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        comtrans_id      path      int     true   "Commercial Transport ID"
// @Param        videos  formData  file    true   "Commercial transport videos (max 1)"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /comtrans/{comtrans_id}/videos [post]
func (h *ComtransHandler) CreateComtransVideos(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid commercial transport ID"),
		})

	}

	form, _ := c.MultipartForm()

	if form == nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		})

	}

	videos := form.File["videos"]

	if len(videos) > 1 {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file(s)"),
		})

	}
	// path, err := pkg.SaveVideos(videos[0], config.ENV.STATIC_PATH+"comtrans/"+idStr+"/videos") // if have ffmpeg on server
	path, err := files.SaveOriginal(videos[0], config.ENV.STATIC_PATH+"comtrans/"+idStr+"/videos")

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  err,
		})

	}

	data := h.service.CreateComtransVideos(ctx, id, path)
	return utils.FiberResponse(c, data)
}

// DeleteComtransImage godoc
// @Summary      Delete commercial transport image
// @Description  Deletes an image from a commercial transport
// @Tags         comtrans
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        comtrans_id      path      int     true   "Commercial Transport ID"
// @Param        image_id      path      int     true   "Image ID"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /comtrans/{comtrans_id}/images/{image_id} [delete]
func (h *ComtransHandler) DeleteComtransImage(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid commercial transport ID"),
		})
	}

	imageIDStr := c.Params("image_id")
	imageID, err := strconv.Atoi(imageIDStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid image ID"),
		})
	}

	data := h.service.DeleteComtransImage(ctx, id, imageID)

	return utils.FiberResponse(c, data)
}

// DeleteComtransVideo godoc
// @Summary      Delete commercial transport video
// @Description  Deletes a video from a commercial transport
// @Tags         comtrans
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        comtrans_id      path      int     true   "Commercial Transport ID"
// @Param        video_id      path      int     true   "Video ID"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /comtrans/{comtrans_id}/videos/{video_id} [delete]
func (h *ComtransHandler) DeleteComtransVideo(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid commercial transport ID"),
		})
	}

	videoIDStr := c.Params("video_id")
	videoID, err := strconv.Atoi(videoIDStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid video ID"),
		})
	}

	data := h.service.DeleteComtransVideo(ctx, id, videoID)

	return utils.FiberResponse(c, data)
}

// GetComtransByID godoc
// @Summary      Get commercial transport by ID
// @Description  Returns a commercial transport by its ID
// @Tags         comtrans
// @Security     BearerAuth
// @Produce      json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.GetComtransResponse
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id} [get]
func (h *ComtransHandler) GetComtransByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	lang := c.Locals("lang").(string)
	data := h.service.GetComtransByID(ctx, id, userID, lang)
	return utils.FiberResponse(c, data)
}

// GetEditComtransByID godoc
// @Summary      Get Edit commercial transport by ID
// @Description  Returns a commercial transport by its ID for editing
// @Tags         comtrans
// @Security     BearerAuth
// @Produce      json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.GetComtransResponse
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id}/edit [get]
func (h *ComtransHandler) GetEditComtransByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	lang := c.Locals("lang").(string)
	data := h.service.GetEditComtransByID(ctx, id, userID, lang)
	return utils.FiberResponse(c, data)
}

// BuyComtrans godoc
// @Summary      Buy commercial transport
// @Description  Returns a status response message
// @Tags         comtrans
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id}/buy [post]
func (h *ComtransHandler) BuyComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	data := h.service.BuyComtrans(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// DontSellComtrans godoc
// @Summary      Set commercial transport as not for sale
// @Description  Updates commercial transport status to not for sale
// @Tags         comtrans
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id}/dont-sell [post]
func (h *ComtransHandler) DontSellComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	data := h.service.DontSellComtrans(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// SellComtrans godoc
// @Summary      Set commercial transport for sale
// @Description  Updates commercial transport status to for sale
// @Tags         comtrans
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id}/sell [post]
func (h *ComtransHandler) SellComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	data := h.service.SellComtrans(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// DeleteComtrans godoc
// @Summary      Delete commercial transport
// @Description  Deletes a commercial transport and its associated files
// @Tags         comtrans
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Commercial Transport ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /comtrans/{id} [delete]
func (h *ComtransHandler) DeleteComtrans(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("commercial transport id must be integer"),
		})
	}

	// Create directory path for file cleanup
	dir := config.ENV.STATIC_PATH + "comtrans/" + idStr

	data := h.service.DeleteComtrans(ctx, id, dir)
	return utils.FiberResponse(c, data)
}
