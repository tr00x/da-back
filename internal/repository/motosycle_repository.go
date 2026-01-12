package repository

import (
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/pkg/auth"
)

type MotorcycleRepository struct {
	config *config.Config
	db     *pgxpool.Pool
}

func NewMotorcycleRepository(config *config.Config, db *pgxpool.Pool) *MotorcycleRepository {
	return &MotorcycleRepository{config, db}
}

func (r *MotorcycleRepository) GetMotorcycleCategories(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.GetMotorcycleCategoriesResponse, error) {
	data := make([]model.GetMotorcycleCategoriesResponse, 0)
	q := `
		SELECT id, ` + nameColumn + ` FROM moto_categories
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category model.GetMotorcycleCategoriesResponse
		err = rows.Scan(&category.ID, &category.Name)

		if err != nil {
			return nil, err
		}

		data = append(data, category)
	}

	return data, nil
}

func (r *MotorcycleRepository) GetMotorcycleParameters(ctx *fasthttp.RequestCtx, categoryID string, nameColumn string) ([]model.GetMotorcycleParametersResponse, error) {
	data := make([]model.GetMotorcycleParametersResponse, 0)
	q := `
		SELECT 
			moto_parameters.id,
			moto_parameters.` + nameColumn + `,
			json_agg(
				json_build_object(
					'id', moto_parameter_values.id,
					'name', moto_parameter_values.` + nameColumn + `
				)
			) as values
		FROM moto_parameters
		LEFT JOIN moto_parameter_values ON moto_parameters.id = moto_parameter_values.moto_parameter_id
		WHERE moto_parameters.moto_category_id = $1
		GROUP BY moto_parameters.id, moto_parameters.` + nameColumn + `
	`

	rows, err := r.db.Query(ctx, q, categoryID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var parameter model.GetMotorcycleParametersResponse
		err = rows.Scan(&parameter.ID, &parameter.Name, &parameter.Values)

		if err != nil {
			return nil, err
		}

		data = append(data, parameter)
	}

	return data, nil
}

func (r *MotorcycleRepository) GetMotorcycleBrands(ctx *fasthttp.RequestCtx, categoryID string, nameColumn string) ([]model.GetMotorcycleBrandsResponse, error) {
	data := make([]model.GetMotorcycleBrandsResponse, 0)
	q := `
		SELECT id, ` + nameColumn + `, $2 || image as image FROM moto_brands
		WHERE moto_category_id = $1
	`

	rows, err := r.db.Query(ctx, q, categoryID, r.config.IMAGE_BASE_URL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var brand model.GetMotorcycleBrandsResponse
		err = rows.Scan(&brand.ID, &brand.Name, &brand.Image)

		if err != nil {
			return nil, err
		}

		data = append(data, brand)
	}

	return data, nil
}

func (r *MotorcycleRepository) GetMotorcycleModelsByBrandID(ctx *fasthttp.RequestCtx, categoryID, brandID, nameColumn string) ([]model.GetMotorcycleModelsResponse, error) {
	data := make([]model.GetMotorcycleModelsResponse, 0)
	q := `
		SELECT id, ` + nameColumn + ` FROM moto_models
		WHERE moto_brand_id = $1
	`

	rows, err := r.db.Query(ctx, q, brandID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var model model.GetMotorcycleModelsResponse
		err = rows.Scan(&model.ID, &model.Name)

		if err != nil {
			return nil, err
		}

		data = append(data, model)
	}

	return data, nil
}

func (r *MotorcycleRepository) CreateMotorcycle(ctx *fasthttp.RequestCtx, req model.CreateMotorcycleRequest, userID int) (model.SuccessWithId, error) {
	data := model.SuccessWithId{}
	parameters := req.Parameters
	req.Parameters = nil

	keys, values, args := auth.BuildParams(req)

	q := `
		INSERT INTO motorcycles ( 
			` + strings.Join(keys, ", ") + `,
			user_id
		) VALUES (
			` + strings.Join(values, ", ") + `,
			$` + strconv.Itoa(len(keys)+1) + `
		) returning id
	`
	var id int
	args = append(args, userID)
	err := r.db.QueryRow(ctx, q, args...).Scan(&id)

	if err != nil {
		return data, err
	}

	for i := range parameters {

		q := `
			INSERT INTO motorcycle_parameters (motorcycle_id, moto_parameter_id, moto_parameter_value_id)
			VALUES ($1, $2, $3)
		`
		_, err = r.db.Exec(ctx, q, id, parameters[i].ParameterID, parameters[i].ValueID)

		if err != nil {
			return data, err
		}
	}

	data.Message = "Motorcycle created successfully"
	data.Id = id

	return data, err
}

func (r *MotorcycleRepository) GetMotorcycles(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.GetMotorcyclesResponse, error) {
	data := make([]model.GetMotorcyclesResponse, 0)
	q := `
		select 
			mcs.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $1 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			mcs.engine,
			mcs.power,
			mcs.year,
			mcs.number_of_cycles,
			mcs.odometer,
			mcs.crash,
			mcs.not_cleared,
			mcs.owners,
			mcs.date_of_purchase,
			mcs.warranty_date,
			mcs.ptc,
			mcs.vin_code,
			mcs.certificate,
			mcs.description,
			mcs.can_look_coordinate,
			mcs.phone_number,
			mcs.refuse_dealers_calls,
			mcs.only_chat,
			mcs.protect_spam,
			mcs.verified_buyers,
			mcs.contact_person,
			mcs.email,
			mcs.price,
			mcs.price_type,
			mcs.status,
			mcs.updated_at,
			mcs.created_at,
			mocs.` + nameColumn + ` as moto_category,
			mbs.` + nameColumn + ` as moto_brand,
			mms.` + nameColumn + ` as moto_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN mcs.user_id = 1 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from motorcycles mcs
		left join profiles pf on pf.user_id = mcs.user_id
		left join moto_categories mocs on mocs.id = mcs.moto_category_id
		left join moto_brands mbs on mbs.id = mcs.moto_brand_id
		left join moto_models mms on mms.id = mcs.moto_model_id
		left join fuel_types fts on fts.id = mcs.fuel_type_id
		left join cities cs on cs.id = mcs.city_id
		left join colors cls on cls.id = mcs.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', mcp.moto_parameter_id,
					'parameter_value_id', mpv.id,
					'parameter', mp.` + nameColumn + `,
					'parameter_value', mpv.` + nameColumn + `
				)
			) AS parameters
			FROM motorcycle_parameters mcp
			LEFT JOIN moto_parameters mp ON mp.id = mcp.moto_parameter_id
			LEFT JOIN moto_parameter_values mpv ON mpv.id = mcp.moto_parameter_value_id
			WHERE mcp.motorcycle_id = mcs.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $1 || image as image
				FROM moto_images
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $1 || video as video
				FROM moto_videos
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) v
		) videos ON true;

	`

	rows, err := r.db.Query(ctx, q, r.config.IMAGE_BASE_URL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var motorcycle model.GetMotorcyclesResponse
		err = rows.Scan(
			&motorcycle.ID, &motorcycle.Owner, &motorcycle.Engine, &motorcycle.Power, &motorcycle.Year,
			&motorcycle.NumberOfCycles, &motorcycle.Odometer, &motorcycle.Crash, &motorcycle.NotCleared,
			&motorcycle.Owners, &motorcycle.DateOfPurchase, &motorcycle.WarrantyDate, &motorcycle.PTC,
			&motorcycle.VinCode, &motorcycle.Certificate, &motorcycle.Description, &motorcycle.CanLookCoordinate,
			&motorcycle.PhoneNumber, &motorcycle.RefuseDealersCalls, &motorcycle.OnlyChat,
			&motorcycle.ProtectSpam, &motorcycle.VerifiedBuyers, &motorcycle.ContactPerson,
			&motorcycle.Email, &motorcycle.Price, &motorcycle.PriceType, &motorcycle.Status,
			&motorcycle.UpdatedAt, &motorcycle.CreatedAt, &motorcycle.MotoCategory, &motorcycle.MotoBrand,
			&motorcycle.MotoModel, &motorcycle.FuelType, &motorcycle.City, &motorcycle.Color, &motorcycle.MyCar,
			&motorcycle.Parameters, &motorcycle.Images, &motorcycle.Videos)

		if err != nil {
			return nil, err
		}

		data = append(data, motorcycle)
	}

	return data, nil
}

func (r *MotorcycleRepository) CreateMotorcycleImages(ctx *fasthttp.RequestCtx, motorcycleID int, images []string) error {

	if len(images) == 0 {
		return nil
	}

	q := `
		INSERT INTO moto_images (moto_id, image) VALUES ($1, $2)
	`

	for i := range images {
		_, err := r.db.Exec(ctx, q, motorcycleID, images[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MotorcycleRepository) CreateMotorcycleVideos(ctx *fasthttp.RequestCtx, motorcycleID int, video string) error {

	q := `
		INSERT INTO moto_videos (moto_id, video) VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, q, motorcycleID, video)
	if err != nil {
		return err
	}

	return err
}

func (r *MotorcycleRepository) DeleteMotorcycleImage(ctx *fasthttp.RequestCtx, motorcycleID int, imageID int) error {
	q := `
		DELETE FROM moto_images WHERE moto_id = $1 AND id = $2
	`

	_, err := r.db.Exec(ctx, q, motorcycleID, imageID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MotorcycleRepository) DeleteMotorcycleVideo(ctx *fasthttp.RequestCtx, motorcycleID int, videoID int) error {
	q := `
		DELETE FROM moto_videos WHERE moto_id = $1 AND id = $2
	`

	_, err := r.db.Exec(ctx, q, motorcycleID, videoID)
	if err != nil {
		return err
	}

	return nil
}

func (r *MotorcycleRepository) GetMotorcycleByID(ctx *fasthttp.RequestCtx, motorcycleID, userID int, nameColumn string) (model.GetMotorcyclesResponse, error) {
	var motorcycle model.GetMotorcyclesResponse
	q := `
		WITH updated AS (
			UPDATE motorcycles
			SET view_count = COALESCE(view_count, 0) + 1
			WHERE id = $1
			RETURNING *
		)
		select 
			mcs.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $3 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			mcs.engine,
			mcs.power,
			mcs.year,
			mcs.number_of_cycles,
			mcs.odometer,
			mcs.crash,
			mcs.not_cleared,
			mcs.owners,
			mcs.date_of_purchase,
			mcs.warranty_date,
			mcs.ptc,
			mcs.vin_code,
			mcs.certificate,
			mcs.description,
			mcs.can_look_coordinate,
			mcs.phone_number,
			mcs.refuse_dealers_calls,
			mcs.only_chat,
			mcs.protect_spam,
			mcs.verified_buyers,
			mcs.contact_person,
			mcs.email,
			mcs.price,
			mcs.price_type,
			mcs.status,
			mcs.updated_at,
			mcs.created_at,
			mocs.` + nameColumn + ` as moto_category,
			mbs.` + nameColumn + ` as moto_brand,
			mms.` + nameColumn + ` as moto_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN mcs.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from updated mcs
		left join profiles pf on pf.user_id = mcs.user_id
		left join moto_categories mocs on mocs.id = mcs.moto_category_id
		left join moto_brands mbs on mbs.id = mcs.moto_brand_id
		left join moto_models mms on mms.id = mcs.moto_model_id
		left join fuel_types fts on fts.id = mcs.fuel_type_id
		left join cities cs on cs.id = mcs.city_id
		left join colors cls on cls.id = mcs.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', mcp.moto_parameter_id,
					'parameter_value_id', mpv.id,
					'parameter', mp.` + nameColumn + `,
					'parameter_value', mpv.` + nameColumn + `
				)
			) AS parameters
			FROM motorcycle_parameters mcp
			LEFT JOIN moto_parameters mp ON mp.id = mcp.moto_parameter_id
			LEFT JOIN moto_parameter_values mpv ON mpv.id = mcp.moto_parameter_value_id
			WHERE mcp.motorcycle_id = mcs.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $3 || image as image
				FROM moto_images
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $3 || video as video
				FROM moto_videos
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		WHERE mcs.id = $1;
	`

	err := r.db.QueryRow(ctx, q, motorcycleID, userID, r.config.IMAGE_BASE_URL).Scan(
		&motorcycle.ID, &motorcycle.Owner, &motorcycle.Engine, &motorcycle.Power, &motorcycle.Year,
		&motorcycle.NumberOfCycles, &motorcycle.Odometer, &motorcycle.Crash, &motorcycle.NotCleared,
		&motorcycle.Owners, &motorcycle.DateOfPurchase, &motorcycle.WarrantyDate, &motorcycle.PTC,
		&motorcycle.VinCode, &motorcycle.Certificate, &motorcycle.Description, &motorcycle.CanLookCoordinate,
		&motorcycle.PhoneNumber, &motorcycle.RefuseDealersCalls, &motorcycle.OnlyChat,
		&motorcycle.ProtectSpam, &motorcycle.VerifiedBuyers, &motorcycle.ContactPerson,
		&motorcycle.Email, &motorcycle.Price, &motorcycle.PriceType, &motorcycle.Status,
		&motorcycle.UpdatedAt, &motorcycle.CreatedAt, &motorcycle.MotoCategory, &motorcycle.MotoBrand,
		&motorcycle.MotoModel, &motorcycle.FuelType, &motorcycle.City, &motorcycle.Color, &motorcycle.MyCar,
		&motorcycle.Parameters, &motorcycle.Images, &motorcycle.Videos)

	return motorcycle, err
}

func (r *MotorcycleRepository) GetEditMotorcycleByID(ctx *fasthttp.RequestCtx, motorcycleID, userID int, nameColumn string) (model.GetMotorcyclesResponse, error) {
	var motorcycle model.GetMotorcyclesResponse
	q := `
		select 
			mcs.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $3 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			mcs.engine,
			mcs.power,
			mcs.year,
			mcs.number_of_cycles,
			mcs.odometer,
			mcs.crash,
			mcs.not_cleared,
			mcs.owners,
			mcs.date_of_purchase,
			mcs.warranty_date,
			mcs.ptc,
			mcs.vin_code,
			mcs.certificate,
			mcs.description,
			mcs.can_look_coordinate,
			mcs.phone_number,
			mcs.refuse_dealers_calls,
			mcs.only_chat,
			mcs.protect_spam,
			mcs.verified_buyers,
			mcs.contact_person,
			mcs.email,
			mcs.price,
			mcs.price_type,
			mcs.status,
			mcs.updated_at,
			mcs.created_at,
			mocs.` + nameColumn + ` as moto_category,
			mbs.` + nameColumn + ` as moto_brand,
			mms.` + nameColumn + ` as moto_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN mcs.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from motorcycles mcs
		left join profiles pf on pf.user_id = mcs.user_id
		left join moto_categories mocs on mocs.id = mcs.moto_category_id
		left join moto_brands mbs on mbs.id = mcs.moto_brand_id
		left join moto_models mms on mms.id = mcs.moto_model_id
		left join fuel_types fts on fts.id = mcs.fuel_type_id
		left join cities cs on cs.id = mcs.city_id
		left join colors cls on cls.id = mcs.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', mcp.moto_parameter_id,
					'parameter_value_id', mpv.id,
					'parameter', mp.` + nameColumn + `,
					'parameter_value', mpv.` + nameColumn + `
				)
			) AS parameters
			FROM motorcycle_parameters mcp
			LEFT JOIN moto_parameters mp ON mp.id = mcp.moto_parameter_id
			LEFT JOIN moto_parameter_values mpv ON mpv.id = mcp.moto_parameter_value_id
			WHERE mcp.motorcycle_id = mcs.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $3 || image as image
				FROM moto_images
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $3 || video as video
				FROM moto_videos
				WHERE moto_id = mcs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		WHERE mcs.id = $1 AND mcs.user_id = $2;
	`

	err := r.db.QueryRow(ctx, q, motorcycleID, userID, r.config.IMAGE_BASE_URL).Scan(
		&motorcycle.ID, &motorcycle.Owner, &motorcycle.Engine, &motorcycle.Power, &motorcycle.Year,
		&motorcycle.NumberOfCycles, &motorcycle.Odometer, &motorcycle.Crash, &motorcycle.NotCleared,
		&motorcycle.Owners, &motorcycle.DateOfPurchase, &motorcycle.WarrantyDate, &motorcycle.PTC,
		&motorcycle.VinCode, &motorcycle.Certificate, &motorcycle.Description, &motorcycle.CanLookCoordinate,
		&motorcycle.PhoneNumber, &motorcycle.RefuseDealersCalls, &motorcycle.OnlyChat,
		&motorcycle.ProtectSpam, &motorcycle.VerifiedBuyers, &motorcycle.ContactPerson,
		&motorcycle.Email, &motorcycle.Price, &motorcycle.PriceType, &motorcycle.Status,
		&motorcycle.UpdatedAt, &motorcycle.CreatedAt, &motorcycle.MotoCategory, &motorcycle.MotoBrand,
		&motorcycle.MotoModel, &motorcycle.FuelType, &motorcycle.City, &motorcycle.Color, &motorcycle.MyCar,
		&motorcycle.Parameters, &motorcycle.Images, &motorcycle.Videos)

	return motorcycle, err
}

func (r *MotorcycleRepository) BuyMotorcycle(ctx *fasthttp.RequestCtx, motorcycleID, userID int) error {
	q := `
		UPDATE motorcycles 
		SET status = 2,
			user_id = $1
		WHERE id = $2 AND status = 3 -- 3 is on sale
	`

	_, err := r.db.Exec(ctx, q, userID, motorcycleID)
	return err
}

func (r *MotorcycleRepository) DontSellMotorcycle(ctx *fasthttp.RequestCtx, motorcycleID, userID int) error {
	q := `
		UPDATE motorcycles 
		SET status = 2 -- 2 is not sale
		WHERE id = $1 AND status = 3 -- 3 is on sale
			AND user_id = $2
	`

	_, err := r.db.Exec(ctx, q, motorcycleID, userID)
	return err
}

func (r *MotorcycleRepository) SellMotorcycle(ctx *fasthttp.RequestCtx, motorcycleID, userID int) error {
	q := `
		UPDATE motorcycles 
		SET status = 3 -- 3 is on sale
		WHERE id = $1 AND status = 2 -- 2 is not sale 
			AND user_id = $2
	`

	_, err := r.db.Exec(ctx, q, motorcycleID, userID)
	return err
}

func (r *MotorcycleRepository) DeleteMotorcycle(ctx *fasthttp.RequestCtx, motorcycleID int) error {
	q := `
		DELETE FROM motorcycles WHERE id = $1
	`

	_, err := r.db.Exec(ctx, q, motorcycleID)
	return err
}
