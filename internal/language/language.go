package language

import (
	"embed"
	"fmt"
	"io/fs"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Info struct {
	Code, Name string
}

type LanguageFile struct {
	Info   Info     `yaml:"info"`
	Values Language `yaml:"values"`
}

type QandA struct {
	Question string `yaml:"q"`
	Answer   string `yaml:"a"`
}

type Language struct {
	PlaceEnterLabel string `yaml:"place_enter_label"`
	Blurb           string `yaml:"blurb"`
	ImageBy         string `yaml:"image_by"`

	AboutQandA []QandA `yaml:"about_q_and_a"`

	// Search Page
	AddressNotFound         string `yaml:"address_not_found"`
	ActualAddressLabel      string `yaml:"actual_address_label"`
	CoordinateLabel         string `yaml:"coordinate_label"`
	SafetyScoreLabel        string `yaml:"safety_score_label"`
	SafetyScoreDescription  string `yaml:"safety_score_description"`
	UniversityDistanceLabel string `yaml:"university_distance_label"`
}

//go:embed **.yaml
var languageFiles embed.FS

func Languages() (map[Info]Language, error) {
	files, err := fs.ReadDir(languageFiles, ".")
	if err != nil {
		return nil, fmt.Errorf("unable to list language directory: %w", err)
	}

	langs := make(map[Info]Language, len(files))

	for _, file := range files {
		var data LanguageFile

		f, err := languageFiles.Open(file.Name())
		if err != nil {
			return nil, fmt.Errorf("unable to open file %q: %w", file.Name(), err)
		}

		if err := yaml.NewDecoder(f).Decode(&data); err != nil {
			return nil, fmt.Errorf("unable to decode %q: %w", file.Name(), err)
		}

		// Check are all fields specified
		v := reflect.ValueOf(data.Values)
		for i := 0; i < v.NumField(); i++ {
			if field := v.Field(i); field.IsZero() {
				return nil, fmt.Errorf("empty value %s.%s", data.Info.Code, v.Type().Field(i).Name)
			}
		}

		langs[data.Info] = data.Values
	}

	return langs, nil
}
