package offerrecords

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type OfferRecord struct {
	Id             int       `form:"id" json:"id,omitempty"`
	Item           string    `form:"item" json:"item,omitempty"`
	Priceper       float32   `form:"priceper" json:"priceper,omitempty"`
	Unit           string    `form:"unit" json:"unit,omitempty"`
	Duration_start time.Time `form:"duration_start" json:"duration_start,omitempty"`
	Duration_end   time.Time `form:"duration_end" json:"duration_end,omitempty"`
	Brand          string    `form:"brand" json:"brand,omitempty"`
	Store          string    `form:"store" json:"store,omitempty"`
}

func (offerrecord *OfferRecord) valid() bool {
	return (len(offerrecord.Item) < 21 &&
		offerrecord.Priceper >= 0 &&
		len(offerrecord.Unit) < 11 &&
		len(offerrecord.Brand) < 31 &&
		len(offerrecord.Store) < 16)
}

func executeQuery(db *sql.DB, sql string) ([]OfferRecord, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return []OfferRecord{}, err
	}

	records := []OfferRecord{}

	for rows.Next() {
		var id int
		var item string
		var priceper float32
		var unit string
		var duration_start time.Time
		var duration_end time.Time
		var brand string
		var store string

		err = rows.Scan(&id, &item, &priceper, &unit, &duration_start, &duration_end, &brand, &store)
		if err != nil {
			return []OfferRecord{}, err
		}

		records = append(records, OfferRecord{
			Id:             id,
			Item:           item,
			Priceper:       priceper,
			Unit:           unit,
			Duration_start: duration_start,
			Duration_end:   duration_end,
			Brand:          brand,
			Store:          store})
	}

	return records, nil
}

func fetchRecordsWithParams(db *sql.DB, params map[string]string) ([]OfferRecord, error) {
	firstparam := true
	//itemname
	itemname := params["item"]
	//pricerange
	maxprice := params["priceper"]
	//datespan
	startdate := params["startdate"]
	enddate := params["enddate"]
	//brand
	brand := params["brand"]
	//store
	store := params["store"]
	sql := "SELECT * FROM offers"
	if itemname != "" {
		sql += " WHERE itemname COLLATE latin1_general_ci LIKE '" + itemname + "'"
		firstparam = false
	}
	if maxprice != "" {
		if firstparam {
			sql += "WHERE priceper <= " + maxprice
		} else {
			sql += " AND priceper <= " + maxprice
		}
	}

	if startdate != "" {
		if firstparam {
			sql += "WHERE duration_end >= '" + startdate + "'"
		} else {
			sql += " AND duration_end >= '" + startdate + "'"
		}
	}

	if enddate != "" {
		if firstparam {
			sql += "WHERE duration_start <= '" + enddate + "'"
		} else {
			sql += " AND duration_start <= '" + enddate + "'"
		}
	}

	if brand != "" {
		if firstparam {
			sql += "WHERE brand COLLATE latin1_general_ci LIKE '" + brand + "'"
		} else {
			sql += " AND brand COLLATE latin1_general_ci LIKE '" + brand + "'"
		}
	}

	if store != "" {
		if firstparam {
			sql += "WHERE store COLLATE latin1_general_ci LIKE '" + store + "'"
		} else {
			sql += " AND store COLLATE latin1_general_ci LIKE '" + store + "'"
		}
	}

	fmt.Println(sql)
	records, err := executeQuery(db, sql)
	return records, err
}

func fetchAllOfferRecords(db *sql.DB) ([]OfferRecord, error) {
	records, err := executeQuery(db, "SELECT * FROM offers")
	return records, err
}

func insertOfferRecord(db *sql.DB, offerrecord OfferRecord) error {
	//Create the insert statement
	insertStatement, err := db.Prepare("INSERT INTO offers (itemname, priceper, unit, duration_start, duration_end, brand, store) VALUES( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		log.Println("Error: " + err.Error())
	}
	defer insertStatement.Close()

	//execute the prepared statement
	_, err = insertStatement.Exec(
		offerrecord.Item,
		offerrecord.Priceper,
		offerrecord.Unit,
		offerrecord.Duration_start,
		offerrecord.Duration_end,
		offerrecord.Brand,
		offerrecord.Store)
	if err != nil {
		return err
	}

	return nil
}
