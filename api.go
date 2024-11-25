package main

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type receipt_install_success_response struct {
	ID string `json:"id"`
}

type receipt_query_success_response struct {
	Points uint64 `json:"points"`
}

var acceptableRetailer validator.Func = func(fl validator.FieldLevel) bool {
	// XXX: This is probably slower than maintaining a global, but retailer_re is not thread-safe
	// 			what is the accepted standard for recording a lint as a unit test?
	retailer_re := regexp.MustCompile(`^[\w\s-&]+$`) // re2 probably ok with
	retailer_tested, ok := fl.Field().Interface().(string)
	if ok {
		return retailer_re.MatchString(retailer_tested) // XXX: \s implies multi-line retailers, and probably some other weird stuff.
		// XXX: re2 doesn't directly support 'horizontal whitespcae' character that would remediate.
	} else {
		return false
	}
}

var acceptableDescription validator.Func = func(fl validator.FieldLevel) bool {
	// XXX: This is probably slower than maintaining a global, but retailer_re is not thread-safe
	// 			what is the accepted standard for recording a lint as a unit test?
	retailer_re := regexp.MustCompile(`^[\w\s-]+$`) // re2 probably ok with
	retailer_tested, ok := fl.Field().Interface().(string)
	if ok {
		return retailer_re.MatchString(retailer_tested) // XXX: \s implies multi-line retailers, and probably some other weird stuff.
		// XXX: re2 doesn't directly support 'horizontal whitespcae' character that would remediate.
	} else {
		return false
	}
}

var acceptablePurchaseDate validator.Func = func(fl validator.FieldLevel) bool {
	date_tested, ok := fl.Field().Interface().(string)
	if ok {
		// XXX: Spec is ambiguous regarding MM/DD
		if _, err := time.Parse(time.DateOnly, date_tested); err == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

var acceptablePurchaseTime validator.Func = func(fl validator.FieldLevel) bool {
	time_tested, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := time.Parse("15:04", time_tested); err == nil {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

var acceptablePrice validator.Func = func(fl validator.FieldLevel) bool {
	price_re := regexp.MustCompile(`^\d+[.]\d{2}$`) // re2 probably ok with
	price_tested, ok := fl.Field().Interface().(string)
	if ok {
		return price_re.MatchString(price_tested)
	} else {
		return false
	}
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// XXX: Missing one of these is a runtime error
		//      struct tags want a compiler
		v.RegisterValidation("acceptableRetailer", acceptableRetailer)
		v.RegisterValidation("acceptablePurchaseDate", acceptablePurchaseDate)
		v.RegisterValidation("acceptablePurchaseTime", acceptablePurchaseTime)
		v.RegisterValidation("acceptablePrice", acceptablePrice) // XXX: mismatching these is also a sort of runtime error.
		v.RegisterValidation("acceptableDescription", acceptableDescription)
	}
}
