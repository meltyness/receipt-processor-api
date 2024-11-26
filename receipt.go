package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"
)

type Receipt struct {
	Receipt     *ReceiptContent
	ReceiptUUID string
	Points      uint64
}

// JSON Binding Validators listed below are defined in api.go
type ReceiptContent struct {
	Retailer     string `json:"retailer" binding:"required,max=1024,acceptableRetailer"`
	PurchaseDate string `json:"purchaseDate" binding:"required,max=1024,acceptablePurchaseDate"`
	PurchaseTime string `json:"purchaseTime" binding:"required,max=1024,acceptablePurchaseTime"`
	Items        []Item `json:"items" binding:"required,max=256,min=1,dive"`
	Total        string `json:"total" binding:"required,max=1024,acceptablePrice"`
}

// XXX: These custom struct tags really want a compiler
//
//	... really? https://go.dev/wiki/Well-known-struct-tags
type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required,max=1024,acceptableDescription"`
	Price            string `json:"price" binding:"required,max=1024,acceptablePrice"`
}

func NewReceipt(contents *ReceiptContent, new_uuid string) (*Receipt, error) {
	// XXX: This is probably slower than maintaining globals, but these are not thread-safe
	// 			what is the accepted standard for recording a lint as a unit test?
	// XXX: Caution: don't call Longest and reconfigure these, ever.
	alpha_re := regexp.MustCompile(`([[:alnum:]]){1}`)
	even_re := regexp.MustCompile(`[.]00$`)
	quarter_re := regexp.MustCompile(`[.](00|25|50|75)$`)

	// TODO: validate the receipt total, and any other things that are possible to happen
	score := big.NewInt(0)

	score = score.Add(score, big.NewInt(int64(len(alpha_re.FindAllStringSubmatchIndex(contents.Retailer, -1)))))

	if even_re.MatchString(contents.Total) {
		score = score.Add(score, big.NewInt(50))
	}

	// XXX: This is also ambiguous! Common language would imply either for this req
	//      -- second example confirms.
	if quarter_re.MatchString(contents.Total) {
		score = score.Add(score, big.NewInt(25))
	}

	// >> 1 is the same as floor div 2
	// limited to 256 items.
	total_items := len(contents.Items)
	score = score.Add(score, big.NewInt(int64(5*uint64(total_items>>1))))

	for i := range contents.Items { // (a receipt contains at least one item)
		specific_item := contents.Items[i]
		trimmed_length := len(strings.TrimSpace(specific_item.ShortDescription))
		if (trimmed_length % 3) == 0 { // If the trimmed length of the item description is a multiple of 3
			// (close only counts in aiml, stocks, and graphics)
			price_apa_float, _, _ := big.ParseFloat(strings.ReplaceAll(specific_item.Price, ".", ""), 10, 0, big.ToZero)
			price_apa_int, _ := price_apa_float.Int(nil)
			quo, rem := price_apa_int.DivMod(price_apa_int, big.NewInt(500), big.NewInt(0)) // multiply the price by 0.2 and
			if rem.Cmp(big.NewInt(0)) != 0 {                                                // round up to the nearest integer (presuming, dollars)
				// (from context, when operating on price floats)
				quo.Add(quo, big.NewInt(1))
			}
			// PC_LOAD_LETTER

			score = score.Add(score, quo) // the result is the number of points earned
		}
	}

	if cal_date, err := time.Parse(time.DateOnly, contents.PurchaseDate); err == nil {
		purchased_dom := cal_date.Day()
		if (purchased_dom % 2) == 1 {
			score = score.Add(score, big.NewInt(6))
		}
		purchased_yyyy := cal_date.Year()
		purchased_month := cal_date.Month()
		purchased_month_obj := time.Month(purchased_month)
		bonus_t1 := time.Date(purchased_yyyy, purchased_month_obj, purchased_dom, 14, 0, 0, 0, time.UTC)
		bonus_t2 := time.Date(purchased_yyyy, purchased_month_obj, purchased_dom, 16, 0, 0, 0, time.UTC)

		if clock_time, err := time.Parse("15:04", contents.PurchaseTime); err == nil {
			purchased_time := time.Date(purchased_yyyy, purchased_month_obj, purchased_dom, clock_time.Hour(), clock_time.Minute(), 0, 0, time.UTC)
			if purchased_time.Compare(bonus_t1) == 1 && purchased_time.Compare(bonus_t2) == -1 { // spec is exclusive range
				score = score.Add(score, big.NewInt(10))
			}
		}
	}

	if score.IsUint64() {
		return &Receipt{
			Receipt:     contents,
			ReceiptUUID: new_uuid,
			Points:      score.Uint64(),
		}, nil
	} else {
		return nil, fmt.Errorf("wrapping score from receipt, %s invalid", new_uuid)
	}
}
