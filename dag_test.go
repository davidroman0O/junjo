package junjo

import (
	"fmt"
	"testing"

	"github.com/davidroman0O/junjo/dag"
	"github.com/davidroman0O/junjo/memory"
	"github.com/davidroman0O/junjo/types"
)

// go test -timeout 30s -v -count=1 -run ^TestEdges$ .
func TestEdges(t *testing.T) {
	storage := memory.NewMemoryStorage()
	var err error
	jj := New(storage)

	var vertexA *types.Owner
	if vertexA, err = jj.CreateOwner("vertexA"); err != nil {
		t.Error(err)
		return
	}

	var unitWorkA *types.TaskDefinition
	if unitWorkA, err = jj.CreateTaskDefinition(
		"start",
		vertexA.Key,
	); err != nil {
		t.Error(err)
		return
	}

	var vertexB *types.Owner
	if vertexB, err = jj.CreateOwner("vertexB"); err != nil {
		t.Error(err)
		return
	}

	var unitWorkB *types.TaskDefinition
	if unitWorkB, err = jj.CreateTaskDefinition(
		"end",
		vertexB.Key,
	); err != nil {
		t.Error(err)
		return
	}

	unitDag := types.NewWorkUnitDag(storage)

	defA := unitDag.AddTaskDefinition(unitWorkA)
	defB := unitDag.AddTaskDefinition(unitWorkB)

	source, target := unitDag.ConnectDef(defA, defB)

	available := unitDag.AvailableNodeUnit()
	if len(available) == 0 {
		t.Error("should have one")
		return
	}

	var ancestorsSource []dag.Vertex
	var ancestorsTarget []dag.Vertex

	if ancestorsSource, err = unitDag.Graph().
		ImmediateAncestors(source); err != nil {
		t.Error(err)
		return
	}
	if ancestorsTarget, err = unitDag.Graph().
		ImmediateAncestors(target); err != nil {
		t.Error(err)
		return
	}

	if len(ancestorsSource) != 0 {
		t.Error(fmt.Errorf("should not be different than zero"))
	}

	if len(ancestorsTarget) == 0 {
		t.Error(fmt.Errorf("should should be at least one"))
	}
}

// // go test -timeout 30s -v -count=1 -run ^TestGeneral$ ./dag
// func TestGeneral(t *testing.T) {

// 	memoryStorage := memory.NewMemoryStorage()

// 	// ok
// 	btl, _ := memoryStorage.CreateOwner("BTL")
// 	network, _ := memoryStorage.CreateOwner("Network")
// 	ipxe, _ := memoryStorage.CreateOwner("iPXE")
// 	netbox, _ := memoryStorage.CreateOwner("netbox")

// 	provisioningUnit, _ := memoryStorage.CreateTaskDefinition(
// 		"provisioning",
// 		btl.Key,
// 		types.WithTaskDefDescription("provisioning a machine"),
// 		types.WithTaskDefIdentifier("metal.provisioning"))
// 	portConfigUnit, _ := memoryStorage.CreateTaskDefinition(
// 		"portConfig",
// 		network.Key,
// 		types.WithTaskDefDescription("portConfig a machine"),
// 		types.WithTaskDefIdentifier("network.portConfig"))
// 	ipxeUnit, _ := memoryStorage.CreateTaskDefinition(
// 		"ipxe",
// 		ipxe.Key,
// 		types.WithTaskDefDescription("ipxe a machine"),
// 		types.WithTaskDefIdentifier("metal.ipxe"))
// 	netboxUpdateDeviceStagedUnit, _ := memoryStorage.CreateTaskDefinition(
// 		"netboxUpdateDeviceStaged",
// 		netbox.Key,
// 		types.WithTaskDefDescription("netboxUpdateDevice a machine"),
// 		types.WithTaskDefIdentifier("netbox.update.dcimdevice"))
// 	netboxUpdateDeviceInventoryUnit, _ := memoryStorage.CreateTaskDefinition(
// 		"netboxUpdateDeviceInventory",
// 		netbox.Key,
// 		types.WithTaskDefDescription("netboxUpdateDevice a machine"),
// 		types.WithTaskDefIdentifier("netbox.update.dcimdevice"))
// 	logNotification, _ := memoryStorage.CreateTaskDefinition(
// 		"log",
// 		btl.Key,
// 		types.WithTaskDefDescription("log machine"),
// 		types.WithTaskDefIdentifier("metal.log"))

// 	dag := NewWorkUnitDag(memoryStorage)

// 	netboxUpdateInventory := dag.AddTaskUnit(
// 		WithUnitDescription(netboxUpdateDeviceInventoryUnit),
// 		WithStatus(types.NoneStatus))
// 	netboxUpdateStaged := dag.AddTaskUnit(
// 		WithUnitDescription(netboxUpdateDeviceStagedUnit),
// 		WithStatus(types.NoneStatus))
// 	logThing := dag.AddTaskUnit(
// 		WithUnitDescription(logNotification),
// 		WithStatus(types.NoneStatus))
// 	provisioningUbuntu := dag.AddTaskUnit(
// 		WithUnitDescription(provisioningUnit),
// 		WithStatus(types.NoneStatus),
// 		WithData(map[string]string{
// 			"machine": "uuid",
// 		}))
// 	ipxeBoot := dag.AddTaskUnit(
// 		WithUnitDescription(ipxeUnit),
// 		WithStatus(types.NoneStatus))
// 	networkConfigure := dag.AddTaskUnit(
// 		WithUnitDescription(portConfigUnit),
// 		WithStatus(types.NoneStatus))

// 	dag.Connect(netboxUpdateInventory, ipxeBoot)
// 	dag.Connect(ipxeBoot, provisioningUbuntu)
// 	dag.Connect(provisioningUbuntu, logThing)
// 	dag.Connect(logThing, netboxUpdateStaged)
// 	dag.Connect(provisioningUbuntu, networkConfigure)
// 	dag.Connect(networkConfigure, netboxUpdateStaged)

// 	dag.Print()

// 	fmt.Println("====")

// 	// fmt.Println(dag.AvailableNodeUnit())
// 	thing := netboxUpdateInventory.(*NodeUnit)
// 	thing.status = types.SuccessStatus
// 	ipxeBootThing := ipxeBoot.(*NodeUnit)
// 	ipxeBootThing.status = types.SuccessStatus
// 	// provisioningUbuntuThiong := provisioningUbuntu.(*NodeUnit)
// 	// provisioningUbuntuThiong.status = types.SuccessStatus

// 	yes, err := dag.CanChangeStatus(networkConfigure)
// 	fmt.Println(yes, err)

// 	yes, err = dag.CanChangeStatusWithOwner(networkConfigure, btl.Key)
// 	fmt.Println(yes, err)

// 	nodes := dag.AvailableNodeUnit()
// 	fmt.Println("len ", len(nodes))
// 	for idx, node := range nodes {
// 		fmt.Println(idx, node.GetDescription().Identifier, node)
// 	}
// }
