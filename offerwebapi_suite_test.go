package main_test

import (
	"fmt"
	"time"

	. "github.com/dr0l3/offerwebapi/offerrecords"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOfferwebapi(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offerwebapi Suite")
}

func defineFactories() {
	gory.Define("offerrecord", OfferRecord{},
		func(factory gory.Factory) {
			duration_start, err := time.Parse("2006.01.02", "2016.01.01")
			if err != nil {
				fmt.Println("Error in timeparsing: " + err.Error())
			}
			duration_end, _ := time.Parse("2006.01.02", "2016.02.02")
			if err != nil {
				fmt.Println("Error in timeparsing: " + err.Error())
			}
			factory["Item"] = "rock"
			factory["Unit"] = "kg"
			factory["Duration_start"] = duration_start
			factory["Duration_end"] = duration_end
			factory["Brand"] = "luxury"
			factory["Store"] = "Netto"
			factory["Priceper"] = gory.Sequence(
				func(n int) interface{} {
					return float32(n)
				})
		})

	gory.Define("offerrecordNegativePrice", OfferRecord{},
		func(factory gory.Factory) {
			duration_start, _ := time.Parse("2016.02.01", "2016.01.01")
			duration_end, _ := time.Parse("2016.02.01", "2016.02.02")
			factory["Item"] = "rock"
			factory["Unit"] = "kg"
			factory["Duration_start"] = duration_start
			factory["Duration_end"] = duration_end
			factory["Brand"] = "luxury"
			factory["Store"] = "Netto"
			factory["Priceper"] = float32(-2.0)
		})
}
