package model

type User struct {
	Name     string
	Email    string
	Password string
	ID       int64
}

type DeleteCarImageRequest struct {
	Image string `json:"image" validate:"required"`
}

type ThirdPartyCreateCarRequest struct {
	// new
	// BodyTypeID     int      `json:"body_type_id" validate:"required"`
	PhoneNumbers   []string `json:"phone_numbers" validate:"required"`
	Wheel          *bool    `json:"wheel" validate:"required"` // left true, right false
	Description    string   `json:"description"`
	CityID         int      `json:"city_id" validate:"required"`
	BrandID        int      `json:"brand_id" validate:"required"`
	ModelID        int      `json:"model_id" validate:"required"`
	ModificationID int      `json:"modification_id" validate:"required"`
	Year           int      `json:"year" validate:"required"`
	Odometer       int      `json:"odometer" validate:"required"`
	Price          int      `json:"price" validate:"required"`
	ColorID        int      `json:"color_id" validate:"required"`
	Owners         int      `json:"owners"`
	TradeIn        int      `json:"trade_in" validate:"required"`
	New            bool     `json:"new"`
	Crash          bool     `json:"crash"`

	//
	// OwnershipTypeId int    `json:"ownership_type_id"`
	// Credit          bool   `json:"credit"`
	// DoorCount       int    `json:"door_count"`
	// InteriorColorID int `json:"interior_color_id"`
	// Negotiable      bool `json:"negotiable"`
	// ModificationID  int  `json:"modification_id"`
	// MileageKM       int    `json:"mileage_km"`
	// GenerationID int `json:"generation_id" validate:"required"`
}

type CreateCarRequest struct {
	// new
	// BodyTypeID     int      `json:"body_type_id" validate:"required"`
	PhoneNumbers   []string `json:"phone_numbers" validate:"required"`
	Wheel          *bool    `json:"wheel" validate:"required"` // left true, right false
	Description    string   `json:"description"`
	VinCode        string   `json:"vin_code" validate:"required"`
	CityID         int      `json:"city_id" validate:"required"`
	BrandID        int      `json:"brand_id" validate:"required"`
	ModelID        int      `json:"model_id" validate:"required"`
	ModificationID int      `json:"modification_id" validate:"required"`
	Year           int      `json:"year" validate:"required"`
	Odometer       int      `json:"odometer" validate:"required"`
	Price          int      `json:"price" validate:"required"`
	ColorID        int      `json:"color_id" validate:"required"`
	Owners         int      `json:"owners"`
	TradeIn        int      `json:"trade_in" validate:"required"`
	New            bool     `json:"new"`
	Crash          bool     `json:"crash"`

	//
	// OwnershipTypeId int    `json:"ownership_type_id"`
	// Credit          bool   `json:"credit"`
	// DoorCount       int    `json:"door_count"`
	// InteriorColorID int `json:"interior_color_id"`
	// Negotiable      bool `json:"negotiable"`
	// ModificationID  int  `json:"modification_id"`
	// MileageKM       int    `json:"mileage_km"`
	// GenerationID int `json:"generation_id" validate:"required"`
}

type DeleteCarVideoRequest struct {
	Video string `json:"video" validate:"required"`
}

type UpdateCarRequest struct {
	PhoneNumbers   []string `json:"phone_numbers" validate:"required"`
	Wheel          *bool    `json:"wheel" validate:"required"` // left true, right false
	Description    string   `json:"description"`
	ID             int      `json:"id" validate:"required"`
	CityID         int      `json:"city_id" validate:"required"`
	BrandID        int      `json:"brand_id" validate:"required"`
	ModificationID int      `json:"modification_id" validate:"required"`
	ModelID        int      `json:"model_id" validate:"required"`
	Year           int      `json:"year" validate:"required"`
	Odometer       int      `json:"odometer" validate:"required"`
	Price          int      `json:"price" validate:"required"`
	ColorID        int      `json:"color_id" validate:"required"`
	Owners         int      `json:"owners" validate:"required"`
	TradeIn        int      `json:"trade_in" validate:"required"`
	New            bool     `json:"new"`
	Crash          bool     `json:"crash"`
}

type UpdateProfileRequest struct {
	Username          string            `json:"username" validate:"required,min=3,max=20"`
	Google            string            `json:"google"`
	Birthday          string            `json:"birthday"`
	AboutMe           string            `json:"about_me"`
	DrivingExperience int               `json:"driving_experience"`
	Notification      bool              `json:"notification"`
	CityID            int               `json:"city_id"`
	Contacts          map[string]string `json:"contacts"`
	Address           string            `json:"address"`
	PhoneNumber       string            `json:"phone_number"`
	Email             string            `json:"email"`
	// todo: add city
}

// Admin request/response models
type CreateCityRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}
