package utils

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type IaaSClientMocked struct {
	GetSecurityGroupRuleFails bool
	GetSecurityGroupRuleResp  *iaas.SecurityGroupRule
	GetSecurityGroupFails     bool
	GetSecurityGroupResp      *iaas.SecurityGroup
	GetPublicIpFails          bool
	GetPublicIpResp           *iaas.PublicIp
	GetServerFails            bool
	GetServerResp             *iaas.Server
	GetVolumeFails            bool
	GetVolumeResp             *iaas.Volume
	GetNetworkFails           bool
	GetNetworkResp            *iaas.Network
	GetNetworkAreaFails       bool
	GetNetworkAreaResp        *iaas.NetworkArea
	GetAttachedProjectsFails  bool
	GetAttachedProjectsResp   *iaas.ProjectListResponse
	GetNetworkAreaRangeFails  bool
	GetNetworkAreaRangeResp   *iaas.NetworkRange
	GetImageFails             bool
	GetImageResp              *iaas.Image
	GetAffinityGroupsFails    bool
	GetAffinityGroupResp      *iaas.AffinityGroup
	GetBackupFails            bool
	GetBackupResp             *iaas.Backup
	GetSnapshotFails          bool
	GetSnapshotResp           *iaas.Snapshot
}

func (m *IaaSClientMocked) GetAffinityGroupExecute(_ context.Context, _, _ string) (*iaas.AffinityGroup, error) {
	if m.GetAffinityGroupsFails {
		return nil, fmt.Errorf("could not get affinity groups")
	}
	return m.GetAffinityGroupResp, nil
}

func (m *IaaSClientMocked) GetSecurityGroupRuleExecute(_ context.Context, _, _, _ string) (*iaas.SecurityGroupRule, error) {
	if m.GetSecurityGroupRuleFails {
		return nil, fmt.Errorf("could not get security group rule")
	}
	return m.GetSecurityGroupRuleResp, nil
}

func (m *IaaSClientMocked) GetSecurityGroupExecute(_ context.Context, _, _ string) (*iaas.SecurityGroup, error) {
	if m.GetSecurityGroupFails {
		return nil, fmt.Errorf("could not get security group")
	}
	return m.GetSecurityGroupResp, nil
}

func (m *IaaSClientMocked) GetPublicIPExecute(_ context.Context, _, _ string) (*iaas.PublicIp, error) {
	if m.GetPublicIpFails {
		return nil, fmt.Errorf("could not get public ip")
	}
	return m.GetPublicIpResp, nil
}

func (m *IaaSClientMocked) GetServerExecute(_ context.Context, _, _ string) (*iaas.Server, error) {
	if m.GetServerFails {
		return nil, fmt.Errorf("could not get server")
	}
	return m.GetServerResp, nil
}

func (m *IaaSClientMocked) GetVolumeExecute(_ context.Context, _, _ string) (*iaas.Volume, error) {
	if m.GetVolumeFails {
		return nil, fmt.Errorf("could not get volume")
	}
	return m.GetVolumeResp, nil
}

func (m *IaaSClientMocked) GetNetworkExecute(_ context.Context, _, _ string) (*iaas.Network, error) {
	if m.GetNetworkFails {
		return nil, fmt.Errorf("could not get network")
	}
	return m.GetNetworkResp, nil
}

func (m *IaaSClientMocked) GetNetworkAreaExecute(_ context.Context, _, _ string) (*iaas.NetworkArea, error) {
	if m.GetNetworkAreaFails {
		return nil, fmt.Errorf("could not get network area")
	}
	return m.GetNetworkAreaResp, nil
}

func (m *IaaSClientMocked) ListNetworkAreaProjectsExecute(_ context.Context, _, _ string) (*iaas.ProjectListResponse, error) {
	if m.GetAttachedProjectsFails {
		return nil, fmt.Errorf("could not get attached projects")
	}
	return m.GetAttachedProjectsResp, nil
}

func (m *IaaSClientMocked) GetNetworkAreaRangeExecute(_ context.Context, _, _, _ string) (*iaas.NetworkRange, error) {
	if m.GetNetworkAreaRangeFails {
		return nil, fmt.Errorf("could not get network range")
	}
	return m.GetNetworkAreaRangeResp, nil
}

func (m *IaaSClientMocked) GetImageExecute(_ context.Context, _, _ string) (*iaas.Image, error) {
	if m.GetImageFails {
		return nil, fmt.Errorf("could not get image")
	}
	return m.GetImageResp, nil
}

func (m *IaaSClientMocked) GetBackupExecute(_ context.Context, _, _ string) (*iaas.Backup, error) {
	if m.GetBackupFails {
		return nil, fmt.Errorf("could not get backup")
	}
	return m.GetBackupResp, nil
}

func (m *IaaSClientMocked) GetSnapshotExecute(_ context.Context, _, _ string) (*iaas.Snapshot, error) {
	if m.GetSnapshotFails {
		return nil, fmt.Errorf("could not get snapshot")
	}
	return m.GetSnapshotResp, nil
}
func TestGetSecurityGroupRuleName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.SecurityGroupRule
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.SecurityGroupRule{
					Ethertype: utils.Ptr("IPv6"),
					Direction: utils.Ptr("ingress"),
				},
			},
			want: "IPv6, ingress",
		},
		{
			name: "get security group rule fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetSecurityGroupRuleFails: tt.args.getInstanceFails,
				GetSecurityGroupRuleResp:  tt.args.getInstanceResp,
			}
			got, err := GetSecurityGroupRuleName(context.Background(), m, "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityGroupRuleName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecurityGroupRuleName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSecurityGroupName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.SecurityGroup
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.SecurityGroup{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get security group fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "name in response is nil",
			args: args{
				getInstanceResp: &iaas.SecurityGroup{
					Name: nil,
				},
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetSecurityGroupFails: tt.args.getInstanceFails,
				GetSecurityGroupResp:  tt.args.getInstanceResp,
			}
			got, err := GetSecurityGroupName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSecurityGroupName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSecurityGroupName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPublicIp(t *testing.T) {
	type args struct {
		getPublicIpFails bool
		getPublicIpResp  *iaas.PublicIp
	}
	tests := []struct {
		name                   string
		args                   args
		wantPublicIp           string
		wantAssociatedResource string
		wantErr                bool
	}{
		{
			name: "base",
			args: args{
				getPublicIpResp: &iaas.PublicIp{
					Ip:               utils.Ptr("1.2.3.4"),
					NetworkInterface: iaas.NewNullableString(utils.Ptr("5.6.7.8")),
				},
			},
			wantPublicIp:           "1.2.3.4",
			wantAssociatedResource: "5.6.7.8",
		},
		{
			name: "get public ip fails",
			args: args{
				getPublicIpFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetPublicIpFails: tt.args.getPublicIpFails,
				GetPublicIpResp:  tt.args.getPublicIpResp,
			}
			gotPublicIP, gotAssociatedResource, err := GetPublicIP(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPublicIP != tt.wantPublicIp {
				t.Errorf("GetPublicIP() = %v, want public IP %v", gotPublicIP, tt.wantPublicIp)
			}
			if gotAssociatedResource != tt.wantAssociatedResource {
				t.Errorf("GetPublicIP() = %v, want associated resource %v", gotAssociatedResource, tt.wantAssociatedResource)
			}
		})
	}
}

func TestGetServerName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Server
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.Server{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get server fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetServerFails: tt.args.getInstanceFails,
				GetServerResp:  tt.args.getInstanceResp,
			}
			got, err := GetServerName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServerName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetServerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVolumeName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Volume
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.Volume{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get volume fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "name in response is nil",
			args: args{
				getInstanceResp: &iaas.Volume{
					Name: nil,
				},
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetVolumeFails: tt.args.getInstanceFails,
				GetVolumeResp:  tt.args.getInstanceResp,
			}
			got, err := GetVolumeName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVolumeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.Network
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.Network{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get network fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "name in response is nil",
			args: args{
				getInstanceResp: &iaas.Network{
					Name: nil,
				},
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkFails: tt.args.getInstanceFails,
				GetNetworkResp:  tt.args.getInstanceResp,
			}
			got, err := GetNetworkName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkAreaName(t *testing.T) {
	type args struct {
		getInstanceFails bool
		getInstanceResp  *iaas.NetworkArea
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getInstanceResp: &iaas.NetworkArea{
					Name: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get network area fails",
			args: args{
				getInstanceFails: true,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "response is nil",
			args: args{
				getInstanceResp:  nil,
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "name in response is nil",
			args: args{
				getInstanceResp: &iaas.NetworkArea{
					Name: nil,
				},
				getInstanceFails: false,
			},
			wantErr: true,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkAreaFails: tt.args.getInstanceFails,
				GetNetworkAreaResp:  tt.args.getInstanceResp,
			}
			got, err := GetNetworkAreaName(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkAreaName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkAreaName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListAttachedProjects(t *testing.T) {
	type args struct {
		getAttachedProjectsFails bool
		getAttachedProjectsResp  *iaas.ProjectListResponse
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getAttachedProjectsResp: &iaas.ProjectListResponse{
					Items: &[]string{"test"},
				},
			},
			want: []string{"test"},
		},
		{
			name: "get attached projects fails",
			args: args{
				getAttachedProjectsFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetAttachedProjectsFails: tt.args.getAttachedProjectsFails,
				GetAttachedProjectsResp:  tt.args.getAttachedProjectsResp,
			}
			got, err := ListAttachedProjects(context.Background(), m, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttachedProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("GetAttachedProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkRangePrefix(t *testing.T) {
	type args struct {
		getNetworkAreaRangeFails bool
		getNetworkAreaRangeResp  *iaas.NetworkRange
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				getNetworkAreaRangeResp: &iaas.NetworkRange{
					Prefix: utils.Ptr("test"),
				},
			},
			want: "test",
		},
		{
			name: "get network area range fails",
			args: args{
				getNetworkAreaRangeFails: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &IaaSClientMocked{
				GetNetworkAreaRangeFails: tt.args.getNetworkAreaRangeFails,
				GetNetworkAreaRangeResp:  tt.args.getNetworkAreaRangeResp,
			}
			got, err := GetNetworkRangePrefix(context.Background(), m, "", "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkRangePrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNetworkRangePrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRouteFromAPIResponse(t *testing.T) {
	type args struct {
		prefix  string
		nexthop string
		routes  *[]iaas.Route
	}
	tests := []struct {
		name    string
		args    args
		want    iaas.Route
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes: &[]iaas.Route{
					{
						Prefix:  utils.Ptr("1.1.1.0/24"),
						Nexthop: utils.Ptr("1.1.1.1"),
					},
					{
						Prefix:  utils.Ptr("2.2.2.0/24"),
						Nexthop: utils.Ptr("2.2.2.2"),
					},
					{
						Prefix:  utils.Ptr("3.3.3.0/24"),
						Nexthop: utils.Ptr("3.3.3.3"),
					},
				},
			},
			want: iaas.Route{
				Prefix:  utils.Ptr("1.1.1.0/24"),
				Nexthop: utils.Ptr("1.1.1.1"),
			},
		},
		{
			name: "not found",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes: &[]iaas.Route{
					{
						Prefix:  utils.Ptr("2.2.2.0/24"),
						Nexthop: utils.Ptr("2.2.2.2"),
					},
					{
						Prefix:  utils.Ptr("3.3.3.0/24"),
						Nexthop: utils.Ptr("3.3.3.3"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				prefix:  "1.1.1.0/24",
				nexthop: "1.1.1.1",
				routes:  &[]iaas.Route{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRouteFromAPIResponse(tt.args.prefix, tt.args.nexthop, tt.args.routes)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRouteFromAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRouteFromAPIResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNetworkRangeFromAPIResponse(t *testing.T) {
	type args struct {
		prefix        string
		networkRanges *[]iaas.NetworkRange
	}
	tests := []struct {
		name    string
		args    args
		want    iaas.NetworkRange
		wantErr bool
	}{
		{
			name: "base",
			args: args{
				prefix: "1.1.1.0/24",
				networkRanges: &[]iaas.NetworkRange{
					{
						Prefix: utils.Ptr("1.1.1.0/24"),
					},
					{
						Prefix: utils.Ptr("2.2.2.0/24"),
					},
					{
						Prefix: utils.Ptr("3.3.3.0/24"),
					},
				},
			},
			want: iaas.NetworkRange{
				Prefix: utils.Ptr("1.1.1.0/24"),
			},
		},
		{
			name: "not found",
			args: args{
				prefix: "1.1.1.0/24",
				networkRanges: &[]iaas.NetworkRange{
					{
						Prefix: utils.Ptr("2.2.2.0/24"),
					},
					{
						Prefix: utils.Ptr("3.3.3.0/24"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				prefix:        "1.1.1.0/24",
				networkRanges: &[]iaas.NetworkRange{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNetworkRangeFromAPIResponse(tt.args.prefix, tt.args.networkRanges)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkRangeFromAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNetworkRangeFromAPIResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetImageName(t *testing.T) {
	tests := []struct {
		name      string
		imageResp *iaas.Image
		imageErr  bool
		want      string
		wantErr   bool
	}{
		{
			name:      "successful retrieval",
			imageResp: &iaas.Image{Name: utils.Ptr("test-image")},
			want:      "test-image",
			wantErr:   false,
		},
		{
			name:     "error on retrieval",
			imageErr: true,
			wantErr:  true,
		},
		{
			name:      "response is nil",
			imageErr:  false,
			imageResp: nil,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "name in response is nil",
			imageErr:  false,
			imageResp: &iaas.Image{Name: nil},
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &IaaSClientMocked{
				GetImageFails: tt.imageErr,
				GetImageResp:  tt.imageResp,
			}
			got, err := GetImageName(context.Background(), client, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAffinityGroupName(t *testing.T) {
	tests := []struct {
		name         string
		affinityResp *iaas.AffinityGroup
		affinityErr  bool
		want         string
		wantErr      bool
	}{
		{
			name:         "successful retrieval",
			affinityResp: &iaas.AffinityGroup{Name: utils.Ptr("test-affinity")},
			want:         "test-affinity",
			wantErr:      false,
		},
		{
			name:        "error on retrieval",
			affinityErr: true,
			wantErr:     true,
		},
		{
			name:        "response is nil",
			affinityErr: false,
			affinityResp: &iaas.AffinityGroup{
				Name: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name:        "affinity group name in response is nil",
			affinityErr: false,
			affinityResp: &iaas.AffinityGroup{
				Name: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client := &IaaSClientMocked{
				GetAffinityGroupsFails: tt.affinityErr,
				GetAffinityGroupResp:   tt.affinityResp,
			}
			got, err := GetAffinityGroupName(ctx, client, "", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAffinityGroupName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAffinityGroupName() = %v, want %v", got, tt.want)
			}
		})
	}
}
