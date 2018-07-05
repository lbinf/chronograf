package chronograf

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/influxdata/platform"
)

// CellService represents a service for managing cell data.
type CellService interface {
	// FindCellByplatform.ID returns a single cell by platform.ID.
	FindCellByID(ctx context.Context, id platform.ID) (*Cell, error)

	// FindCells returns a list of cells that match filter and the total count of matching cells.
	// Additional options provide pagination & sorting.
	FindCells(ctx context.Context, filter CellFilter) ([]*Cell, int, error)

	// CreateCell creates a new cell and sets b.platform.ID with the new identifier.
	CreateCell(ctx context.Context, b *Cell) error

	// UpdateCell updates a single cell with changeset.
	// Returns the new cell state after update.
	UpdateCell(ctx context.Context, id platform.ID, upd CellUpdate) (*Cell, error)

	// DeleteCell removes a cell by platform.ID.
	DeleteCell(ctx context.Context, id platform.ID) error
}

// CellUpdate is a struct for updating cells.
type CellUpdate struct {
	Name          *string       `json:"name"`
	Visualization Visualization `json:"visualization"`
}

// CellFilter represents a set of filter that restrict the returned results.
type CellFilter struct {
	ID             *platform.ID
	OrganizationID *platform.ID
	Organization   *string
}

// Cell holds positional and visual information for a cell.
type Cell struct {
	CellContents
	Visualization Visualization
}

type CellContents struct {
	ID   platform.ID `json:"id"`
	Name string      `json:"name"`
}

type Visualization interface {
	Visualization()
}

func UnmarshalVisualizationJSON(b []byte) (Visualization, error) {
	var v struct {
		B json.RawMessage `json:"visualization"`
	}

	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}

	var t struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(v.B, &t); err != nil {
		return nil, err
	}

	var vis Visualization
	switch t.Type {
	case "chronograf-v1":
		var qv V1Visualization
		if err := json.Unmarshal(v.B, &qv); err != nil {
			return nil, err
		}
		vis = qv
	default:
		return nil, fmt.Errorf("unknown type %v", t.Type)
	}

	return vis, nil
}

func MarshalVisualizationJSON(v Visualization) ([]byte, error) {
	var s interface{}
	switch vis := v.(type) {
	case V1Visualization:
		s = struct {
			Type string `json:"type"`
			V1Visualization
		}{
			Type:            "chronograf-v1",
			V1Visualization: vis,
		}
	default:
		return nil, fmt.Errorf("unsupported type")
	}
	return json.Marshal(s)
}

func (c Cell) MarshalJSON() ([]byte, error) {
	vis, err := MarshalVisualizationJSON(c.Visualization)
	if err != nil {
		return nil, err
	}

	return json.Marshal(struct {
		CellContents
		Visualization json.RawMessage `json:"visualization"`
	}{
		CellContents:  c.CellContents,
		Visualization: vis,
	})
}

func (c *Cell) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &c.CellContents); err != nil {
		return err
	}

	v, err := UnmarshalVisualizationJSON(b)
	if err != nil {
		return err
	}
	c.Visualization = v
	return nil
}

type V1Visualization struct {
	Queries []DashboardQuery `json:"queries"`
	Axes    map[string]Axis  `json:"axes"`
	// TODO: chronograf will have to use visualizationType rather than type
	Type          string           `json:"visualizationType"`
	CellColors    []CellColor      `json:"colors"`
	Legend        Legend           `json:"legend"`
	TableOptions  TableOptions     `json:"tableOptions,omitempty"`
	FieldOptions  []RenamableField `json:"fieldOptions"`
	TimeFormat    string           `json:"timeFormat"`
	DecimalPlaces DecimalPlaces    `json:"decimalPlaces"`
}

func (V1Visualization) Visualization() {}

/////////////////////////////
// Old Chronograf Types
/////////////////////////////

// DashboardQuery includes state for the query builder.  This is a transition
// struct while we move to the full InfluxQL AST
type DashboardQuery struct {
	Command     string      `json:"query"`                 // Command is the query itself
	Label       string      `json:"label,omitempty"`       // Label is the Y-Axis label for the data
	Range       *Range      `json:"range,omitempty"`       // Range is the default Y-Axis range for the data
	QueryConfig QueryConfig `json:"queryConfig,omitempty"` // QueryConfig represents the query state that is understood by the data explorer
	Source      string      `json:"source"`                // Source is the optional URI to the data source for this queryConfig
	Shifts      []TimeShift `json:"-"`                     // Shifts represents shifts to apply to an influxql query's time range.  Clients expect the shift to be in the generated QueryConfig
}

// Range represents an upper and lower bound for data
type Range struct {
	Upper int64 `json:"upper"` // Upper is the upper bound
	Lower int64 `json:"lower"` // Lower is the lower bound
}

// QueryConfig represents UI query from the data explorer
type QueryConfig struct {
	ID              string              `json:"id,omitempty"`
	Database        string              `json:"database"`
	Measurement     string              `json:"measurement"`
	RetentionPolicy string              `json:"retentionPolicy"`
	Fields          []Field             `json:"fields"`
	Tags            map[string][]string `json:"tags"`
	GroupBy         GroupBy             `json:"groupBy"`
	AreTagsAccepted bool                `json:"areTagsAccepted"`
	Fill            string              `json:"fill,omitempty"`
	RawText         *string             `json:"rawText"`
	Range           *DurationRange      `json:"range"`
	Shifts          []TimeShift         `json:"shifts"`
}

// TimeShift represents a shift to apply to an influxql query's time range
type TimeShift struct {
	Label    string `json:"label"`    // Label user facing description
	Unit     string `json:"unit"`     // Unit influxql time unit representation i.e. ms, s, m, h, d
	Quantity string `json:"quantity"` // Quantity number of units
}

// Field represent influxql fields and functions from the UI
type Field struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Alias string      `json:"alias"`
	Args  []Field     `json:"args,omitempty"`
}

// GroupBy represents influxql group by tags from the UI
type GroupBy struct {
	Time string   `json:"time"`
	Tags []string `json:"tags"`
}

// DurationRange represents the lower and upper durations of the query config
type DurationRange struct {
	Upper string `json:"upper"`
	Lower string `json:"lower"`
}

// Axis represents the visible extents of a visualization
type Axis struct {
	Bounds       []string `json:"bounds"` // bounds are an arbitrary list of client-defined strings that specify the viewport for a cell
	LegacyBounds [2]int64 `json:"-"`      // legacy bounds are for testing a migration from an earlier version of axis
	Label        string   `json:"label"`  // label is a description of this Axis
	Prefix       string   `json:"prefix"` // Prefix represents a label prefix for formatting axis values
	Suffix       string   `json:"suffix"` // Suffix represents a label suffix for formatting axis values
	Base         string   `json:"base"`   // Base represents the radix for formatting axis values
	Scale        string   `json:"scale"`  // Scale is the axis formatting scale. Supported: "log", "linear"
}

// CellColor represents the encoding of data into visualizations
type CellColor struct {
	ID    string `json:"id"`    // ID is the unique id of the cell color
	Type  string `json:"type"`  // Type is how the color is used. Accepted (min,max,threshold)
	Hex   string `json:"hex"`   // Hex is the hex number of the color
	Name  string `json:"name"`  // Name is the user-facing name of the hex color
	Value string `json:"value"` // Value is the data value mapped to this color
}

// Legend represents the encoding of data into a legend
type Legend struct {
	Type        string `json:"type,omitempty"`
	Orientation string `json:"orientation,omitempty"`
}

// TableOptions is a type of options for a DashboardCell with type Table
type TableOptions struct {
	VerticalTimeAxis bool           `json:"verticalTimeAxis"`
	SortBy           RenamableField `json:"sortBy"`
	Wrapping         string         `json:"wrapping"`
	FixFirstColumn   bool           `json:"fixFirstColumn"`
}

// RenamableField is a column/row field in a DashboardCell of type Table
type RenamableField struct {
	InternalName string `json:"internalName"`
	DisplayName  string `json:"displayName"`
	Visible      bool   `json:"visible"`
}

// DecimalPlaces indicates whether decimal places should be enforced, and how many digits it should show.
type DecimalPlaces struct {
	IsEnforced bool  `json:"isEnforced"`
	Digits     int32 `json:"digits"`
}
