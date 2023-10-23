package junjo

import (
	"fmt"
	"testing"
	"time"

	"github.com/davidroman0O/junjo/dag"
	"github.com/davidroman0O/junjo/memory"
	"github.com/davidroman0O/junjo/types"
)

// Let's try to modelize the topic "cooking a meal" of a kitchen which orchestrate different owner of `UnitDescription`
// https://en.wikipedia.org/wiki/Brigade_de_cuisine
// go test -timeout 30s -v -count=1 -run ^TestBrigadeDeCuisine$ .
func TestBrigadeDeCuisine(t *testing.T) {

	now := time.Now()

	storage := memory.NewMemoryStorage()

	jj := New(storage)
	var err error

	var topic *types.Topic
	if topic, err = jj.CreateTopic(
		"Meal",
		types.WithTopicDescription("cooking a meal")); err != nil {
		t.Error(err)
		return
	}

	// fmt.Println(topic)

	var chief *types.Owner
	if chief, err = jj.CreateOwner(
		"Chief",
		types.WithOwnerDescription("Chief of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var deputy *types.Owner
	if deputy, err = jj.CreateOwner(
		"Deputy",
		types.WithOwnerDescription("Deputy of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var sauceMaker *types.Owner
	if sauceMaker, err = jj.CreateOwner(
		"Sauce Maker",
		types.WithOwnerDescription("Sauce Maker of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var cook *types.Owner
	if cook, err = jj.CreateOwner(
		"Cook",
		types.WithOwnerDescription("Cook of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var juniorCook *types.Owner
	if juniorCook, err = jj.CreateOwner(
		"Junior Cook",
		types.WithOwnerDescription("Junior cook of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var dishwaher *types.Owner
	if dishwaher, err = jj.CreateOwner(
		"Dishwaher",
		types.WithOwnerDescription("Dishwaher of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var barker *types.Owner
	if barker, err = jj.CreateOwner(
		"Barker",
		types.WithOwnerDescription("Barker of the kitchen")); err != nil {
		t.Error(err)
		return
	}

	var unitWorkBarkerGiveNewMeal *types.TaskDefinition
	if unitWorkBarkerGiveNewMeal, err = jj.CreateTaskDefinition(
		"give new meal",
		barker.Key,
		types.WithTaskDefDescription("new meal to be prepared"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkChiefMonitorMeal *types.TaskDefinition
	if unitWorkChiefMonitorMeal, err = jj.CreateTaskDefinition(
		"look at meal order",
		chief.Key,
		types.WithTaskDefDescription("just looking don't mind me"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkDeputyCheckMealProgress *types.TaskDefinition
	if unitWorkDeputyCheckMealProgress, err = jj.CreateTaskDefinition(
		"verify meal progression",
		deputy.Key,
		types.WithTaskDefDescription("often have to check the quality of the preparation"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkSauceMakerCookSauce *types.TaskDefinition
	if unitWorkSauceMakerCookSauce, err = jj.CreateTaskDefinition(
		"cook sauce",
		sauceMaker.Key,
		types.WithTaskDefDescription("sauce marker will make a great sauce, not all the time though"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkJuniorCookVegetables *types.TaskDefinition
	if unitWorkJuniorCookVegetables, err = jj.CreateTaskDefinition(
		"cook vegetables",
		juniorCook.Key,
		types.WithTaskDefDescription("junior will make vegetables"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkCookCookMeat *types.TaskDefinition
	if unitWorkCookCookMeat, err = jj.CreateTaskDefinition(
		"cook meat",
		cook.Key,
		types.WithTaskDefDescription("cook will make meat"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkDeputyAssembleMeal *types.TaskDefinition
	if unitWorkDeputyAssembleMeal, err = jj.CreateTaskDefinition(
		"assemble meal",
		deputy.Key,
		types.WithTaskDefDescription("assemble the meal with sauce, meat (else) and vegetables"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkDeputyNotifyBarker *types.TaskDefinition
	if unitWorkDeputyNotifyBarker, err = jj.CreateTaskDefinition(
		"notify meal complete",
		deputy.Key,
		types.WithTaskDefDescription("barker need to be ready to take the meal"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkerBarkerDistributeMeal *types.TaskDefinition
	if unitWorkerBarkerDistributeMeal, err = jj.CreateTaskDefinition(
		"give meal in room",
		barker.Key,
		types.WithTaskDefDescription("barker will give meal to client"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkBarkerWaitMealEaten *types.TaskDefinition
	if unitWorkBarkerWaitMealEaten, err = jj.CreateTaskDefinition(
		"waiting until eaten",
		barker.Key,
		types.WithTaskDefDescription("barker will monitor the completion of the meal"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkBarkerTakePlateToDishwasher *types.TaskDefinition
	if unitWorkBarkerTakePlateToDishwasher, err = jj.CreateTaskDefinition(
		"take plate to dishwasher",
		barker.Key,
		types.WithTaskDefDescription("barker will take the plate and give it to dishwasher"),
	); err != nil {
		t.Error(err)
		return
	}

	var unitWorkDishwasherCleanPlate *types.TaskDefinition
	if unitWorkDishwasherCleanPlate, err = jj.CreateTaskDefinition(
		"clean plate",
		dishwaher.Key,
		types.WithTaskDefDescription("dishwasher clean the plate"),
	); err != nil {
		t.Error(err)
		return
	}

	var workUnitDag *types.WorkUnitDag
	workUnitDag = jj.CreateDagTaskUnits()

	// Maybe we should return a function that will create a new vertex each time we connect for description
	defBarkerGiveNewMeal := workUnitDag.AddTaskDefinition(unitWorkBarkerGiveNewMeal)
	defChiefMonitorMeal := workUnitDag.AddTaskDefinition(unitWorkChiefMonitorMeal)
	defDeputyCheckMealProgress := workUnitDag.AddTaskDefinition(unitWorkDeputyCheckMealProgress)
	defSauceMakerCookSauce := workUnitDag.AddTaskDefinition(unitWorkSauceMakerCookSauce)
	defJuniorCookVegetables := workUnitDag.AddTaskDefinition(unitWorkJuniorCookVegetables)
	defCookCookMeat := workUnitDag.AddTaskDefinition(unitWorkCookCookMeat)
	defDeputyAssembleMeal := workUnitDag.AddTaskDefinition(unitWorkDeputyAssembleMeal)
	defDeputyNotifyBarker := workUnitDag.AddTaskDefinition(unitWorkDeputyNotifyBarker)
	deferBarkerDistributeMeal := workUnitDag.AddTaskDefinition(unitWorkerBarkerDistributeMeal)
	defBarkerWaitMealEaten := workUnitDag.AddTaskDefinition(unitWorkBarkerWaitMealEaten)
	defBarkerTakePlateToDishwasher := workUnitDag.AddTaskDefinition(unitWorkBarkerTakePlateToDishwasher)
	defDishwasherCleanPlate := workUnitDag.AddTaskDefinition(unitWorkDishwasherCleanPlate)

	// Barker create a new meal in the DAG but the Chief want to check it first
	giveMeal, monitorOrder := workUnitDag.ConnectDef(defBarkerGiveNewMeal, defChiefMonitorMeal)

	// Oncethe chief finished to look at the meal, the deputy will check it progress
	deputyCheckMealProgress := workUnitDag.AssociateDef(monitorOrder, defDeputyCheckMealProgress)

	sauceMakerCookSauce := workUnitDag.AssociateDef(deputyCheckMealProgress, defSauceMakerCookSauce)   // then the sauce can start
	juniorCookVegetables := workUnitDag.AssociateDef(deputyCheckMealProgress, defJuniorCookVegetables) // then the vege can start
	cookCookMeat := workUnitDag.AssociateDef(deputyCheckMealProgress, defCookCookMeat)                 // then the cook can start

	// Create new vertex for new meal check progression while waiting for all cooks to finish their work
	deputySecondCheckProgress := workUnitDag.MConnectDef(defDeputyCheckMealProgress, sauceMakerCookSauce, juniorCookVegetables, cookCookMeat)

	// One it's fine, let's assemble the meal
	deputyAssembleMeal := workUnitDag.AssociateDef(deputySecondCheckProgress, defDeputyAssembleMeal)

	notifyBarker := workUnitDag.AssociateDef(deputyAssembleMeal, defDeputyNotifyBarker)
	mealToRoom := workUnitDag.AssociateDef(notifyBarker, deferBarkerDistributeMeal)
	waitMeal := workUnitDag.AssociateDef(mealToRoom, defBarkerWaitMealEaten)
	dish := workUnitDag.AssociateDef(waitMeal, defBarkerTakePlateToDishwasher)

	// final state of the DAG
	workUnitDag.AssociateDef(dish, defDishwasherCleanPlate)

	// print just for fun
	// printProgression(workUnitDag)

	available := workUnitDag.AvailableNodeUnit()
	if len(available) > 1 || len(available) == 0 {
		fmt.Println("instead have", len(available))
		t.Error(fmt.Errorf("should be 1"))
		return
	}

	// Check if description is correct
	if available[0].Definition.Key != unitWorkBarkerGiveNewMeal.Key {
		fmt.Println("instead have", available[0].Definition.Key)
		t.Error(fmt.Errorf("should be giving meal"))
		return
	}

	// From DAG generate the TaskUnits we need to create a task
	var taskUnits []*types.TaskUnit
	if taskUnits, err = workUnitDag.ToTaskUnits(); err != nil {
		t.Error(err)
		return
	}

	if len(taskUnits) == 0 {
		t.Error(fmt.Errorf("should have task units"))
		return
	}

	ids := []types.TaskUnitID{}
	// You need to create and store those draft TaskUnits
	if ids, err = jj.CreateTaskUnits(taskUnits); err != nil {
		t.Error(err)
		return
	}

	var task *types.Task
	if task, err = jj.CreateTask(); err != nil {
		t.Error(err)
		return
	}

	var job *types.Job
	if job, err = jj.CreateJob(); err != nil {
		t.Error(err)
		return
	}

	// Attach drafts DAG to Task
	if err := jj.AssignTaskUnits(task.Key, ids); err != nil {
		t.Error(err)
		return
	}

	// Attach draft Task to Job
	if err := jj.AssignTask(job.Key, task.Key); err != nil {
		t.Error(err)
		return
	}

	// Attach Job to Topic for processing
	if err := jj.AssignJob(topic.Key, job.Key); err != nil {
		t.Error(err)
		return
	}

	// fmt.Println("task", task.JobID)

	var ancestorsGiveMeal []dag.Vertex
	if ancestorsGiveMeal, err = workUnitDag.Graph().ImmediateAncestors(giveMeal); err != nil {
		t.Error(err)
		return
	}

	if len(ancestorsGiveMeal) != 0 {
		fmt.Println(ancestorsGiveMeal)
		t.Error(fmt.Errorf("should be zero"))
		return
	}

	// storage.Print()

	//	Now let's create a Job + Task + TaskUnits
	// fmt.Println("first work to do", available[0].Definition.Name, "for", available[0].Definition.OwnerID)
	// fmt.Println("search work for", barker.Key)
	inboxBarker := []types.InboxTaskUnit{}
	if inboxBarker, err = jj.GetInbox(barker.Key); err != nil {
		t.Error(err)
		return
	}

	// fmt.Println("inbox", len(inboxBarker)) // should be 1 though
	if len(inboxBarker) == 0 {
		t.Error(fmt.Errorf("should be one"))
		return
	}
	fmt.Println(time.Since(now))
}

func printProgression(workUnitDag *types.WorkUnitDag) {
	available := workUnitDag.AvailableNodeUnit()
	for i := 0; i < len(available); i++ {
		fmt.Println(available[i].Definition.Name)
	}
}
