package offerrecords

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
}

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"/item/":          "List of all records for all items",
		"/item/:itemname": "List of all records for itemname",
		"/insert":         "insert a record"})
	return
}

func allItemsHandler(c *gin.Context) {
	dbCon := c.MustGet("databaseconnection").(*DatabaseConnection)
	records, err := fetchAllOfferRecords(dbCon.Db)
	if err != nil {
		c.JSON(400, gin.H{"error": "error during fetching of records: " + err.Error()})
		return
	}
	c.JSON(200, records)
}

func itemHandler(c *gin.Context) {
	dbCon := c.MustGet("databaseconnection").(*DatabaseConnection)
	params := make(map[string]string)
	params["item"] = strings.Replace(c.Param("itemname"), "*", "%", -1)
	params["priceper"] = c.Query("maxprice")
	params["startdate"] = c.Query("startdate")
	params["enddate"] = c.Query("enddate")
	params["store"] = strings.Replace(c.Query("store"), "*", "%", -1)
	params["brand"] = strings.Replace(c.Query("brand"), "*", "%", -1)
	record, err := fetchRecordsWithParams(dbCon.Db, params)
	if err != nil {
		c.JSON(400, gin.H{"error": "error during fetching of records: " + err.Error()})
		return
	}
	c.JSON(200, record)
}

func insertHandler(c *gin.Context) {
	dbCon := c.MustGet("databaseconnection").(*DatabaseConnection)
	var recievedRecord OfferRecord
	err := c.Bind(&recievedRecord)
	if err != nil {
		c.JSON(400, gin.H{"error": "error during parsing of parameters: " + err.Error()})
		return
	}

	if !recievedRecord.valid() {
		c.JSON(400, gin.H{"error": "invalid record"})
		return
	}

	err = insertOfferRecord(dbCon.Db, recievedRecord)

	if err != nil {
		c.JSON(400, gin.H{"error": "error during insert: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success"})
}

func NewServer(databaseconnection *DatabaseConnection) Server {
	router := Server{gin.Default()}

	database := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("databaseconnection", databaseconnection)
		}
	}

	router.Use(database())

	router.GET("/", indexHandler)
	router.GET("/item/", allItemsHandler)
	router.POST("/insert", insertHandler)
	router.GET("/item/:itemname", itemHandler)
	/*r.Get("/item/:itemname", "get offers for an item", itemHandler)
	r.Get("/item/", "get list of items", itemsHandler)
	r.Get("/store/:storename", "get offers for a store", storeHandler)
	r.Get("/store/", "get list of stores", storesHandler)
	r.Post("/insert", "insert an offerline", insertHandler)*/

	return router
}
