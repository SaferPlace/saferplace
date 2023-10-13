package saferplace

import "testing"

func TestParseComponent(t *testing.T) {
	// Empty string in the response assumes an error
	testCases := map[string]Component{
		"unknown": Component(""),
		"":        Component(""),

		"consumer": ConsumerComponent,
		"review":   ReviewComponent,
		"report":   ReportComponent,
		"uploader": UploaderComponent,
		"viewer":   ViewerComponent,
	}

	for in, want := range testCases {
		t.Run(in, func(t *testing.T) {
			if got, err := ParseComponent(in); (err != nil && got != "") || got != want {
				t.Errorf("ParseComponent(%s) = %s, want %s", in, got, want)
			}
		})
	}
}
