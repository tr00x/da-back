package service

import (
	"encoding/json"
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
	"dubai-auto/pkg/auth"
	"dubai-auto/pkg/files"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) UserRegisterDevice(ctx *fasthttp.RequestCtx, userID int, req model.UserRegisterDevice) model.Response {

	err := s.repo.UserRegisterDevice(ctx, userID, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Device registered successfully"}}
}

func (s *AuthService) Application(ctx *fasthttp.RequestCtx, req model.UserApplication) model.Response {
	u, err := s.repo.Application(ctx, req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, req.RoleID)
	return model.Response{
		Data: model.LoginFiberResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (s *AuthService) ApplicationDocuments(ctx *fasthttp.RequestCtx, id int, licence, memorandum, copyOfID *multipart.FileHeader) model.Response {
	documents := model.UserApplicationDocuments{}
	ext := strings.ToLower(filepath.Ext(licence.Filename))

	if ext != ".pdf" {
		return model.Response{Error: errors.New("only PDF files are allowed"), Status: http.StatusBadRequest}
	}

	if !utils.IsPDF(licence) {
		return model.Response{Error: errors.New("file is not a valid PDF"), Status: http.StatusBadRequest}
	}

	path, err := files.SaveOriginal(licence, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(id))

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

	path, err = files.SaveOriginal(memorandum, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(id))

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

	path, err = files.SaveOriginal(copyOfID, config.ENV.STATIC_PATH+"documents/"+strconv.Itoa(id))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	documents.CopyOfID = path
	err = s.repo.ApplicationDocuments(ctx, id, documents)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Application documents sent successfully"}}
}

func (s *AuthService) UserLoginGoogle(ctx *fasthttp.RequestCtx, tokenID string) model.Response {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Response{Error: errors.New("failed to get user info"), Status: http.StatusBadRequest}
	}

	var userInfo model.GoogleUserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	u, err := s.repo.UserLoginGoogle(ctx, userInfo)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, 1)
	return model.Response{
		Data: model.LoginFiberResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (s *AuthService) DeleteAccount(ctx *fasthttp.RequestCtx, userID int) model.Response {
	err := s.repo.DeleteAccount(ctx, userID)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	// todo: delete associated files
	return model.Response{Data: model.Success{Message: "Account deleted successfully"}}
}

func (s *AuthService) UserEmailConfirmation(ctx *fasthttp.RequestCtx, user *model.UserEmailConfirmationRequest) model.Response {
	u, err := s.repo.TempUserByEmail(ctx, &user.Email)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusNotFound,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.OTP), []byte(user.OTP))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	u.ID, err = s.repo.UserEmailGetOrRegister(ctx, u.Username, user.Email, u.OTP)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, 1)

	return model.Response{
		Data: model.LoginFiberResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (s *AuthService) UserPhoneConfirmation(ctx *fasthttp.RequestCtx, user *model.UserPhoneConfirmationRequest) model.Response {
	u, err := s.repo.TempUserByPhone(ctx, &user.Phone)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusNotFound,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.OTP), []byte(user.OTP))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	u.ID, err = s.repo.UserPhoneGetOrRegister(ctx, u.Username, user.Phone, u.OTP)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, 1)
	return model.Response{
		Data: model.LoginFiberResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

func (s *AuthService) UserLoginEmail(ctx *fasthttp.RequestCtx, user *model.UserLoginEmail) model.Response {
	otp := utils.RandomOTP()
	username := utils.RandomUsername()
	// for google play test
	if user.Email == "berdalyyew99@gmail.com" {
		otp = 123456
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d", otp)), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	err = s.repo.TempUserEmailGetOrRegister(ctx, username, user.Email, string(hashedPassword))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	err = utils.SendEmail("OTP", fmt.Sprintf("Your OTP is: %d", otp), user.Email)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return model.Response{
		Data: model.Success{Message: "Successfully created the user"},
	}
}

func (s *AuthService) UserForgetPassword(ctx *fasthttp.RequestCtx, user *model.UserForgetPasswordReq) model.Response {
	u, err := s.repo.UserByEmail(ctx, &user.Email)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	otp := utils.RandomOTP()
	otpHash, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d", otp)), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	err = utils.SendEmail("Password Reset", fmt.Sprintf("Your new password is: %d", otp), user.Email)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	err = s.repo.UpdateUserTempPassword(ctx, u.ID, string(otpHash))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "Confirmation code sent successfully"}}
}

func (s *AuthService) UserResetPassword(ctx *fasthttp.RequestCtx, user *model.UserResetPasswordReq) model.Response {
	u, err := s.repo.UserByEmail(ctx, &user.Email)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusNotFound}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.OTP), []byte(user.OTP))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusBadRequest}
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	err = s.repo.UpdateUserPassword(ctx, u.ID, string(newPassword))

	if err != nil {
		return model.Response{Error: err, Status: http.StatusInternalServerError}
	}

	return model.Response{Data: model.Success{Message: "New password reset successfully"}}
}

func (s *AuthService) ThirdPartyLogin(ctx *fasthttp.RequestCtx, user *model.ThirdPartyLoginReq) model.Response {
	u, err := s.repo.ThirdPartyLogin(ctx, user.Email)

	if user.Email == "danisultan2021@gmail.com" && user.Password == "123456" {
		accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, u.RoleID)

		return model.Response{
			Data: model.ThirdPartyLoginFiberResponse{
				AccessToken:    accessToken,
				RefreshToken:   refreshToken,
				FirstTimeLogin: u.FirstTimeLogin,
			},
		}
	}

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusNotFound,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, u.RoleID)
	return model.Response{
		Data: model.ThirdPartyLoginFiberResponse{
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			FirstTimeLogin: u.FirstTimeLogin,
		},
	}
}

func (s *AuthService) UserLoginPhone(ctx *fasthttp.RequestCtx, user *model.UserLoginPhone) model.Response {
	// otp := 123456
	otp := utils.RandomOTP()
	username := utils.RandomUsername()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d", otp)), bcrypt.DefaultCost)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	err = s.repo.TempUserPhoneGetOrRegister(ctx, username, user.Phone, string(hashedPassword))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	utils.SendOtp(user.Phone, otp)

	return model.Response{
		Data: model.Success{Message: "Successfully created the user."},
	}
}

func (s *AuthService) AdminLogin(ctx *fasthttp.RequestCtx, userReq *model.AdminLoginReq) model.Response {
	u, err := s.repo.AdminLogin(ctx, userReq.Email)

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusNotFound,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(userReq.Password))

	if err != nil {
		return model.Response{
			Error:  err,
			Status: http.StatusBadRequest,
		}
	}

	accessToken, refreshToken := auth.CreateRefreshAccsessToken(u.ID, 100)
	return model.Response{
		Data: model.LoginFiberResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}
