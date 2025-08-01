package utils

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"golang.org/x/mod/semver"

	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
)

// The number of replicas is enforced by the API according to the instance type
var instanceTypeToReplicas = map[string]int64{
	"Single":  1,
	"Replica": 3,
	"Sharded": 9,
}

type MongoDBFlexClient interface {
	ListVersionsExecute(ctx context.Context, projectId, region string) (*mongodbflex.ListVersionsResponse, error)
	GetInstanceExecute(ctx context.Context, projectId, instanceId, region string) (*mongodbflex.GetInstanceResponse, error)
	GetUserExecute(ctx context.Context, projectId, instanceId, userId, region string) (*mongodbflex.GetUserResponse, error)
	ListRestoreJobsExecute(ctx context.Context, projectId string, instanceId, region string) (*mongodbflex.ListRestoreJobsResponse, error)
}

func AvailableInstanceTypes() []string {
	instanceTypes := make([]string, len(instanceTypeToReplicas))
	i := 0
	for k := range instanceTypeToReplicas {
		instanceTypes[i] = k
		i++
	}
	// Dict keys aren't iterated in a consistent order
	// So we sort the array to make the output consistent
	slices.Sort(instanceTypes)
	return instanceTypes
}

func GetInstanceReplicas(instanceType string) (int64, error) {
	numReplicas, ok := instanceTypeToReplicas[instanceType]
	if !ok {
		return 0, fmt.Errorf("invalid instance type: %v", instanceType)
	}
	return numReplicas, nil
}

func GetInstanceType(numReplicas int64) (string, error) {
	for k, v := range instanceTypeToReplicas {
		if v == numReplicas {
			return k, nil
		}
	}
	return "", fmt.Errorf("invalid number of replicas: %v", numReplicas)
}

func ValidateFlavorId(flavorId string, flavors *[]mongodbflex.InstanceFlavor) error {
	if flavors == nil {
		return fmt.Errorf("nil flavors")
	}

	for _, f := range *flavors {
		if f.Id != nil && strings.EqualFold(*f.Id, flavorId) {
			return nil
		}
	}

	return &errors.DatabaseInvalidFlavorError{
		Service: "mongodbflex",
		Details: fmt.Sprintf("You provided flavor ID '%s', which is invalid.", flavorId),
	}
}

func ValidateStorage(storageClass *string, storageSize *int64, storages *mongodbflex.ListStoragesResponse, flavorId string) error {
	if storages == nil {
		return fmt.Errorf("nil storages")
	}

	if storageSize != nil {
		if *storageSize < *storages.StorageRange.Min || *storageSize > *storages.StorageRange.Max {
			return fmt.Errorf("%s", fmt.Sprintf("You provided storage size '%d', which is invalid. The valid range is %d-%d.", *storageSize, *storages.StorageRange.Min, *storages.StorageRange.Max))
		}
	}

	if storageClass == nil {
		return nil
	}

	for _, sc := range *storages.StorageClasses {
		if strings.EqualFold(*storageClass, sc) {
			return nil
		}
	}
	return &errors.DatabaseInvalidStorageError{
		Service:  "mongodbflex",
		Details:  fmt.Sprintf("You provided storage class '%s', which is invalid.", *storageClass),
		FlavorId: flavorId,
	}
}

func LoadFlavorId(cpu, ram int64, flavors *[]mongodbflex.InstanceFlavor) (*string, error) {
	if flavors == nil {
		return nil, fmt.Errorf("nil flavors")
	}

	availableFlavors := ""
	for _, f := range *flavors {
		if f.Id == nil || f.Cpu == nil || f.Memory == nil {
			continue
		}
		if *f.Cpu == cpu && *f.Memory == ram {
			return f.Id, nil
		}
		availableFlavors = fmt.Sprintf("%s\n- %d CPU, %d GB RAM", availableFlavors, *f.Cpu, *f.Cpu)
	}
	return nil, &errors.DatabaseInvalidFlavorError{
		Service: "mongodbflex",
		Details: "You provided an invalid combination for CPU and RAM.",
	}
}

func GetLatestMongoDBVersion(ctx context.Context, apiClient MongoDBFlexClient, projectId, region string) (string, error) {
	resp, err := apiClient.ListVersionsExecute(ctx, projectId, region)
	if err != nil {
		return "", fmt.Errorf("get MongoDB versions: %w", err)
	}
	versions := *resp.Versions

	latestVersion := "0"
	for i := range versions {
		oldSemVer := fmt.Sprintf("v%s", latestVersion)
		newSemVer := fmt.Sprintf("v%s", versions[i])
		if semver.Compare(newSemVer, oldSemVer) != 1 {
			continue
		}
		latestVersion = versions[i]
	}
	if latestVersion == "0" {
		return "", fmt.Errorf("no MongoDB versions found")
	}
	return latestVersion, nil
}

func GetInstanceName(ctx context.Context, apiClient MongoDBFlexClient, projectId, instanceId, region string) (string, error) {
	resp, err := apiClient.GetInstanceExecute(ctx, projectId, instanceId, region)
	if err != nil {
		return "", fmt.Errorf("get MongoDB Flex instance: %w", err)
	}
	return *resp.Item.Name, nil
}

func GetUserName(ctx context.Context, apiClient MongoDBFlexClient, projectId, instanceId, userId, region string) (string, error) {
	resp, err := apiClient.GetUserExecute(ctx, projectId, instanceId, userId, region)
	if err != nil {
		return "", fmt.Errorf("get MongoDB Flex user: %w", err)
	}
	return *resp.Item.Username, nil
}

func GetRestoreStatus(backupId string, restoreJobs *mongodbflex.ListRestoreJobsResponse) string {
	state := "-"
	if restoreJobs.Items == nil {
		return state
	}

	restoreJobsSlice := *restoreJobs.Items

	// sort array by descending date
	slices.SortFunc(restoreJobsSlice, func(i, j mongodbflex.RestoreInstanceStatus) int {
		// swap elements to sort by descending order
		return cmp.Compare(*j.Date, *i.Date)
	})

	for _, restoreJob := range *restoreJobs.Items {
		if *restoreJob.BackupID == backupId {
			state = *restoreJob.Status
			break
		}
	}
	return state
}
