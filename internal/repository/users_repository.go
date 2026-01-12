package repository

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
	"dubai-auto/pkg/auth"
)

type UserRepository struct {
	config *config.Config
	db     *pgxpool.Pool
}

func NewUserRepository(config *config.Config, db *pgxpool.Pool) *UserRepository {
	return &UserRepository{config, db}
}

func (r *UserRepository) GetMyCars(ctx *fasthttp.RequestCtx, userID, limit, lastID, status int, nameColumn string) ([]model.GetMyCarsResponse, error) {
	cars := make([]model.GetMyCarsResponse, 0)
	q := `
		with vs as (        
			select 
				vs.id,
				'car' as type,
				bs.` + nameColumn + ` as brand,
				ms.` + nameColumn + ` as model,
				vs.year,
				vs.price,
				vs.credit,
				vs.status,
				vs.created_at,
				images.images,
				vs.view_count,
				true as my_car,
				vs.crash
			from vehicles vs
			left join brands bs on vs.brand_id = bs.id
			left join models ms on vs.model_id = ms.id
			LEFT JOIN LATERAL (
				SELECT json_agg(img.image) AS images
				FROM (
					SELECT $2 || image as image
					FROM images
					WHERE vehicle_id = vs.id
					ORDER BY created_at DESC
				) img
			) images ON true
			where vs.user_id = $1 and status = $3
			order by vs.id desc
		),
		cms as (
			select
				cm.id,
				'comtran' as type,
				cbs.` + nameColumn + ` as brand,
				cms.` + nameColumn + ` as model,
				cm.year,
				cm.price,
				cm.credit,
				cm.status,
				cm.created_at,
				images.images,
				cm.view_count,
				true as my_car,
				cm.crash
			from comtrans cm
			left join com_brands cbs on cbs.id = cm.comtran_brand_id
			left join com_models cms on cms.id = cm.comtran_model_id
			LEFT JOIN LATERAL (
				SELECT json_agg(img.image) AS images
				FROM (
					SELECT $2 || image as image
					FROM comtran_images
					WHERE comtran_id = cm.id
					ORDER BY created_at DESC
				) img
			) images ON true
			where cm.user_id = $1 and cm.status = $3
		),
		mts as (
			select
				mt.id,
				'motorcycle' as type,
				mbs.` + nameColumn + ` as brand,
				mms.` + nameColumn + ` as model,
				mt.year,
				mt.price,
				mt.credit,
				mt.status,
				mt.created_at,
				mt.view_count,
				images.images,
				true as my_car,
				mt.crash
			from motorcycles mt
			left join moto_brands mbs on mbs.id = mt.moto_brand_id
			left join moto_models mms on mms.id = mt.moto_model_id
			LEFT JOIN LATERAL (
				SELECT json_agg(img.image) AS images
				FROM (
					SELECT $2 || image as image
					FROM moto_images
					WHERE moto_id = mt.id
					ORDER BY created_at DESC
				) img
			) images ON true
			where mt.user_id = $1 and mt.status = $3
		)
		-- Union all three CTEs
		select 
			id, type, brand, model, 
			year, price, credit, 
			status, created_at, 
			view_count, images, my_car, 
			crash 
		from vs
		union all
		select 
			id, type, brand, model, 
			year, price, credit, 
			status, created_at, 
			view_count, images, my_car, 
			crash 
		from cms
		union all
		select 
			id, type, brand, model, 
			year, price, credit, 
			status, created_at, 
			view_count, images, my_car, 
			crash 
		from mts
		order by created_at desc;

	`
	rows, err := r.db.Query(ctx, q, userID, r.config.IMAGE_BASE_URL, status)

	if err != nil {
		return cars, err
	}
	defer rows.Close()

	for rows.Next() {
		var car model.GetMyCarsResponse
		if err := rows.Scan(
			&car.ID,
			&car.Type,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.Price,
			&car.Credit,
			&car.Status,
			&car.CreatedAt,
			&car.ViewCount,
			&car.Images,
			&car.MyCar,
			&car.Crash,
		); err != nil {
			return cars, err
		}
		cars = append(cars, car)
	}
	return cars, err
}

func (r *UserRepository) Cancel(ctx *fasthttp.RequestCtx, carID *int) error {
	q := `
		delete from vehicles where id = $1
	`
	_, err := r.db.Exec(ctx, q, *carID)
	return err
}

func (r *UserRepository) DeleteCar(ctx *fasthttp.RequestCtx, carID *int) error {
	q := `
		delete from vehicles where id = $1
	`
	_, err := r.db.Exec(ctx, q, *carID)
	return err
}

func (r *UserRepository) DontSell(ctx *fasthttp.RequestCtx, carID, userID *int) error {
	q := `
		update vehicles 
			set status = 2 -- 2 is not sale
		where id = $1 and status = 3 -- 3 is on sale
			and user_id = $2
	`

	_, err := r.db.Exec(ctx, q, *carID, *userID)
	return err
}

func (r *UserRepository) Sell(ctx *fasthttp.RequestCtx, carID, userID *int) error {
	q := `
		update vehicles 
			set status = 3 -- 3 is on sale
		where id = $1 and status = 2 -- 2 is not sale 
			and user_id = $2
	`
	_, err := r.db.Exec(ctx, q, *carID, *userID)
	return err
}

func (r *UserRepository) GetBrands(ctx *fasthttp.RequestCtx, text string, nameColumn string) ([]*model.GetBrandsResponse, error) {
	q := `
		SELECT 
			brands.id, 
			brands.` + nameColumn + ` as name, 
			$2 || logo, 
			count(m.id) as model_count 
		FROM brands 
		left join models m on m.brand_id = brands.id
		WHERE brands.` + nameColumn + ` ILIKE '%' || $1 || '%'
		group by brands.id
	`
	rows, err := r.db.Query(ctx, q, text, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var brands = make([]*model.GetBrandsResponse, 0)

	for rows.Next() {
		var brand model.GetBrandsResponse

		if err := rows.Scan(&brand.ID, &brand.Name, &brand.Logo, &brand.ModelCount); err != nil {
			return nil, err
		}
		brands = append(brands, &brand)
	}
	return brands, err
}

func (r *UserRepository) GetProfile(ctx *fasthttp.RequestCtx, userID int, nameColumn string) (model.GetProfileResponse, error) {
	q := `
		select 
			us.id,
			us.email,
			us.phone,
			ps.driving_experience,
			ps.notification,
			ps.username,
			ps.google,
			ps.birthday,
			ps.about_me,
			ps.registered_by,
			ps.contacts,
			ps.address,
			json_build_object(
				'id', cs.id,
				'name', cs.` + nameColumn + `
			) as city
		from users us
		left join profiles as ps on ps.user_id = us.id
		left join cities as cs on cs.id = ps.city_id
		where us.id = $1;

	`
	var pf model.GetProfileResponse
	var contactsJSON []byte
	err := r.db.QueryRow(ctx, q, userID).Scan(&pf.ID, &pf.Email, &pf.Phone,
		&pf.DrivingExperience, &pf.Notification, &pf.Username, &pf.Google, &pf.Birthday, &pf.AboutMe, &pf.RegisteredBy, &contactsJSON, &pf.Address, &pf.City)

	if err == nil && len(contactsJSON) > 0 {
		if err := json.Unmarshal(contactsJSON, &pf.Contacts); err != nil {
			return pf, err
		}
	}

	return pf, err
}

func (r *UserRepository) UpdateProfile(ctx *fasthttp.RequestCtx, userID int, profile *model.UpdateProfileRequest) error {

	q := `
	UPDATE users 
	SET username = $2, 
	phone = $3, email = $4
	WHERE id = $1
	`
	_, err := r.db.Exec(ctx, q, userID, profile.Username, profile.PhoneNumber, profile.Email)

	if err != nil {
		return err
	}

	profile.PhoneNumber = ""
	profile.Email = ""

	keys, _, args := auth.BuildParams(profile)
	// Handle contacts map separately - BuildParams will skip maps, so we add it manually
	var contactsJSON []byte
	var hasContacts bool
	if profile.Contacts != nil {
		var err error
		contactsJSON, err = json.Marshal(profile.Contacts)

		if err != nil {
			return err
		}
		hasContacts = true
	}

	// Remove contacts from keys/args if BuildParams somehow included it (shouldn't happen, but be safe)
	for i := len(keys) - 1; i >= 0; i-- {
		if keys[i] == "contacts" {
			keys = append(keys[:i], keys[i+1:]...)
			args = append(args[:i], args[i+1:]...)
		}
	}

	// Add contacts if provided
	if hasContacts {
		keys = append(keys, "contacts")
		args = append(args, contactsJSON)
	}

	if len(keys) == 0 {
		return nil // No fields to update
	}

	// Build dynamic SET clause
	var setClause []string

	for i, key := range keys {
		setClause = append(setClause, fmt.Sprintf("%s = $%d", key, i+1))
	}

	setClause = append(setClause, "last_active_date = NOW()")
	args = append(args, userID)

	q = fmt.Sprintf(`
		UPDATE profiles 
		SET %s
		WHERE user_id = $%d
	`, strings.Join(setClause, ", "), len(args))

	_, err = r.db.Exec(ctx, q, args...)

	return err
}

func (r *UserRepository) GetFilterBrands(ctx *fasthttp.RequestCtx, text string, nameColumn string) (model.GetFilterBrandsResponse, error) {
	var brand model.GetFilterBrandsResponse
	q := `
		with popular as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `, 
						'logo', $2 || logo, 
						'model_count', model_count 
					)
				) as popular_brands
			FROM brands 
			WHERE name ILIKE '%' || $1 || '%' and popular = true
		), all_brands as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `, 
						'logo', $2 || logo, 
						'model_count', model_count 
					)
				) as all_brands
			FROM brands 
			WHERE name ILIKE '%' || $1 || '%'
		)
		select 
			pp.popular_brands,
			ab.all_brands
		from popular as pp
		left join all_brands as ab on true;

	`
	err := r.db.QueryRow(ctx, q, text, r.config.IMAGE_BASE_URL).Scan(&brand.PopularBrands, &brand.AllBrands)

	return brand, err
}

func (r *UserRepository) GetCities(ctx *fasthttp.RequestCtx, text, nameColumn string) ([]model.GetCitiesResponse, error) {
	var cities = make([]model.GetCitiesResponse, 0)
	var city model.GetCitiesResponse
	q := `
		select 
			c.id, 
			c.` + nameColumn + ` as name,
			json_agg(
				json_build_object(
					'id', r.id,
					'name', r.` + nameColumn + `
				)
			) as regions
		from cities c
		left join regions r on r.city_id = c.id
		where c.` + nameColumn + ` ILIKE '%' || $1 || '%'
		group by c.id, c.` + nameColumn + `;
	`
	rows, err := r.db.Query(ctx, q, text)

	if err != nil {
		return cities, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&city.ID, &city.Name, &city.Regions)

		if err != nil {
			return cities, err
		}
		cities = append(cities, city)
	}
	return cities, err
}

func (r *UserRepository) GetModelsByBrandID(ctx *fasthttp.RequestCtx, brandID int64, text string, nameColumn string) ([]model.Model, error) {
	q := `
			SELECT 
				id, 
				` + nameColumn + ` as name
			FROM models 
			WHERE 
				brand_id = $1 AND 
				name ILIKE '%' || $2 || '%'
		`
	rows, err := r.db.Query(ctx, q, brandID, text)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	models := make([]model.Model, 0)

	for rows.Next() {
		var model model.Model

		if err := rows.Scan(&model.ID, &model.Name); err != nil {
			return nil, err
		}
		models = append(models, model)
	}

	return models, err
}

func (r *UserRepository) GetFilterModelsByBrandID(ctx *fasthttp.RequestCtx, brandID int64, text, nameColumn string) (model.GetFilterModelsResponse, error) {
	responseModel := model.GetFilterModelsResponse{}
	q := `
		with popular as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `
					)
				) as popular_models
			FROM models 
			WHERE brand_id = $1 AND (name ILIKE '%' || $2 || '%' or name_ru ILIKE '%' || $2 || '%') and popular = true
		), all_models as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `
					)
				) as all_models
			FROM models 
			WHERE brand_id = $1 AND (name ILIKE '%' || $2 || '%' or name_ru ILIKE '%' || $2 || '%')
		)
		select 
			pp.popular_models,
			ms.all_models
		from popular as pp
		left join all_models as ms on true;
		`
	err := r.db.QueryRow(ctx, q, brandID, text).Scan(&responseModel.PopularModels, &responseModel.AllModels)

	if err != nil {
		return responseModel, err
	}

	return responseModel, err
}

func (r *UserRepository) GetFilterModelsByBrands(ctx *fasthttp.RequestCtx, brands []int, text, nameColumn string) (model.GetFilterModelsResponse, error) {
	responseModel := model.GetFilterModelsResponse{}
	q := `
		with popular as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `
					)
				) as popular_models
			FROM models 
			WHERE brand_id = any ($1) AND (name ILIKE '%' || $2 || '%' or name_ru ILIKE '%' || $2 || '%') and popular = true
		), all_models as (
			SELECT 
				json_agg(
					json_build_object(
						'id', id, 
						'name', ` + nameColumn + `
					)
				) as all_models
			FROM models 
			WHERE brand_id = any ($1) AND (name ILIKE '%' || $2 || '%' or name_ru ILIKE '%' || $2 || '%')
		)
		select 
			pp.popular_models,
			ms.all_models
		from popular as pp
		left join all_models as ms on true;
	`
	err := r.db.QueryRow(ctx, q, brands, text).Scan(&responseModel.PopularModels, &responseModel.AllModels)

	if err != nil {
		return responseModel, err
	}

	return responseModel, err
}

func (r *UserRepository) GetGenerationsByModelID(ctx *fasthttp.RequestCtx, modelID int, wheel bool, year, bodyTypeID, nameColumn string) ([]model.Generation, error) {
	q := `
		select
			gs.id,
			gs.` + nameColumn + ` as name,
			$5 || gs.image,
			gs.start_year,
			gs.end_year,
			json_agg(
				json_build_object(
					'id', gms.id,
					'engine', es.` + nameColumn + `, 
					'fuel_type', fts.` + nameColumn + `, 
					'drivetrain', ds.` + nameColumn + `, 
					'transmission', ts.` + nameColumn + `
				)
			) as modifications
		from generations gs
		join generation_modifications gms on gms.generation_id = gs.id
		left join engines es on es.id = gms.engine_id
		left join fuel_types fts on fts.id = gms.fuel_type_id
		left join drivetrains ds on ds.id = gms.drivetrain_id
		left join transmissions ts on ts.id = gms.transmission_id
		where 
			gs.model_id = $1 and gs.start_year <= $2 and gs.end_year >= $2
			and gms.body_type_id = $3 and gs.wheel = $4
		group by gs.id;
	`
	rows, err := r.db.Query(ctx, q, modelID, year, bodyTypeID, wheel, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	generations := make([]model.Generation, 0)

	for rows.Next() {
		var generation model.Generation
		if err = rows.Scan(&generation.ID, &generation.Name,
			&generation.Image, &generation.StartYear, &generation.EndYear,
			&generation.Modifications,
		); err != nil {
			return nil, err
		}
		generations = append(generations, generation)
	}
	return generations, err
}

func (r *UserRepository) GetGenerationsByModels(ctx *fasthttp.RequestCtx, models []int, nameColumn string) ([]model.Generation, error) {

	q := `
		select
			gs.id,
			gs.` + nameColumn + ` as name,
			$2 || gs.image,
			gs.start_year,
			gs.end_year,
			json_agg(
				json_build_object(
					'id', gms.id,
					'engine', es.` + nameColumn + `, 
					'fuel_type', fts.` + nameColumn + `, 
					'drivetrain', ds.` + nameColumn + `, 
					'transmission', ts.` + nameColumn + `
				)
			) as modifications
		from generations gs
		join generation_modifications gms on gms.generation_id = gs.id
		left join engines es on es.id = gms.engine_id
		left join fuel_types fts on fts.id = gms.fuel_type_id
		left join drivetrains ds on ds.id = gms.drivetrain_id
		left join transmissions ts on ts.id = gms.transmission_id
		where 
			gs.model_id = any ($1)
		group by gs.id;
	`
	rows, err := r.db.Query(ctx, q, models, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	generations := make([]model.Generation, 0)

	for rows.Next() {
		var generation model.Generation
		if err = rows.Scan(&generation.ID, &generation.Name,
			&generation.Image, &generation.StartYear, &generation.EndYear,
			&generation.Modifications,
		); err != nil {
			return nil, err
		}
		generations = append(generations, generation)
	}
	return generations, err
}

// todo: after remove the array response, return an object
func (r *UserRepository) GetYearsByModelID(ctx *fasthttp.RequestCtx, modelID int64, wheel bool) ([]*int, error) {
	q := `
		SELECT 
			array_agg(y ORDER BY y) AS years
		FROM (
			SELECT DISTINCT generate_series(start_year, end_year) AS y
			FROM generations
			WHERE model_id = $1 AND wheel = $2
		) AS years_series;
	`
	var years []*int
	err := r.db.QueryRow(ctx, q, modelID, wheel).Scan(&years)

	return years, err
}

// todo: after remove the array response, return an object
func (r *UserRepository) GetYearsByModels(ctx *fasthttp.RequestCtx, models []int, wheel bool) ([]*int, error) {
	q := `
		SELECT 
			array_agg(y ORDER BY y) AS years
		FROM (
			SELECT DISTINCT generate_series(start_year, end_year) AS y
			FROM generations
			WHERE model_id = any ($1) AND wheel = $2
		) AS years_series;
	`
	var years []*int
	err := r.db.QueryRow(ctx, q, models, wheel).Scan(&years)

	return years, err
}

func (r *UserRepository) GetBodysByModelID(ctx *fasthttp.RequestCtx, modelID int, wheel bool, year string, nameColumn string) ([]model.BodyType, error) {
	q := `
		select DISTINCT ON (bts.id)
			bts.id,
			bts.` + nameColumn + ` as name,
			$4 || bts.image
		from generation_modifications gms
		left join body_types bts on bts.id = gms.body_type_id
		where gms.generation_id in (
			select 
				gs.id 
			from generations gs 
			WHERE gs.start_year <= $1 AND gs.end_year >= $1 
					and gs.model_id = $2 and gs.wheel = $3
		)
	`

	rows, err := r.db.Query(ctx, q, year, modelID, wheel, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	bodyTypes := make([]model.BodyType, 0)

	for rows.Next() {
		var bodyType model.BodyType

		if err := rows.Scan(&bodyType.ID, &bodyType.Name, &bodyType.Image); err != nil {
			return nil, err
		}

		bodyTypes = append(bodyTypes, bodyType)
	}

	return bodyTypes, err
}

func (r *UserRepository) GetBodyTypes(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.BodyType, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name, $1 || image as image FROM body_types
	`

	rows, err := r.db.Query(ctx, q, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	bodyTypes := make([]model.BodyType, 0)

	for rows.Next() {
		var bodyType model.BodyType

		if err := rows.Scan(&bodyType.ID, &bodyType.Name, &bodyType.Image); err != nil {
			return nil, err
		}

		bodyTypes = append(bodyTypes, bodyType)
	}
	return bodyTypes, err
}

func (r *UserRepository) GetTransmissions(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.Transmission, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name FROM transmissions
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	transmissions := make([]model.Transmission, 0)

	for rows.Next() {
		var transmission model.Transmission
		if err := rows.Scan(&transmission.ID, &transmission.Name); err != nil {
			return nil, err
		}
		transmissions = append(transmissions, transmission)
	}
	return transmissions, err
}

func (r *UserRepository) GetEngines(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.Engine, error) {
	q := `
		SELECT id, ` + nameColumn + ` FROM engines
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	engines := make([]model.Engine, 0)

	for rows.Next() {
		var engine model.Engine
		if err := rows.Scan(&engine.ID, &engine.Value); err != nil {
			return nil, err
		}
		engines = append(engines, engine)
	}
	return engines, err
}

func (r *UserRepository) GetDrivetrains(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.Drivetrain, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name FROM drivetrains
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	drivetrains := make([]model.Drivetrain, 0)

	for rows.Next() {
		var drivetrain model.Drivetrain
		if err := rows.Scan(&drivetrain.ID, &drivetrain.Name); err != nil {
			return nil, err
		}
		drivetrains = append(drivetrains, drivetrain)
	}
	return drivetrains, err
}

func (r *UserRepository) GetFuelTypes(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.FuelType, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name FROM fuel_types
	`

	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	fuelTypes := make([]model.FuelType, 0)

	for rows.Next() {
		var fuelType model.FuelType
		if err := rows.Scan(&fuelType.ID, &fuelType.Name); err != nil {
			return nil, err
		}
		fuelTypes = append(fuelTypes, fuelType)
	}
	return fuelTypes, err
}

func (r *UserRepository) GetColors(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.Color, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name, $1 || image as image FROM colors
	`

	rows, err := r.db.Query(ctx, q, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	colors := make([]model.Color, 0)

	for rows.Next() {
		var color model.Color

		if err := rows.Scan(&color.ID, &color.Name, &color.Image); err != nil {
			return nil, err
		}
		colors = append(colors, color)
	}
	return colors, err
}

func (r *UserRepository) GetCountries(ctx *fasthttp.RequestCtx, nameColumn string) ([]model.Country, error) {
	q := `
		SELECT id, ` + nameColumn + ` as name, country_code, $1 || flag FROM countries
	`

	rows, err := r.db.Query(ctx, q, r.config.IMAGE_BASE_URL)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	countries := make([]model.Country, 0)

	for rows.Next() {
		var country model.Country

		if err := rows.Scan(&country.ID, &country.Name, &country.CountryCode, &country.Flag); err != nil {
			return nil, err
		}

		countries = append(countries, country)
	}
	return countries, err
}

func (r *UserRepository) GetHome(ctx *fasthttp.RequestCtx, userID int, nameColumn string) (model.Home, error) {
	home := model.Home{}
	cars := make([]model.GetCarsResponse, 0)

	q := `
		select 
			vs.id,
			bs.` + nameColumn + ` as brand,
			rs.` + nameColumn + ` as region,
			cs.` + nameColumn + ` as city,
			cls.` + nameColumn + ` as color,
			ms.` + nameColumn + ` as model,
			ts.` + nameColumn + ` as transmission,
			es.` + nameColumn + ` as engine,
			ds.` + nameColumn + ` as drive,
			bts.` + nameColumn + ` as body_type,
			fts.` + nameColumn + ` as fuel_type,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			vs.updated_at,
			images,
			vs.phone_numbers,
			vs.view_count,
			CASE
				WHEN vs.user_id = $1 THEN TRUE
				ELSE FALSE
			END AS my_car
		from vehicles vs
		left join generation_modifications gms on gms.id = vs.modification_id
		left join colors cls on vs.color_id = cls.id
		left join brands bs on vs.brand_id = bs.id
		left join regions rs on vs.region_id = rs.id
		left join cities cs on vs.city_id = cs.id
		left join models ms on vs.model_id = ms.id
		left join transmissions ts on gms.transmission_id = ts.id
		left join engines es on gms.engine_id = es.id
		left join drivetrains ds on gms.drivetrain_id = ds.id
		left join body_types bts on gms.body_type_id = bts.id
		left join fuel_types fts on gms.fuel_type_id = fts.id
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $2 || image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		where vs.status = 3
		order by vs.id desc limit 4
	`

	rows, err := r.db.Query(ctx, q, userID, r.config.IMAGE_BASE_URL)

	if err != nil {
		return home, err
	}

	defer rows.Close()

	for rows.Next() {
		var car model.GetCarsResponse

		if err := rows.Scan(
			&car.ID, &car.Brand, &car.Region, &car.City, &car.Color, &car.Model, &car.Transmission, &car.Engine,
			&car.Drivetrain, &car.BodyType, &car.FuelType, &car.Year, &car.Price, &car.Mileage, &car.VinCode,
			&car.Credit, &car.New, &car.Status, &car.CreatedAt,
			&car.UpdatedAt, &car.Images, &car.PhoneNumbers, &car.ViewCount, &car.MyCar,
		); err != nil {
			return home, err
		}

		cars = append(cars, car)
	}
	home.Popular = cars
	return home, nil
}

func (r *UserRepository) GetCars(ctx *fasthttp.RequestCtx, userID int,
	targetUserID string,
	brands, models, regions, cities, generations, transmissions,
	engines, drivetrains, body_types, fuel_types, ownership_types, colors []string,
	year_from, year_to, credit, price_from, price_to, tradeIn, owners,
	crash, odometer string, new, wheel *bool, limit, lastID int, nameColumn string) ([]model.GetCarsResponse, error) {
	var qWhere string
	var qValues []any
	qValues = append(qValues, userID)
	var i = 1

	if len(brands) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND bs.id = ANY($%d)", i)
		qValues = append(qValues, brands)
	}

	if len(models) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND ms.id = ANY($%d)", i)
		qValues = append(qValues, models)
	}

	if len(regions) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND rs.id = ANY($%d)", i)
		qValues = append(qValues, regions)
	}

	if len(cities) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND cs.id = ANY($%d)", i)
		qValues = append(qValues, cities)
	}

	if len(transmissions) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND ts.id = ANY($%d)", i)
		qValues = append(qValues, transmissions)
	}

	if len(engines) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND es.id = ANY($%d)", i)
		qValues = append(qValues, engines)
	}

	if len(drivetrains) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND ds.id = ANY($%d)", i)
		qValues = append(qValues, drivetrains)
	}

	if len(body_types) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND bts.id = ANY($%d)", i)
		qValues = append(qValues, body_types)
	}

	if len(fuel_types) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND fts.id = ANY($%d)", i)
		qValues = append(qValues, fuel_types)
	}

	if len(generations) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND gms.generation_id = ANY($%d)", i)
		qValues = append(qValues, generations)
	}

	if len(colors) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.color_id = ANY($%d)", i)
		qValues = append(qValues, colors)
	}

	if len(ownership_types) > 0 {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.ownership_type_id = ANY($%d) ", i)
		qValues = append(qValues, ownership_types)
	}

	if year_from != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.year >= $%d", i)
		qValues = append(qValues, year_from)
	}

	if year_to != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.year <= $%d", i)
		qValues = append(qValues, year_to)
	}

	if tradeIn != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.trade_in = $%d", i)
		qValues = append(qValues, tradeIn)
	}

	if owners != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.owners = $%d", i)
		qValues = append(qValues, owners)
	}

	if crash != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.crash = $%d", i)
		qValues = append(qValues, crash)
	}

	if credit != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.credit = $%d", i)
		qValues = append(qValues, true)
	}

	if credit != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.credit = $%d", i)
		qValues = append(qValues, true)
	}

	if price_from != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.price >= $%d", i)
		qValues = append(qValues, price_from)
	}

	if price_to != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.price <= $%d", i)
		qValues = append(qValues, price_to)
	}

	if new != nil {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.new = $%d", i)
		qValues = append(qValues, new)
	}

	if wheel != nil {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.wheel = $%d", i)
		qValues = append(qValues, wheel)
	}

	if odometer != "" {
		i += 1
		qWhere += fmt.Sprintf(" AND vs.odometer <= $%d", i)
		qValues = append(qValues, odometer)
	}

	if targetUserID != "" {
		targetUserIDInt, err := strconv.Atoi(targetUserID)

		if err != nil {
			return nil, err
		}

		i += 1
		qWhere += fmt.Sprintf(" AND vs.user_id = $%d", i)
		qValues = append(qValues, targetUserIDInt)
	}

	cars := make([]model.GetCarsResponse, 0)
	q := `
		select 
			vs.id,
			bs.` + nameColumn + ` as brand,
			rs.` + nameColumn + ` as region,
			cs.` + nameColumn + ` as city,
			cls.` + nameColumn + ` as color,
			ms.` + nameColumn + ` as model,
			ts.` + nameColumn + ` as transmission,
			es.` + nameColumn + ` as engine,
			ds.` + nameColumn + ` as drive,
			bts.` + nameColumn + ` as body_type,
			fts.` + nameColumn + ` as fuel_type,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			vs.trade_in,
			vs.owners,
			vs.crash,
			vs.updated_at,
			images.images,
			videos.videos,
			vs.phone_numbers,
			vs.view_count,
			CASE
				WHEN vs.user_id = $1 THEN TRUE
				ELSE FALSE
			END AS my_car,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', '` + r.config.IMAGE_BASE_URL + `' || pf.avatar,
				'role_id', u.role_id,
				'contacts', pf.contacts
			) as owner,
			vs.description,
			CASE 
				WHEN ul.vehicle_id IS NOT NULL THEN true
				ELSE false
			END AS liked
		from vehicles vs
		left join generation_modifications gms on gms.id = vs.modification_id
		left join colors cls on vs.color_id = cls.id
		left join profiles pf on pf.user_id = vs.user_id
		left join users u on u.id = vs.user_id
		left join brands bs on vs.brand_id = bs.id
		left join regions rs on vs.region_id = rs.id
		left join cities cs on vs.city_id = cs.id
		left join models ms on vs.model_id = ms.id
		left join transmissions ts on gms.transmission_id = ts.id
		left join engines es on gms.engine_id = es.id
		left join drivetrains ds on gms.drivetrain_id = ds.id
		left join body_types bts on gms.body_type_id = bts.id
		left join fuel_types fts on gms.fuel_type_id = fts.id
		left join user_likes ul on ul.vehicle_id = vs.id AND ul.user_id = $1
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || video as video
				FROM videos
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		where vs.status = 3 and vs.id > ` + strconv.Itoa(lastID) + `
		` + qWhere + `
		order by vs.id desc
		limit ` + strconv.Itoa(limit) + `
	`
	rows, err := r.db.Query(ctx, q, qValues...)

	if err != nil {
		return cars, err
	}

	defer rows.Close()
	for rows.Next() {
		var car model.GetCarsResponse
		if err := rows.Scan(
			&car.ID, &car.Brand, &car.Region, &car.City, &car.Color, &car.Model, &car.Transmission, &car.Engine,
			&car.Drivetrain, &car.BodyType, &car.FuelType, &car.Year, &car.Price, &car.Mileage, &car.VinCode,
			&car.Credit, &car.New, &car.Status, &car.CreatedAt, &car.TradeIn, &car.Owners, &car.Crash,
			&car.UpdatedAt, &car.Images, &car.Videos, &car.PhoneNumbers, &car.ViewCount, &car.MyCar, &car.Owner, &car.Description, &car.Liked,
		); err != nil {
			return cars, err
		}
		cars = append(cars, car)
	}
	return cars, err
}

func (r *UserRepository) GetPriceRecommendation(ctx *fasthttp.RequestCtx, filter model.GetPriceRecommendationRequest) ([]int, error) {
	keys, _, args := auth.BuildParams(filter)
	qWhere := ""

	for i, key := range keys {
		if i == 0 {
			qWhere += fmt.Sprintf("vs.%s = $%d", key, i+1)
		} else {
			qWhere += fmt.Sprintf(" AND vs.%s = $%d", key, i+1)
		}
	}

	q := `
		with prices as (
			select 
				price
			from vehicles vs 
			where 
				` + qWhere + `
			order by updated_at desc
			limit 10
		),
		ordered as (
			select 
			price
			from prices
			order by price desc
		)
		select 
			json_agg(price) as prices
		from ordered;
	`

	var prices []int
	err := r.db.QueryRow(ctx, q, args...).Scan(&prices)
	return prices, err
}

func (r *UserRepository) GetCarByID(ctx *fasthttp.RequestCtx, carID, userID int, nameColumn string) (model.GetCarsResponse, error) {
	car := model.GetCarsResponse{}
	q := `
		WITH updated AS (
			UPDATE vehicles
			SET view_count = view_count + 1
			WHERE id = $1
			RETURNING *
		)
		SELECT 
			vs.id,
			bs.` + nameColumn + ` as brand,
			rs.` + nameColumn + ` as region,
			cs.` + nameColumn + ` as city,
			cls.` + nameColumn + ` as color,
			ms.` + nameColumn + ` as model,
			ts.` + nameColumn + ` as transmission,
			es.` + nameColumn + ` as engine,
			ds.` + nameColumn + ` as drive,
			bts.` + nameColumn + ` as body_type,
			fts.` + nameColumn + ` as fuel_type,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			vs.trade_in,
			vs.owners,
			vs.updated_at,
			images,
			videos,
			vs.phone_numbers,
			vs.view_count,
			CASE
				WHEN vs.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', $3 || pf.avatar,
				'contacts', pf.contacts
			) as owner,
			 vs.description,
			CASE 
				WHEN ul.vehicle_id IS NOT NULL THEN true
				ELSE false
			END AS liked
		FROM updated vs
		LEFT JOIN generation_modifications gms ON gms.id = vs.modification_id
		LEFT JOIN profiles pf on pf.user_id = vs.user_id
		LEFT JOIN colors cls ON vs.color_id = cls.id
		LEFT JOIN brands bs ON vs.brand_id = bs.id
		LEFT JOIN regions rs ON vs.region_id = rs.id
		LEFT JOIN cities cs ON vs.city_id = cs.id
		LEFT JOIN models ms ON vs.model_id = ms.id
		LEFT JOIN transmissions ts ON gms.transmission_id = ts.id
		LEFT JOIN engines es ON gms.engine_id = es.id
		LEFT JOIN drivetrains ds ON gms.drivetrain_id = ds.id
		LEFT JOIN body_types bts ON gms.body_type_id = bts.id
		LEFT JOIN fuel_types fts ON gms.fuel_type_id = fts.id
		left join user_likes ul on ul.vehicle_id = vs.id AND ul.user_id = $2
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || video as video
				FROM videos
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		WHERE vs.id = $1;
	`

	err := r.db.QueryRow(ctx, q, carID, userID, r.config.IMAGE_BASE_URL).Scan(
		&car.ID, &car.Brand, &car.Region, &car.City, &car.Color, &car.Model, &car.Transmission, &car.Engine,
		&car.Drivetrain, &car.BodyType, &car.FuelType, &car.Year, &car.Price, &car.Mileage, &car.VinCode,
		&car.Credit, &car.New, &car.Status, &car.CreatedAt, &car.TradeIn, &car.Owners,
		&car.UpdatedAt, &car.Images, &car.Videos, &car.PhoneNumbers, &car.ViewCount, &car.MyCar, &car.Owner, &car.Description, &car.Liked,
	)

	return car, err
}

func (r *UserRepository) GetEditCarByID(ctx *fasthttp.RequestCtx, carID, userID int, nameColumn string) (model.GetEditCarsResponse, error) {
	car := model.GetEditCarsResponse{}

	q := `
		select 
			vs.id,
			json_build_object(
				'id', bs.id,
				'name', bs.` + nameColumn + `,
				'logo', $3 || bs.logo,
				'model_count', bs.model_count
			) as brand,
			json_build_object(
				'id', rs.id,
				'name', rs.` + nameColumn + `
			) as region,
			json_build_object(
				'id', cs.id,
				'name', cs.` + nameColumn + `
			) as city,
			json_build_object(
				'id', ms.id,
				'name', ms.` + nameColumn + `
			) as model,
			json_build_object(
				'id', mfs.id,
				'engine', es.` + nameColumn + `,
				'fuel_type', fts.` + nameColumn + `,
				'drivetrain', ds.` + nameColumn + `,
				'transmission', ts.` + nameColumn + `
			) as modification,
			json_build_object(
				'id', cls.id,
				'name', cls.` + nameColumn + `,
				'image', cls.image
			) as color,
			json_build_object(
				'id', bts.id,
				'name', bts.` + nameColumn + `,
				'image', bts.image
			) as body_type,
			json_build_object(
				'id', gs.id,
				'name', gs.` + nameColumn + `,
				'image', gs.image,
				'start_year', gs.start_year,
				'end_year', gs.end_year
			) as generation,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.wheel,
			vs.trade_in,
			vs.crash,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			images,
			videos,
			vs.phone_numbers,
			vs.view_count,
			vs.description,
			CASE
				WHEN vs.user_id = $2 THEN TRUE
				ELSE FALSE
			END AS my_car,
			vs.owners
		from vehicles vs
		left join colors cls on vs.color_id = cls.id
		left join generation_modifications mfs on mfs.id = vs.modification_id
		left join generations gs on gs.id = mfs.generation_id
		left join body_types bts on bts.id = mfs.body_type_id
		left join engines es on es.id = mfs.engine_id
		left join transmissions ts on ts.id = mfs.transmission_id
		left join drivetrains ds on ds.id = mfs.drivetrain_id
		left join fuel_types fts on fts.id = mfs.fuel_type_id
		left join brands bs on vs.brand_id = bs.id
		left join regions rs on vs.region_id = rs.id
		left join cities cs on vs.city_id = cs.id
		left join models ms on vs.model_id = ms.id
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT '` + r.config.IMAGE_BASE_URL + `' || video as video
				FROM videos
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		where vs.id = $1 and vs.user_id = $2;
	`
	err := r.db.QueryRow(ctx, q, carID, userID, r.config.IMAGE_BASE_URL).Scan(
		&car.ID, &car.Brand, &car.Region, &car.City, &car.Model, &car.Modification,
		&car.Color, &car.BodyType, &car.Generation, &car.Year, &car.Price, &car.Odometer, &car.VinCode,
		&car.Wheel, &car.TradeIN, &car.Crash,
		&car.Credit, &car.New, &car.Status, &car.CreatedAt, &car.Images, &car.Videos, &car.PhoneNumbers,
		&car.ViewCount, &car.Description, &car.MyCar, &car.Owners,
	)

	return car, err
}

func (r *UserRepository) BuyCar(ctx *fasthttp.RequestCtx, carID, userID int) error {

	q := `
		update vehicles 
			set status = 2,
				user_id = $1
		where id = $2 and status = 3 -- 3 is on sale
	`

	_, err := r.db.Exec(ctx, q, userID, carID)

	return err
}

func (r *UserRepository) CreateCar(ctx *fasthttp.RequestCtx, car *model.CreateCarRequest, userID int) (int, error) {

	keys, values, args := auth.BuildParams(car)

	q := `
		INSERT INTO vehicles 
			(
				` + strings.Join(keys, ", ") + `
				, user_id
			) VALUES (
				` + strings.Join(values, ", ") + `,
				$` + strconv.Itoa(len(keys)+1) + `
			) RETURNING id
	`
	var id int
	args = append(args, userID)
	err := r.db.QueryRow(ctx, q, args...).Scan(&id)

	return id, err
}

func (r *UserRepository) UpdateCar(ctx *fasthttp.RequestCtx, car *model.UpdateCarRequest, userID int) error {
	// First check if the car belongs to the user
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM vehicles WHERE id = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, checkQuery, car.ID, userID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("car not found or access denied")
	}

	keys, _, args := auth.BuildParams(car)

	var updateFields []string
	var updateArgs []any
	updateArgs = append(updateArgs, car.ID)

	paramIndex := 2
	for i, key := range keys {
		if key != "id" && key != "user_id" {
			updateFields = append(updateFields, fmt.Sprintf("%s = $%d", key, paramIndex))
			updateArgs = append(updateArgs, args[i])
			paramIndex++
		}
	}

	if len(updateFields) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	q := `
		UPDATE vehicles 
		SET ` + strings.Join(updateFields, ", ") + `, updated_at = NOW()
		WHERE id = $1 AND user_id = $` + fmt.Sprintf("%d", paramIndex)

	updateArgs = append(updateArgs, userID)

	_, err = r.db.Exec(ctx, q, updateArgs...)
	return err
}

func (r *UserRepository) CarLike(ctx *fasthttp.RequestCtx, carID, userID *int) error {

	q := `
		INSERT INTO user_likes(user_id, vehicle_id) values ($2, $1)
	`
	_, err := r.db.Exec(ctx, q, carID, userID)
	return err
}

func (r *UserRepository) RemoveLike(ctx *fasthttp.RequestCtx, carID, userID *int) error {

	q := `
		delete from user_likes where vehicle_id = $1 and user_id = $2
	`
	_, err := r.db.Exec(ctx, q, carID, userID)
	return err
}

func (r *UserRepository) Likes(ctx *fasthttp.RequestCtx, userID *int, nameColumn string) ([]model.GetCarsResponse, error) {
	cars := make([]model.GetCarsResponse, 0)
	q := `
		select 
			vs.id,
			bs.` + nameColumn + ` as brand,
			rs.` + nameColumn + ` as region,
			cs.` + nameColumn + ` as city,
			cls.` + nameColumn + ` as color,
			ms.` + nameColumn + ` as model,
			ts.` + nameColumn + ` as transmission,
			es.` + nameColumn + ` as engine,
			ds.` + nameColumn + ` as drive,
			bts.` + nameColumn + ` as body_type,
			fts.` + nameColumn + ` as fuel_type,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			vs.trade_in,
			vs.owners,
			vs.updated_at,
			images.images,
			videos.videos,
			vs.phone_numbers, 
			vs.view_count,
			true as my_car,
			vs.description,
			CASE 
				WHEN ul.vehicle_id IS NOT NULL THEN true
				ELSE false
			END AS liked
		from vehicles vs
		left join generation_modifications gms on gms.id = vs.modification_id
		left join colors cls on vs.color_id = cls.id
		left join brands bs on vs.brand_id = bs.id
		left join regions rs on vs.region_id = rs.id
		left join cities cs on vs.city_id = cs.id
		left join models ms on vs.model_id = ms.id
		left join transmissions ts on gms.transmission_id = ts.id
		left join engines es on gms.engine_id = es.id
		left join drivetrains ds on gms.drivetrain_id = ds.id
		left join body_types bts on gms.body_type_id = bts.id
		left join fuel_types fts on gms.fuel_type_id = fts.id
		inner join user_likes ul on ul.vehicle_id = vs.id AND ul.user_id = $1
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT $2 || image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT $2 || video as video
				FROM videos
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		where vs.status = 3
		order by vs.id desc
	`
	rows, err := r.db.Query(ctx, q, *userID, r.config.IMAGE_BASE_URL)

	if err != nil {
		return cars, err
	}

	defer rows.Close()

	for rows.Next() {
		var car model.GetCarsResponse

		if err := rows.Scan(
			&car.ID, &car.Brand, &car.Region, &car.City, &car.Color, &car.Model, &car.Transmission, &car.Engine,
			&car.Drivetrain, &car.BodyType, &car.FuelType, &car.Year, &car.Price, &car.Mileage, &car.VinCode,
			&car.Credit, &car.New, &car.Status, &car.CreatedAt, &car.TradeIn, &car.Owners,
			&car.UpdatedAt, &car.Images, &car.Videos, &car.PhoneNumbers, &car.ViewCount, &car.MyCar, &car.Description, &car.Liked,
		); err != nil {
			return cars, err
		}
		cars = append(cars, car)
	}
	return cars, err
}

func (r *UserRepository) CreateCarImages(ctx *fasthttp.RequestCtx, carID int, images []string) error {

	if len(images) == 0 {
		return nil
	}

	q := `
		INSERT INTO images (vehicle_id, image) VALUES ($1, $2)
	`

	for i := range images {
		_, err := r.db.Exec(ctx, q, carID, images[i])

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) CreateCarVideos(ctx *fasthttp.RequestCtx, carID int, video string) error {
	q := `
		INSERT INTO videos (vehicle_id, video) VALUES ($1, $2)
	`
	_, err := r.db.Exec(ctx, q, carID, video)
	return err
}

func (r *UserRepository) CreateMessageFile(ctx *fasthttp.RequestCtx, senderID int, filePath string) error {
	q := `
		INSERT INTO message_files (sender_id, file_path) VALUES ($1, $2)
	`
	_, err := r.db.Exec(ctx, q, senderID, filePath)
	return err
}

func (r *UserRepository) DeleteCarImage(ctx *fasthttp.RequestCtx, carID int, imagePath string) error {
	q := `DELETE FROM images WHERE vehicle_id = $1 AND image = $2`
	_, err := r.db.Exec(ctx, q, carID, imagePath)
	return err
}

func (r *UserRepository) DeleteCarVideo(ctx *fasthttp.RequestCtx, carID int, videoPath string) error {
	q := `DELETE FROM videos WHERE vehicle_id = $1 AND video = $2`
	_, err := r.db.Exec(ctx, q, carID, videoPath)
	return err
}

// GetUserByRoleAndID fetches a single user by role_id and user id
func (r *UserRepository) GetUserByRoleAndID(ctx *fasthttp.RequestCtx, userID int, nameColumn string) (model.ThirdPartyGetProfileRes, error) {

	q := `
		select
			about_me,
			contacts,
			address,
			coordinates,
			$2 || avatar,
			$2 || banner,
			company_name,
			message,
			vat_number,
			company_types.` + nameColumn + ` as company_type,
			activity_fields.` + nameColumn + ` as activity_field,
			profiles.created_at,
            destinations.countries as destinations
		from users 
        left join profiles on profiles.user_id = users.id
        left join (
                SELECT json_agg(
                        json_build_object(
                        'from_country', json_build_object(
                            'id', cf.id,
                            'name', cf.` + nameColumn + `,
                            'flag', cf.flag
                        ),
                        'to_country', json_build_object(
                            'id', ct.id,
                            'name', ct.` + nameColumn + `,
                            'flag', ct.flag
                        )
                        )
                    ) AS countries
                FROM user_destinations ds
                LEFT JOIN countries cf ON cf.id = ds.from_id
                LEFT JOIN countries ct ON ct.id = ds.to_id
                WHERE ds.user_id = $1
        ) destinations on true
		left join company_types on company_types.id = profiles.company_type_id
		left join activity_fields on activity_fields.id = profiles.activity_field_id
        where users.id = $1;
	`
	var profile model.ThirdPartyGetProfileRes
	var contactsJSON []byte
	err := r.db.QueryRow(ctx, q, userID, r.config.IMAGE_BASE_URL).Scan(
		&profile.AboutUs, &contactsJSON,
		&profile.Address,
		&profile.Coordinates, &profile.Avatar,
		&profile.Banner,
		&profile.CompanyName, &profile.Message,
		&profile.VATNumber, &profile.CompanyType,
		&profile.ActivityField,
		&profile.Registered,
		&profile.Destinations,
	)

	if err != nil {
		return profile, err
	}

	if len(contactsJSON) > 0 {
		if err := json.Unmarshal(contactsJSON, &profile.Contacts); err != nil {
			return profile, err
		}
	}

	q = `
		select 
			email,
			phone,
			role_id
		from users 
		where id = $1
	`
	err = r.db.QueryRow(ctx, q, userID).Scan(&profile.Email, &profile.Phone, &profile.RoleID)
	return profile, err
}

func (r *UserRepository) GetThirdPartyUsers(ctx *fasthttp.RequestCtx, roleID, fromID, toID int, search string) ([]model.ThirdPartyUserResponse, error) {
	qWhere := " where u.role_id = $1 "

	if search != "" {
		qWhere = fmt.Sprintf(" %s AND u.username ILIKE '%%%s%%' ", qWhere, search)
	}

	if fromID > 0 && toID > 0 {
		qWhere = fmt.Sprintf(" right join user_destinations ds on ds.user_id = u.id %s AND ds.from_id = %d AND ds.to_id = %d ", qWhere, fromID, toID)
	}

	q := `
		select
			u.id,
			p.company_name,
			p.created_at,
			$2 || p.avatar
		from users u
        left join profiles p on p.user_id = u.id
         %s;
	`
	var users []model.ThirdPartyUserResponse
	var user model.ThirdPartyUserResponse
	rows, err := r.db.Query(ctx, fmt.Sprintf(q, qWhere), roleID, r.config.IMAGE_BASE_URL)

	if err != nil {
		return users, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.Registered, &user.Avatar); err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, err
}
