package model

import "time"

type GetPriceRecommendationRequest struct {
	BrandID        string `json:"brand_id"`
	ModelID        string `json:"model_id"`
	Year           string `json:"year"`
	ModificationID string `json:"modification_id"`
	Odometer       string `json:"odometer"`
	CityID         string `json:"city_id"`
}

type GetPriceRecommendationResponse struct {
	MinPrice int `json:"min_price"`
	MaxPrice int `json:"max_price"`
	AvgPrice int `json:"avg_price"`
}

type Brand struct {
	ID         *int    `json:"id"`
	Name       *string `json:"name"`
	Logo       *string `json:"logo"`
	ModelCount *int    `json:"model_count"`
}

type GetBrandsResponse struct {
	Name       string  `json:"name"`
	Logo       *string `json:"logo"`
	ID         int     `json:"id"`
	ModelCount int     `json:"model_count"`
}

type GetProfileResponse struct {
	Birthday          *time.Time        `json:"birthday"`
	Email             *string           `json:"email"`
	Phone             *string           `json:"phone"`
	Username          *string           `json:"username"`
	Google            *string           `json:"google"`
	RegisteredBy      *string           `json:"registered_by"`
	City              *City             `json:"city"`
	Contacts          map[string]string `json:"contacts"`
	Address           *string           `json:"address"`
	AboutMe           *string           `json:"about_me"`
	DrivingExperience *int              `json:"driving_experience"`
	ID                int               `json:"id"`
	Notification      *bool             `json:"notification"`
}

type GetFilterBrandsResponse struct {
	PopularBrands []Brand `json:"popular_brands"`
	AllBrands     []Brand `json:"all_brands"`
}

type Region struct {
	Name *string `json:"name"`
	ID   *int    `json:"id"`
}

type GetCitiesResponse struct {
	Regions []Region `json:"regions"`
	Name    *string  `json:"name"`
	ID      int      `json:"id"`
}

type GetModificationsResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Model struct {
	Name *string `json:"name"`
	ID   *int    `json:"id"`
}

type GetFilterModelsResponse struct {
	PopularModels []Model `json:"popular_models"`
	AllModels     []Model `json:"all_models"`
}

type GetYearsResponse struct {
	Years []*int `json:"years"`
}

type Modification struct {
	Engine       *string `json:"engine"`
	FuelType     *string `json:"fuel_type"`
	Drivetrain   *string `json:"drivetrain"`
	Transmission *string `json:"transmission"`
	ID           *int    `json:"id"`
}

type Generation struct {
	Modifications []Modification `json:"modifications"`
	Name          string         `json:"name"`
	Image         string         `json:"image"`
	ID            int            `json:"id"`
	StartYear     int            `json:"start_year"`
	EndYear       int            `json:"end_year"`
}

type BodyType struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
	ID    *int    `json:"id"`
}

type Transmission struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Engine struct {
	Value string `json:"value"`
	ID    int    `json:"id"`
}

type Drivetrain struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type FuelType struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Color struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
	ID    *int    `json:"id"`
}

type Home struct {
	Popular []GetCarsResponse `json:"popular"`
}

type Owner struct {
	Avatar   *string           `json:"avatar"`
	Username *string           `json:"username"`
	Id       *int              `json:"id"`
	RoleID   *int              `json:"role_id"`
	Contacts map[string]string `json:"contacts"`
}

type GetCarsResponse struct {
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Images       *[]string  `json:"images"`
	Videos       *[]string  `json:"videos"`
	PhoneNumbers *[]string  `json:"phone_numbers"`
	Owner        *Owner     `json:"owner"`
	Brand        *string    `json:"brand"`
	Region       *string    `json:"region"`
	City         *string    `json:"city"`
	Transmission *string    `json:"transmission"`
	Engine       *string    `json:"engine"`
	Drivetrain   *string    `json:"drivetrain"`
	FuelType     *string    `json:"fuel_type"`
	VinCode      *string    `json:"vin_code"`
	Color        *string    `json:"color"`
	Description  *string    `json:"description"`
	Status       *int       `json:"status"`
	TradeIn      *int       `json:"trade_in"`
	Owners       *int       `json:"owners"`
	Mileage      *int       `json:"mileage"` // todo: change it to odometer
	Model        string     `json:"model"`
	BodyType     string     `json:"body_type"`
	ID           int        `json:"id"`
	Year         int        `json:"year"`
	Price        int        `json:"price"`
	ViewCount    int        `json:"view_count"`
	Credit       *bool      `json:"credit"`
	New          *bool      `json:"new"`
	Crash        *bool      `json:"crash"`
	MyCar        *bool      `json:"my_car"`
	Liked        *bool      `json:"liked"`
}

type GetMyCarsResponse struct {
	Type      string     `json:"type"`
	Model     *string    `json:"model"`
	CreatedAt *time.Time `json:"created_at"`
	Images    *[]string  `json:"images"`
	Brand     *string    `json:"brand"`
	Status    *int       `json:"status"`
	TradeIn   *int       `json:"trade_in"`
	Year      *int       `json:"year"`
	Price     *int       `json:"price"`
	ViewCount *int       `json:"view_count"`
	Credit    *bool      `json:"credit"`
	New       *bool      `json:"new"`
	Crash     *bool      `json:"crash"`
	MyCar     *bool      `json:"my_car"`
	ID        int        `json:"id"`
}

type City struct {
	Name *string `json:"name"`
	ID   *int    `json:"id"`
}

type EditCarGeneration struct {
	Name      string `json:"name"`
	Image     string `json:"image"`
	ID        int    `json:"id"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year"`
}

type GetEditCarsResponse struct {
	CreatedAt    *time.Time         `json:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at"`
	Images       *[]string          `json:"images"`
	Videos       *[]string          `json:"videos"`
	PhoneNumbers *[]string          `json:"phone_numbers"`
	Brand        *Brand             `json:"brand"`
	Region       *Region            `json:"region"`
	City         *City              `json:"city"`
	Model        *Model             `json:"model"`
	Modification *Modification      `json:"modification"`
	BodyType     *BodyType          `json:"body_type"`
	Generation   *EditCarGeneration `json:"generation"`
	Color        *Color             `json:"color"`
	Year         *int               `json:"year"`
	Price        *int               `json:"price"`
	Odometer     *int               `json:"odometer"`
	Status       *int               `json:"status"`
	ViewCount    *int               `json:"view_count"`
	TradeIN      *int               `json:"trade_id"`
	Owners       *int               `json:"owners"`
	VinCode      *string            `json:"vin_code"`
	Description  *string            `json:"description"`
	ID           int                `json:"id"`
	Credit       *bool              `json:"credit"`
	New          *bool              `json:"new"`
	MyCar        *bool              `json:"my_car"`
	Wheel        *bool              `json:"wheel"`
	Crash        *bool              `json:"crash"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	VerifiedEmail bool   `json:"verified_email"`
}

type AdminProfileResponse struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	ID          int      `json:"id"`
}

type AdminResponse struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	Status      int      `json:"status"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type AdminApplicationResponse struct {
	LicenceIssueDate  *time.Time `json:"licence_issue_date"`
	LicenceExpiryDate *time.Time `json:"licence_expiry_date"`
	CreatedAt         *time.Time `json:"created_at"`
	CompanyName       *string    `json:"company_name"`
	FullName          *string    `json:"full_name"`
	Email             *string    `json:"email"`
	Phone             *string    `json:"phone"`
	Status            *string    `json:"status"`
	ID                *int       `json:"id"`
}

type AdminApplicationByIDResponse struct {
	LicenceIssueDate  *time.Time `json:"licence_issue_date"`
	LicenceExpiryDate *time.Time `json:"licence_expiry_date"`
	CreatedAt         *time.Time `json:"created_at"`
	CompanyName       *string    `json:"company_name"`
	FullName          *string    `json:"full_name"`
	Email             *string    `json:"email"`
	Phone             *string    `json:"phone"`
	Status            *string    `json:"status"`
	CopyOFIDURL       *string    `json:"copy_of_id_url"`
	MemorandumURL     *string    `json:"memorandum_url"`
	LicenceURL        *string    `json:"licence_url"`
	CompanyType       *string    `json:"company_type"`
	ActivityField     *string    `json:"activity_field"`
	VATNumber         *string    `json:"vat_number"`
	Address           *string    `json:"address"`
	ID                *int       `json:"id"`
}

// Admin response models
type AdminCityResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type CompanyType struct {
	Name   string `json:"name"`
	NameRu string `json:"name_ru"`
	NameAe string `json:"name_ae"`
	ID     int    `json:"id"`
}

type AdminCountryResponse struct {
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	NameRu      string    `json:"name_ru"`
	NameAe      string    `json:"name_ae"`
	CountryCode string    `json:"country_code"`
	Flag        string    `json:"flag"`
	ID          int       `json:"id"`
}

type Country struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Flag        string `json:"flag"`
}

type AdminBrandResponse struct {
	UpdatedAt  *time.Time `json:"updated_at"`
	Name       *string    `json:"name"`
	NameRu     *string    `json:"name_ru"`
	NameAe     *string    `json:"name_ae"`
	Logo       *string    `json:"logo"`
	ID         *int       `json:"id"`
	ModelCount *int       `json:"model_count"`
	Popular    *bool      `json:"popular"`
}

type AdminModelResponse struct {
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	NameRu      string    `json:"name_ru"`
	NameAe      string    `json:"name_ae"`
	BrandName   string    `json:"brand_name"`
	BrandNameRu string    `json:"brand_name_ru"`
	ID          int       `json:"id"`
	BrandID     int       `json:"brand_id"`
	Popular     bool      `json:"popular"`
}

type AdminBodyTypeResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	Image     string    `json:"image"`
	ID        int       `json:"id"`
}

type AdminTransmissionResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminEngineResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminDrivetrainResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminFuelTypeResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminServiceTypeResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
}

type AdminServiceResponse struct {
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	ServiceTypeName string    `json:"service_type_name"`
	ID              int       `json:"id"`
	ServiceTypeID   int       `json:"service_type_id"`
}

type AdminGenerationResponse struct {
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	NameRu      string    `json:"name_ru"`
	NameAe      string    `json:"name_ae"`
	ModelName   string    `json:"model_name"`
	Image       string    `json:"image"`
	ModelNameRu string    `json:"model_name_ru"`
	ID          int       `json:"id"`
	ModelID     int       `json:"model_id"`
	StartYear   int       `json:"start_year"`
	EndYear     int       `json:"end_year"`
	Wheel       bool      `json:"wheel"`
}

type AdminGenerationModificationResponse struct {
	BodyTypeName       string `json:"body_type_name"`
	BodyTypeNameRu     string `json:"body_type_name_ru"`
	EngineName         string `json:"engine_name"`
	EngineNameRu       string `json:"engine_name_ru"`
	FuelTypeName       string `json:"fuel_type_name"`
	FuelTypeNameRu     string `json:"fuel_type_name_ru"`
	DrivetrainName     string `json:"drivetrain_name"`
	DrivetrainNameRu   string `json:"drivetrain_name_ru"`
	TransmissionName   string `json:"transmission_name"`
	TransmissionNameRu string `json:"transmission_name_ru"`
	ID                 int    `json:"id"`
	GenerationID       int    `json:"generation_id"`
	BodyTypeID         int    `json:"body_type_id"`
	EngineID           int    `json:"engine_id"`
	FuelTypeID         int    `json:"fuel_type_id"`
	DrivetrainID       int    `json:"drivetrain_id"`
	TransmissionID     int    `json:"transmission_id"`
}

type AdminConfigurationResponse struct {
	BodyTypeName string `json:"body_type_name"`
	ID           int    `json:"id"`
	BodyTypeID   int    `json:"body_type_id"`
	GenerationID int    `json:"generation_id"`
}

type AdminColorResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	Image     string    `json:"image"`
	ID        int       `json:"id"`
}

type AdminMotoCategoryResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminMotoBrandResponse struct {
	CreatedAt          time.Time `json:"created_at"`
	Name               string    `json:"name"`
	NameRu             string    `json:"name_ru"`
	NameAe             string    `json:"name_ae"`
	Image              string    `json:"image"`
	MotoCategoryName   string    `json:"moto_category_name"`
	MotoCategoryNameRu string    `json:"moto_category_name_ru"`
	ID                 int       `json:"id"`
	MotoCategoryID     int       `json:"moto_category_id"`
}

type AdminMotoModelResponse struct {
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	NameRu          string    `json:"name_ru"`
	NameAe          string    `json:"name_ae"`
	MotoBrandName   string    `json:"moto_brand_name"`
	MotoBrandNameRu string    `json:"moto_brand_name_ru"`
	ID              int       `json:"id"`
	MotoBrandID     int       `json:"moto_brand_id"`
}

type AdminMotoParameterResponse struct {
	CreatedAt          time.Time `json:"created_at"`
	Name               string    `json:"name"`
	NameRu             string    `json:"name_ru"`
	NameAe             string    `json:"name_ae"`
	MotoCategoryName   string    `json:"moto_category_name"`
	MotoCategoryNameRu string    `json:"moto_category_name_ru"`
	ID                 int       `json:"id"`
	MotoCategoryID     int       `json:"moto_category_id"`
}

type AdminMotoParameterValueResponse struct {
	CreatedAt       time.Time `json:"created_at"`
	Name            string    `json:"name"`
	NameRu          string    `json:"name_ru"`
	NameAe          string    `json:"name_ae"`
	ID              int       `json:"id"`
	MotoParameterID int       `json:"moto_parameter_id"`
}

type AdminMotoCategoryParameterResponse struct {
	CreatedAt           time.Time `json:"created_at"`
	MotoParameterName   string    `json:"moto_parameter_name"`
	MotoParameterNameRu string    `json:"moto_parameter_name_ru"`
	MotoCategoryID      int       `json:"moto_category_id"`
	MotoParameterID     int       `json:"moto_parameter_id"`
}

type AdminComtransCategoryResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	NameRu    string    `json:"name_ru"`
	NameAe    string    `json:"name_ae"`
	ID        int       `json:"id"`
}

type AdminComtransBrandResponse struct {
	CreatedAt              time.Time `json:"created_at"`
	Name                   string    `json:"name"`
	NameRu                 string    `json:"name_ru"`
	NameAe                 string    `json:"name_ae"`
	Image                  string    `json:"image"`
	ComtransCategoryName   string    `json:"comtrans_category_name"`
	ComtransCategoryNameRu string    `json:"comtrans_category_name_ru"`
	ID                     int       `json:"id"`
	ComtransCategoryID     int       `json:"comtrans_category_id"`
}

type AdminComtransModelResponse struct {
	CreatedAt           time.Time `json:"created_at"`
	Name                string    `json:"name"`
	NameRu              string    `json:"name_ru"`
	NameAe              string    `json:"name_ae"`
	ComtransBrandName   string    `json:"comtrans_brand_name"`
	ComtransBrandNameRu string    `json:"comtrans_brand_name_ru"`
	ID                  int       `json:"id"`
	ComtransBrandID     int       `json:"comtrans_brand_id"`
}

type AdminComtransParameterResponse struct {
	CreatedAt              time.Time `json:"created_at"`
	Name                   string    `json:"name"`
	NameRu                 string    `json:"name_ru"`
	NameAe                 string    `json:"name_ae"`
	ComtransCategoryName   string    `json:"comtrans_category_name"`
	ComtransCategoryNameRu string    `json:"comtrans_category_name_ru"`
	ID                     int       `json:"id"`
	ComtransCategoryID     int       `json:"comtrans_category_id"`
}

type AdminComtransParameterValueResponse struct {
	CreatedAt           time.Time `json:"created_at"`
	Name                string    `json:"name"`
	NameRu              string    `json:"name_ru"`
	NameAe              string    `json:"name_ae"`
	ID                  int       `json:"id"`
	ComtransParameterID int       `json:"comtrans_parameter_id"`
}

type AdminComtransCategoryParameterResponse struct {
	CreatedAt               time.Time `json:"created_at"`
	ComtransParameterName   string    `json:"comtrans_parameter_name"`
	ComtransParameterNameRu string    `json:"comtrans_parameter_name_ru"`
	ID                      int       `json:"id"`
	ComtransCategoryID      int       `json:"comtrans_category_id"`
	ComtransParameterID     int       `json:"comtrans_parameter_id"`
}

type ThirdPartProfileDestinationsRes struct {
	FromCountry *Country `json:"from_country"`
	ToCountry   *Country `json:"to_country"`
}

type ThirdPartyGetProfileRes struct {
	Registered    *time.Time                         `json:"registered"`
	Destinations  *[]ThirdPartProfileDestinationsRes `json:"destinations"`
	CompanyName   *string                            `json:"company_name"`
	AboutUs       *string                            `json:"about_us"`
	Email         *string                            `json:"email"`
	Contacts      map[string]string                  `json:"contacts"`
	Phone         *string                            `json:"phone"`
	Address       *string                            `json:"address"`
	Coordinates   *string                            `json:"coordinates"`
	Avatar        *string                            `json:"avatar"`
	Banner        *string                            `json:"banner"`
	Message       *string                            `json:"message"`
	VATNumber     *string                            `json:"vat_number"`
	CompanyType   *string                            `json:"company_type"`
	ActivityField *string                            `json:"activity_field"`
	UserID        *int                               `json:"user_id"`
	RoleID        *int                               `json:"role_id"`
}

type ThirdPartyGetRegistrationDataRes struct {
	CompanyTypes   []Model `json:"company_types"`
	ActivityFields []Model `json:"activity_fields"`
}

type LogistDestinationResponse struct {
	ID        int      `json:"id"`
	From      *Country `json:"from"`
	To        *Country `json:"to"`
	CreatedAt string   `json:"created_at"`
}

// User role responses (brokers, logists, services)
type ThirdPartyUserResponse struct {
	Registered *time.Time `json:"registered"`
	Username   *string    `json:"username"`
	Avatar     *string    `json:"avatar"`
	ID         int        `json:"id"`
}
