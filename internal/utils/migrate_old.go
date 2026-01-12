package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

func ExcelMigrateOld(db *pgxpool.Pool) error {
	err := MarksOld("docs/1.marks.xlsx", db)
	fmt.Println(err)
	err = ModelsOld("docs/2.models.xlsx", db)
	fmt.Println(err)
	err = GenerationsOld("docs/3.generations.xlsx", db)
	fmt.Println(err)
	err = ConfigurationsOld("docs/4.configurations.xlsx", db)
	fmt.Println(err)
	err = ComplectationsOld("docs/5.complectations.xlsx", db)
	fmt.Println(err)
	db.Exec(context.Background(), `
		delete from models where id not in (
			3232,
			3233,
			3234
		);
	`)
	return nil
}

func MarksOld(filePath string, db *pgxpool.Pool) error {
	f, err := excelize.OpenFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	q := `
		insert into brands(
			id, name, logo, popular
		) values ($1, $2, $3, true)
	`

	for i := range rows {
		if i == 0 || rows[i][0] != "310" {
			continue
		}

		id, _ := strconv.Atoi(rows[i][0])
		_, err = db.Exec(context.Background(), q, id, rows[i][2], "/images/logo/audi.png")

		if err != nil {
			continue
		}
	}
	return err
}

func ModelsOld(filePath string, db *pgxpool.Pool) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	q := `
		insert into models(
			id, brand_id, name
		) values ($1, $2, $3)
	`

	for i := range rows {
		if i == 0 || rows[i][6] != "310" {
			continue
		}

		id, _ := strconv.Atoi(rows[i][0])
		brand_id, _ := strconv.Atoi(rows[i][6])
		_, err = db.Exec(context.Background(), q, id, brand_id, rows[i][2])

		if err != nil {
			continue
			// return err
		}
	}
	return err
}

func GenerationsOld(filePath string, db *pgxpool.Pool) error {
	f, err := excelize.OpenFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	q := `
		insert into generations (
			id, name, start_year, end_year, image, model_id, wheel
		) values 
		(
			$1, $2, $3, $4, $5, $6, $7
		)
	`

	for i := range rows {

		if i == 0 || rows[i][9] != "310" {
			continue
		}

		if i == 479 {
			break
		}
		// inner_id, _ := strconv.Atoi(rows[i][1])
		// autoru_mark_id, _ := strconv.Atoi(rows[i][9])
		id, _ := strconv.Atoi(rows[i][1]) // inner_id
		autoru_model_id, _ := strconv.Atoi(rows[i][8])
		year_from, _ := strconv.Atoi(rows[i][3])
		year_to, _ := strconv.Atoi(rows[i][4])
		wheel := true
		if i%4 == 0 {
			wheel = false
		}
		_, err = db.Exec(context.Background(), q,
			id, rows[i][2], year_from, year_to, rows[i][7][2:], autoru_model_id, wheel)

		if err != nil {
			continue
		}
	}
	return nil
}

func ConfigurationsOld(filePath string, db *pgxpool.Pool) error {
	f, err := excelize.OpenFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	q := `
		insert into configurations (
			id, generation_id, body_type_id
		) values (
		 	$1, $2, $3
		)
	`
	qBodyType := `
		select
			id 
		from body_types 
		where 
			name ilike $1;
	`

	qBodyTypeInsert := `
		insert into body_types (name, image)
		values ($1, 'empty')
	`

	for i := range rows {

		if i == 0 || rows[i][9] != "310" {
			continue
		}
		bodyTypeID := 0
		db.QueryRow(context.Background(), qBodyType, rows[i][4]).Scan(&bodyTypeID)

		if bodyTypeID == 0 {
			db.QueryRow(context.Background(), qBodyTypeInsert, rows[i][4]).Scan(&bodyTypeID)
		}

		id, _ := strconv.Atoi(rows[i][1]) // inner_id
		// autoru_generation_id, _ := strconv.Atoi(rows[i][8])
		autoru_generation_inner_id, _ := strconv.Atoi(rows[i][11])
		// autoru_mark_id, _ := strconv.Atoi(rows[i][9])

		_, err = db.Exec(context.Background(), q,
			id, autoru_generation_inner_id, bodyTypeID)

		if err != nil {
			continue
		}
	}
	return nil
}

type TechInfoMainOld struct {
	ID     string `json:"id"`
	Entity []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Units string `json:"units,omitempty"`
		Value string `json:"value"`
	} `json:"entity"`
}

func ComplectationsOld(filePath string, db *pgxpool.Pool) error {
	f, err := excelize.OpenFile(filePath)

	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	if sheetName == "" {
		return fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := f.GetRows(sheetName)

	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	qGenerationID := `
		select generation_id, body_type_id from configurations where id = $1
	`

	q := `
		insert into generation_modifications (
			generation_id, body_type_id, engine_id, 
			fuel_type_id, drivetrain_id, transmission_id
		) values (
		 	$1, $2, $3, 
			$4, $5, $6
		)
	`
	for i := range rows {

		if i == 0 {
			continue
		}

		var data struct {
			TechInfoMain TechInfoMainOld `json:"tech_info_main"`
		}

		// break

		generationID := 0
		bodyTypeID := 0
		engineID := 0
		fuelTypeID := 0
		drivetrainID := 0
		transmissionID := 0

		id, _ := strconv.Atoi(rows[i][10])
		db.QueryRow(context.Background(), qGenerationID, id).Scan(&generationID, &bodyTypeID)

		if err := json.Unmarshal([]byte(rows[i][4]), &data); err != nil {
			return err
		}
		// var cnt = 0
		for i := range data.TechInfoMain.Entity {

			if data.TechInfoMain.Entity[i].ID == "displacement" {
				_ = db.QueryRow(context.Background(), `
					insert into engines (value) values ($1)
					on conflict(value)
					DO UPDATE SET value = EXCLUDED.value
					returning id
				`, data.TechInfoMain.Entity[i].Value).Scan(&engineID)
			}

			if data.TechInfoMain.Entity[i].ID == "engine_type" {
				db.QueryRow(context.Background(), `
					insert into fuel_types (name) values ($1)
					on conflict(name)
					DO UPDATE SET name = EXCLUDED.name
					returning id
				`, data.TechInfoMain.Entity[i].Value).Scan(&fuelTypeID)
			}

			if data.TechInfoMain.Entity[i].ID == "gear_type" {
				db.QueryRow(context.Background(), `
					insert into drivetrains (name) values ($1)
					on conflict(name)
					DO UPDATE SET name = EXCLUDED.name
					returning id
				`, data.TechInfoMain.Entity[i].Value).Scan(&drivetrainID)
			}

			if data.TechInfoMain.Entity[i].ID == "transmission" {
				db.QueryRow(context.Background(), `
					insert into transmissions (name) values ($1)
					on conflict(name)
					DO UPDATE SET name = EXCLUDED.name
					returning id
				`, data.TechInfoMain.Entity[i].Value).Scan(&transmissionID)
			}
		}

		if generationID == 0 || transmissionID == 0 || engineID == 0 || drivetrainID == 0 || fuelTypeID == 0 || bodyTypeID == 0 {
			continue
		}

		_, err = db.Exec(context.Background(), q, generationID, bodyTypeID,
			engineID, fuelTypeID, drivetrainID, transmissionID)

		if err != nil {
			continue
		}
	}
	return nil
}
