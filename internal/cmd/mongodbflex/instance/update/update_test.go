package update

import (
	"context"
	"fmt"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

const (
	testRegion = "eu02"
)

type testCtxKey struct{}

var testCtx = context.WithValue(context.Background(), testCtxKey{}, "foo")
var testClient = &mongodbflex.APIClient{}

type mongoDBFlexClientMocked struct {
	listFlavorsFails  bool
	listFlavorsResp   *mongodbflex.ListFlavorsResponse
	listStoragesFails bool
	listStoragesResp  *mongodbflex.ListStoragesResponse
	getInstanceFails  bool
	getInstanceResp   *mongodbflex.GetInstanceResponse
}

func (c *mongoDBFlexClientMocked) PartialUpdateInstance(ctx context.Context, projectId, instanceId, region string) mongodbflex.ApiPartialUpdateInstanceRequest {
	return testClient.PartialUpdateInstance(ctx, projectId, instanceId, region)
}

func (c *mongoDBFlexClientMocked) GetInstanceExecute(_ context.Context, _, _, _ string) (*mongodbflex.GetInstanceResponse, error) {
	if c.getInstanceFails {
		return nil, fmt.Errorf("get instance failed")
	}
	return c.getInstanceResp, nil
}

func (c *mongoDBFlexClientMocked) ListStoragesExecute(_ context.Context, _, _, _ string) (*mongodbflex.ListStoragesResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list storages failed")
	}
	return c.listStoragesResp, nil
}

func (c *mongoDBFlexClientMocked) ListFlavorsExecute(_ context.Context, _, _ string) (*mongodbflex.ListFlavorsResponse, error) {
	if c.listFlavorsFails {
		return nil, fmt.Errorf("list flavors failed")
	}
	return c.listFlavorsResp, nil
}

var testProjectId = uuid.NewString()
var testInstanceId = uuid.NewString()
var testFlavorId = uuid.NewString()

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testInstanceId,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureRequiredFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureStandardFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		globalflags.ProjectIdFlag: testProjectId,
		globalflags.RegionFlag:    testRegion,
		flavorIdFlag:              testFlavorId,
		instanceNameFlag:          "example-name",
		aclFlag:                   "0.0.0.0/0",
		backupScheduleFlag:        "0 0 * * *",
		storageClassFlag:          "class",
		storageSizeFlag:           "10",
		versionFlag:               "5.0",
		typeFlag:                  "Single",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureRequiredInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId: testInstanceId,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureStandardInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		GlobalFlagModel: &globalflags.GlobalFlagModel{
			ProjectId: testProjectId,
			Region:    testRegion,
			Verbosity: globalflags.VerbosityDefault,
		},
		InstanceId:     testInstanceId,
		FlavorId:       utils.Ptr(testFlavorId),
		InstanceName:   utils.Ptr("example-name"),
		ACL:            utils.Ptr([]string{"0.0.0.0/0"}),
		BackupSchedule: utils.Ptr("0 0 * * *"),
		StorageClass:   utils.Ptr("class"),
		StorageSize:    utils.Ptr(int64(10)),
		Version:        utils.Ptr("5.0"),
		Type:           utils.Ptr("Single"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *mongodbflex.ApiPartialUpdateInstanceRequest)) mongodbflex.ApiPartialUpdateInstanceRequest {
	request := testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion)
	request = request.PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{})
	for _, mod := range mods {
		mod(&request)
	}
	return request
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description   string
		argValues     []string
		flagValues    map[string]string
		aclValues     []string
		isValid       bool
		expectedModel *inputModel
	}{
		{
			description: "no values",
			argValues:   []string{},
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     false,
		},
		{
			description: "only instance and project ids",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureRequiredFlagValues(),

			isValid: false,
		},
		{
			description:   "all values with flavor id",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureStandardFlagValues(),
			isValid:       true,
			expectedModel: fixtureStandardInputModel(),
		},
		{
			description: "all values with cpu and ram",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, flavorIdFlag)
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.FlavorId = nil
				model.CPU = utils.Ptr(int64(2))
				model.RAM = utils.Ptr(int64(4))
			}),
		},
		{
			description: "project id missing",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				delete(flagValues, globalflags.ProjectIdFlag)
			}),
			isValid: false,
		},
		{
			description: "project id invalid 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "project id invalid 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[globalflags.ProjectIdFlag] = "invalid-uuid"
			}),
			isValid: false,
		},
		{
			description: "instance id invalid 1",
			argValues:   []string{""},
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "instance id invalid 2",
			argValues:   []string{"invalid-uuid"},
			flagValues:  fixtureRequiredFlagValues(),
			isValid:     false,
		},
		{
			description: "invalid with flavor ID, CPU and RAM",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[flavorIdFlag] = testFlavorId
				flagValues[cpuFlag] = "2"
				flagValues[ramFlag] = "4"
			}),
			isValid: false,
		},
		{
			description: "invalid with flavor ID and CPU",
			argValues:   fixtureArgValues(),
			flagValues: fixtureRequiredFlagValues(func(flagValues map[string]string) {
				flagValues[flavorIdFlag] = testFlavorId
				flagValues[cpuFlag] = "2"
			}),
			isValid: false,
		},
		{
			description: "no acl flag",
			argValues:   fixtureArgValues(),
			flagValues: fixtureStandardFlagValues(func(flagValues map[string]string) {
				delete(flagValues, aclFlag)
			}),
			isValid: true,
			expectedModel: fixtureStandardInputModel(func(model *inputModel) {
				model.ACL = nil
			}),
		},
		{
			description: "repeated acl flags",
			argValues:   fixtureArgValues(),
			flagValues:  fixtureRequiredFlagValues(),
			aclValues:   []string{"198.51.100.14/24", "198.51.100.14/32"},
			isValid:     true,
			expectedModel: fixtureRequiredInputModel(func(model *inputModel) {
				model.ACL = utils.Ptr([]string{"198.51.100.14/24", "198.51.100.14/32"})
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(&params.CmdParams{Printer: p})
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			for _, value := range tt.aclValues {
				err := cmd.Flags().Set(aclFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", aclFlag, value, err)
				}
			}

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing flags: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	tests := []struct {
		description       string
		model             *inputModel
		expectedRequest   mongodbflex.ApiPartialUpdateInstanceRequest
		getInstanceFails  bool
		getInstanceResp   *mongodbflex.GetInstanceResponse
		listFlavorsFails  bool
		listFlavorsResp   *mongodbflex.ListFlavorsResponse
		listStoragesFails bool
		listStoragesResp  *mongodbflex.ListStoragesResponse
		isValid           bool
	}{
		{
			description:     "no values",
			model:           fixtureRequiredInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "update flavor from id",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.FlavorId = utils.Ptr(testFlavorId)
			}),
			isValid: true,
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: &[]mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					FlavorId: utils.Ptr(testFlavorId),
				}),
		},
		{
			description: "update flavor from cpu and ram",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.CPU = utils.Ptr(int64(2))
				model.RAM = utils.Ptr(int64(4))
			}),
			isValid: true,
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: &[]mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					FlavorId: utils.Ptr(testFlavorId),
				}),
		},
		{
			description: "update storage class only",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr("class")
			}),
			isValid: true,
			getInstanceResp: &mongodbflex.GetInstanceResponse{
				Item: &mongodbflex.Instance{
					Flavor: &mongodbflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					Storage: &mongodbflex.Storage{
						Class: utils.Ptr("class"),
					},
				}),
		},
		{
			description: "update storage class and size",
			model: fixtureRequiredInputModel(func(model *inputModel) {
				model.StorageClass = utils.Ptr("class")
				model.StorageSize = utils.Ptr(int64(10))
			}),
			isValid: true,
			getInstanceResp: &mongodbflex.GetInstanceResponse{
				Item: &mongodbflex.Instance{
					Flavor: &mongodbflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			expectedRequest: testClient.PartialUpdateInstance(testCtx, testProjectId, testInstanceId, testRegion).
				PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
					Storage: &mongodbflex.Storage{
						Class: utils.Ptr("class"),
						Size:  utils.Ptr(int64(10)),
					},
				}),
		},
		{
			description: "get flavors fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "flavor id not found",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.CPU = utils.Ptr(int64(5))
					model.RAM = utils.Ptr(int64(9))
				},
			),
			listFlavorsResp: &mongodbflex.ListFlavorsResponse{
				Flavors: &[]mongodbflex.InstanceFlavor{
					{
						Id:     utils.Ptr(testFlavorId),
						Cpu:    utils.Ptr(int64(2)),
						Memory: utils.Ptr(int64(4)),
					},
					{
						Id:     utils.Ptr("other-flavor"),
						Cpu:    utils.Ptr(int64(1)),
						Memory: utils.Ptr(int64(8)),
					},
				},
			},
			isValid: false,
		},
		{
			description: "get instance fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("class")
				},
			),
			getInstanceFails: true,
			isValid:          false,
		},
		{
			description: "get storages fails",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.FlavorId = nil
					model.CPU = utils.Ptr(int64(2))
					model.RAM = utils.Ptr(int64(4))
				},
			),
			listFlavorsFails: true,
			isValid:          false,
		},
		{
			description: "invalid storage class",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageClass = utils.Ptr("non-existing-class")
				},
			),
			getInstanceResp: &mongodbflex.GetInstanceResponse{
				Item: &mongodbflex.Instance{
					Flavor: &mongodbflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
		{
			description: "invalid storage size",
			model: fixtureRequiredInputModel(
				func(model *inputModel) {
					model.StorageSize = utils.Ptr(int64(9))
				},
			),
			getInstanceResp: &mongodbflex.GetInstanceResponse{
				Item: &mongodbflex.Instance{
					Flavor: &mongodbflex.Flavor{
						Id: utils.Ptr(testFlavorId),
					},
				},
			},
			listStoragesResp: &mongodbflex.ListStoragesResponse{
				StorageClasses: &[]string{"class"},
				StorageRange: &mongodbflex.StorageRange{
					Min: utils.Ptr(int64(10)),
					Max: utils.Ptr(int64(100)),
				},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			client := &mongoDBFlexClientMocked{
				getInstanceFails:  tt.getInstanceFails,
				getInstanceResp:   tt.getInstanceResp,
				listFlavorsFails:  tt.listFlavorsFails,
				listFlavorsResp:   tt.listFlavorsResp,
				listStoragesFails: tt.listStoragesFails,
				listStoragesResp:  tt.listStoragesResp,
			}
			request, err := buildRequest(testCtx, tt.model, client)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(tt.expectedRequest),
				cmpopts.EquateComparable(testCtx),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResult(t *testing.T) {
	type args struct {
		outputFormat  string
		async         bool
		instanceLabel string
		resp          *mongodbflex.UpdateInstanceResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "empty response",
			args: args{
				resp: &mongodbflex.UpdateInstanceResponse{},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(&params.CmdParams{Printer: p})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResult(p, tt.args.outputFormat, tt.args.async, tt.args.instanceLabel, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
