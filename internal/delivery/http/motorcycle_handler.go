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

type MotorcycleHandler struct {
	service   *service.MotorcycleService
	validator *auth.Validator
}

func NewMotorcycleHandler(service *service.MotorcycleService, validator *auth.Validator) *MotorcycleHandler {
	return &MotorcycleHandler{service, validator}
}

// GetMotorcycleCategories godoc
// @Summary Get motorcycle categories
// @Description Get motorcycle categories
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Success 200 {array} model.GetMotorcycleCategoriesResponse
// @Failure 500 {object} model.ResultMessage
// @Router /motorcycles/categories [get]
func (h *MotorcycleHandler) GetMotorcycleCategories(c *fiber.Ctx) error {

	ctx := c.Context()
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetMotorcycleCategories(ctx, lang))

}

// GetMotorcycleParameters godoc
// @Summary Get motorcycle parameters
// @Description Get motorcycle parameters
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Success 200 {array} model.GetMotorcycleParametersResponse
// @Failure 500 {object} model.ResultMessage
// @Router /motorcycles/categories/{category_id}/parameters [get]
func (h *MotorcycleHandler) GetMotorcycleParameters(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetMotorcycleParameters(ctx, categoryID, lang))
}

// GetMotorcycleBrands godoc
// @Summary Get motorcycle brands
// @Description Get motorcycle brands
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Success 200 {array} model.GetMotorcycleBrandsResponse
// @Failure 500 {object} model.ResultMessage
// @Router /motorcycles/categories/{category_id}/brands [get]
func (h *MotorcycleHandler) GetMotorcycleBrands(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetMotorcycleBrands(ctx, categoryID, lang))

}

// GetMotorcycleModelsByBrandID godoc
// @Summary Get motorcycle models by brand ID
// @Description Get motorcycle models by brand ID
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param category_id path string true "Category ID"
// @Param brand_id path string true "Brand ID"
// @Success 200 {array} model.GetMotorcycleModelsResponse
// @Failure 500 {object} model.ResultMessage
// @Router /motorcycles/categories/{category_id}/brands/{brand_id}/models [get]
func (h *MotorcycleHandler) GetMotorcycleModelsByBrandID(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryID := c.Params("category_id")
	brandID := c.Params("brand_id")
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetMotorcycleModelsByBrandID(ctx, categoryID, brandID, lang))
}

// GetMotorcycles godoc
// @Summary Get motorcycles
// @Description Get motorcycles
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Success 200 {array} model.GetMotorcyclesResponse
// @Failure 500 {object} model.ResultMessage
// @Router /motorcycles [get]
func (h *MotorcycleHandler) GetMotorcycles(c *fiber.Ctx) error {
	ctx := c.Context()
	lang := c.Locals("lang").(string)
	return utils.FiberResponse(c, h.service.GetMotorcycles(ctx, lang))

}

// CreateMotorcycle godoc
// @Summary Create motorcycle
// @Description Create motorcycle
// @Tags motorcycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param motorcycle body model.CreateMotorcycleRequest true "Motorcycle"
// @Success 200 {object} model.SuccessWithId
// @Failure 500 {object} model.ResultMessage
// @Failure 400 {object} model.ResultMessage
// @Router /motorcycles [post]
func (h *MotorcycleHandler) CreateMotorcycle(c *fiber.Ctx) error {
	ctx := c.Context()
	var motorcycle model.CreateMotorcycleRequest

	if err := c.BodyParser(&motorcycle); err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  err,
		})
	}

	userID := c.Locals("id").(int)

	createdMotorcycle := h.service.CreateMotorcycle(ctx, motorcycle, userID)

	return utils.FiberResponse(c, createdMotorcycle)
}

// CreateMotorcycleImages godoc
// @Summary      Upload motorcycle images
// @Description  Uploads images for a motorcycle (max 10 files)
// @Tags         motorcycles
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        motorcycle_id      path      int     true   "Motorcycle ID"
// @Param        images  formData  file    true   "Motorcycle images (max 10)"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /motorcycles/{motorcycle_id}/images [post]
func (h *MotorcycleHandler) CreateMotorcycleImages(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid motorcycle ID"),
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

	paths, status, err := files.SaveFiles(images, config.ENV.STATIC_PATH+"motorcycles/"+strconv.Itoa(id), config.ENV.DEFAULT_IMAGE_WIDTHS)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: status,
			Error:  err,
		})

	}

	data := h.service.CreateMotorcycleImages(ctx, id, paths)
	return utils.FiberResponse(c, data)
}

// CreateMotorcycleVideos godoc
// @Summary      Upload motorcycle videos
// @Description  Uploads videos for a motorcycle (max 1 files)
// @Tags         motorcycles
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        motorcycle_id      path      int     true   "Motorcycle ID"
// @Param        videos  formData  file    true   "Motorcycle videos (max 10)"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /motorcycles/{motorcycle_id}/videos [post]
func (h *MotorcycleHandler) CreateMotorcycleVideos(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid motorcycle ID"),
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
	// path, err := pkg.SaveVideos(videos[0], config.ENV.STATIC_PATH+"motorcycles/"+idStr+"/videos") // if have ffmpeg on server
	path, err := files.SaveOriginal(videos[0], config.ENV.STATIC_PATH+"motorcycles/"+idStr+"/videos")

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  err,
		})

	}

	data := h.service.CreateMotorcycleVideos(ctx, id, path)
	return utils.FiberResponse(c, data)
}

// DeleteMotorcycleImage godoc
// @Summary      Delete motorcycle image
// @Description  Deletes an image from a motorcycle
// @Tags         motorcycles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        motorcycle_id      path      int     true   "Motorcycle ID"
// @Param        image_id      path      int     true   "Image ID"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /motorcycles/{motorcycle_id}/images/{image_id} [delete]
func (h *MotorcycleHandler) DeleteMotorcycleImage(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid motorcycle ID"),
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

	data := h.service.DeleteMotorcycleImage(ctx, id, imageID)

	return utils.FiberResponse(c, data)
}

// DeleteMotorcycleVideo godoc
// @Summary      Delete motorcycle video
// @Description  Deletes a video from a motorcycle
// @Tags         motorcycles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        motorcycle_id      path      int     true   "Motorcycle ID"
// @Param        video_id      path      int     true   "Video ID"
// @Success      200     {object}  model.Success
// @Failure      400     {object}  model.ResultMessage
// @Failure      401     {object}  auth.ErrorResponse
// @Failure	 	 403  	 {object}  auth.ErrorResponse
// @Failure      404     {object}  model.ResultMessage
// @Failure      500     {object}  model.ResultMessage
// @Router       /motorcycles/{motorcycle_id}/videos/{video_id} [delete]
func (h *MotorcycleHandler) DeleteMotorcycleVideo(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("invalid motorcycle ID"),
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

	data := h.service.DeleteMotorcycleVideo(ctx, id, videoID)

	return utils.FiberResponse(c, data)
}

// GetMotorcycleByID godoc
// @Summary      Get motorcycle by ID
// @Description  Returns a motorcycle by its ID
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.GetMotorcyclesResponse
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id} [get]
func (h *MotorcycleHandler) GetMotorcycleByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	lang := c.Locals("lang").(string)
	data := h.service.GetMotorcycleByID(ctx, id, userID, lang)
	return utils.FiberResponse(c, data)
}

// GetEditMotorcycleByID godoc
// @Summary      Get Edit motorcycle by ID
// @Description  Returns a motorcycle by its ID for editing
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Security 	 BearerAuth
// @Param   Accept-Language  header  string  false  "Language"
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.GetMotorcyclesResponse
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id}/edit [get]
func (h *MotorcycleHandler) GetEditMotorcycleByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	lang := c.Locals("lang").(string)
	data := h.service.GetEditMotorcycleByID(ctx, id, userID, lang)
	return utils.FiberResponse(c, data)
}

// BuyMotorcycle godoc
// @Summary      Buy motorcycle
// @Description  Returns a status response message
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id}/buy [post]
func (h *MotorcycleHandler) BuyMotorcycle(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	data := h.service.BuyMotorcycle(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// DontSellMotorcycle godoc
// @Summary      Set motorcycle as not for sale
// @Description  Updates motorcycle status to not for sale
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id}/dont-sell [post]
func (h *MotorcycleHandler) DontSellMotorcycle(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	data := h.service.DontSellMotorcycle(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// SellMotorcycle godoc
// @Summary      Set motorcycle for sale
// @Description  Updates motorcycle status to for sale
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id}/sell [post]
func (h *MotorcycleHandler) SellMotorcycle(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	userID := c.Locals("id").(int)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	data := h.service.SellMotorcycle(ctx, id, userID)
	return utils.FiberResponse(c, data)
}

// DeleteMotorcycle godoc
// @Summary      Delete motorcycle
// @Description  Deletes a motorcycle and its associated files
// @Tags         motorcycles
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Motorcycle ID"
// @Success      200  {object}  model.Success
// @Failure      400  {object}  model.ResultMessage
// @Failure      401  {object}  auth.ErrorResponse
// @Failure      403  {object}  auth.ErrorResponse
// @Failure      404  {object}  model.ResultMessage
// @Failure      500  {object}  model.ResultMessage
// @Router       /motorcycles/{id} [delete]
func (h *MotorcycleHandler) DeleteMotorcycle(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return utils.FiberResponse(c, model.Response{
			Status: 400,
			Error:  errors.New("motorcycle id must be integer"),
		})
	}

	// Create directory path for file cleanup
	dir := config.ENV.STATIC_PATH + "motorcycles/" + idStr

	data := h.service.DeleteMotorcycle(ctx, id, dir)
	return utils.FiberResponse(c, data)
}
