package service

import (
	"errors"

	"github.com/valyala/fasthttp"

	"dubai-auto/internal/model"
	"dubai-auto/internal/repository"
	"dubai-auto/pkg/files"
)

type ComtransService struct {
	repository *repository.ComtransRepository
}

func NewComtransService(repository *repository.ComtransRepository) *ComtransService {
	return &ComtransService{repository}
}

func (s *ComtransService) GetComtransCategories(ctx *fasthttp.RequestCtx, lang string) model.Response {
	data, err := s.repository.GetComtransCategories(ctx, lang)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) GetComtransParameters(ctx *fasthttp.RequestCtx, categoryID string, lang string) model.Response {
	data, err := s.repository.GetComtransParameters(ctx, categoryID, lang)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) GetComtransBrands(ctx *fasthttp.RequestCtx, categoryID string, lang string) model.Response {
	data, err := s.repository.GetComtransBrands(ctx, categoryID, lang)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) GetComtransModelsByBrandID(ctx *fasthttp.RequestCtx, categoryID string, brandID string, lang string) model.Response {
	data, err := s.repository.GetComtransModelsByBrandID(ctx, categoryID, brandID, lang)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) CreateComtrans(ctx *fasthttp.RequestCtx, comtrans model.CreateComtransRequest, userID int) model.Response {
	data, err := s.repository.CreateComtrans(ctx, comtrans, userID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) GetComtrans(ctx *fasthttp.RequestCtx, lang string) model.Response {
	data, err := s.repository.GetComtrans(ctx, lang)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	return model.Response{
		Status: 200,
		Data:   data,
	}
}

func (s *ComtransService) CreateComtransImages(ctx *fasthttp.RequestCtx, comtransID int, images []string) model.Response {

	if len(images) == 0 {
		return model.Response{
			Status: 400,
			Error:  errors.New("images are required"),
		}
	}
	err := s.repository.CreateComtransImages(ctx, comtransID, images)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Commercial transport images created successfully"},
	}
}

func (s *ComtransService) CreateComtransVideos(ctx *fasthttp.RequestCtx, comtransID int, video string) model.Response {
	err := s.repository.CreateComtransVideos(ctx, comtransID, video)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Commercial transport videos created successfully"},
	}
}

func (s *ComtransService) DeleteComtransImage(ctx *fasthttp.RequestCtx, comtransID int, imageID int) model.Response {

	err := s.repository.DeleteComtransImage(ctx, comtransID, imageID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	// todo: delete image if exists
	return model.Response{
		Data: model.Success{Message: "Commercial transport image deleted successfully"},
	}
}

func (s *ComtransService) DeleteComtransVideo(ctx *fasthttp.RequestCtx, comtransID int, videoID int) model.Response {

	err := s.repository.DeleteComtransVideo(ctx, comtransID, videoID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	// todo: delete video if exists
	return model.Response{
		Data: model.Success{Message: "Commercial transport video deleted successfully"},
	}
}

func (s *ComtransService) GetComtransByID(ctx *fasthttp.RequestCtx, comtransID, userID int, lang string) model.Response {
	comtrans, err := s.repository.GetComtransByID(ctx, comtransID, userID, lang)
	if err != nil {
		return model.Response{
			Status: 404,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   comtrans,
	}
}

func (s *ComtransService) GetEditComtransByID(ctx *fasthttp.RequestCtx, comtransID, userID int, lang string) model.Response {
	comtrans, err := s.repository.GetEditComtransByID(ctx, comtransID, userID, lang)

	if err != nil {
		return model.Response{
			Status: 404,
			Error:  err,
		}
	}

	return model.Response{
		Status: 200,
		Data:   comtrans,
	}
}

func (s *ComtransService) BuyComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) model.Response {
	err := s.repository.BuyComtrans(ctx, comtransID, userID)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	return model.Response{
		Status: 200,
		Data:   model.Success{Message: "Successfully bought the commercial transport"},
	}
}

func (s *ComtransService) DontSellComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) model.Response {
	err := s.repository.DontSellComtrans(ctx, comtransID, userID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   model.Success{Message: "Successfully updated commercial transport status to not for sale"},
	}
}

func (s *ComtransService) SellComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) model.Response {
	err := s.repository.SellComtrans(ctx, comtransID, userID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	return model.Response{
		Status: 200,
		Data:   model.Success{Message: "Successfully updated commercial transport status to for sale"},
	}
}

func (s *ComtransService) DeleteComtrans(ctx *fasthttp.RequestCtx, comtransID int, dir string) model.Response {
	err := s.repository.DeleteComtrans(ctx, comtransID)
	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	// todo: delete associated files
	if dir != "" {
		files.RemoveFolder(dir)
	}
	return model.Response{
		Status: 200,
		Data:   model.Success{Message: "Successfully deleted commercial transport"},
	}
}
