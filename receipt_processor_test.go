package main

import (
	"github.com/google/uuid"
	"testing"
)

func TestTesting(t *testing.T) {
	//t.Errorf("No this is fine")
}

func TestAlsoTesting(t *testing.T) {
	// t.Errorf("No this is fine")
}

func TestExampleCase1(t *testing.T) {
	new_receipt := &ReceiptContent{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []Item{
			{
				ShortDescription: "Mountain Dew 12PK",
				Price:            "6.49",
			},
			{
				ShortDescription: "Emils Cheese Pizza",
				Price:            "12.25",
			},
			{
				ShortDescription: "Knorr Creamy Chicken",
				Price:            "1.26",
			},
			{
				ShortDescription: "Doritos Nacho Cheese",
				Price:            "3.35",
			},
			{
				ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				Price:            "12.00",
			},
		},
		Total: "35.35",
	}
	new_uuid := uuid.NewString()
	installed_receipt := NewReceipt(new_receipt, new_uuid)

	if installed_receipt.Points != (28) {
		t.Errorf("Tally incorrect, found installed_receipt.Points %d want 28", installed_receipt.Points)
	}
}

func TestExampleCase2(t *testing.T) {
	new_receipt := &ReceiptContent{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Items: []Item{
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
			{
				ShortDescription: "Gatorade",
				Price:            "2.25",
			},
		},
		Total: "9.00",
	}
	new_uuid := uuid.NewString()
	installed_receipt := NewReceipt(new_receipt, new_uuid)

	if installed_receipt.Points != (109) {
		t.Errorf("Tally incorrect, found installed_receipt.Points %d want 109", installed_receipt.Points)
	}
}
