package usecases

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"

	helpers "gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	responses "gitlab.com/mims-api-service/responses"
)

// The function should correctly parse the input CSV data and create RoadConditionSurveyM objects with the correct fields.
func TestParseCSVData(t *testing.T) {
	// Create a sample CSV data
	csvData := MockRoadConditionCSV()

	// Create a sample full geometry
	fullGeom := models.FullGeom{
		KmStart: 1000.0,
		KmEnd:   3000.0,
		Geom:    "LINESTRING(0 0, 2 2)",
	}

	// Call the function under test
	res, err := RoadConditionAnalysisData(csvData, 0.1, fullGeom)

	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the RoadConditionSurveyM objects are created correctly
	expected := []models.RoadConditionSurveyM{
		{
			KmStart:     1800.0,
			KmEnd:       1825.0,
			IRI:         helpers.Float64Ptr(5.6),
			IFI:         helpers.Float64Ptr(3.2),
			MPD:         helpers.Float64Ptr(2.1),
			RUT:         helpers.Float64Ptr(1.8),
			TheGeom:     "ST_LineSubstring('LINESTRING(0 0, 2 2)', 0.0, 0.5)",
			ImgFilepath: "image1.jpg",
		},
		{
			KmStart:     1825.0,
			KmEnd:       1850.0,
			IRI:         helpers.Float64Ptr(4.8),
			IFI:         helpers.Float64Ptr(2.9),
			MPD:         helpers.Float64Ptr(1.7),
			RUT:         helpers.Float64Ptr(1.5),
			TheGeom:     "ST_LineSubstring('LINESTRING(0 0, 2 2)', 0.5, 1.0)",
			ImgFilepath: "image2.jpg",
		},
	}

	if !reflect.DeepEqual(res.RoadConditionSurveyM, expected) {
		t.Errorf("Unexpected RoadConditionSurveyM objects. Got %v, want %v", res.RoadConditionSurveyM, expected)
	}
}

// The function should correctly calculate the IRI, IFI, MPD, and RUT averages for the entire road.
func TestCalculateAveragesForEntireRoad(t *testing.T) {
	// Create a sample CSV data
	csvData := []models.RoadConditionCSV{
		{KMStart: 0, KMEnd: 1, IRI: helpers.Float64Ptr(5.6), IFI: helpers.Float64Ptr(3.2), MPD: helpers.Float64Ptr(2.1), RUT: helpers.Float64Ptr(1.8), ImgFilepath: "image1.jpg"},
		{KMStart: 1, KMEnd: 2, IRI: helpers.Float64Ptr(4.8), IFI: helpers.Float64Ptr(2.9), MPD: helpers.Float64Ptr(1.7), RUT: helpers.Float64Ptr(1.5), ImgFilepath: "image2.jpg"},
	}

	// Create a sample full geometry
	fullGeom := models.FullGeom{
		KmStart: 0,
		KmEnd:   2,
		Geom:    "LINESTRING(0 0, 2 2)",
	}

	// Call the function under test
	res, err := RoadConditionAnalysisData(csvData, 0.1, fullGeom)

	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the averages are calculated correctly for the entire road
	expected := responses.RoadData{
		IriAverage: helpers.Float64Ptr(4.9),
		IfiAverage: helpers.Float64Ptr(3.05),
		MpdAverage: helpers.Float64Ptr(1.9),
		RutAverage: helpers.Float64Ptr(1.65),
	}

	if !reflect.DeepEqual(res.RoadData, expected) {
		t.Errorf("Unexpected RoadData. Got %v, want %v", res.RoadData, expected)
	}
}

// The function should correctly calculate the IRI, IFI, MPD, and RUT averages for each 100m segment of the road.
func TestCalculateAveragesFor100mSegments(t *testing.T) {
	// Create a sample CSV data
	csvData := []models.RoadConditionCSV{
		{KMStart: 0, KMEnd: 1, IRI: helpers.Float64Ptr(5.6), IFI: helpers.Float64Ptr(3.2), MPD: helpers.Float64Ptr(2.1), RUT: helpers.Float64Ptr(1.8), ImgFilepath: "image1.jpg"},
		{KMStart: 1, KMEnd: 2, IRI: helpers.Float64Ptr(4.8), IFI: helpers.Float64Ptr(2.9), MPD: helpers.Float64Ptr(1.7), RUT: helpers.Float64Ptr(1.5), ImgFilepath: "image2.jpg"},
	}

	// Create a sample full geometry
	fullGeom := models.FullGeom{
		KmStart: 0,
		KmEnd:   2,
		Geom:    "LINESTRING(0 0, 2 2)",
	}

	// Call the function under test
	res, err := RoadConditionAnalysisData(csvData, 0.1, fullGeom)

	// Check if there is no error
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check if the averages are calculated correctly for each 100m segment
	expected := []models.RoadConditionSurvey100M{
		{
			KmStart: 0.1,
			KmEnd:   1,
			IRI:     helpers.CalculateRoadCondition(res.RoadData.IriKm, res.RoadData.DividerCountIri),
			IFI:     helpers.CalculateRoadCondition(res.RoadData.IfiKm, res.RoadData.DividerCountIfi),
			MPD:     helpers.CalculateRoadCondition(res.RoadData.MpdKm, res.RoadData.DividerCountMpd),
			RUT:     helpers.CalculateRoadCondition(res.RoadData.RutKm, res.RoadData.DividerCountRut),
		},
		{
			KmStart: 1,
			KmEnd:   2,
			IRI:     nil,
			IFI:     nil,
			MPD:     nil,
			RUT:     nil,
		},
	}

	if !reflect.DeepEqual(res.RoadConditionSurvey100M, expected) {
		t.Errorf("Unexpected RoadConditionSurvey100M objects. Got %v, want %v", res.RoadConditionSurvey100M, expected)
	}
}

// The function should return an error if the input CSV data is empty.
func TestEmptyCSVData(t *testing.T) {
	// Create an empty CSV data
	csvData := []models.RoadConditionCSV{}

	// Create a sample full geometry
	fullGeom := models.FullGeom{
		KmStart: 0,
		KmEnd:   2,
		Geom:    "LINESTRING(0 0, 2 2)",
	}

	// Call the function under test
	_, err := RoadConditionAnalysisData(csvData, 0.1, fullGeom)

	// Check if the error is returned
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

func MockRoadConditionCSV() []models.RoadConditionCSV {
	mockJSON := `[
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1800,
        "km_end": 1825,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1825,
        "km_end": 1850,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1850,
        "km_end": 1875,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1875,
        "km_end": 1900,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1900,
        "km_end": 1925,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1925,
        "km_end": 1950,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1950,
        "km_end": 1975,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1975,
        "km_end": 1995,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "CC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 1995,
        "km_end": 2000,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "AC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 2000,
        "km_end": 2025,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "AC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 2025,
        "km_end": 2050,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "AC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 2050,
        "km_end": 2075,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "AC",
        "img_filepath": ""
    },
    {
        "road_id": 1,
        "road_code": "SRINAKARIN ROAD - KLONG SONG TON  NUN DISTRICT",
        "name": "ถนนศรีนครินทร์ - แขวงคลองสองต้นนุ่น",
        "km_start": 2075,
        "km_end": 2100,
        "iri": 0,
        "mpd": 0,
        "rut": 0,
        "ifi": 0,
        "survey_type": "AC",
        "img_filepath": ""
    }
]`
	var data []models.RoadConditionCSV
	err := json.Unmarshal([]byte(mockJSON), &data)
	if err != nil {
		log.Fatalf("Error parsing mock JSON: %s", err)
	}

	return data
}
