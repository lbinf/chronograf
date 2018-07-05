package chronograf_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/chronograf/v2"
	"github.com/influxdata/platform"
)

func TestCell_MarshalJSON(t *testing.T) {
	type args struct {
		cell chronograf.Cell
	}
	type wants struct {
		json string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			args: args{
				cell: chronograf.Cell{
					CellContents: chronograf.CellContents{
						ID:   platform.ID("0"), // This ends up being id 30 encoded
						Name: "hello",
					},
					Visualization: chronograf.V1Visualization{
						Type: "line",
					},
				},
			},
			wants: wants{
				json: `
{
  "id": "30",
  "name": "hello",
  "visualization": {
    "type": "chronograf-v1",
    "queries": null,
    "axes": null,
    "visualizationType": "line",
    "colors": null,
    "legend": {},
    "tableOptions": {
      "verticalTimeAxis": false,
      "sortBy": {
        "internalName": "",
        "displayName": "",
        "visible": false
      },
      "wrapping": "",
      "fixFirstColumn": false
    },
    "fieldOptions": null,
    "timeFormat": "",
    "decimalPlaces": {
      "isEnforced": false,
      "digits": 0
    }
  }
}
`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.MarshalIndent(tt.args.cell, "", "  ")
			if err != nil {
				t.Fatalf("error marshalling json")
			}

			eq, err := jsonEqual(string(b), tt.wants.json)
			if err != nil {
				t.Fatalf("error marshalling json")
			}
			if !eq {
				t.Errorf("JSON did not match\nexpected:%s\ngot:\n%s\n", tt.wants.json, string(b))
			}
		})
	}
}

func jsonEqual(s1, s2 string) (eq bool, err error) {
	var o1, o2 interface{}

	if err = json.Unmarshal([]byte(s1), &o1); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(s2), &o2); err != nil {
		return
	}

	return cmp.Equal(o1, o2), nil
}
