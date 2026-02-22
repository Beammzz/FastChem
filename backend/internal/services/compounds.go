package services

import "github.com/takumi/fastchem/internal/models"

// Predefined compounds with known oxidation numbers for each element
var compounds = []models.Compound{
	{
		Formula: "H₂O",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 1},
		},
	},
	{
		Formula: "NaCl",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 1},
			{Symbol: "Cl", OxidationNumber: -1, Count: 1},
		},
	},
	{
		Formula: "CO₂",
		Elements: []models.CompoundElement{
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 2},
		},
	},
	{
		Formula: "NH₃",
		Elements: []models.CompoundElement{
			{Symbol: "N", OxidationNumber: -3, Count: 1},
			{Symbol: "H", OxidationNumber: +1, Count: 3},
		},
	},
	{
		Formula: "H₂SO₄",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "KMnO₄",
		Elements: []models.CompoundElement{
			{Symbol: "K", OxidationNumber: +1, Count: 1},
			{Symbol: "Mn", OxidationNumber: +7, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
	{
		Formula: "CaCO₃",
		Elements: []models.CompoundElement{
			{Symbol: "Ca", OxidationNumber: +2, Count: 1},
			{Symbol: "C", OxidationNumber: +4, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "Fe₂O₃",
		Elements: []models.CompoundElement{
			{Symbol: "Fe", OxidationNumber: +3, Count: 2},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "HNO₃",
		Elements: []models.CompoundElement{
			{Symbol: "H", OxidationNumber: +1, Count: 1},
			{Symbol: "N", OxidationNumber: +5, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 3},
		},
	},
	{
		Formula: "Na₂SO₄",
		Elements: []models.CompoundElement{
			{Symbol: "Na", OxidationNumber: +1, Count: 2},
			{Symbol: "S", OxidationNumber: +6, Count: 1},
			{Symbol: "O", OxidationNumber: -2, Count: 4},
		},
	},
}
