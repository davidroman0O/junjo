package memory

import (
	"testing"

	"github.com/davidroman0O/junjo/types"
)

// go test -timeout 30s -v -count=1 -run ^TestGeneralTest$ ./memory
func TestGeneralTest(t *testing.T) {
	storage := NewMemoryStorage()

	// ok
	topicProvisioning, _ := storage.CreateTopic("provisioning")

	// ok
	btl, _ := storage.CreateOwner("BTL")
	network, _ := storage.CreateOwner("Network")
	ipxe, _ := storage.CreateOwner("iPXE")
	netbox, _ := storage.CreateOwner("netbox")

	// ok
	provisioningUnit, _ := storage.CreateTaskDefinition(
		"provisioning",
		btl.Key,
		types.WithTaskDefDescription("provisioning a machine"),
		types.WithTaskDefIdentifier("metal.provisioning"))
	portConfigUnit, _ := storage.CreateTaskDefinition(
		"portConfig",
		network.Key,
		types.WithTaskDefDescription("portConfig a machine"),
		types.WithTaskDefIdentifier("network.portConfig"))
	ipxeUnit, _ := storage.CreateTaskDefinition(
		"ipxe",
		ipxe.Key,
		types.WithTaskDefDescription("ipxe a machine"),
		types.WithTaskDefIdentifier("metal.ipxe"))
	netboxUpdateDeviceUnit, _ := storage.CreateTaskDefinition(
		"netboxUpdateDevice",
		netbox.Key,
		types.WithTaskDefDescription("netboxUpdateDevice a machine"),
		types.WithTaskDefIdentifier("netbox.update.dcimdevice"))

	task, _ := storage.CreateTask()

	unitNetboxDcim, _ := storage.CreateTaskUnit(netboxUpdateDeviceUnit.Key, []types.TaskUnitID{})

	unitIpxe, _ := storage.CreateTaskUnit(ipxeUnit.Key, []types.TaskUnitID{
		unitNetboxDcim.Key,
	})

	unitProv, _ := storage.CreateTaskUnit(provisioningUnit.Key, []types.TaskUnitID{
		unitIpxe.Key,
		unitNetboxDcim.Key,
	})

	unitiPort, _ := storage.CreateTaskUnit(portConfigUnit.Key, []types.TaskUnitID{
		unitProv.Key,
	})

	unitNetbox, _ := storage.CreateTaskUnit(netboxUpdateDeviceUnit.Key, []types.TaskUnitID{
		unitiPort.Key,
	})

	if err := storage.CreateUnitOnTask(task.Key, []types.TaskUnit{
		*unitNetboxDcim,
		*unitIpxe,
		*unitProv,
		*unitiPort,
		*unitNetbox,
	}); err != nil {
		t.Error(err)
	}

	if _, err := storage.CreateJobWithTasks(topicProvisioning.Key, []types.Task{
		*task,
	}); err != nil {
		t.Error(err)
	}

	storage.PrintTree()

	if err := storage.UpdateTaskUnitStatus(unitIpxe.Key, types.ProgressStatus, nil); err != nil {
		t.Error(err)
	}

	storage.PrintTree()

	if err := storage.UpdateTaskUnitStatus(unitProv.Key, types.ProgressStatus, nil); err != nil {
		t.Error(err)
	}

	storage.PrintTree()
	storage.PrintDAG()
}
