package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"dubai-auto/internal/config"
	"dubai-auto/internal/delivery/http"
	"dubai-auto/internal/repository"
	"dubai-auto/internal/service"
	"dubai-auto/pkg/auth"
)

// TODO: Hemmesini interface cykar!

func SetupAdminRoutes(r fiber.Router, config *config.Config, db *pgxpool.Pool, validator *auth.Validator) {
	adminRepository := repository.NewAdminRepository(config, db)
	adminService := service.NewAdminService(adminRepository)
	adminHandler := http.NewAdminHandler(adminService, validator)

	// profile routes
	profile := r.Group("/profile")
	{
		profile.Get("/", adminHandler.GetProfile)
	}

	// admin users CRUD
	users := r.Group("/users")
	{
		users.Get("/", adminHandler.GetAdmins)
		users.Get("/:id", adminHandler.GetAdmin)
		users.Post("/", adminHandler.CreateAdmin)
		users.Put("/:id", adminHandler.UpdateAdmin)
		users.Delete("/:id", adminHandler.DeleteAdmin)
	}

	// countries routes
	countries := r.Group("/countries")
	{
		countries.Get("/", adminHandler.GetCountries)
		countries.Post("/", adminHandler.CreateCountry)
		countries.Post("/:id/images", adminHandler.CreateCountryImage)
		countries.Put("/:id", adminHandler.UpdateCountry)
		countries.Delete("/:id", adminHandler.DeleteCountry)
	}

	// Application routes
	// todo: update to third-party routes
	{
		r.Get("/applications", adminHandler.GetApplications)
		r.Post("/applications", adminHandler.CreateApplication)
		r.Get("/applications/:id", adminHandler.GetApplication)
		r.Post("/applications/:id/documents", adminHandler.CreateApplicationDocuments)
		r.Post("/applications/:id/accept", adminHandler.AcceptApplication)
		r.Post("/applications/:id/reject", adminHandler.RejectApplication)
	}

	// Cities routes
	cities := r.Group("/cities")
	{
		cities.Get("/", adminHandler.GetCities)
		cities.Post("/", adminHandler.CreateCity)
		cities.Put("/:id", adminHandler.UpdateCity)
		cities.Delete("/:id", adminHandler.DeleteCity)
	}

	// Regions routes
	{
		cities.Get("/:city_id/regions", adminHandler.GetRegions)
		cities.Post("/:city_id/regions", adminHandler.CreateRegion)
		cities.Put("/:city_id/regions/:id", adminHandler.UpdateRegion)
		cities.Delete("/:city_id/regions/:id", adminHandler.DeleteRegion)
	}

	// Brands routes
	brands := r.Group("/brands")
	{
		brands.Get("/", adminHandler.GetBrands)
		brands.Post("/", adminHandler.CreateBrand)
		brands.Post("/:id/images", adminHandler.CreateBrandImage)
		brands.Put("/:id", adminHandler.UpdateBrand)
		brands.Delete("/:id", adminHandler.DeleteBrand)
	}

	// Models routes
	{
		brands.Get("/:brand_id/models", adminHandler.GetModels)
		brands.Post("/:brand_id/models", adminHandler.CreateModel)
		brands.Get("/:brand_id/models/:model_id/generations", adminHandler.GetGenerationsByModel)
		brands.Put("/:brand_id/models/:id", adminHandler.UpdateModel)
		brands.Delete("/:brand_id/models/:id", adminHandler.DeleteModel)
	}

	// Body Types routes
	bodyTypes := r.Group("/body-types")
	{
		bodyTypes.Get("/", adminHandler.GetBodyTypes)
		bodyTypes.Post("/", adminHandler.CreateBodyType)
		bodyTypes.Post("/:id", adminHandler.CreateBodyTypeImage)
		bodyTypes.Put("/:id", adminHandler.UpdateBodyType)
		bodyTypes.Delete("/:id", adminHandler.DeleteBodyType)
		bodyTypes.Delete("/:id/images", adminHandler.DeleteBodyTypeImage)
	}

	// Transmissions routes
	transmissions := r.Group("/transmissions")
	{
		transmissions.Get("/", adminHandler.GetTransmissions)
		transmissions.Post("/", adminHandler.CreateTransmission)
		transmissions.Put("/:id", adminHandler.UpdateTransmission)
		transmissions.Delete("/:id", adminHandler.DeleteTransmission)
	}

	// Engines routes
	engines := r.Group("/engines")
	{
		engines.Get("/", adminHandler.GetEngines)
		engines.Post("/", adminHandler.CreateEngine)
		engines.Put("/:id", adminHandler.UpdateEngine)
		engines.Delete("/:id", adminHandler.DeleteEngine)
	}

	// Drivetrains routes
	drivetrains := r.Group("/drivetrains")
	{
		drivetrains.Get("/", adminHandler.GetDrivetrains)
		drivetrains.Post("/", adminHandler.CreateDrivetrain)
		drivetrains.Put("/:id", adminHandler.UpdateDrivetrain)
		drivetrains.Delete("/:id", adminHandler.DeleteDrivetrain)
	}

	// Fuel Types routes
	fuelTypes := r.Group("/fuel-types")
	{
		fuelTypes.Get("/", adminHandler.GetFuelTypes)
		fuelTypes.Post("/", adminHandler.CreateFuelType)
		fuelTypes.Put("/:id", adminHandler.UpdateFuelType)
		fuelTypes.Delete("/:id", adminHandler.DeleteFuelType)
	}
	// Generations routes
	generations := r.Group("/generations")
	{
		generations.Get("/", adminHandler.GetGenerations)
		generations.Post("/", adminHandler.CreateGeneration)
		generations.Put("/:id", adminHandler.UpdateGeneration)
		generations.Post("/:id/images", adminHandler.CreateGenerationImage)
		generations.Delete("/:id", adminHandler.DeleteGeneration)
		// generations.Delete("/:id/images", adminHandler.DeleteGenerationImage)
	}

	// Generation Modifications routes
	{
		generations.Get("/:generation_id/", adminHandler.GetGenerationModifications)
		generations.Post("/:generation_id/", adminHandler.CreateGenerationModification)
		generations.Put("/:generation_id/:id", adminHandler.UpdateGenerationModification)
		generations.Delete("/:generation_id/:id", adminHandler.DeleteGenerationModification)
	}

	// Colors routes
	colors := r.Group("/colors")
	{
		colors.Get("/", adminHandler.GetColors)
		colors.Post("/", adminHandler.CreateColor)
		colors.Post("/:id/images", adminHandler.CreateColorImage)
		colors.Put("/:id", adminHandler.UpdateColor)
		colors.Delete("/:id", adminHandler.DeleteColor)
	}

	// Moto Categories routes
	motoCategories := r.Group("/moto-categories")
	{
		motoCategories.Get("/", adminHandler.GetMotoCategories)
		motoCategories.Get("/:id/brands", adminHandler.GetMotoBrandsByCategoryID)
		motoCategories.Post("/", adminHandler.CreateMotoCategory)
		motoCategories.Put("/:id", adminHandler.UpdateMotoCategory)
		motoCategories.Delete("/:id", adminHandler.DeleteMotoCategory)
	}

	// Moto Category Parameters routes
	{
		motoCategories.Get("/:category_id/parameters", adminHandler.GetMotoCategoryParameters)
		motoCategories.Post("/:category_id/parameters", adminHandler.CreateMotoCategoryParameter)
		motoCategories.Put("/:category_id/parameters/:parameter_id", adminHandler.UpdateMotoCategoryParameter)
		motoCategories.Delete("/:category_id/parameters/:parameter_id", adminHandler.DeleteMotoCategoryParameter)
	}

	// Moto Brands routes
	motoBrands := r.Group("/moto-brands")
	{
		motoBrands.Get("/", adminHandler.GetMotoBrands)
		motoBrands.Get("/:id/models", adminHandler.GetMotoModelsByBrandID)
		motoBrands.Post("/", adminHandler.CreateMotoBrand)
		motoBrands.Post("/:id/images", adminHandler.CreateMotoBrandImage)
		motoBrands.Put("/:id", adminHandler.UpdateMotoBrand)
		motoBrands.Delete("/:id", adminHandler.DeleteMotoBrand)
	}

	// Moto Models routes
	motoModels := r.Group("/moto-models")
	{
		motoModels.Get("/", adminHandler.GetMotoModels)
		motoModels.Post("/", adminHandler.CreateMotoModel)
		motoModels.Put("/:id", adminHandler.UpdateMotoModel)
		motoModels.Delete("/:id", adminHandler.DeleteMotoModel)
	}

	// Moto Parameters routes
	motoParameters := r.Group("/moto-parameters")
	{
		motoParameters.Get("/", adminHandler.GetMotoParameters)
		motoParameters.Post("/", adminHandler.CreateMotoParameter)
		motoParameters.Put("/:id", adminHandler.UpdateMotoParameter)
		motoParameters.Delete("/:id", adminHandler.DeleteMotoParameter)
	}

	// Moto Parameter Values routes
	{
		motoParameters.Get("/:moto_param_id/values", adminHandler.GetMotoParameterValues)
		motoParameters.Post("/:moto_param_id/values", adminHandler.CreateMotoParameterValue)
		motoParameters.Put("/:moto_param_id/values/:id", adminHandler.UpdateMotoParameterValue)
		motoParameters.Delete("/:moto_param_id/values/:id", adminHandler.DeleteMotoParameterValue)
	}

	// Comtrans Categories routes
	comtransCategories := r.Group("/comtrans-categories")
	{
		comtransCategories.Get("/", adminHandler.GetComtransCategories)
		comtransCategories.Get("/:id/brands", adminHandler.GetComtransBrandsByCategoryID)
		comtransCategories.Post("/", adminHandler.CreateComtransCategory)
		comtransCategories.Put("/:id", adminHandler.UpdateComtransCategory)
		comtransCategories.Delete("/:id", adminHandler.DeleteComtransCategory)
	}

	// Comtrans Category Parameters routes
	{
		comtransCategories.Get("/:category_id/parameters", adminHandler.GetComtransCategoryParameters)
		comtransCategories.Post("/:category_id/parameters", adminHandler.CreateComtransCategoryParameter)
		comtransCategories.Put("/:category_id/parameters/:id", adminHandler.UpdateComtransCategoryParameter)
		comtransCategories.Delete("/:category_id/parameters/:id", adminHandler.DeleteComtransCategoryParameter)
	}

	// Comtrans Brands routes
	comtransBrands := r.Group("/comtrans-brands")
	{
		comtransBrands.Get("/", adminHandler.GetComtransBrands)
		comtransBrands.Get("/:id/models", adminHandler.GetComtransModelsByBrandID)
		comtransBrands.Post("/", adminHandler.CreateComtransBrand)
		comtransBrands.Post("/:id/images", adminHandler.CreateComtransBrandImage)
		comtransBrands.Put("/:id", adminHandler.UpdateComtransBrand)
		comtransBrands.Delete("/:id", adminHandler.DeleteComtransBrand)
	}

	// Comtrans Models routes
	comtransModels := r.Group("/comtrans-models")
	{
		comtransModels.Get("/", adminHandler.GetComtransModels)
		comtransModels.Post("/", adminHandler.CreateComtransModel)
		comtransModels.Put("/:id", adminHandler.UpdateComtransModel)
		comtransModels.Delete("/:id", adminHandler.DeleteComtransModel)
	}

	// Comtrans Parameters routes
	comtransParameters := r.Group("/comtrans-parameters")
	{
		comtransParameters.Get("/", adminHandler.GetComtransParameters)
		comtransParameters.Post("/", adminHandler.CreateComtransParameter)
		comtransParameters.Put("/:id", adminHandler.UpdateComtransParameter)
		comtransParameters.Delete("/:id", adminHandler.DeleteComtransParameter)
	}

	// Comtrans Parameter Values routes
	{
		comtransParameters.Get("/:parameter_id/values", adminHandler.GetComtransParameterValues)
		comtransParameters.Post("/:parameter_id/values", adminHandler.CreateComtransParameterValue)
		comtransParameters.Put("/:parameter_id/values/:id", adminHandler.UpdateComtransParameterValue)
		comtransParameters.Delete("/:parameter_id/values/:id", adminHandler.DeleteComtransParameterValue)
	}

	// company types
	companyTypes := r.Group("/company-types")
	{
		companyTypes.Get("/", adminHandler.GetCompanyTypes)
		companyTypes.Get("/:id", adminHandler.GetCompanyType)
		companyTypes.Post("/", adminHandler.CreateCompanyType)
		companyTypes.Put("/:id", adminHandler.UpdateCompanyType)
		companyTypes.Delete("/:id", adminHandler.DeleteCompanyType)
	}

	// activity fields
	activityFields := r.Group("/activity-fields")
	{
		activityFields.Get("/", adminHandler.GetActivityFields)
		activityFields.Get("/:id", adminHandler.GetActivityField)
		activityFields.Post("/", adminHandler.CreateActivityField)
		activityFields.Put("/:id", adminHandler.UpdateActivityField)
		activityFields.Delete("/:id", adminHandler.DeleteActivityField)
	}

}
