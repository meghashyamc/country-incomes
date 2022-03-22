package db

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/meghashyamc/country-incomes/models"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var dbClient *countryIncomesDB

type countryIncomesDB struct {
	db *gorm.DB
}

func PGDSN(dbhost, dbusername, dbpassword, dbname, dbport string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbhost, dbusername, dbpassword, dbname, dbport)
}

func Get() (*countryIncomesDB, error) {
	if dbClient != nil {
		return dbClient, nil
	}

	dsn := PGDSN("localhost", os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"), os.Getenv("PG_DB"), os.Getenv("PG_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Info("failed to open DB: %v", err)
		return nil, err
	}

	dbClient = &countryIncomesDB{
		db: db,
	}

	return dbClient, nil
}

func (db *countryIncomesDB) GetParityFactor(countryFromISO, countryToISO string, year int) (float64, error) {
	countryParityFactors := []models.CountryParityFactor{}
	aMonthAgo := time.Now().Add(-1 * 24 * 30 * time.Hour)
	result := db.db.Where("country_from_iso=? AND country_to_iso=? AND (year BETWEEN ? AND ?) AND updated_at >= ?", countryFromISO, countryToISO, 2020, year, aMonthAgo).Select("parity_factor").Order("year desc").Find(&countryParityFactors)
	if result.Error != nil {
		return 0.0, result.Error
	}
	if result.RowsAffected == 0 || len(countryParityFactors) == 0 {
		return 0.0, errors.New(fmt.Sprintf("could not get parity factor between %s (from) and %s (to)", countryFromISO, countryToISO))
	}
	return countryParityFactors[0].ParityFactor, nil
}

func (db *countryIncomesDB) GetGDPPerCapitaPPP(countryISO string, year int) (int, error) {
	countries := []models.Country{}
	aMonthAgo := time.Now().Add(-1 * 24 * 30 * time.Hour)

	result := db.db.Where("iso=? AND (year BETWEEN ? AND ?) AND updated_at >= ?", countryISO, 2020, year, aMonthAgo).Select("gdp_percapita_ppp").Order("year desc").Find(&countries)
	if result.Error != nil {
		return 0.0, result.Error
	}
	if result.RowsAffected == 0 || len(countries) == 0 {
		return 0.0, errors.New(fmt.Sprintf("could not get GDP per capita (PPP) for %s", countryISO))
	}
	return countries[0].GDPPerCapitaPPP, nil
}

func (db *countryIncomesDB) InsertGDPPerCapitaPPP(countryName, countryISO string, year, gdpPerCapitaPPP int) error {
	timeNowUTC := time.Now().UTC()
	country := &models.Country{CountryName: countryName, ISO: countryISO, Year: year, GDPPerCapitaPPP: gdpPerCapitaPPP, CreatedAt: timeNowUTC, UpdatedAt: timeNowUTC}
	if err := db.db.Clauses(clause.OnConflict{OnConstraint: "iso_gdp_per_capita_year", DoUpdates: clause.AssignmentColumns([]string{"updated_at"})}).Create(&country).Error; err != nil {
		return err
	}

	return nil
}
func (db *countryIncomesDB) InsertParityFactor(countryFromISO, countryToISO string, year int, parityFactor float64) error {
	timeNowUTC := time.Now().UTC()
	countryParityFactor := &models.CountryParityFactor{CountryFromISO: countryFromISO, CountryToISO: countryToISO, Year: year, ParityFactor: parityFactor, CreatedAt: timeNowUTC, UpdatedAt: timeNowUTC}
	if err := db.db.Clauses(clause.OnConflict{OnConstraint: "country_from_country_to_year", DoUpdates: clause.AssignmentColumns([]string{"updated_at"})}).Create(&countryParityFactor).Error; err != nil {
		return err
	}

	return nil
}
