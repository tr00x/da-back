package utils

import (
	"bytes"
	"context"
	"dubai-auto/internal/config"
	"dubai-auto/pkg/files"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

func MigrateV2(filePath string, db *pgxpool.Pool) error {
	fmt.Println("Migrating v2...")
	f, err := excelize.OpenFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}

	defer f.Close()
	fmt.Println("Excel file opened...")
	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	fmt.Println("Rows fetched...")

	for i := 1; i < len(rows); i++ {

		fmt.Println("Processing row", i)
		brandID, err := getBrandID(rows[i][1], rows[i][2], rows[i][3], rows[i][0]+".png", db)

		if err != nil {
			fmt.Println("Error getting brand ID:", err)
			return err
		}

		modelID, err := getModelID(rows[i][9], rows[i][10], rows[i][11], brandID, db)

		if err != nil {
			fmt.Println("Error getting model ID:", err)
			return err
		}

		generationNameRu := rows[i][16]
		generationName := rows[i][17]
		generationNameAe := rows[i][18]

		if generationNameRu == "" {
			generationNameRu = rows[i][27]
			generationName = rows[i][28]
			generationNameAe = rows[i][29]
		}

		generationID, exists, err := getGenerationID(generationNameRu, generationName, generationNameAe, rows[i][19], rows[i][20], rows[i][36], modelID, db)

		if err != nil {
			fmt.Println("Error getting generation ID:", err)
			return err
		}

		if !exists {
			err = helperResizeImage(generationID, "./generations/"+rows[i][21]+"_main.jpg", config.ENV.DEFAULT_IMAGE_WIDTHS, db)

			if err != nil {
				fmt.Println("Error resizing image:", err)
				return err
			}
		}

		bodyTypeID, err := getBodyTypeID(rows[i][22], rows[i][23], rows[i][24], db)

		if err != nil {
			fmt.Println("Error getting body type ID:", err)
			return err
		}

		engineID, err := getEngineID(rows[i][73], db)

		if err != nil {
			fmt.Println("Error getting engine ID:", err)
			return err
		}

		horsePowerID, err := getHorsePowerID(helperGetHorsePower(rows[i][75]), db)

		if err != nil {
			fmt.Println("Error getting horse power ID:", err)
			return err
		}

		transmissionID, err := getTransmissionID(rows[i][51], rows[i][52], rows[i][53], db)

		if err != nil {
			fmt.Println("Error getting transmission ID:", err)
			return err
		}

		drivetrainID, err := getDrivetrainID(rows[i][55], rows[i][56], rows[i][57], db)

		if err != nil {
			fmt.Println("Error getting drivetrain ID:", err)
			return err
		}

		fuelTypeID, err := getFuelTypeID(rows[i][69], rows[i][70], rows[i][71], db)

		if err != nil {
			fmt.Println("Error getting fuel type ID:", err)
			return err
		}

		_, err = getGenerationModificationID(generationID, horsePowerID, bodyTypeID, engineID, fuelTypeID, drivetrainID, transmissionID, db)

		if err != nil {
			fmt.Println("Error getting generation modification ID:", err)
			return err
		}

	}
	return err
}

func getBrandID(name, name_ru, name_ae, logoFileName string, db *pgxpool.Pool) (int, error) {
	q := `
		select id from brands where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name).Scan(&id)

	if err == pgx.ErrNoRows {

		q = `
			insert into brands (name_ru, name_ae, name) values ($2, $3, $1) returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae).Scan(&id)

		if err != nil {
			return 0, err
		}

		newName := uuid.NewString()
		readerFile, err := os.Open("./logos/" + logoFileName)

		if err != nil {
			return 0, err
		}

		defer readerFile.Close()
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, readerFile)

		err = os.WriteFile(
			"./images/logos/"+newName+".png",
			buf.Bytes(),
			os.ModePerm,
		)

		if err != nil {
			return 0, err
		}

		q = "UPDATE brands SET logo = $1 WHERE id = $2"
		_, err = db.Exec(context.Background(), q, "/images/logos/"+newName+".png", id)

		return id, err

	}

	return id, err
}

func getModelID(name, name_ru, name_ae string, brandID int, db *pgxpool.Pool) (int, error) {
	q := `
		select 
			id 
		from models 
		where name = $1 and brand_id = $2
	`
	var id int
	err := db.QueryRow(context.Background(), q, name, brandID).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into models (name, name_ru, name_ae, brand_id) values ($1, $2, $3, $4) returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae, brandID).Scan(&id)
		return id, err
	}

	return id, err
}

func getBodyTypeID(name_ru, name, name_ae string, db *pgxpool.Pool) (int, error) {
	q := `
		select id from body_types where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into body_types (name, name_ru, name_ae, image) values ($1, $2, $3, '') returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae).Scan(&id)
		return id, err
	}
	return id, err
}

func getTransmissionID(name_ru, name, name_ae string, db *pgxpool.Pool) (int, error) {
	q := `
		select id from transmissions where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name).Scan(&id)
	if err == pgx.ErrNoRows {
		q = `
			insert into transmissions (name, name_ru, name_ae) values ($1, $2, $3) returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae).Scan(&id)
		return id, err
	}
	return id, err
}

func getEngineID(name string, db *pgxpool.Pool) (int, error) {
	name = helperSm3ToL(name)
	name_en := fmt.Sprintf("%s %s", name, "L")
	name_ru := fmt.Sprintf("%s %s", name, "Л")
	name_ae := fmt.Sprintf("%s %s", name, "L")

	q := `
		select id from engines where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name_en).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into engines (name, name_ru, name_ae) values ($1, $2, $3) returning id
		`
		err = db.QueryRow(context.Background(), q, name_en, name_ru, name_ae).Scan(&id)
		return id, err
	}

	return id, err
}

func getDrivetrainID(name, name_ru, name_ae string, db *pgxpool.Pool) (int, error) {
	q := `
		select id from drivetrains where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into drivetrains (name, name_ru, name_ae) values ($1, $2, $3) returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae).Scan(&id)
		return id, err
	}
	return id, err
}

func getFuelTypeID(name_ru, name, name_ae string, db *pgxpool.Pool) (int, error) {
	q := `
		select id from fuel_types where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name).Scan(&id)
	if err == pgx.ErrNoRows {
		q = `
			insert into fuel_types (name, name_ru, name_ae) values ($1, $2, $3) returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae).Scan(&id)
		return id, err
	}
	return id, err
}

func getGenerationID(name_ru, name, name_ae, from, to, wheelStr string, modelID int, db *pgxpool.Pool) (int, bool, error) {
	wheel := wheelStr == "Левый"

	q := `
		select id from generations where name_ru = $1 and model_id = $2
	`
	var id int
	err := db.QueryRow(context.Background(), q, name, modelID).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into generations (
				name, name_ru, name_ae, 
				model_id, start_year, end_year, 
				wheel, image
			) 
			values ($1, $2, $3, $4, $5, $6, $7, '') returning id
		`
		err = db.QueryRow(context.Background(), q, name, name_ru, name_ae, modelID, from, to, wheel).Scan(&id)
		return id, false, err
	}

	return id, true, err
}

func getGenerationModificationID(
	generationID, horsePowerID, bodyTypeID, engineID, fuelTypeID, drivetrainID, transmissionID int,
	db *pgxpool.Pool) (int, error) {
	q := `
		select id from generation_modifications 
		where 
			generation_id = $1 and body_type_id = $2 and 
			engine_id = $3 and fuel_type_id = $4 and 
			drivetrain_id = $5 and transmission_id = $6 and 
			horse_power_id = $7
	`
	var id int
	err := db.QueryRow(context.Background(), q, generationID, bodyTypeID,
		engineID, fuelTypeID,
		drivetrainID, transmissionID,
		horsePowerID).Scan(&id)

	if err == pgx.ErrNoRows {
		q = `
			insert into generation_modifications (
				generation_id, body_type_id, engine_id, fuel_type_id, drivetrain_id, transmission_id, horse_power_id
			) values ($1, $2, $3, $4, $5, $6, $7) returning id
		`
		err = db.QueryRow(context.Background(), q, generationID, bodyTypeID, engineID, fuelTypeID, drivetrainID, transmissionID, horsePowerID).Scan(&id)
		return id, err
	}
	return id, err
}

func getHorsePowerID(name string, db *pgxpool.Pool) (int, error) {
	name_en := fmt.Sprintf("%s h.p", name)
	name_ru := fmt.Sprintf("%s л.с.", name)
	name_ae := fmt.Sprintf("%s حصان", name)

	q := `
		select id from horse_powers where name = $1
	`
	var id int
	err := db.QueryRow(context.Background(), q, name_en).Scan(&id)
	if err == pgx.ErrNoRows {
		q = `
			insert into horse_powers (name, name_ru, name_ae) values ($1, $2, $3) returning id
		`
		err = db.QueryRow(context.Background(), q, name_en, name_ru, name_ae).Scan(&id)
		return id, err
	}
	return id, err
}

func helperGetHorsePower(str string) string {
	idx := strings.Index(str, "л.с.")

	if idx != -1 {
		return str[:idx+len("л.с.")]
	}

	return str
}

func helperSm3ToL(str string) string {
	cm3, err := strconv.Atoi(str)

	if err != nil {
		return ""
	}

	litres := float64(cm3) / 1000.0
	return fmt.Sprintf("%.1f", litres)
}

func helperResizeImage(generationID int, imagePath string, widths []uint, db *pgxpool.Pool) error {
	newName := uuid.NewString()
	readerFile, err := os.Open(imagePath)

	if err != nil {
		return err
	}
	defer readerFile.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, readerFile)

	err = os.WriteFile(
		"./images/generations/"+newName,
		buf.Bytes(),
		os.ModePerm,
	)

	if err != nil {
		fmt.Println("Error writing image:", err)
		return err
	}

	for _, width := range widths {

		err := files.ResizeImage("./images/generations/"+newName, width)

		if err != nil {
			fmt.Println("Error resizing image:", err)
			return err
		}

	}
	err = os.Remove("./images/generations/" + newName)

	if err != nil {
		fmt.Println("failed to remove temp file: ", err)
		return err
	}

	err = updateGenerationImagePath(generationID, "/images/generations/"+newName, db)

	return err
}

func updateGenerationImagePath(generationID int, path string, db *pgxpool.Pool) error {
	q := "UPDATE generations SET image = $1 WHERE id = $2"
	_, err := db.Exec(context.Background(), q, path, generationID)
	return err
}
