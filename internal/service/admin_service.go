package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/utils"
	"dubai-auto/pkg/files"
)

type AdminService struct {
	repo *repository.AdminRepository
}

func NewAdminService(repo *repository.AdminRepository) *AdminService {
	return &AdminService{repo}
}

// Admin service methods
func (s *AdminService) CreateAdmin(ctx *fasthttp.RequestCtx, req *model.CreateAdminRequest) model.Response {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	req.Password = string(hashedPassword)
	userID, err := s.repo.CreateAdmin(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: userID, Message: "Admin created successfully"}}
}

func (s *AdminService) GetAdmins(ctx *fasthttp.RequestCtx) model.Response {
	admins, err := s.repo.GetAdmins(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: admins}
}

func (s *AdminService) GetAdmin(ctx *fasthttp.RequestCtx, id int) model.Response {
	admin, err := s.repo.GetAdmin(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}
	return model.Response{Data: admin}
}

func (s *AdminService) UpdateAdmin(ctx *fasthttp.RequestCtx, id int, req *model.UpdateAdminRequest) model.Response {
	// Hash password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		if err != nil {
			return model.Response{Error: err, Status: http.StatusInternalServerError}
		}
		req.Password = string(hashedPassword)
	}

	err := s.repo.UpdateAdmin(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Admin updated successfully"}}
}

func (s *AdminService) DeleteAdmin(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteAdmin(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete admin files if exists
	return model.Response{Data: model.Success{Message: "Admin deleted successfully"}}
}

// Profile service methods
func (s *AdminService) GetProfile(ctx *fasthttp.RequestCtx, id int) model.Response {
	profile, err := s.repo.GetProfile(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: profile}
}

// Application service methods
func (s *AdminService) GetApplications(ctx *fasthttp.RequestCtx, qRole, qStatus, limit, lastID, search string) model.Response {
	lastIDInt, limitInt := utils.CheckLastIDLimit(lastID, limit, "")
	qRoleInt, err := strconv.Atoi(qRole)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	qStatusInt, err := strconv.Atoi(qStatus)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	if (qRoleInt > model.ROLE_COUNT || qRoleInt < 1) || (qStatusInt > model.APPLICATION_STATUS_COUNT || qStatusInt < 1) {
		return model.Response{Error: errors.New("invalid role or status"), Status: http.StatusBadRequest}
	}

	applications, err := s.repo.GetApplications(ctx, qRoleInt, qStatusInt, limitInt, lastIDInt, search)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: applications}
}

func (s *AdminService) CreateApplication(ctx *fasthttp.RequestCtx, req model.UserApplication) model.Response {

	if req.Password == "" {
		req.Password = fmt.Sprintf("%d", utils.RandomOTP())
	}

	err := utils.SendEmail("Password", fmt.Sprintf("Your password is: %s", req.Password), req.Email)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	id, err := s.repo.CreateApplication(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Application created successfully"}}
}

func (s *AdminService) GetApplication(ctx *fasthttp.RequestCtx, idStr string, qStatus string) model.Response {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	qStatusInt, err := strconv.Atoi(qStatus)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	if qStatusInt > model.APPLICATION_STATUS_COUNT || qStatusInt < 1 {
		return model.Response{Error: errors.New("invalid status"), Status: http.StatusBadRequest}
	}

	application, err := s.repo.GetApplication(ctx, id, qStatusInt)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: application}
}

func (s *AdminService) CreateApplicationDocuments(ctx *fasthttp.RequestCtx, userID int, licence, memorandum, copyOfID *multipart.FileHeader) model.Response {
	documents := model.UserApplicationDocuments{}
	ext := strings.ToLower(filepath.Ext(licence.Filename))

	if ext != ".pdf" {
		return model.Response{Error: errors.New("only PDF files are allowed"), Status: http.StatusBadRequest}
	}

	if !utils.IsPDF(licence) {
		return model.Response{Error: errors.New("file is not a valid PDF"), Status: http.StatusBadRequest}
	}

	path, err := files.SaveOriginal(licence, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(userID))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	documents.Licence = path
	ext = strings.ToLower(filepath.Ext(memorandum.Filename))

	if ext != ".pdf" {
		return model.Response{Error: errors.New("only PDF files are allowed"), Status: http.StatusBadRequest}
	}

	if !utils.IsPDF(memorandum) {
		return model.Response{Error: errors.New("file is not a valid PDF"), Status: http.StatusBadRequest}
	}

	path, err = files.SaveOriginal(memorandum, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(userID))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	documents.Memorandum = path
	ext = strings.ToLower(filepath.Ext(copyOfID.Filename))

	if ext != ".pdf" {
		return model.Response{Error: errors.New("only PDF files are allowed"), Status: http.StatusBadRequest}
	}

	if !utils.IsPDF(copyOfID) {
		return model.Response{Error: errors.New("file is not a valid PDF"), Status: http.StatusBadRequest}
	}

	path, err = files.SaveOriginal(copyOfID, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(userID))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	documents.CopyOfID = path
	err = s.repo.CreateApplicationDocuments(ctx, userID, documents)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Application documents sent successfully"}}
}

func (s *AdminService) AcceptApplication(ctx *fasthttp.RequestCtx, id int, req model.AcceptApplicationRequest) model.Response {

	if req.Password == "" {
		req.Password = fmt.Sprintf("%d", utils.RandomOTP())
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	email, err := s.repo.AcceptApplication(ctx, id, string(cryptedPassword))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	go func() {
		err = utils.SendEmail("Application accepted", "Your application has been accepted. Please login to your account to continue. Your password is: "+req.Password, email)

		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return model.Response{Data: model.Success{Message: "Application accepted successfully"}}
}

func (s *AdminService) RejectApplication(ctx *fasthttp.RequestCtx, idStr string, qStatus string, qMessage string) model.Response {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	qStatusInt, err := strconv.Atoi(qStatus)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	// todo: remove files/folders after reject
	// todo: send email to user
	email, err := s.repo.RejectApplication(ctx, id, qStatusInt)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	go utils.SendEmail("Application rejected", "Your application has been rejected. Reason: "+qMessage, email)
	return model.Response{Data: model.Success{Message: "Application rejected successfully"}}
}

// Cities service methods
func (s *AdminService) GetCities(ctx *fasthttp.RequestCtx) model.Response {
	cities, err := s.repo.GetCities(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: cities}
}

func (s *AdminService) CreateCity(ctx *fasthttp.RequestCtx, req *model.CreateNameRequest) model.Response {
	id, err := s.repo.CreateCity(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "City created successfully"}}
}

func (s *AdminService) UpdateCity(ctx *fasthttp.RequestCtx, id int, req *model.CreateNameRequest) model.Response {
	err := s.repo.UpdateCity(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "City updated successfully"}}
}

func (s *AdminService) DeleteCity(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteCity(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "City deleted successfully"}}
}

// Company Types service methods
func (s *AdminService) GetCompanyTypes(ctx *fasthttp.RequestCtx) model.Response {
	items, err := s.repo.GetCompanyTypes(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: items}
}

func (s *AdminService) GetCompanyType(ctx *fasthttp.RequestCtx, id int) model.Response {
	item, err := s.repo.GetCompanyType(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: item}
}

func (s *AdminService) CreateCompanyType(ctx *fasthttp.RequestCtx, req *model.CreateCompanyTypeRequest) model.Response {
	id, err := s.repo.CreateCompanyType(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Company type created successfully"}}
}

func (s *AdminService) UpdateCompanyType(ctx *fasthttp.RequestCtx, id int, req *model.CreateCompanyTypeRequest) model.Response {
	err := s.repo.UpdateCompanyType(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Company type updated successfully"}}
}

func (s *AdminService) DeleteCompanyType(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteCompanyType(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Company type deleted successfully"}}
}

// Activity Fields service methods
func (s *AdminService) GetActivityFields(ctx *fasthttp.RequestCtx) model.Response {
	items, err := s.repo.GetActivityFields(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: items}
}

func (s *AdminService) GetActivityField(ctx *fasthttp.RequestCtx, id int) model.Response {
	item, err := s.repo.GetActivityField(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: item}
}

func (s *AdminService) CreateActivityField(ctx *fasthttp.RequestCtx, req *model.CreateCompanyTypeRequest) model.Response {
	id, err := s.repo.CreateActivityField(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Activity field created successfully"}}
}

func (s *AdminService) UpdateActivityField(ctx *fasthttp.RequestCtx, id int, req *model.CreateCompanyTypeRequest) model.Response {
	err := s.repo.UpdateActivityField(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Activity field updated successfully"}}
}

func (s *AdminService) DeleteActivityField(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteActivityField(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Activity field deleted successfully"}}
}

// Brands service methods
func (s *AdminService) GetBrands(ctx *fasthttp.RequestCtx) model.Response {
	brands, err := s.repo.GetBrands(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: brands}
}

func (s *AdminService) CreateBrand(ctx *fasthttp.RequestCtx, req *model.CreateBrandRequest) model.Response {
	id, err := s.repo.CreateBrand(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Brand created successfully"}}
}

func (s *AdminService) CreateBrandImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	path, err := files.SaveOriginal(image[0], config.ENV.STATIC_PATH+"logos/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	err = s.repo.CreateBrandImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Brand image created successfully"}}
}

func (s *AdminService) UpdateBrand(ctx *fasthttp.RequestCtx, id int, req *model.CreateBrandRequest) model.Response {
	err := s.repo.UpdateBrand(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Brand updated successfully"}}
}

func (s *AdminService) DeleteBrand(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteBrand(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Brand deleted successfully"}}
}

// Models service methods
func (s *AdminService) GetModels(ctx *fasthttp.RequestCtx, brand_id int) model.Response {
	models, err := s.repo.GetModels(ctx, brand_id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: models}
}

func (s *AdminService) CreateModel(ctx *fasthttp.RequestCtx, brand_id int, req *model.CreateModelRequest) model.Response {
	id, err := s.repo.CreateModel(ctx, brand_id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Model created successfully"}}
}

func (s *AdminService) UpdateModel(ctx *fasthttp.RequestCtx, id int, req *model.UpdateModelRequest) model.Response {
	err := s.repo.UpdateModel(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Model updated successfully"}}
}

func (s *AdminService) DeleteModel(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteModel(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Model deleted successfully"}}
}

// Body Types service methods
func (s *AdminService) GetBodyTypes(ctx *fasthttp.RequestCtx) model.Response {
	bodyTypes, err := s.repo.GetBodyTypes(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: bodyTypes}
}

func (s *AdminService) CreateBodyType(ctx *fasthttp.RequestCtx, req *model.CreateBodyTypeRequest) model.Response {
	id, err := s.repo.CreateBodyType(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Body type created successfully"}}
}

func (s *AdminService) CreateBodyTypeImage(ctx *fasthttp.RequestCtx, id int, path string) model.Response {
	err := s.repo.CreateBodyTypeImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Body type image created successfully"}}
}

func (s *AdminService) UpdateBodyType(ctx *fasthttp.RequestCtx, id int, req *model.CreateBodyTypeRequest) model.Response {
	err := s.repo.UpdateBodyType(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Body type updated successfully"}}
}

func (s *AdminService) DeleteBodyType(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteBodyType(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Body type deleted successfully"}}
}

func (s *AdminService) DeleteBodyTypeImage(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteBodyTypeImage(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Body type image deleted successfully"}}
}

// Transmissions service methods
func (s *AdminService) GetTransmissions(ctx *fasthttp.RequestCtx) model.Response {
	transmissions, err := s.repo.GetTransmissions(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: transmissions}
}

func (s *AdminService) CreateTransmission(ctx *fasthttp.RequestCtx, req *model.CreateTransmissionRequest) model.Response {
	id, err := s.repo.CreateTransmission(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Transmission created successfully"}}
}

func (s *AdminService) UpdateTransmission(ctx *fasthttp.RequestCtx, id int, req *model.CreateTransmissionRequest) model.Response {
	err := s.repo.UpdateTransmission(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Transmission updated successfully"}}
}

func (s *AdminService) DeleteTransmission(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteTransmission(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Transmission deleted successfully"}}
}

// Engines service methods
func (s *AdminService) GetEngines(ctx *fasthttp.RequestCtx) model.Response {
	engines, err := s.repo.GetEngines(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: engines}
}

func (s *AdminService) CreateEngine(ctx *fasthttp.RequestCtx, req *model.CreateEngineRequest) model.Response {
	id, err := s.repo.CreateEngine(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Engine created successfully"}}
}

func (s *AdminService) UpdateEngine(ctx *fasthttp.RequestCtx, id int, req *model.CreateEngineRequest) model.Response {
	err := s.repo.UpdateEngine(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Engine updated successfully"}}
}

func (s *AdminService) DeleteEngine(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteEngine(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Engine deleted successfully"}}
}

// Regions service methods
func (s *AdminService) GetRegions(ctx *fasthttp.RequestCtx, city_id int) model.Response {
	regions, err := s.repo.GetRegions(ctx, city_id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: regions}
}

func (s *AdminService) CreateRegion(ctx *fasthttp.RequestCtx, city_id int, req *model.CreateNameRequest) model.Response {
	id, err := s.repo.CreateRegion(ctx, city_id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Region created successfully"}}
}

func (s *AdminService) UpdateRegion(ctx *fasthttp.RequestCtx, id int, req *model.CreateNameRequest) model.Response {
	err := s.repo.UpdateRegion(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Region updated successfully"}}
}

func (s *AdminService) DeleteRegion(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteRegion(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Region deleted successfully"}}
}

// Drivetrains service methods
func (s *AdminService) GetDrivetrains(ctx *fasthttp.RequestCtx) model.Response {
	drivetrains, err := s.repo.GetDrivetrains(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: drivetrains}
}

func (s *AdminService) CreateDrivetrain(ctx *fasthttp.RequestCtx, req *model.CreateDrivetrainRequest) model.Response {
	id, err := s.repo.CreateDrivetrain(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Drivetrain created successfully"}}
}

func (s *AdminService) UpdateDrivetrain(ctx *fasthttp.RequestCtx, id int, req *model.CreateDrivetrainRequest) model.Response {
	err := s.repo.UpdateDrivetrain(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Drivetrain updated successfully"}}
}

func (s *AdminService) DeleteDrivetrain(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteDrivetrain(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Drivetrain deleted successfully"}}
}

// Fuel Types service methods
func (s *AdminService) GetFuelTypes(ctx *fasthttp.RequestCtx) model.Response {
	fuelTypes, err := s.repo.GetFuelTypes(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: fuelTypes}
}

func (s *AdminService) CreateFuelType(ctx *fasthttp.RequestCtx, req *model.CreateFuelTypeRequest) model.Response {
	id, err := s.repo.CreateFuelType(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Fuel type created successfully"}}
}

func (s *AdminService) UpdateFuelType(ctx *fasthttp.RequestCtx, id int, req *model.CreateFuelTypeRequest) model.Response {
	err := s.repo.UpdateFuelType(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Fuel type updated successfully"}}
}

func (s *AdminService) DeleteFuelType(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteFuelType(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Fuel type deleted successfully"}}
}

// Generations service methods
func (s *AdminService) GetGenerations(ctx *fasthttp.RequestCtx) model.Response {
	generations, err := s.repo.GetGenerations(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: generations}
}

func (s *AdminService) GetGenerationsByModel(ctx *fasthttp.RequestCtx, brandId, modelId int) model.Response {
	// First validate that the model belongs to the specified brand
	isValid, err := s.repo.ValidateModelBelongsToBrand(ctx, modelId, brandId)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	if !isValid {
		return model.Response{
			Error:  errors.New("model does not belong to the specified brand"),
			Status: http.StatusBadRequest,
		}
	}

	generations, err := s.repo.GetGenerationsByModel(ctx, modelId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: generations}
}

func (s *AdminService) CreateGeneration(ctx *fasthttp.RequestCtx, req *model.CreateGenerationRequest) model.Response {
	id, err := s.repo.CreateGeneration(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Generation created successfully"}}
}

func (s *AdminService) UpdateGeneration(ctx *fasthttp.RequestCtx, id int, req *model.UpdateGenerationRequest) model.Response {
	err := s.repo.UpdateGeneration(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Generation updated successfully"}}
}

func (s *AdminService) DeleteGeneration(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteGeneration(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Generation deleted successfully"}}
}

func (s *AdminService) CreateGenerationImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	paths, status, err := files.SaveFiles(image, config.ENV.STATIC_PATH+"cars/generations/"+strconv.Itoa(id), config.ENV.DEFAULT_IMAGE_WIDTHS)

	if err != nil {
		return model.Response{
			Status: status,
			Error:  err,
		}
	}

	// todo: delete old image if exists
	err = s.repo.CreateGenerationImage(ctx, id, paths)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Generation image created successfully"}}
}

// Generation Modifications service methods
func (s *AdminService) GetGenerationModifications(ctx *fasthttp.RequestCtx, generationId int) model.Response {
	generationModifications, err := s.repo.GetGenerationModifications(ctx, generationId)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: generationModifications}
}

func (s *AdminService) CreateGenerationModification(ctx *fasthttp.RequestCtx, generationId int, req *model.CreateGenerationModificationRequest) model.Response {
	id, err := s.repo.CreateGenerationModification(ctx, generationId, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Generation modification created successfully"}}
}

func (s *AdminService) UpdateGenerationModification(ctx *fasthttp.RequestCtx, generationId int, id int, req *model.UpdateGenerationModificationRequest) model.Response {
	err := s.repo.UpdateGenerationModification(ctx, generationId, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Generation modification updated successfully"}}
}

func (s *AdminService) DeleteGenerationModification(ctx *fasthttp.RequestCtx, generationId int, id int) model.Response {
	err := s.repo.DeleteGenerationModification(ctx, generationId, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Generation modification deleted successfully"}}
}

// Colors service methods
func (s *AdminService) GetColors(ctx *fasthttp.RequestCtx) model.Response {
	colors, err := s.repo.GetColors(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: colors}
}

func (s *AdminService) CreateColor(ctx *fasthttp.RequestCtx, req *model.CreateColorRequest) model.Response {
	id, err := s.repo.CreateColor(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Color created successfully"}}
}

func (s *AdminService) CreateColorImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	path, err := files.SaveOriginal(image[0], config.ENV.STATIC_PATH+"colors/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	err = s.repo.CreateColorImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Color image created successfully"}}
}

func (s *AdminService) UpdateColor(ctx *fasthttp.RequestCtx, id int, req *model.UpdateColorRequest) model.Response {
	err := s.repo.UpdateColor(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Color updated successfully"}}
}

func (s *AdminService) DeleteColor(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteColor(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete image if exists
	return model.Response{Data: model.Success{Message: "Color deleted successfully"}}
}

// Moto Categories service methods
func (s *AdminService) GetMotoCategories(ctx *fasthttp.RequestCtx) model.Response {
	motoCategories, err := s.repo.GetMotoCategories(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoCategories}
}

func (s *AdminService) CreateMotoCategory(ctx *fasthttp.RequestCtx, req *model.CreateMotoCategoryRequest) model.Response {
	id, err := s.repo.CreateMotoCategory(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto category created successfully"}}
}

func (s *AdminService) GetMotoBrandsByCategoryID(ctx *fasthttp.RequestCtx, id int) model.Response {
	motoBrands, err := s.repo.GetMotoBrandsByCategoryID(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoBrands}
}

func (s *AdminService) UpdateMotoCategory(ctx *fasthttp.RequestCtx, id int, req *model.UpdateMotoCategoryRequest) model.Response {
	err := s.repo.UpdateMotoCategory(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto category updated successfully"}}
}

func (s *AdminService) DeleteMotoCategory(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteMotoCategory(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto category deleted successfully"}}
}

// Moto Brands service methods
func (s *AdminService) GetMotoBrands(ctx *fasthttp.RequestCtx) model.Response {
	motoBrands, err := s.repo.GetMotoBrands(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoBrands}
}

func (s *AdminService) GetMotoModelsByBrandID(ctx *fasthttp.RequestCtx, id int) model.Response {
	motoModels, err := s.repo.GetMotoModelsByBrandID(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoModels}
}

func (s *AdminService) CreateMotoBrand(ctx *fasthttp.RequestCtx, req *model.CreateMotoBrandRequest) model.Response {
	id, err := s.repo.CreateMotoBrand(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto brand created successfully"}}
}

func (s *AdminService) CreateMotoBrandImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	path, err := files.SaveOriginal(image[0], config.ENV.STATIC_PATH+"logos/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	err = s.repo.CreateMotoBrandImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto brand image created successfully"}}
}

func (s *AdminService) UpdateMotoBrand(ctx *fasthttp.RequestCtx, id int, req *model.UpdateMotoBrandRequest) model.Response {
	err := s.repo.UpdateMotoBrand(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto brand updated successfully"}}
}

func (s *AdminService) DeleteMotoBrand(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteMotoBrand(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Moto brand deleted successfully"}}
}

// Moto Models service methods
func (s *AdminService) GetMotoModels(ctx *fasthttp.RequestCtx) model.Response {
	motoModels, err := s.repo.GetMotoModels(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoModels}
}

func (s *AdminService) CreateMotoModel(ctx *fasthttp.RequestCtx, req *model.CreateMotoModelRequest) model.Response {
	id, err := s.repo.CreateMotoModel(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto model created successfully"}}
}

func (s *AdminService) UpdateMotoModel(ctx *fasthttp.RequestCtx, id int, req *model.UpdateMotoModelRequest) model.Response {
	err := s.repo.UpdateMotoModel(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto model updated successfully"}}
}

func (s *AdminService) DeleteMotoModel(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteMotoModel(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Moto model deleted successfully"}}
}

// Moto Parameters service methods
func (s *AdminService) GetMotoParameters(ctx *fasthttp.RequestCtx) model.Response {
	motoParameters, err := s.repo.GetMotoParameters(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoParameters}
}

func (s *AdminService) CreateMotoParameter(ctx *fasthttp.RequestCtx, req *model.CreateMotoParameterRequest) model.Response {
	id, err := s.repo.CreateMotoParameter(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto parameter created successfully"}}
}

func (s *AdminService) UpdateMotoParameter(ctx *fasthttp.RequestCtx, id int, req *model.UpdateMotoParameterRequest) model.Response {
	err := s.repo.UpdateMotoParameter(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto parameter updated successfully"}}
}

func (s *AdminService) DeleteMotoParameter(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteMotoParameter(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto parameter deleted successfully"}}
}

// Moto Parameter Values service methods
func (s *AdminService) GetMotoParameterValues(ctx *fasthttp.RequestCtx, motoParamId int) model.Response {
	motoParameterValues, err := s.repo.GetMotoParameterValues(ctx, motoParamId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoParameterValues}
}

func (s *AdminService) CreateMotoParameterValue(ctx *fasthttp.RequestCtx, motoParamId int, req *model.CreateMotoParameterValueRequest) model.Response {
	id, err := s.repo.CreateMotoParameterValue(ctx, motoParamId, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Moto parameter value created successfully"}}
}

func (s *AdminService) UpdateMotoParameterValue(ctx *fasthttp.RequestCtx, motoParamId int, id int, req *model.UpdateMotoParameterValueRequest) model.Response {
	err := s.repo.UpdateMotoParameterValue(ctx, motoParamId, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto parameter value updated successfully"}}
}

func (s *AdminService) DeleteMotoParameterValue(ctx *fasthttp.RequestCtx, motoParamId int, id int) model.Response {
	err := s.repo.DeleteMotoParameterValue(ctx, motoParamId, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto parameter value deleted successfully"}}
}

// Moto Category Parameters service methods
func (s *AdminService) GetMotoCategoryParameters(ctx *fasthttp.RequestCtx, categoryId int) model.Response {
	motoCategoryParameters, err := s.repo.GetMotoCategoryParameters(ctx, categoryId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: motoCategoryParameters}
}

func (s *AdminService) CreateMotoCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, req *model.CreateMotoCategoryParameterRequest) model.Response {
	parameterId, err := s.repo.CreateMotoCategoryParameter(ctx, categoryId, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: parameterId, Message: "Moto category parameter created successfully"}}
}

func (s *AdminService) UpdateMotoCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, parameterId int, req *model.UpdateMotoCategoryParameterRequest) model.Response {
	err := s.repo.UpdateMotoCategoryParameter(ctx, categoryId, parameterId, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto category parameter updated successfully"}}
}

func (s *AdminService) DeleteMotoCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, parameterId int) model.Response {
	err := s.repo.DeleteMotoCategoryParameter(ctx, categoryId, parameterId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Moto category parameter deleted successfully"}}
}

// Comtrans Categories service methods
func (s *AdminService) GetComtransCategories(ctx *fasthttp.RequestCtx) model.Response {
	comtransCategories, err := s.repo.GetComtransCategories(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransCategories}
}

func (s *AdminService) GetComtransBrandsByCategoryID(ctx *fasthttp.RequestCtx, categoryId int) model.Response {
	comtransBrands, err := s.repo.GetComtransBrandsByCategoryID(ctx, categoryId)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransBrands}
}

func (s *AdminService) CreateComtransCategory(ctx *fasthttp.RequestCtx, req *model.CreateComtransCategoryRequest) model.Response {
	id, err := s.repo.CreateComtransCategory(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans category created successfully"}}
}

func (s *AdminService) UpdateComtransCategory(ctx *fasthttp.RequestCtx, id int, req *model.UpdateComtransCategoryRequest) model.Response {
	err := s.repo.UpdateComtransCategory(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans category updated successfully"}}
}

func (s *AdminService) DeleteComtransCategory(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteComtransCategory(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans category deleted successfully"}}
}

// Comtrans Brands service methods
func (s *AdminService) GetComtransBrands(ctx *fasthttp.RequestCtx) model.Response {
	comtransBrands, err := s.repo.GetComtransBrands(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransBrands}
}

func (s *AdminService) GetComtransModelsByBrandID(ctx *fasthttp.RequestCtx, id int) model.Response {
	comtransModels, err := s.repo.GetComtransModelsByBrandID(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransModels}
}

func (s *AdminService) CreateComtransBrand(ctx *fasthttp.RequestCtx, req *model.CreateComtransBrandRequest) model.Response {
	id, err := s.repo.CreateComtransBrand(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans brand created successfully"}}
}

func (s *AdminService) CreateComtransBrandImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	path, err := files.SaveOriginal(image[0], config.ENV.STATIC_PATH+"logos/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	err = s.repo.CreateComtransBrandImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans brand image created successfully"}}
}

func (s *AdminService) UpdateComtransBrand(ctx *fasthttp.RequestCtx, id int, req *model.UpdateComtransBrandRequest) model.Response {
	err := s.repo.UpdateComtransBrand(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans brand updated successfully"}}
}

func (s *AdminService) DeleteComtransBrand(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteComtransBrand(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	// todo: delete image files if exists
	return model.Response{Data: model.Success{Message: "Comtrans brand deleted successfully"}}
}

// Comtrans Models service methods
func (s *AdminService) GetComtransModels(ctx *fasthttp.RequestCtx) model.Response {
	comtransModels, err := s.repo.GetComtransModels(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransModels}
}

func (s *AdminService) CreateComtransModel(ctx *fasthttp.RequestCtx, req *model.CreateComtransModelRequest) model.Response {
	id, err := s.repo.CreateComtransModel(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans model created successfully"}}
}

func (s *AdminService) UpdateComtransModel(ctx *fasthttp.RequestCtx, id int, req *model.UpdateComtransModelRequest) model.Response {
	err := s.repo.UpdateComtransModel(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans model updated successfully"}}
}

func (s *AdminService) DeleteComtransModel(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteComtransModel(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Comtrans model deleted successfully"}}
}

// Comtrans Parameters service methods
func (s *AdminService) GetComtransParameters(ctx *fasthttp.RequestCtx) model.Response {
	comtransParameters, err := s.repo.GetComtransParameters(ctx)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransParameters}
}

func (s *AdminService) CreateComtransParameter(ctx *fasthttp.RequestCtx, req *model.CreateComtransParameterRequest) model.Response {
	id, err := s.repo.CreateComtransParameter(ctx, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans parameter created successfully"}}
}

func (s *AdminService) UpdateComtransParameter(ctx *fasthttp.RequestCtx, id int, req *model.UpdateComtransParameterRequest) model.Response {
	err := s.repo.UpdateComtransParameter(ctx, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans parameter updated successfully"}}
}

func (s *AdminService) DeleteComtransParameter(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteComtransParameter(ctx, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans parameter deleted successfully"}}
}

// Comtrans Parameter Values service methods
func (s *AdminService) GetComtransParameterValues(ctx *fasthttp.RequestCtx, parameterId int) model.Response {
	comtransParameterValues, err := s.repo.GetComtransParameterValues(ctx, parameterId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransParameterValues}
}

func (s *AdminService) CreateComtransParameterValue(ctx *fasthttp.RequestCtx, parameterId int, req *model.CreateComtransParameterValueRequest) model.Response {
	id, err := s.repo.CreateComtransParameterValue(ctx, parameterId, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans parameter value created successfully"}}
}

func (s *AdminService) UpdateComtransParameterValue(ctx *fasthttp.RequestCtx, parameterId int, id int, req *model.UpdateComtransParameterValueRequest) model.Response {
	err := s.repo.UpdateComtransParameterValue(ctx, parameterId, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans parameter value updated successfully"}}
}

func (s *AdminService) DeleteComtransParameterValue(ctx *fasthttp.RequestCtx, parameterId int, id int) model.Response {
	err := s.repo.DeleteComtransParameterValue(ctx, parameterId, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans parameter value deleted successfully"}}
}

// Comtrans Category Parameters service methods
func (s *AdminService) GetComtransCategoryParameters(ctx *fasthttp.RequestCtx, categoryId int) model.Response {
	comtransCategoryParameters, err := s.repo.GetComtransCategoryParameters(ctx, categoryId)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: comtransCategoryParameters}
}

func (s *AdminService) CreateComtransCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, req *model.CreateComtransCategoryParameterRequest) model.Response {
	id, err := s.repo.CreateComtransCategoryParameter(ctx, categoryId, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Comtrans category parameter created successfully"}}
}

func (s *AdminService) UpdateComtransCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, id int, req *model.UpdateComtransCategoryParameterRequest) model.Response {
	err := s.repo.UpdateComtransCategoryParameter(ctx, categoryId, id, req)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans category parameter updated successfully"}}
}

func (s *AdminService) DeleteComtransCategoryParameter(ctx *fasthttp.RequestCtx, categoryId int, id int) model.Response {
	err := s.repo.DeleteComtransCategoryParameter(ctx, categoryId, id)
	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}
	return model.Response{Data: model.Success{Message: "Comtrans category parameter deleted successfully"}}
}

// Countries service methods
func (s *AdminService) GetCountries(ctx *fasthttp.RequestCtx) model.Response {
	countries, err := s.repo.GetCountries(ctx)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: countries}
}

func (s *AdminService) CreateCountry(ctx *fasthttp.RequestCtx, req *model.CreateNameRequest) model.Response {
	id, err := s.repo.CreateCountry(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Country created successfully"}}
}

func (s *AdminService) CreateCountryImage(ctx *fasthttp.RequestCtx, form *multipart.Form, id int) model.Response {

	if form == nil {
		return model.Response{
			Status: 400,
			Error:  errors.New("didn't upload the files"),
		}
	}

	image := form.File["image"]

	if len(image) > 1 {
		return model.Response{
			Status: 400,
			Error:  errors.New("must load maximum 1 file"),
		}
	}

	path, err := files.SaveOriginal(image[0], config.ENV.STATIC_PATH+"countries/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{
			Status: 400,
			Error:  err,
		}
	}

	err = s.repo.CreateCountryImage(ctx, id, path)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.SuccessWithId{Id: id, Message: "Country flag image created successfully"}}
}

func (s *AdminService) UpdateCountry(ctx *fasthttp.RequestCtx, id int, req *model.CreateNameRequest) model.Response {
	err := s.repo.UpdateCountry(ctx, id, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Country updated successfully"}}
}

func (s *AdminService) DeleteCountry(ctx *fasthttp.RequestCtx, id int) model.Response {
	err := s.repo.DeleteCountry(ctx, id)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Country deleted successfully"}}
}
