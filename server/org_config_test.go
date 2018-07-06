package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/influxdata/chronograf"
	"github.com/influxdata/chronograf/log"
	"github.com/influxdata/chronograf/mocks"
	"github.com/influxdata/chronograf/organizations"
)

func TestOrganizationConfig(t *testing.T) {
	type args struct {
		organizationID string
	}
	type fields struct {
		organizationConfigStore chronograf.OrganizationConfigStore
	}
	type wants struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		wants  wants
	}{
		{
			name: "Get organization configuration",
			args: args{
				organizationID: "default",
			},
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "default":
							return &chronograf.OrganizationConfig{
								OrganizationID: "default",
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "time",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "visibility",
													Value: "hidden",
												},
											},
										},
										{
											Name:     "severity",
											Position: 1,
											Encodings: []chronograf.ColumnEncoding{

												{
													Type:  "visibility",
													Value: "visible",
												},
												{
													Type:  "label",
													Value: "icon",
												},
												{
													Type:  "label",
													Value: "text",
												},
											},
										},
										{
											Name:     "timestamp",
											Position: 2,
											Encodings: []chronograf.ColumnEncoding{

												{
													Type:  "visibility",
													Value: "visible",
												},
											},
										},
										{
											Name:     "message",
											Position: 3,
											Encodings: []chronograf.ColumnEncoding{

												{
													Type:  "visibility",
													Value: "visible",
												},
											},
										},
										{
											Name:     "facility",
											Position: 4,
											Encodings: []chronograf.ColumnEncoding{

												{
													Type:  "visibility",
													Value: "visible",
												},
											},
										},
										{
											Name:     "procid",
											Position: 5,
											Encodings: []chronograf.ColumnEncoding{

												{
													Type:  "visibility",
													Value: "visible",
												},
												{
													Type:  "displayName",
													Value: "Proc ID",
												},
											},
										},
										{
											Name:     "appname",
											Position: 6,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "visibility",
													Value: "visible",
												},
												{
													Type:  "displayName",
													Value: "Application",
												},
											},
										},
										{
											Name:     "host",
											Position: 7,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "visibility",
													Value: "visible",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
				},
			},
			wants: wants{
				statusCode:  200,
				contentType: "application/json",
				body:        `{"links":{"self":"/chronograf/v1/org_config"},"organization":"default","logViewer":{"columns":[{"name":"time","position":0,"encodings":[{"type":"visibility","value":"hidden"}]},{"name":"severity","position":1,"encodings":[{"type":"visibility","value":"visible"},{"type":"label","value":"icon"},{"type":"label","value":"text"}]},{"name":"timestamp","position":2,"encodings":[{"type":"visibility","value":"visible"}]},{"name":"message","position":3,"encodings":[{"type":"visibility","value":"visible"}]},{"name":"facility","position":4,"encodings":[{"type":"visibility","value":"visible"}]},{"name":"procid","position":5,"encodings":[{"type":"visibility","value":"visible"},{"type":"displayName","value":"Proc ID"}]},{"name":"appname","position":6,"encodings":[{"type":"visibility","value":"visible"},{"type":"displayName","value":"Application"}]},{"name":"host","position":7,"encodings":[{"type":"visibility","value":"visible"}]}]}}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store: &mocks.Store{
					OrganizationConfigStore: tt.fields.organizationConfigStore,
				},
				Logger: log.New(log.DebugLevel),
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://any.url", nil)
			ctx := context.WithValue(r.Context(), organizations.ContextKey, tt.args.organizationID)
			r = r.WithContext(ctx)

			s.OrganizationConfig(w, r)

			resp := w.Result()
			content := resp.Header.Get("Content-Type")
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != tt.wants.statusCode {
				t.Errorf("%q. OrganizationConfig() = %v, want %v", tt.name, resp.StatusCode, tt.wants.statusCode)
			}
			if tt.wants.contentType != "" && content != tt.wants.contentType {
				t.Errorf("%q. OrganizationConfig() = %v, want %v", tt.name, content, tt.wants.contentType)
			}
			if eq, _ := jsonEqual(string(body), tt.wants.body); tt.wants.body != "" && !eq {
				t.Errorf("%q. OrganizationConfig() = \n***%v***\n,\nwant\n***%v***", tt.name, string(body), tt.wants.body)
			}
		})
	}
}

func TestLogViewerOrganizationConfig(t *testing.T) {
	type args struct {
		organizationID string
	}
	type fields struct {
		organizationConfigStore chronograf.OrganizationConfigStore
	}
	type wants struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name   string
		args   args
		fields fields
		wants  wants
	}{
		{
			name: "Get log viewer configuration",
			args: args{
				organizationID: "default",
			},
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "default":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "severity",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "color",
													Value: "emergency",
													Name:  "ruby",
												},
												{
													Type:  "color",
													Value: "info",
													Name:  "rainforest",
												},
												{
													Type:  "displayName",
													Value: "Log Severity",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
				},
			},
			wants: wants{
				statusCode:  200,
				contentType: "application/json",
				body:        `{"links":{"self":"/chronograf/v1/org_config/logviewer"},"columns":[{"name":"severity","position":0,"encodings":[{"type":"color","value":"emergency","name":"ruby"},{"type":"color","value":"info","name":"rainforest"},{"type":"displayName","value":"Log Severity"}]}]}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store: &mocks.Store{
					OrganizationConfigStore: tt.fields.organizationConfigStore,
				},
				Logger: log.New(log.DebugLevel),
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://any.url", nil)
			ctx := context.WithValue(r.Context(), organizations.ContextKey, tt.args.organizationID)
			r = r.WithContext(ctx)

			s.LogViewerOrganizationConfig(w, r)

			resp := w.Result()
			content := resp.Header.Get("Content-Type")
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != tt.wants.statusCode {
				t.Errorf("%q. Config() = %v, want %v", tt.name, resp.StatusCode, tt.wants.statusCode)
			}
			if tt.wants.contentType != "" && content != tt.wants.contentType {
				t.Errorf("%q. Config() = %v, want %v", tt.name, content, tt.wants.contentType)
			}
			if eq, _ := jsonEqual(string(body), tt.wants.body); tt.wants.body != "" && !eq {
				t.Errorf("%q. Config() = \n***%v***\n,\nwant\n***%v***", tt.name, string(body), tt.wants.body)
			}
		})
	}
}

func TestReplaceLogViewerOrganizationConfig(t *testing.T) {
	type fields struct {
		organizationConfigStore chronograf.OrganizationConfigStore
	}
	type args struct {
		payload        interface{} // expects JSON serializable struct
		organizationID string
	}
	type wants struct {
		statusCode  int
		contentType string
		body        string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		wants  wants
	}{
		{
			name: "Set log viewer configuration",
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "1337":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "severity",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "color",
													Value: "info",
													Name:  "rainforest",
												},
												{
													Type:  "visibility",
													Value: "visible",
												},
												{
													Type:  "label",
													Value: "icon",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
					UpdateF: func(ctx context.Context, target *chronograf.OrganizationConfig) error {
						return nil
					},
				},
			},
			args: args{
				payload: chronograf.LogViewerConfig{
					Columns: []chronograf.LogViewerColumn{
						{
							Name:     "severity",
							Position: 1,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "color",
									Value: "info",
									Name:  "pineapple",
								},
								{
									Type:  "color",
									Value: "emergency",
									Name:  "ruby",
								},
								{
									Type:  "visibility",
									Value: "visible",
								},
								{
									Type:  "label",
									Value: "icon",
								},
							},
						},
						{
							Name:     "messages",
							Position: 0,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "displayName",
									Value: "Log Messages",
								},
								{
									Type:  "visibility",
									Value: "visible",
								},
							},
						},
					},
				},
				organizationID: "1337",
			},
			wants: wants{
				statusCode:  200,
				contentType: "application/json",
				body:        `{"links":{"self":"/chronograf/v1/org_config/logviewer"},"columns":[{"name":"severity","position":1,"encodings":[{"type":"color","value":"info","name":"pineapple"},{"type":"color","value":"emergency","name":"ruby"},{"type":"visibility","value":"visible"},{"type":"label","value":"icon"}]},{"name":"messages","position":0,"encodings":[{"type":"displayName","value":"Log Messages"},{"type":"visibility","value":"visible"}]}]}`,
			},
		},
		{
			name: "Set invalid log viewer configuration – empty",
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "1337":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "severity",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "color",
													Value: "info",
													Name:  "rainforest",
												},
												{
													Type:  "label",
													Value: "icon",
												},
												{
													Type:  "visibility",
													Value: "visible",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
					UpdateF: func(ctx context.Context, target *chronograf.OrganizationConfig) error {
						return nil
					},
				},
			},
			args: args{
				payload: chronograf.LogViewerConfig{
					Columns: []chronograf.LogViewerColumn{},
				},
				organizationID: "1337",
			},
			wants: wants{
				statusCode:  400,
				contentType: "application/json",
				body:        `{"code":400,"message":"Invalid log viewer config: must have at least 1 column"}`,
			},
		},
		{
			name: "Set invalid log viewer configuration - duplicate column name",
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "1337":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "procid",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "visibility",
													Value: "hidden",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
					UpdateF: func(ctx context.Context, target *chronograf.OrganizationConfig) error {
						return nil
					},
				},
			},
			args: args{
				payload: chronograf.LogViewerConfig{
					Columns: []chronograf.LogViewerColumn{
						{
							Name:     "procid",
							Position: 0,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "visibility",
									Value: "hidden",
								},
							},
						},
						{
							Name:     "procid",
							Position: 1,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "visibility",
									Value: "hidden",
								},
							},
						},
					},
				},
				organizationID: "1337",
			},
			wants: wants{
				statusCode:  400,
				contentType: "application/json",
				body:        `{"code":400,"message":"Invalid log viewer config: Duplicate column name procid"}`,
			},
		},
		{
			name: "Set invalid log viewer configuration - multiple columns with same position value",
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "1337":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "procid",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "visibility",
													Value: "hidden",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
					UpdateF: func(ctx context.Context, target *chronograf.OrganizationConfig) error {
						return nil
					},
				},
			},
			args: args{
				payload: chronograf.LogViewerConfig{
					Columns: []chronograf.LogViewerColumn{
						{
							Name:     "procid",
							Position: 0,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "visibility",
									Value: "hidden",
								},
							},
						},
						{
							Name:     "timestamp",
							Position: 0,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "visibility",
									Value: "hidden",
								},
							},
						},
					},
				},
				organizationID: "1337",
			},
			wants: wants{
				statusCode:  400,
				contentType: "application/json",
				body:        `{"code":400,"message":"Invalid log viewer config: Multiple columns with same position value"}`,
			},
		},
		{
			name: "Set invalid log viewer configuration – no visibility",
			fields: fields{
				organizationConfigStore: &mocks.OrganizationConfigStore{
					FindOrCreateF: func(ctx context.Context, orgID string) (*chronograf.OrganizationConfig, error) {
						switch orgID {
						case "1337":
							return &chronograf.OrganizationConfig{
								LogViewer: chronograf.LogViewerConfig{
									Columns: []chronograf.LogViewerColumn{
										{
											Name:     "severity",
											Position: 0,
											Encodings: []chronograf.ColumnEncoding{
												{
													Type:  "color",
													Value: "info",
													Name:  "rainforest",
												},
												{
													Type:  "label",
													Value: "icon",
												},
											},
										},
									},
								},
							}, nil
						default:
							return nil, chronograf.ErrOrganizationConfigFindOrCreateFailed
						}
					},
					UpdateF: func(ctx context.Context, target *chronograf.OrganizationConfig) error {
						return nil
					},
				},
			},
			args: args{
				payload: chronograf.LogViewerConfig{
					Columns: []chronograf.LogViewerColumn{
						{
							Name:     "severity",
							Position: 1,
							Encodings: []chronograf.ColumnEncoding{
								{
									Type:  "color",
									Value: "info",
									Name:  "pineapple",
								},
								{
									Type:  "color",
									Value: "emergency",
									Name:  "ruby",
								},
								{
									Type:  "label",
									Value: "icon",
								},
							},
						},
					},
				},
				organizationID: "1337",
			},
			wants: wants{
				statusCode:  400,
				contentType: "application/json",
				body:        `{"code":400,"message":"Invalid log viewer config: missing visibility encoding in column severity"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store: &mocks.Store{
					OrganizationConfigStore: tt.fields.organizationConfigStore,
				},
				Logger: log.New(log.DebugLevel),
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://any.url", nil)
			ctx := context.WithValue(r.Context(), organizations.ContextKey, tt.args.organizationID)
			r = r.WithContext(ctx)
			buf, _ := json.Marshal(tt.args.payload)
			r.Body = ioutil.NopCloser(bytes.NewReader(buf))

			s.ReplaceLogViewerOrganizationConfig(w, r)

			resp := w.Result()
			content := resp.Header.Get("Content-Type")
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != tt.wants.statusCode {
				t.Errorf("%q. Config() = %v, want %v", tt.name, resp.StatusCode, tt.wants.statusCode)
			}
			if tt.wants.contentType != "" && content != tt.wants.contentType {
				t.Errorf("%q. Config() = %v, want %v", tt.name, content, tt.wants.contentType)
			}
			if eq, _ := jsonEqual(string(body), tt.wants.body); tt.wants.body != "" && !eq {
				t.Errorf("%q. Config() = \n***%v***\n,\nwant\n***%v***", tt.name, string(body), tt.wants.body)
			}
		})
	}
}
