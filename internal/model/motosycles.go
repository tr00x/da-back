package model

import "time"

// RESPONSES
type GetMotorcycleCategoriesResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GetMotorcycleParameterValuesResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type GetMotorcycleParametersResponse struct {
	Values []GetMotorcycleParameterValuesResponse `json:"values"`
	Name   string                                 `json:"name"`
	ID     int                                    `json:"id"`
}

type GetMotorcycleBrandsResponse struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	ID    int    `json:"id"`
}

type GetMotorcycleModelsResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// REQUESTS
type CreateMotorcycleParameterRequest struct {
	ParameterID int `json:"parameter_id" validate:"required"`
	ValueID     int `json:"value_id" validate:"required"`
}

type CreateMotorcycleRequest struct {
	Parameters         []CreateMotorcycleParameterRequest `json:"parameters"`
	Crash              *bool                              `json:"crash"`
	NotCleared         *bool                              `json:"not_cleared"`
	PTC                *bool                              `json:"ptc"`
	RefuseDealersCalls *bool                              `json:"refuse_dealers_calls"`
	OnlyChat           *bool                              `json:"only_chat"`
	ProtectSpam        *bool                              `json:"protect_spam"`
	VerifiedBuyers     *bool                              `json:"verified_buyers"`
	MotoCategoryID     string                             `json:"moto_category_id" validate:"required"`
	BrandID            string                             `json:"moto_brand_id" validate:"required"`
	ModelID            string                             `json:"moto_model_id" validate:"required"`
	DateOfPurchase     string                             `json:"date_of_purchase"`
	WarrantyDate       string                             `json:"warranty_date"`
	VinCode            string                             `json:"vin_code" validate:"required"`
	Certificate        string                             `json:"certificate"`
	Description        string                             `json:"description"`
	CanLookCoordinate  string                             `json:"can_look_coordinate"`
	PhoneNumber        string                             `json:"phone_number" validate:"required"`
	ContactPerson      string                             `json:"contact_person"`
	Email              string                             `json:"email"`
	PriceType          string                             `json:"price_type" validate:"required,oneof=USD AED RUB EUR"`
	FuelTypeID         int                                `json:"fuel_type_id" validate:"required"`
	CityID             int                                `json:"city_id" validate:"required"`
	ColorID            int                                `json:"color_id" validate:"required"`
	Engine             int                                `json:"engine"`
	Power              int                                `json:"power"`
	Year               int                                `json:"year" validate:"required"`
	NumberOfCycles     int                                `json:"number_of_cycles"`
	Odometer           int                                `json:"odometer"`
	Owners             int                                `json:"owners"`
	Price              int                                `json:"price" validate:"required"`
}

// Owner represents the motorcycle owner information
type MotorcycleOwner struct {
	Contacts map[string]string `json:"contacts"`
	Username string            `json:"username"`
	Avatar   string            `json:"avatar"`
	ID       int               `json:"id"`
}

// MotorcycleParameter represents a motorcycle parameter with its value
type MotorcycleParameter struct {
	Parameter        string `json:"parameter"`
	ParameterValue   string `json:"parameter_value"`
	ParameterID      int    `json:"parameter_id"`
	ParameterValueID int    `json:"parameter_value_id"`
}

type GetMotorcyclesResponse struct {
	UpdatedAt          time.Time             `json:"updated_at"`
	CreatedAt          time.Time             `json:"created_at"`
	Parameters         []MotorcycleParameter `json:"parameters"`
	Images             []string              `json:"images"`
	Videos             []string              `json:"videos"`
	Owner              MotorcycleOwner       `json:"owner"`
	Crash              *bool                 `json:"crash"`
	NotCleared         *bool                 `json:"not_cleared"`
	PTC                *bool                 `json:"ptc"`
	RefuseDealersCalls *bool                 `json:"refuse_dealers_calls"`
	OnlyChat           *bool                 `json:"only_chat"`
	ProtectSpam        *bool                 `json:"protect_spam"`
	VerifiedBuyers     *bool                 `json:"verified_buyers"`
	DateOfPurchase     string                `json:"date_of_purchase"`
	WarrantyDate       string                `json:"warranty_date"`
	VinCode            string                `json:"vin_code"`
	Certificate        string                `json:"certificate"`
	Description        string                `json:"description"`
	CanLookCoordinate  string                `json:"can_look_coordinate"`
	PhoneNumber        string                `json:"phone_number"`
	ContactPerson      string                `json:"contact_person"`
	Email              string                `json:"email"`
	PriceType          string                `json:"price_type"`
	Status             string                `json:"status"`
	MotoCategory       string                `json:"moto_category"`
	MotoBrand          string                `json:"moto_brand"`
	MotoModel          string                `json:"moto_model"`
	FuelType           string                `json:"fuel_type"`
	City               string                `json:"city"`
	Color              string                `json:"color"`
	ID                 int                   `json:"id"`
	Engine             int                   `json:"engine"`
	Power              int                   `json:"power"`
	Year               int                   `json:"year"`
	NumberOfCycles     int                   `json:"number_of_cycles"`
	Odometer           int                   `json:"odometer"`
	Owners             int                   `json:"owners"`
	Price              int                   `json:"price"`
	MyCar              bool                  `json:"my_car"`
}
