package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/valyala/fasthttp"

	"dubai-auto/internal/config"
	"dubai-auto/internal/model"
)

type AuthRepository struct {
	config *config.Config
	db     *pgxpool.Pool
}

func NewAuthRepository(config *config.Config, db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{config, db}
}

func (r *AuthRepository) UserRegisterDevice(ctx *fasthttp.RequestCtx, userID int, req model.UserRegisterDevice) error {
	q := `
		INSERT INTO user_tokens (user_id, device_id, device_type, device_token)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET device_type = EXCLUDED.device_type, device_token = EXCLUDED.device_token
	`
	_, err := r.db.Exec(ctx, q, userID, req.DeviceID, req.DeviceType, req.DeviceToken)
	return err
}

func (r *AuthRepository) UserLoginGoogle(ctx *fasthttp.RequestCtx, claims model.GoogleUserInfo) (model.UserByEmail, error) {
	var userByEmail model.UserByEmail
	q := `
		INSERT INTO users (email, password, username)
		VALUES ($1, 'google', $2)
		ON CONFLICT (email) DO UPDATE
		SET email = EXCLUDED.email, username = EXCLUDED.username
		RETURNING id, role_id;
	`
	row := r.db.QueryRow(ctx, q, claims.Email, claims.Name)
	err := row.Scan(&userByEmail.ID, &userByEmail.RoleID)

	if err != nil {
		return userByEmail, err
	}

	q = `
		INSERT INTO profiles (user_id, username, registered_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING;
	`
	_, err = r.db.Exec(ctx, q, userByEmail.ID, claims.Name, "email")

	return userByEmail, err
}

func (r *AuthRepository) Application(ctx *fasthttp.RequestCtx, req model.UserApplication) (model.UserByEmail, error) {
	var userByEmail model.UserByEmail
	q := `
		INSERT INTO temp_users (email, username, role_id, phone, password, 
			registered_by, company_name, company_type_id, activity_field_id, 
			vat_number, address, licence_issue_date, licence_expiry_date)
		VALUES ($1, $2, $3, $4, 'application', 'application', $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (email) DO UPDATE
		SET role_id = EXCLUDED.role_id, company_name = EXCLUDED.company_name, company_type_id = EXCLUDED.company_type_id, activity_field_id = EXCLUDED.activity_field_id, vat_number = EXCLUDED.vat_number, address = EXCLUDED.address, licence_issue_date = EXCLUDED.licence_issue_date, licence_expiry_date = EXCLUDED.licence_expiry_date
		RETURNING id, role_id;
	`
	row := r.db.QueryRow(
		ctx, q, req.Email, req.FullName, req.RoleID, req.Phone, req.CompanyName,
		req.CompanyTypeID, req.ActivityFieldID, req.VATNumber, req.Address,
		req.LicenceIssueDate, req.LicenceExpiryDate)
	err := row.Scan(&userByEmail.ID, &userByEmail.RoleID)
	return userByEmail, err
}

func (r *AuthRepository) ApplicationDocuments(ctx *fasthttp.RequestCtx, id int, documents model.UserApplicationDocuments) error {
	q := `
		INSERT INTO documents (licence_path, memorandum_path, copy_of_id_path) 
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var documentID int
	err := r.db.QueryRow(ctx, q, documents.Licence, documents.Memorandum, documents.CopyOfID).Scan(&documentID)

	if err != nil {
		return err
	}

	q = `
		UPDATE temp_users SET documents_id = $1 WHERE id = $2
	`
	_, err = r.db.Exec(ctx, q, documentID, id)
	return err
}

func (r *AuthRepository) DeleteAccount(ctx *fasthttp.RequestCtx, userID int) error {
	// Delete user
	_, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	return err
}

func (r *AuthRepository) TempUserByEmail(ctx *fasthttp.RequestCtx, email *string) (*model.UserByEmail, error) {

	query := `
		SELECT id, email, password, username FROM temp_users WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email)

	var u model.UserByEmail
	err := row.Scan(&u.ID, &u.Email, &u.OTP, &u.Username)

	return &u, err
}

func (r *AuthRepository) UserByEmail(ctx *fasthttp.RequestCtx, email *string) (*model.UserByEmail, error) {

	query := `
		SELECT id, temp_password FROM users WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email)

	var u model.UserByEmail
	err := row.Scan(&u.ID, &u.OTP)

	return &u, err
}

func (r *AuthRepository) UpdateUserTempPassword(ctx *fasthttp.RequestCtx, userID int, password string) error {
	q := `
		UPDATE users SET temp_password = $1 WHERE id = $2
	`
	_, err := r.db.Exec(ctx, q, password, userID)
	return err
}

func (r *AuthRepository) UpdateUserPassword(ctx *fasthttp.RequestCtx, userID int, password string) error {
	q := `
		UPDATE users SET password = $1 WHERE id = $2
	`
	_, err := r.db.Exec(ctx, q, password, userID)
	return err
}

func (r *AuthRepository) TempUserByPhone(ctx *fasthttp.RequestCtx, phone *string) (model.UserByPhone, error) {
	query := `
		SELECT id, phone, password, username FROM temp_users WHERE phone = $1
	`
	row := r.db.QueryRow(ctx, query, phone)

	var userByPhone model.UserByPhone
	err := row.Scan(&userByPhone.ID, &userByPhone.Phone, &userByPhone.OTP, &userByPhone.Username)

	return userByPhone, err
}

func (r *AuthRepository) TempUserEmailGetOrRegister(ctx *fasthttp.RequestCtx, username, email, password string) error {
	var userID int
	q := `
		insert into temp_users (email, password, username, registered_by)
		values ($1, $2, $3, 'email')
		on conflict (email)
		do update
		set 
			password = EXCLUDED.password
		returning id
	`
	err := r.db.QueryRow(ctx, q, email, password, username).Scan(&userID)

	return err
}

func (r *AuthRepository) ThirdPartyLogin(ctx *fasthttp.RequestCtx, email string) (model.ThirdPartyLogin, error) {
	var u model.ThirdPartyLogin
	q := `
		SELECT 
			id, 
			password,
			role_id,
			CASE 
				WHEN created_at = updated_at 
				THEN true 
				ELSE false 
			END as first_time_login
		FROM users 
		WHERE email = $1
	`
	err := r.db.QueryRow(ctx, q, email).Scan(&u.ID, &u.Password, &u.RoleID, &u.FirstTimeLogin)

	if err != nil {
		return u, err
	}

	q = `
		update users set
			updated_at = now()
		where id = $1
	`
	_, err = r.db.Exec(ctx, q, u.ID)

	return u, err
}

func (r *AuthRepository) AdminLogin(ctx *fasthttp.RequestCtx, email string) (model.ThirdPartyLogin, error) {
	var u model.ThirdPartyLogin
	query := `
		SELECT 
			id, 
			password
		FROM users 
		WHERE email = $1 and role_id = 0
	`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Password)
	return u, err
}

func (r *AuthRepository) TempUserPhoneGetOrRegister(ctx *fasthttp.RequestCtx, username, phone, password string) error {
	var userID int
	q := `
		insert into temp_users (phone, password, username, registered_by)
		values ($1, $2, $3, 'phone')
		on conflict (phone)
		do update
		set 
			password = EXCLUDED.password
		returning id
	`
	err := r.db.QueryRow(ctx, q, phone, password, username).Scan(&userID)

	return err
}

func (r *AuthRepository) UserEmailGetOrRegister(ctx *fasthttp.RequestCtx, username, email, password string) (int, error) {
	var userID int
	q := `
		insert into users (email, password, username)
		values ($1, $2, $3)
		on conflict (email)
		do update
		set 
			password = EXCLUDED.password
		returning id
	`
	err := r.db.QueryRow(ctx, q, email, password, username).Scan(&userID)

	if err != nil {
		return userID, err
	}

	q = `
		INSERT INTO profiles (user_id, username, registered_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING;
	`
	_, err = r.db.Exec(ctx, q, userID, username, "email")
	return userID, err
}

func (r *AuthRepository) UserPhoneGetOrRegister(ctx *fasthttp.RequestCtx, username, phone, password string) (int, error) {

	var userID int
	q := `
		insert into users (phone, password, username)
		values ($1, $2, $3)
		on conflict (phone)
		do update
		set 
			password = EXCLUDED.password
		returning id
	`
	err := r.db.QueryRow(ctx, q, phone, password, username).Scan(&userID)

	if err != nil {
		return userID, err
	}

	q = `
		INSERT INTO profiles (user_id, username, registered_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO NOTHING;
	`
	_, err = r.db.Exec(ctx, q, userID, username, "phone")

	return userID, err
}
