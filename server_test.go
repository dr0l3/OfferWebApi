package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/dr0l3/offerwebapi/offerrecords"
	"github.com/modocache/gory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func sliceFromJSON(data []byte) []map[string]interface{} {
	var result []map[string]interface{}
	json.Unmarshal(data, &result)
	//fmt.Println(result)
	return result
}

func mapFromJSON(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

func createTestDB() *DatabaseConnection {
	openstring := "root:2308esbjerg@/offerdb_test?parseTime=true"

	db, _ := sql.Open("mysql", openstring)
	truncateStatement, _ := db.Prepare("truncate offers;")
	defer truncateStatement.Close()
	truncateStatement.Exec()

	return &DatabaseConnection{db, "testdb"}
}

func insertOfferRecordInTestServer(db *sql.DB, offerrecord OfferRecord) error {
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

var _ = Describe("Server", func() {
	var dbName string
	var dbCon *DatabaseConnection
	var server Server
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Set up a new server, connected to a test database,
		// before each test.
		dbName = "test_server"
		dbCon = createTestDB()
		server = NewServer(dbCon)

		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("GET /item/:itemname", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/item/rock", nil)
		})

		Context("when no offerrecords exits", func() {
			It("returns a statuscode of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})
		})

		Context("when offerrecords exists", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "water"
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those offerrecords in the body", func() {
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
			})
		})

		Context("when a maxprice is given", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or1.Priceper = 20
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "rock"
				or2.Priceper = 30
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the correct offerrecord in the body", func() {
				request, _ = http.NewRequest("GET", "/item/rock?maxprice=25", nil)
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
				Expect(firstRecord["priceper"]).To(Equal(float64(20)), "Items received: ", records)
			})
		})

		Context("when a start date is given", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or1.Duration_start, _ = time.Parse("2006.01.02", "2016.01.01")
				or1.Duration_end, _ = time.Parse("2006.01.02", "2016.01.01")
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "rock"
				or2.Duration_start, _ = time.Parse("2006.01.02", "2016.02.02")
				or2.Duration_end, _ = time.Parse("2006.01.02", "2016.02.02")
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the correct offerrecord in the body", func() {
				request, _ = http.NewRequest("GET", "/item/rock?startdate=2016.01.15", nil)
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
				expectedTime, _ := time.Parse("2006.01.02", "2016.02.02")
				recievedTime, _ := time.Parse("2006-01-02T15:04:05Z", firstRecord["duration_end"].(string))
				Expect(recievedTime).To(Equal(expectedTime))
			})
		})

		Context("when an end date is given", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or1.Duration_end, _ = time.Parse("2006.01.02", "2016.01.01")
				or1.Duration_start, _ = time.Parse("2006.01.02", "2016.01.01")
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "rock"
				or2.Duration_end, _ = time.Parse("2006.01.02", "2016.02.02")
				or2.Duration_start, _ = time.Parse("2006.01.02", "2016.02.02")
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the correct offerrecord in the body", func() {
				request, _ = http.NewRequest("GET", "/item/rock?enddate=2016.01.15", nil)
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
				expectedTime, _ := time.Parse("2006.01.02", "2016.01.01")
				recievedTime, _ := time.Parse("2006-01-02T15:04:05Z", firstRecord["duration_start"].(string))
				Expect(recievedTime).To(Equal(expectedTime))
			})
		})

		Context("when a brand is specified", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or1.Brand = "luxury"
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "rock"
				or2.Brand = "shitty"
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the correct offerrecord in the body", func() {
				request, _ = http.NewRequest("GET", "/item/rock?brand=*lux*", nil)
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
				Expect(firstRecord["brand"]).To(Equal("luxury"), "Items received: ", records)
			})
		})

		Context("when a store is specified", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or1.Item = "rock"
				or1.Store = "Netto"
				or2 := gory.Build("offerrecord").(*OfferRecord)
				or2.Item = "rock"
				or2.Store = "hole in ground"
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the correct offerrecord in the body", func() {
				request, _ = http.NewRequest("GET", "/item/rock?store=*netto*", nil)
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(1))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"), "Items received: ", records)
				Expect(firstRecord["store"]).To(Equal("Netto"), "Items received: ", records)
			})
		})
	})

	Describe("GET /", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/", nil)
		})

		It("returns a response code of 200", func() {
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
		})

		It("returns an overview of available routes", func() {
			server.ServeHTTP(recorder, request)
			var resultJson interface{}
			json.Unmarshal(recorder.Body.Bytes(), &resultJson)
			Expect(len(resultJson.(map[string]interface{}))).To(Equal(3))
		})

	})

	Describe("GET /item/", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/item/", nil)
		})

		Context("when no offerrecords exits", func() {
			It("returns a statuscode of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})
		})

		Context("when offerrecords exists", func() {
			BeforeEach(func() {
				or1 := gory.Build("offerrecord").(*OfferRecord)
				or2 := gory.Build("offerrecord").(*OfferRecord)
				insertOfferRecordInTestServer(dbCon.Db, *or1)
				insertOfferRecordInTestServer(dbCon.Db, *or2)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns those offerrecords in the body", func() {
				server.ServeHTTP(recorder, request)

				records := sliceFromJSON(recorder.Body.Bytes())
				Expect(len(records)).To(Equal(2))

				firstRecord := records[0]

				Expect(firstRecord["item"]).To(Equal("rock"))
				Expect(firstRecord["unit"]).To(Equal("kg"))
				Expect(firstRecord["brand"]).To(Equal("luxury"))
				Expect(firstRecord["store"]).To(Equal("Netto"))
			})
		})
	})

	Describe("with valid JSON", func() {
		BeforeEach(func() {
			offer := gory.Build("offerrecord").(*OfferRecord)
			fmt.Println(offer)

			offer.Duration_start, _ = time.Parse("2016.01.02", "2016.01.01")
			offer.Duration_end, _ = time.Parse("2016.01.02", "2016.02.02")
			body, err := json.Marshal(*offer)
			if err != nil {
				fmt.Println("Error in marshalling: " + err.Error())
			}
			request, _ = http.NewRequest("POST", "/insert", bytes.NewReader(body))
		})

		It("returns a statuscode of 201", func() {
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(201))
		})

		It("returns an success message", func() {
			server.ServeHTTP(recorder, request)
			resultsJSON := mapFromJSON(recorder.Body.Bytes())
			Expect(resultsJSON["status"]).To(ContainSubstring("success"), "For the item: %s. The error message was: %s", resultsJSON["error"])
		})
	})
})
