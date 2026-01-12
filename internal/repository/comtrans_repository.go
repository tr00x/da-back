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

type ComtransRepository struct {
	config *config.Config
	db     *pgxpool.Pool
}

func NewComtransRepository(config *config.Config, db *pgxpool.Pool) *ComtransRepository {
	return &ComtransRepository{config, db}
}

func (r *ComtransRepository) GetComtransCategories(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.GetComtransCategoriesResponse, error) {
	data := make([]model.GetComtransCategoriesResponse, 0)
	q := `
		SELECT id, ` + nameColumn + ` FROM com_categories
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category model.GetComtransCategoriesResponse
		err = rows.Scan(&category.ID, &category.Name)

		if err != nil {
			return nil, err
		}

		data = append(data, category)
	}

	return data, nil
}

func (r *ComtransRepository) GetComtransParameters(ctx *fasthttp.RequestCtx, categoryID string, nameColumn string) ([]model.GetComtransParametersResponse, error) {
	data := make([]model.GetComtransParametersResponse, 0)
	q := `
		SELECT 
			comtran_parameters.id,
			comtran_parameters.` + nameColumn + `,
			json_agg(
				json_build_object(
					'id', com_parameter_values.id,
					'name', com_parameter_values.` + nameColumn + `
				)
			) as values
		FROM comtran_parameters
		LEFT JOIN com_parameter_values ON comtran_parameters.id = com_parameter_values.comtran_parameter_id
		WHERE comtran_parameters.comtran_category_id = $1
		GROUP BY comtran_parameters.id, comtran_parameters.` + nameColumn + `
	`

	rows, err := r.db.Query(ctx, q, categoryID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var parameter model.GetComtransParametersResponse
		err = rows.Scan(&parameter.ID, &parameter.Name, &parameter.Values)

		if err != nil {
			return nil, err
		}

		data = append(data, parameter)
	}

	return data, nil
}

func (r *ComtransRepository) GetComtransBrands(ctx *fasthttp.RequestCtx, categoryID string, nameColumn string) ([]model.GetComtransBrandsResponse, error) {
	data := make([]model.GetComtransBrandsResponse, 0)
	q := `
		SELECT id, ` + nameColumn + `, $2 || image FROM com_brands
		WHERE comtran_category_id = $1
	`

	rows, err := r.db.Query(ctx, q, categoryID, r.config.IMAGE_BASE_URL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var brand model.GetComtransBrandsResponse
		err = rows.Scan(&brand.ID, &brand.Name, &brand.Image)

		if err != nil {
			return nil, err
		}

		data = append(data, brand)
	}

	return data, nil
}

func (r *ComtransRepository) GetComtransModelsByBrandID(ctx *fasthttp.RequestCtx, categoryID string, brandID string, nameColumn string) ([]model.GetComtransModelsResponse, error) {
	data := make([]model.GetComtransModelsResponse, 0)
	q := `
		SELECT id, ` + nameColumn + ` FROM com_models
		WHERE comtran_brand_id = $1
	`

	rows, err := r.db.Query(ctx, q, brandID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var model model.GetComtransModelsResponse
		err = rows.Scan(&model.ID, &model.Name)

		if err != nil {
			return nil, err
		}

		data = append(data, model)
	}

	return data, nil
}

func (r *ComtransRepository) CreateComtrans(ctx *fasthttp.RequestCtx, req model.CreateComtransRequest, userID int) (model.SuccessWithId, error) {
	data := model.SuccessWithId{}
	parameters := req.Parameters
	req.Parameters = nil

	keys, values, args := auth.BuildParams(req)

	q := `
		INSERT INTO comtrans ( 
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
			INSERT INTO comtran_parameters (comtran_id, comtran_parameter_id, comtran_parameter_value_id)
			VALUES ($1, $2, $3)
		`
		_, err = r.db.Exec(ctx, q, id, parameters[i].ComtransCategoryID, parameters[i].Name)

		if err != nil {
			return data, err
		}
	}

	data.Message = "Commercial transport created successfully"
	data.Id = id

	return data, err
}

func (r *ComtransRepository) GetComtrans(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.GetComtransResponse, error) {
	data := make([]model.GetComtransResponse, 0)
	q := `
		select 
			cts.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $1 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			cts.engine,
			cts.power,
			cts.year,
			cts.number_of_cycles,
			cts.odometer,
			cts.crash,
			cts.not_cleared,
			cts.owners,
			cts.date_of_purchase,
			cts.warranty_date,
			cts.ptc,
			cts.vin_code,
			cts.certificate,
			cts.description,
			cts.can_look_coordinate,
			cts.phone_number,
			cts.refuse_dealers_calls,
			cts.only_chat,
			cts.protect_spam,
			cts.verified_buyers,
			cts.contact_person,
			cts.email,
			cts.price,
			cts.price_type,
			cts.status,
			cts.updated_at,
			cts.created_at,
			cocs.` + nameColumn + ` as comtran_category,
			cbs.` + nameColumn + ` as comtran_brand,
			cms.` + nameColumn + ` as comtran_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN cts.user_id = 1 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from comtrans cts
		left join profiles pf on pf.user_id = cts.user_id
		left join com_categories cocs on cocs.id = cts.comtran_category_id
		left join com_brands cbs on cbs.id = cts.comtran_brand_id
		left join com_models cms on cms.id = cts.comtran_model_id
		left join fuel_types fts on fts.id = cts.fuel_type_id
		left join cities cs on cs.id = cts.city_id
		left join colors cls on cls.id = cts.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', ccp.comtran_parameter_id,
					'parameter_value_id', cpv.id,
					'parameter', cp.` + nameColumn + `,
					'parameter_value', cpv.` + nameColumn + `
				)
			) AS parameters
			FROM comtran_parameters ccp
			LEFT JOIN com_parameters cp ON cp.id = ccp.comtran_parameter_id
			LEFT JOIN com_parameter_values cpv ON cpv.id = ccp.comtran_parameter_value_id
			WHERE ccp.comtran_id = cts.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $1 || image as image
				FROM comtran_images
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $1 || video as video
				FROM comtran_videos
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) v
		) videos ON true;

	`

	rows, err := r.db.Query(ctx, q, r.config.IMAGE_BASE_URL)
	if err != nil {
		return data, err
	}

	defer rows.Close()

	for rows.Next() {
		var comtrans model.GetComtransResponse
		err = rows.Scan(
			&comtrans.ID, &comtrans.Owner, &comtrans.Engine, &comtrans.Power, &comtrans.Year,
			&comtrans.NumberOfCycles, &comtrans.Odometer, &comtrans.Crash, &comtrans.NotCleared,
			&comtrans.Owners, &comtrans.DateOfPurchase, &comtrans.WarrantyDate, &comtrans.PTC,
			&comtrans.VinCode, &comtrans.Certificate, &comtrans.Description, &comtrans.CanLookCoordinate,
			&comtrans.PhoneNumber, &comtrans.RefuseDealersCalls, &comtrans.OnlyChat,
			&comtrans.ProtectSpam, &comtrans.VerifiedBuyers, &comtrans.ContactPerson,
			&comtrans.Email, &comtrans.Price, &comtrans.PriceType, &comtrans.Status,
			&comtrans.UpdatedAt, &comtrans.CreatedAt, &comtrans.ComtranCategory, &comtrans.ComtranBrand,
			&comtrans.ComtranModel, &comtrans.FuelType, &comtrans.City, &comtrans.Color, &comtrans.MyCar,
			&comtrans.Parameters, &comtrans.Images, &comtrans.Videos)

		if err != nil {
			return data, err
		}

		data = append(data, comtrans)
	}

	return data, err
}

func (r *ComtransRepository) CreateComtransImages(ctx *fasthttp.RequestCtx, comtransID int, images []string) error {

	if len(images) == 0 {
		return nil
	}

	q := `
		INSERT INTO comtran_images (comtran_id, image) VALUES ($1, $2)
	`

	for i := range images {
		_, err := r.db.Exec(ctx, q, comtransID, images[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ComtransRepository) CreateComtransVideos(ctx *fasthttp.RequestCtx, comtransID int, video string) error {

	q := `
		INSERT INTO comtran_videos (comtran_id, video) VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, q, comtransID, video)
	if err != nil {
		return err
	}

	return err
}

func (r *ComtransRepository) DeleteComtransImage(ctx *fasthttp.RequestCtx, comtransID int, imageID int) error {
	q := `
		DELETE FROM comtran_images WHERE comtran_id = $1 AND id = $2
	`

	_, err := r.db.Exec(ctx, q, comtransID, imageID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ComtransRepository) DeleteComtransVideo(ctx *fasthttp.RequestCtx, comtransID int, videoID int) error {
	q := `
		DELETE FROM comtran_videos WHERE comtran_id = $1 AND id = $2
	`

	_, err := r.db.Exec(ctx, q, comtransID, videoID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ComtransRepository) GetComtransByID(ctx *fasthttp.RequestCtx, comtransID, userID int, nameColumn string) (model.GetComtransResponse, error) {
	var comtrans model.GetComtransResponse
	q := `
		WITH updated AS (
			UPDATE comtrans
			SET view_count = COALESCE(view_count, 0) + 1
			WHERE id = $1
			RETURNING *
		)
		select 
			cts.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $3 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			cts.engine,
			cts.power,
			cts.year,
			cts.number_of_cycles,
			cts.odometer,
			cts.crash,
			cts.not_cleared,
			cts.owners,
			cts.date_of_purchase,
			cts.warranty_date,
			cts.ptc,
			cts.vin_code,
			cts.certificate,
			cts.description,
			cts.can_look_coordinate,
			cts.phone_number,
			cts.refuse_dealers_calls,
			cts.only_chat,
			cts.protect_spam,
			cts.verified_buyers,
			cts.contact_person,
			cts.email,
			cts.price,
			cts.price_type,
			cts.status,
			cts.updated_at,
			cts.created_at,
			cocs.` + nameColumn + ` as comtran_category,
			cbs.name as comtran_brand,
			cms.` + nameColumn + ` as comtran_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN cts.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from updated cts
		left join profiles pf on pf.user_id = cts.user_id
		left join com_categories cocs on cocs.id = cts.comtran_category_id
		left join com_brands cbs on cbs.id = cts.comtran_brand_id
		left join com_models cms on cms.id = cts.comtran_model_id
		left join fuel_types fts on fts.id = cts.fuel_type_id
		left join cities cs on cs.id = cts.city_id
		left join colors cls on cls.id = cts.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', ccp.comtran_parameter_id,
					'parameter_value_id', cpv.id,
					'parameter', cp.` + nameColumn + `,
					'parameter_value', cpv.` + nameColumn + `
				)
			) AS parameters
			FROM comtran_parameters ccp
			LEFT JOIN com_parameters cp ON cp.id = ccp.comtran_parameter_id
			LEFT JOIN com_parameter_values cpv ON cpv.id = ccp.comtran_parameter_value_id
			WHERE ccp.comtran_id = cts.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $3 || image as image
				FROM comtran_images
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $3 || video as video
				FROM comtran_videos
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		WHERE cts.id = $1;
	`

	err := r.db.QueryRow(ctx, q, comtransID, userID, r.config.IMAGE_BASE_URL).Scan(
		&comtrans.ID, &comtrans.Owner, &comtrans.Engine, &comtrans.Power, &comtrans.Year,
		&comtrans.NumberOfCycles, &comtrans.Odometer, &comtrans.Crash, &comtrans.NotCleared,
		&comtrans.Owners, &comtrans.DateOfPurchase, &comtrans.WarrantyDate, &comtrans.PTC,
		&comtrans.VinCode, &comtrans.Certificate, &comtrans.Description, &comtrans.CanLookCoordinate,
		&comtrans.PhoneNumber, &comtrans.RefuseDealersCalls, &comtrans.OnlyChat,
		&comtrans.ProtectSpam, &comtrans.VerifiedBuyers, &comtrans.ContactPerson,
		&comtrans.Email, &comtrans.Price, &comtrans.PriceType, &comtrans.Status,
		&comtrans.UpdatedAt, &comtrans.CreatedAt, &comtrans.ComtranCategory, &comtrans.ComtranBrand,
		&comtrans.ComtranModel, &comtrans.FuelType, &comtrans.City, &comtrans.Color, &comtrans.MyCar,
		&comtrans.Parameters, &comtrans.Images, &comtrans.Videos)

	return comtrans, err
}

func (r *ComtransRepository) GetEditComtransByID(ctx *fasthttp.RequestCtx, comtransID, userID int, nameColumn string) (model.GetComtransResponse, error) {
	var comtrans model.GetComtransResponse
	q := `
		select 
			cts.id,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $3 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			cts.engine,
			cts.power,
			cts.year,
			cts.number_of_cycles,
			cts.odometer,
			cts.crash,
			cts.not_cleared,
			cts.owners,
			cts.date_of_purchase,
			cts.warranty_date,
			cts.ptc,
			cts.vin_code,
			cts.certificate,
			cts.description,
			cts.can_look_coordinate,
			cts.phone_number,
			cts.refuse_dealers_calls,
			cts.only_chat,
			cts.protect_spam,
			cts.verified_buyers,
			cts.contact_person,
			cts.email,
			cts.price,
			cts.price_type,
			cts.status,
			cts.updated_at,
			cts.created_at,
			cocs.` + nameColumn + ` as comtran_category,
			cbs.name as comtran_brand,
			cms.` + nameColumn + ` as comtran_model,
			fts.` + nameColumn + ` as fuel_type,
			cs.name as city,
			cls.` + nameColumn + ` as color,
			CASE
				WHEN cts.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			ps.parameters,
			images.images,
			videos.videos
		from comtrans cts
		left join profiles pf on pf.user_id = cts.user_id
		left join com_categories cocs on cocs.id = cts.comtran_category_id
		left join com_brands cbs on cbs.id = cts.comtran_brand_id
		left join com_models cms on cms.id = cts.comtran_model_id
		left join fuel_types fts on fts.id = cts.fuel_type_id
		left join cities cs on cs.id = cts.city_id
		left join colors cls on cls.id = cts.color_id
		LEFT JOIN LATERAL (
			SELECT json_agg(
				json_build_object(
					'parameter_id', ccp.comtran_parameter_id,
					'parameter_value_id', cpv.id,
					'parameter', cp.` + nameColumn + `,
					'parameter_value', cpv.` + nameColumn + `
				)
			) AS parameters
			FROM comtran_parameters ccp
			LEFT JOIN com_parameters cp ON cp.id = ccp.comtran_parameter_id
			LEFT JOIN com_parameter_values cpv ON cpv.id = ccp.comtran_parameter_value_id
			WHERE ccp.comtran_id = cts.id
		) ps ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $1 || image as image
				FROM comtran_images
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $1 || video as video
				FROM comtran_videos
				WHERE comtran_id = cts.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		WHERE cts.id = $1 AND cts.user_id = $2;
	`

	err := r.db.QueryRow(ctx, q, comtransID, userID, r.config.IMAGE_BASE_URL).Scan(
		&comtrans.ID, &comtrans.Owner, &comtrans.Engine, &comtrans.Power, &comtrans.Year,
		&comtrans.NumberOfCycles, &comtrans.Odometer, &comtrans.Crash, &comtrans.NotCleared,
		&comtrans.Owners, &comtrans.DateOfPurchase, &comtrans.WarrantyDate, &comtrans.PTC,
		&comtrans.VinCode, &comtrans.Certificate, &comtrans.Description, &comtrans.CanLookCoordinate,
		&comtrans.PhoneNumber, &comtrans.RefuseDealersCalls, &comtrans.OnlyChat,
		&comtrans.ProtectSpam, &comtrans.VerifiedBuyers, &comtrans.ContactPerson,
		&comtrans.Email, &comtrans.Price, &comtrans.PriceType, &comtrans.Status,
		&comtrans.UpdatedAt, &comtrans.CreatedAt, &comtrans.ComtranCategory, &comtrans.ComtranBrand,
		&comtrans.ComtranModel, &comtrans.FuelType, &comtrans.City, &comtrans.Color, &comtrans.MyCar,
		&comtrans.Parameters, &comtrans.Images, &comtrans.Videos)

	return comtrans, err
}

func (r *ComtransRepository) BuyComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) error {
	q := `
		UPDATE comtrans 
		SET status = 2,
			user_id = $1
		WHERE id = $2 AND status = 3 -- 3 is on sale
	`

	_, err := r.db.Exec(ctx, q, userID, comtransID)
	return err
}

func (r *ComtransRepository) DontSellComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) error {
	q := `
		UPDATE comtrans 
		SET status = 2 -- 2 is not sale
		WHERE id = $1 AND status = 3 -- 3 is on sale
			AND user_id = $2
	`

	_, err := r.db.Exec(ctx, q, comtransID, userID)
	return err
}

func (r *ComtransRepository) SellComtrans(ctx *fasthttp.RequestCtx, comtransID, userID int) error {
	q := `
		UPDATE comtrans 
		SET status = 3 -- 3 is on sale
		WHERE id = $1 AND status = 2 -- 2 is not sale 
			AND user_id = $2
	`

	_, err := r.db.Exec(ctx, q, comtransID, userID)
	return err
}

func (r *ComtransRepository) DeleteComtrans(ctx *fasthttp.RequestCtx, comtransID int) error {
	q := `
		DELETE FROM comtrans WHERE id = $1
	`

	_, err := r.db.Exec(ctx, q, comtransID)
	return err
}
