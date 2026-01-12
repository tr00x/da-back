package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/utils"
	"dubai-auto/pkg/files"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetMyCars(ctx *fasthttp.RequestCtx, userID int, limit, lastID, nameColumn string) model.Response {
	lastIDInt, limitInt := utils.CheckLastIDLimit(lastID, limit, "")
	cars, err := s.UserRepository.GetMyCars(ctx, userID, limitInt, lastIDInt, 2, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: cars}
}

func (s *UserService) OnSale(ctx *fasthttp.RequestCtx, userID int, limit, lastID, nameColumn string) model.Response {
	lastIDInt, limitInt := utils.CheckLastIDLimit(lastID, limit, "")
	cars, err := s.UserRepository.GetMyCars(ctx, userID, limitInt, lastIDInt, 3, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: cars}
}

func (s *UserService) Cancel(ctx *fasthttp.RequestCtx, carID *int, dir string) model.Response {
	err := s.UserRepository.Cancel(ctx, carID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	files.RemoveFolder(dir)

	return model.Response{Data: model.Success{Message: "succesfully cancelled"}}
}

func (s *UserService) DeleteCar(ctx *fasthttp.RequestCtx, carID *int, dir string) model.Response {
	err := s.UserRepository.DeleteCar(ctx, carID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	files.RemoveFolder(dir)
	return model.Response{Data: model.Success{Message: "succesfully deleted"}}
}

func (s *UserService) DontSell(ctx *fasthttp.RequestCtx, carID, userID *int) model.Response {
	err := s.UserRepository.DontSell(ctx, carID, userID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "succesfully updated status"}}
}

func (s *UserService) Sell(ctx *fasthttp.RequestCtx, carID, userID *int) model.Response {
	err := s.UserRepository.Sell(ctx, carID, userID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "succesfully updated status"}}
}

func (s *UserService) GetBrands(ctx *fasthttp.RequestCtx, text, nameColumn string) model.Response {
	brands, err := s.UserRepository.GetBrands(ctx, text, nameColumn)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}
	return model.Response{
		Data: brands,
	}
}

func (s *UserService) GetProfile(ctx *fasthttp.RequestCtx, userID int, nameColumn string) model.Response {
	profile, err := s.UserRepository.GetProfile(ctx, userID, nameColumn)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}
	return model.Response{
		Data: profile,
	}
}

func (s *UserService) UpdateProfile(ctx *fasthttp.RequestCtx, userID int, profile *model.UpdateProfileRequest) model.Response {
	err := s.UserRepository.UpdateProfile(ctx, userID, profile)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}
	return model.Response{
		Data: model.Success{Message: "Profile updated successfully"},
	}
}

func (s *UserService) GetFilterBrands(ctx *fasthttp.RequestCtx, text, nameColumn string) model.Response {
	brands, err := s.UserRepository.GetFilterBrands(ctx, text, nameColumn)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}
	return model.Response{
		Data: brands,
	}
}

func (s *UserService) GetCities(ctx *fasthttp.RequestCtx, text, nameColumn string) model.Response {
	cities, err := s.UserRepository.GetCities(ctx, text, nameColumn)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}
	return model.Response{
		Data: cities,
	}
}

func (s *UserService) GetModelsByBrandID(ctx *fasthttp.RequestCtx, brandID int64, text, nameColumn string) model.Response {
	data, err := s.UserRepository.GetModelsByBrandID(ctx, brandID, text, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetFilterModelsByBrandID(ctx *fasthttp.RequestCtx, brandID int64, text, nameColumn string) model.Response {
	data, err := s.UserRepository.GetFilterModelsByBrandID(ctx, brandID, text, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetFilterModelsByBrands(ctx *fasthttp.RequestCtx, brands []int, text, nameColumn string) model.Response {
	data, err := s.UserRepository.GetFilterModelsByBrands(ctx, brands, text, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetGenerationsByModelID(ctx *fasthttp.RequestCtx, modelID int, wheel bool, year, bodyTypeID, nameColumn string) model.Response {
	data, err := s.UserRepository.GetGenerationsByModelID(ctx, modelID, wheel, year, bodyTypeID, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetGenerationsByModels(ctx *fasthttp.RequestCtx, models []int, nameColumn string) model.Response {
	data, err := s.UserRepository.GetGenerationsByModels(ctx, models, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetYearsByModelID(ctx *fasthttp.RequestCtx, modelID int64, wheel bool) model.Response {
	data, err := s.UserRepository.GetYearsByModelID(ctx, modelID, wheel)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetYearsByModels(ctx *fasthttp.RequestCtx, models []int, wheel bool) model.Response {
	data, err := s.UserRepository.GetYearsByModels(ctx, models, wheel)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetBodysByModelID(ctx *fasthttp.RequestCtx, modelID int, wheel bool, year string, nameColumn string) model.Response {
	data, err := s.UserRepository.GetBodysByModelID(ctx, modelID, wheel, year, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: 400}
	}
	return model.Response{Data: data}
}

// func (s *UserService) GetBodysByModels(ctx *fasthttp.RequestCtx, wheel bool, models, years []int) model.Response {
// 	data, err := s.UserRepository.GetBodysByModels(ctx, wheel, models, years)

// 	if err != nil {
// 		return model.Response{Error: err, Status: 400}
// 	}
// 	return model.Response{Data: data}
// }

func (s *UserService) GetBodyTypes(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetBodyTypes(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: data}
}

func (s *UserService) GetTransmissions(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetTransmissions(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetEngines(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetEngines(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetDrivetrains(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetDrivetrains(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetFuelTypes(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetFuelTypes(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetColors(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetColors(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetCountries(ctx *fasthttp.RequestCtx, nameColumn string) model.Response {
	data, err := s.UserRepository.GetCountries(ctx, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetHome(ctx *fasthttp.RequestCtx, userID int, nameColumn string) model.Response {
	data, err := s.UserRepository.GetHome(ctx, userID, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: data}
}

func (s *UserService) GetCars(ctx *fasthttp.RequestCtx, userID int, targetUserID string, brands, models, regions, cities,
	generations, transmissions, engines, drivetrains, body_types, fuel_types, ownership_types, colors []string,
	year_from, year_to, credit, price_from, price_to, tradeIn, owners, crash, odometer string,
	new, wheel *bool, limit, lastID int, nameColumn string) model.Response {

	cars, err := s.UserRepository.GetCars(ctx, userID, targetUserID, brands, models, regions, cities,
		generations, transmissions, engines, drivetrains, body_types, fuel_types,
		ownership_types, colors, year_from, year_to, credit,
		price_from, price_to, tradeIn, owners, crash, odometer, new, wheel, limit, lastID, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: cars}
}

func (s *UserService) GetPriceRecommendation(ctx *fasthttp.RequestCtx, filter model.GetPriceRecommendationRequest) model.Response {
	prices, err := s.UserRepository.GetPriceRecommendation(ctx, filter)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	if len(prices) == 0 {
		return model.Response{Error: errors.New("no prices found"), Status: http.StatusNotFound}
	}

	return model.Response{Data: model.GetPriceRecommendationResponse{MaxPrice: prices[0], MinPrice: prices[len(prices)-1], AvgPrice: prices[len(prices)/2]}}
}

func (s *UserService) GetCarByID(ctx *fasthttp.RequestCtx, carID, userID int, nameColumn string) model.Response {
	car, err := s.UserRepository.GetCarByID(ctx, carID, userID, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	return model.Response{Data: car}
}

func (s *UserService) GetEditCarByID(ctx *fasthttp.RequestCtx, carID, userID int, nameColumn string) model.Response {
	car, err := s.UserRepository.GetEditCarByID(ctx, carID, userID, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	return model.Response{Data: car}
}

func (s *UserService) BuyCar(ctx *fasthttp.RequestCtx, carID, userID int) model.Response {
	err := s.UserRepository.BuyCar(ctx, carID, userID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	return model.Response{Data: model.Success{Message: "successfully buy a car"}}
}

func (s *UserService) CreateCar(ctx *fasthttp.RequestCtx, car *model.CreateCarRequest, userID int) model.Response {

	id, err := s.UserRepository.CreateCar(ctx, car, userID)

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.SuccessWithId{Id: id, Message: "Car created successfully"},
	}
}

func (s *UserService) UpdateCar(ctx *fasthttp.RequestCtx, car *model.UpdateCarRequest, userID int) model.Response {
	err := s.UserRepository.UpdateCar(ctx, car, userID)

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Car updated successfully"},
	}
}

func (s *UserService) CarLike(ctx *fasthttp.RequestCtx, carID, userID *int) model.Response {
	err := s.UserRepository.CarLike(ctx, carID, userID)

	if err != nil {
		return model.Response{
			Status: 409,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Like created successfully"},
	}
}

func (s *UserService) RemoveLike(ctx *fasthttp.RequestCtx, carID, userID *int) model.Response {
	err := s.UserRepository.RemoveLike(ctx, carID, userID)

	if err != nil {
		return model.Response{
			Status: 409,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Like removed successfully"},
	}
}

func (s *UserService) Likes(ctx *fasthttp.RequestCtx, userID *int, nameColumn string) model.Response {
	data, err := s.UserRepository.Likes(ctx, userID, nameColumn)

	if err != nil {
		return model.Response{
			Status: 409,
			Error:  err,
		}
	}

	return model.Response{Data: data}
}

func (s *UserService) CreateCarImages(ctx *fasthttp.RequestCtx, carID int, images []string) model.Response {
	err := s.UserRepository.CreateCarImages(ctx, carID, images)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	// todo: delete old images if exist
	return model.Response{
		Data: model.Success{Message: "Car images created successfully"},
	}
}

func (s *UserService) CreateCarVideos(ctx *fasthttp.RequestCtx, carID int, video string) model.Response {
	err := s.UserRepository.CreateCarVideos(ctx, carID, video)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}
	// todo: delete old video if exist
	return model.Response{
		Data: model.Success{Message: "Car videos created successfully"},
	}
}

func (s *UserService) CreateMessageFile(ctx *fasthttp.RequestCtx, senderID int, filePath string) model.Response {
	err := s.UserRepository.CreateMessageFile(ctx, senderID, filePath)

	if err != nil {
		return model.Response{
			Status: 500,
			Error:  err,
		}
	}

	return model.Response{
		Data: model.Success{Message: config.ENV.IMAGE_BASE_URL + filePath},
	}
}

func (s *UserService) DeleteCarImage(ctx *fasthttp.RequestCtx, carID int, imagePath string) model.Response {
	err := s.UserRepository.DeleteCarImage(ctx, carID, imagePath)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Car image deleted successfully"}}
}

func (s *UserService) DeleteCarVideo(ctx *fasthttp.RequestCtx, carID int, videoPath string) model.Response {
	err := s.UserRepository.DeleteCarVideo(ctx, carID, videoPath)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Car video deleted successfully"}}
}

// GetBrokerByID returns a single broker by ID
func (s *UserService) GetUserByID(ctx *fasthttp.RequestCtx, userID string, nameColumn string) model.Response {

	userIDInt, err := strconv.Atoi(userID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	user, err := s.UserRepository.GetUserByRoleAndID(ctx, userIDInt, nameColumn)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	return model.Response{Data: user}
}

func (s *UserService) GetThirdPartyUsers(ctx *fasthttp.RequestCtx, roleID, fromID, toID, search string) model.Response {
	roleIDInt, _ := strconv.Atoi(roleID)
	fromIDInt, _ := strconv.Atoi(fromID)
	toIDInt, _ := strconv.Atoi(toID)
	user, err := s.UserRepository.GetThirdPartyUsers(ctx, roleIDInt, fromIDInt, toIDInt, search)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	return model.Response{Data: user}
}
