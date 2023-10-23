package memory

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/davidroman0O/junjo/types"
)

/// `MemoryStorage` is an implementation that is intented for single applications or for your unit tests, if you need real databases storage then you should consider using other `junjo` repositories for storage

type MemoryStorage struct {
	mu          sync.RWMutex
	topics      map[types.TopicID]*types.Topic
	jobs        map[types.JobID]*types.Job
	tasks       map[types.TaskID]*types.Task
	units       map[types.TaskUnitID]*types.TaskUnit
	definitions map[types.TaskDefinitionID]*types.TaskDefinition
	owners      map[types.OwnerID]*types.Owner
}

func (ms *MemoryStorage) Print() {
	ownersKeys := make([]types.OwnerID, 0, len(ms.owners))
	for oi := range ms.owners {
		ownersKeys = append(ownersKeys, oi)
	}
	for i := 0; i < len(ownersKeys); i++ {
		fmt.Println("owner", ms.owners[ownersKeys[i]].Key)
	}

	definitionsKeys := make([]types.TaskDefinitionID, 0, len(ms.definitions))
	for oi := range ms.definitions {
		definitionsKeys = append(definitionsKeys, oi)
	}
	for i := 0; i < len(definitionsKeys); i++ {
		fmt.Println("definition", ms.definitions[definitionsKeys[i]].Key)
	}

	unitsKeys := make([]types.TaskUnitID, 0, len(ms.units))
	for oi := range ms.units {
		unitsKeys = append(unitsKeys, oi)
	}
	for i := 0; i < len(unitsKeys); i++ {
		fmt.Println("unit", ms.units[unitsKeys[i]].Key)
	}

	tasksKeys := make([]types.TaskID, 0, len(ms.tasks))
	for oi := range ms.tasks {
		tasksKeys = append(tasksKeys, oi)
	}
	for i := 0; i < len(tasksKeys); i++ {
		fmt.Println("task", ms.tasks[tasksKeys[i]].Key)
	}

	jobsKeys := make([]types.JobID, 0, len(ms.jobs))
	for oi := range ms.jobs {
		jobsKeys = append(jobsKeys, oi)
	}
	for i := 0; i < len(jobsKeys); i++ {
		fmt.Println("job", ms.jobs[jobsKeys[i]].Key)
	}

	topicsKeys := make([]types.TopicID, 0, len(ms.topics))
	for oi := range ms.topics {
		topicsKeys = append(topicsKeys, oi)
	}
	for i := 0; i < len(topicsKeys); i++ {
		fmt.Println("topic", ms.topics[topicsKeys[i]].Key)
	}
}

func (ms *MemoryStorage) PrintTree() {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var builder strings.Builder

	builder.WriteString("Topics:\n")

	topics, _ := ms.GetTopics() // Assuming this function exists and works

	topicKeys := make([]types.TopicID, 0, len(ms.topics))
	for ti := range ms.topics {
		topicKeys = append(topicKeys, ti)
	}

	for i := 0; i < len(topicKeys); i++ {
		topic := topics[i] // Assuming topics is indexed by integers and corresponds with topicKeys
		builder.WriteString(fmt.Sprintf("- %s\n  Jobs:\n", topic.Name))

		jobKeys := make([]types.JobID, 0, len(topic.Jobs))
		for k := range topic.Jobs {
			jobKeys = append(jobKeys, k)
		}

		for j := 0; j < len(jobKeys); j++ {
			job := topic.Jobs[jobKeys[j]]
			builder.WriteString(fmt.Sprintf("  - %s (%s)\n    Tasks:\n", jobKeys[j], job.Status))

			taskKeys := make([]types.TaskID, 0, len(job.Tasks))
			for k := range job.Tasks {
				taskKeys = append(taskKeys, k)
			}

			for k := 0; k < len(taskKeys); k++ {
				task := job.Tasks[taskKeys[k]]
				builder.WriteString(fmt.Sprintf("    - %s (%s)\n      Task Units:\n", taskKeys[k], task.Status))

				// Get a topologically sorted order of TaskUnits
				sortedUnits, err := types.TopologicalSort(task.TaskUnits)
				if err != nil {
					builder.WriteString(fmt.Sprintf("      Error: %s\n", err))
					continue
				}

				for _, unitID := range sortedUnits {
					unit := task.TaskUnits[unitID]
					builder.WriteString(fmt.Sprintf("      - %s (%s)\n", unitID, unit.Status))
					desc, _ := ms.GetTaskDefinition(unit.TaskDefinitionID) // Assuming this function exists and works
					owner, _ := ms.GetOwner(desc.OwnerID)                  // Assuming this function exists and works
					builder.WriteString(fmt.Sprintf("        - %s (%s)\n", desc.Name, owner.Name))
				}
			}
		}
	}

	fmt.Println(builder.String())
}

func (ms *MemoryStorage) PrintDAG() {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var builder strings.Builder

	builder.WriteString("Topics:\n")

	// Assuming you have the rest of the code to fetch and loop through topics and jobs.
	topics, err := ms.GetTopics() // Assuming this function exists and works.
	if err != nil {
		// If there's an error retrieving topics, log it and return.
		builder.WriteString(fmt.Sprintf("Error retrieving topics: %s\n", err))
		fmt.Println(builder.String())
		return
	}

	for _, topic := range topics {
		builder.WriteString(fmt.Sprintf("- %s\n", topic.Name))
		builder.WriteString("  Jobs:\n")

		// Loop through jobs.
		for jobID, job := range topic.Jobs {
			builder.WriteString(fmt.Sprintf("  - %s (%s)\n", jobID, job.Status))
			builder.WriteString("    Tasks:\n")

			// Loop through tasks.
			for taskID, task := range job.Tasks {
				builder.WriteString(fmt.Sprintf("    - %s (%s)\n      Task Units:\n", taskID, task.Status))

				// Get a topologically sorted order of TaskUnits.
				sortedUnits, err := types.TopologicalSort(task.TaskUnits)
				if err != nil {
					builder.WriteString(fmt.Sprintf("      Error: %s\n", err))
					continue
				}

				unitToDepth := make(map[types.TaskUnitID]int)
				for _, unitID := range sortedUnits {
					unit := task.TaskUnits[unitID]
					depth := 0

					// Find the maximum depth from dependencies.
					for _, depID := range unit.DependsOnIDs {
						if d, found := unitToDepth[depID]; found {
							if d+1 > depth {
								depth = d + 1
							}
						}
					}

					unitToDepth[unitID] = depth
					indent := strings.Repeat("    ", depth+1) // Additional indentation for visual hierarchy.

					// Fetch additional details like description and owner.
					desc, errDesc := ms.GetTaskDefinition(unit.TaskDefinitionID) // Assuming this function exists and works.
					owner, errOwner := ms.GetOwner(desc.OwnerID)                 // Assuming this function exists and works.

					if errDesc != nil || errOwner != nil {
						builder.WriteString(fmt.Sprintf("%s- %s (%s)\n", indent, unitID, unit.Status))
						// Handle errors or continue.
						if errDesc != nil {
							builder.WriteString(fmt.Sprintf("Error retrieving unit description: %s\n", errDesc))
						}
						if errOwner != nil {
							builder.WriteString(fmt.Sprintf("Error retrieving owner: %s\n", errOwner))
						}
						continue
					}

					builder.WriteString(fmt.Sprintf("%s- %s (%s): Desc(%s), Owner(%s)\n", indent, unitID, unit.Status, desc.Name, owner.Name))
					// ... [your existing code to print unit details].
				}
			}
		}
	}

	fmt.Println(builder.String())
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		topics:      make(map[types.TopicID]*types.Topic),
		jobs:        make(map[types.JobID]*types.Job),
		tasks:       make(map[types.TaskID]*types.Task),
		units:       make(map[types.TaskUnitID]*types.TaskUnit),
		definitions: make(map[types.TaskDefinitionID]*types.TaskDefinition),
		owners:      make(map[types.OwnerID]*types.Owner),
	}
}

func (ms *MemoryStorage) NewUUID() (string, error) {
	return types.GenerateUUID(), nil
}

// TODO (@droman): too lazy to implement the params but it's here
func (ms *MemoryStorage) GetInbox(ownerID types.OwnerID, params *types.QueryParams) ([]types.InboxTaskUnit, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	var err error
	var inboxUnits []types.InboxTaskUnit

	definitionIDs := []types.TaskDefinitionID{}
	definitions := map[types.TaskDefinitionID]types.TaskDefinition{}

	for _, def := range ms.definitions {
		if def.OwnerID == ownerID {
			definitionIDs = append(definitionIDs, def.Key)
		}
	}

	watchTasksForOwner := []*types.Task{}

	taskKeys := make([]types.TaskID, 0, len(ms.tasks))
	for ti := range ms.tasks {
		taskKeys = append(taskKeys, ti)
	}

	for i := 0; i < len(taskKeys); i++ {
		add := false
		for _, unit := range ms.tasks[taskKeys[i]].TaskUnits {
			for i := 0; i < len(definitionIDs); i++ {
				if unit.TaskDefinitionID == definitionIDs[i] {
					add = true
				}
			}
		}
		if add {
			for _, unit := range ms.tasks[taskKeys[i]].TaskUnits {
				definitions[unit.TaskDefinitionID] = *ms.definitions[unit.TaskDefinitionID]
			}
			watchTasksForOwner = append(watchTasksForOwner, ms.tasks[taskKeys[i]])
		}

	}

	defsKeys := make([]types.TaskDefinitionID, 0, len(definitions))
	for tdi := range definitions {
		defsKeys = append(defsKeys, tdi)
	}

	definitionsFlt := []types.TaskDefinition{}
	for i := 0; i < len(defsKeys); i++ {
		definitionsFlt = append(definitionsFlt, definitions[defsKeys[i]])
	}

	for idxTask := 0; idxTask < len(watchTasksForOwner); idxTask++ {

		if len(watchTasksForOwner[idxTask].JobID) == 0 {
			return nil, fmt.Errorf("critical error a task without JobID")
		}
		if _, ok := ms.jobs[watchTasksForOwner[idxTask].JobID]; !ok {
			continue
		}
		units := []types.TaskUnit{}

		// convert to value (stack)
		for i := 0; i < len(watchTasksForOwner[idxTask].TaskUnitIDs); i++ {
			units = append(units, *watchTasksForOwner[idxTask].TaskUnits[watchTasksForOwner[idxTask].TaskUnitIDs[i]])
		}

		var dag *types.WorkUnitDag

		if dag, err = types.CreateDagFromTaskUnits(ms, units, definitionsFlt); err != nil {
			return nil, err
		}

		workOwner := []types.TaskUnit{}
		workAvailable := dag.AvailableNodeUnitWithOwner(ownerID)

		// this make no sense, you don't have a jobid
		for _, v := range workAvailable {
			workOwner = append(workOwner, *v.Unit)
		}
		if len(workOwner) > 0 {
			inboxUnits = append(inboxUnits, types.InboxTaskUnit{
				TopicID:   ms.jobs[watchTasksForOwner[idxTask].JobID].TopicID,
				JobID:     watchTasksForOwner[idxTask].JobID,
				TaskID:    watchTasksForOwner[idxTask].Key,
				TaskUnits: workOwner,
			})
		}
	}

	return inboxUnits, nil
}

// Helper function to find the job ID associated with a task unit.
func (ms *MemoryStorage) findJobIDByTaskUnit(taskUnitID types.TaskUnitID) types.JobID {
	for _, job := range ms.jobs {
		for _, taskID := range job.TaskIDs {
			task, err := ms.GetTask(taskID)
			if err != nil {
				continue
			}
			if _, ok := task.TaskUnits[taskUnitID]; ok {
				return job.Key
			}
		}
	}
	return ""
}

// Helper function to find the task ID associated with a task unit.
func (ms *MemoryStorage) findTaskIDByTaskUnit(taskUnitID types.TaskUnitID) types.TaskID {
	for _, task := range ms.tasks {
		if _, ok := task.TaskUnits[taskUnitID]; ok {
			return task.Key
		}
	}
	return ""
}

// Assign a drafted `Job` to a `Topic` for processing
func (ms *MemoryStorage) AssignJob(topicID types.TopicID, jobID types.JobID) error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, exists := ms.topics[topicID]; !exists {
		return fmt.Errorf("topic %v doesn't exists", topicID)
	}
	if _, exists := ms.jobs[jobID]; !exists {
		return fmt.Errorf("job %v doesn't exists", jobID)
	}

	// in a real database, you will associate it without retrieving
	ms.topics[topicID].
		Mutate(
			types.WithTopicJobIDs(jobID),
			types.WithTopicJobs(ms.jobs[jobID]))

	return nil
}

// Assign a drafted `Task` to a `Job` for processing
func (ms *MemoryStorage) AssignTask(jobID types.JobID, taskID types.TaskID) error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, exists := ms.tasks[taskID]; !exists {
		return fmt.Errorf("task %v doesn't exists", taskID)
	}
	if _, exists := ms.jobs[jobID]; !exists {
		return fmt.Errorf("job %v doesn't exists", jobID)
	}

	// mutate the task to add jobID
	ms.tasks[taskID].
		Mutate(types.WithTaskJobID(jobID))

	// mutate the job to add the task
	ms.jobs[jobID].
		Mutate(
			types.WithJobTaskIDs(taskID),
			types.WithJobTasks(ms.tasks[taskID]))

	return nil
}

// Assign a drafted `TaskUnit` to a `Task` for processing
func (ms *MemoryStorage) AssignTaskUnits(taskID types.TaskID, ids []types.TaskUnitID) error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, exists := ms.tasks[taskID]; !exists {
		return fmt.Errorf("task %v doesn't exists", taskID)
	}
	for i := 0; i < len(ids); i++ {
		if _, exists := ms.units[ids[i]]; !exists {
			return fmt.Errorf("unit %v doesn't exists", ids[i])
		}
	}

	arr := []*types.TaskUnit{}
	for i := 0; i < len(ids); i++ {
		arr = append(arr, ms.units[ids[i]])
		ms.units[ids[i]].
			Mutate(types.WithTaskUnitTaskID(taskID))
	}
	ms.tasks[taskID].
		Mutate(
			types.WithTaskUnitsIDs(ids...),
			types.WithTaskUnits(arr...))

	return nil
}

func (ms *MemoryStorage) GetOwner(ownerID types.OwnerID) (*types.Owner, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	owner, exists := ms.owners[ownerID]
	if !exists {
		return nil, fmt.Errorf("owner with ID %s not found", ownerID)
	}

	return owner, nil
}

func (ms *MemoryStorage) GetTaskDefinition(id types.TaskDefinitionID) (*types.TaskDefinition, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	unitDesc, exists := ms.definitions[id]
	if !exists {
		return nil, fmt.Errorf("unit description with ID %s not found", id)
	}

	return unitDesc, nil
}

// CreateUnitOnTask associate units on existing task
func (ms *MemoryStorage) CreateUnitOnTask(taskID types.TaskID, units []*types.TaskUnit) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Retrieve the task
	_, exists := ms.tasks[taskID]
	if !exists {
		return fmt.Errorf("task with ID %s does not exist", taskID)
	}

	ids := []types.TaskUnitID{}
	// Create and associate task units
	for i := 0; i < len(units); i++ {

		// Set task and dependencies
		units[i].Mutate(types.WithTaskUnitTaskID(taskID))

		// cummulate keys
		ids = append(ids, units[i].Key)

		// Add the TaskUnit to the storage
		ms.units[units[i].Key] = units[i]
	}

	// Add the TaskUnit to the task
	ms.tasks[taskID].
		Mutate(
			types.WithTaskUnitsIDs(ids...),
			types.WithTaskUnits(units...),
		)

	// task.TaskUnitIDs = append(task.TaskUnitIDs, units[i].Key)
	// task.TaskUnits[units[i].Key] = &units[i]

	return nil
}

// CreateTasksOnJob creates new tasks and associates them with a specific job.
func (ms *MemoryStorage) CreateTasksOnJob(jobID types.JobID, tasks []types.Task) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Retrieve the job
	job, exists := ms.jobs[jobID]
	if !exists {
		return fmt.Errorf("job with ID %s does not exist", jobID)
	}

	// Create and associate tasks
	for _, task := range tasks {
		// Generate a new Task ID
		taskID := types.TaskID(types.GenerateUUID())

		// Set job and description
		task.JobID = jobID
		task.Key = taskID

		// Initialize TaskUnits map
		task.TaskUnits = make(map[types.TaskUnitID]*types.TaskUnit)

		// Add the Task to the job
		job.TaskIDs = append(job.TaskIDs, taskID)
		job.Tasks[taskID] = &task

		// Add the Task to the storage
		ms.tasks[taskID] = &task
	}

	return nil
}

// HasTopic checks if a topic with the given ID exists.
func (s *MemoryStorage) HasTopic(id types.TopicID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.topics[id]
	return exists, nil
}

// HasJob checks if a job with the given ID exists.
func (s *MemoryStorage) HasJob(id types.JobID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.jobs[id]
	return exists, nil
}

// HasTask checks if a task with the given ID exists.
func (s *MemoryStorage) HasTask(id types.TaskID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.tasks[id]
	return exists, nil
}

// GetOwners retrieves all owners.
func (ms *MemoryStorage) GetOwners() ([]types.Owner, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	owners := make([]types.Owner, 0, len(ms.owners))
	for _, owner := range ms.owners {
		owners = append(owners, *owner)
	}

	return owners, nil
}

// GetTaskDefinitions retrieves all unit descriptions.
func (ms *MemoryStorage) GetTaskDefinitions() ([]types.TaskDefinition, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	unitDescriptions := make([]types.TaskDefinition, 0, len(ms.definitions))
	for _, unitDescription := range ms.definitions {
		unitDescriptions = append(unitDescriptions, *unitDescription)
	}

	return unitDescriptions, nil
}

// HasTaskUnit checks if a task unit with the given ID exists.
func (s *MemoryStorage) HasTaskUnit(id types.TaskUnitID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.units[id]
	return exists, nil
}

// HasTaskDefinition checks if a unit description with the given ID exists.
func (s *MemoryStorage) HasTaskDefinition(id types.TaskDefinitionID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.definitions[id]
	return exists, nil
}

// HasOwner checks if an owner with the given ID exists.
func (s *MemoryStorage) HasOwner(id types.OwnerID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.owners[id]
	return exists, nil
}

// CreateOwner creates a new owner
func (ms *MemoryStorage) CreateOwner(name string, cfgs ...types.OwnerConfig) (*types.Owner, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	owner := types.NewOwner(types.OwnerID(uuid), name, cfgs...)

	// we don't want the same id
	if _, exists := ms.owners[owner.Key]; exists {
		return nil, types.ErrOwnerIDAlreadyExists
	}
	// we don't want the same name
	for _, v := range ms.owners {
		if v.Name == name {
			return nil, types.ErrOwnerNameAlreadyExists
		}
	}

	ms.owners[owner.Key] = owner
	return owner, nil
}

func (ms *MemoryStorage) CreateTaskUnits(units []*types.TaskUnit) ([]types.TaskUnitID, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	ids := []types.TaskUnitID{}
	for i := 0; i < len(units); i++ {
		if _, exists := ms.units[units[i].Key]; exists {
			return nil, types.ErrOwnerIDAlreadyExists
		}
		ms.units[units[i].Key] = units[i]
		ids = append(ids, units[i].Key)
	}

	return ids, nil
}

// UpdateOwner updates an owner by ID
func (ms *MemoryStorage) UpdateOwner(ownerID types.OwnerID, name string) (*types.Owner, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	owner, exists := ms.owners[ownerID]
	if !exists {
		return nil, fmt.Errorf("not found")
	}

	owner.Name = name
	return owner, nil
}

// DeprecateOwner deprecates an owner by ID
func (ms *MemoryStorage) DeprecateOwner(ownerID types.OwnerID) (*types.Owner, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	owner, exists := ms.owners[ownerID]
	if !exists {
		return nil, fmt.Errorf("not found")
	}

	delete(ms.owners, ownerID)
	return owner, nil
}

// CreateTaskDefinition creates a new unit description.
func (ms *MemoryStorage) CreateTaskDefinition(name string, ownerID types.OwnerID, cfgs ...types.TaskDefinitionConfig) (*types.TaskDefinition, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	unitDescription := types.NewUnitDescription(types.TaskDefinitionID(uuid), name, ownerID, cfgs...)

	// we don't want the same id
	if _, exists := ms.definitions[unitDescription.Key]; exists {
		return nil, types.ErrOwnerIDAlreadyExists
	}
	// we don't want the same name
	for _, v := range ms.definitions {
		if v.Name == name {
			return nil, types.ErrOwnerNameAlreadyExists
		}
	}

	ms.definitions[unitDescription.Key] = unitDescription
	return unitDescription, nil
}

// UpdateTaskDefinition updates an existing unit description by ID.
func (s *MemoryStorage) UpdateTaskDefinition(id types.TaskDefinitionID, ownerID types.OwnerID, name string, description string, identifier string) (*types.TaskDefinition, error) {
	unitDesc, exists := s.definitions[id]
	if !exists {
		return nil, errors.New("unit description not found")
	}

	// Check if the owner is the same.
	if unitDesc.OwnerID != ownerID {
		return nil, errors.New("invalid owner for unit description")
	}

	unitDesc.Name = name
	unitDesc.Details = description
	unitDesc.Identifier = identifier

	s.definitions[id] = unitDesc
	return unitDesc, nil
}

// DeprecateTaskDefinition deprecates a unit description by ID.
func (s *MemoryStorage) DeprecateTaskDefinition(id types.TaskDefinitionID) error {
	_, exists := s.definitions[id]
	if !exists {
		return errors.New("unit description not found")
	}

	// Check if the unit description is associated with any task units.
	for _, taskUnit := range s.units {
		if taskUnit.TaskDefinitionID == id {
			return errors.New("cannot deprecate unit description associated with task unit(s)")
		}
	}

	delete(s.definitions, id)
	return nil
}

// CancelJob cancels a job by ID.
func (s *MemoryStorage) CancelJob(jobID types.JobID) error {
	job, exists := s.jobs[jobID]
	if !exists {
		return errors.New("job not found")
	}

	// Set the status of the job to "canceled."
	job.Status = types.ErrorStatus
	s.jobs[jobID] = job

	// Update the status of all task units in the job.
	for _, taskID := range job.TaskIDs {
		task, exists := s.tasks[taskID]
		if !exists {
			continue
		}

		for _, taskUnitID := range task.TaskUnitIDs {
			taskUnit, exists := s.units[taskUnitID]
			if !exists {
				continue
			}

			if taskUnit.Status == types.QueuedStatus || taskUnit.Status == types.ProgressStatus {
				taskUnit.Status = types.ErrorStatus
				taskUnit.Error = errors.New("job canceled")
				s.units[taskUnitID] = taskUnit
			}
		}
	}

	return nil
}

// CancelTask cancels a task by ID.
func (s *MemoryStorage) CancelTask(taskID types.TaskID) error {
	task, exists := s.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	// Set the status of the task to "canceled."
	task.Status = types.ErrorStatus
	s.tasks[taskID] = task

	// Update the status of all task units in the task.
	for _, taskUnitID := range task.TaskUnitIDs {
		taskUnit, exists := s.units[taskUnitID]
		if !exists {
			continue
		}

		if taskUnit.Status == types.QueuedStatus || taskUnit.Status == types.ProgressStatus {
			taskUnit.Status = types.ErrorStatus
			taskUnit.Error = errors.New("task canceled")
			s.units[taskUnitID] = taskUnit
		}
	}

	return nil
}

// Create new `Topic`
func (ms *MemoryStorage) CreateTopic(name string, cfgs ...types.TopicConfig) (*types.Topic, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	topic := types.NewTopic(types.TopicID(uuid), name, cfgs...)

	// we don't want the same id
	if _, exists := ms.topics[topic.Key]; exists {
		return nil, types.ErrTopicIDAlreadyExists
	}
	// we don't want the same name
	for _, v := range ms.topics {
		if v.Name == name {
			return nil, types.ErrTopicNameAlreadyExists
		}
	}

	ms.topics[topic.Key] = topic
	return topic, nil
}

func (ms *MemoryStorage) CreateJob(cfgs ...types.JobConfig) (*types.Job, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	job := types.NewJob(types.JobID(uuid), cfgs...)

	// we don't want the same id
	if _, exists := ms.jobs[job.Key]; exists {
		return nil, types.ErrTopicIDAlreadyExists
	}

	ms.jobs[job.Key] = job

	return job, nil
}

// CreateTask creates a new task.
func (ms *MemoryStorage) CreateTask(cfgs ...types.TaskConfig) (*types.Task, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	task := types.NewTask(types.TaskID(uuid), cfgs...)

	// we don't want the same id
	if _, exists := ms.tasks[task.Key]; exists {
		return nil, types.ErrOwnerIDAlreadyExists
	}

	ms.tasks[task.Key] = task
	return task, nil
}

// CreateTaskUnit creates a new TaskUnit with the given dependencies.
func (ms *MemoryStorage) CreateTaskUnit(cfgs ...types.TaskUnitConfig) (*types.TaskUnit, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var uuid string
	var err error
	if uuid, err = ms.NewUUID(); err != nil {
		return nil, err
	}

	// Create the TaskUnit
	taskUnit := types.NewTaskUnit(types.TaskUnitID(uuid), cfgs...)

	// we don't want the same id
	if _, exists := ms.units[taskUnit.Key]; exists {
		return nil, types.ErrOwnerIDAlreadyExists
	}

	// Add the TaskUnit to the storage
	ms.units[taskUnit.Key] = taskUnit

	return taskUnit, nil
}

// GetTasks returns all tasks associated with a job.
func (m *MemoryStorage) GetTasks(jobID types.JobID) ([]types.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[jobID]
	if !ok {
		return nil, fmt.Errorf("job with ID %s not found", jobID)
	}

	tasks := make([]types.Task, 0, len(job.TaskIDs))
	for _, taskID := range job.TaskIDs {
		task, ok := m.tasks[taskID]
		if !ok {
			return nil, fmt.Errorf("task with ID %s not found", taskID)
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (ms *MemoryStorage) GetTopics() ([]types.Topic, error) {
	// Implement the logic to retrieve all topics from memory.
	topics := make([]types.Topic, 0)
	for _, topic := range ms.topics {
		for jid, value := range ms.jobs {
			if topic.Key == value.TopicID {
				job, _ := ms.GetJob(jid)
				topic.Jobs[jid] = job
			}
		}
		topics = append(topics, *topic)
	}
	return topics, nil
}

// // CreateTaskUnits creates new task units for a task
// func (m *MemoryStorage) CreateTaskUnits(taskID types.TaskID, units []types.TaskUnit) ([]types.TaskUnit, error) {
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	task, ok := m.tasks[taskID]
// 	if !ok {
// 		return nil, fmt.Errorf("task with ID %s not found", taskID)
// 	}

// 	var createdUnits []types.TaskUnit

// 	for _, unit := range units {
// 		if _, ok := m.units[unit.Key]; ok {
// 			return nil, fmt.Errorf("task unit with ID %s already exists", unit.Key)
// 		}

// 		unit.TaskID = taskID
// 		m.units[unit.Key] = &unit
// 		task.TaskUnitIDs = append(task.TaskUnitIDs, unit.Key)
// 		task.TaskUnits[unit.Key] = &unit

// 		createdUnits = append(createdUnits, unit)
// 	}

// 	return createdUnits, nil
// }

func (ms *MemoryStorage) GetTopic(id types.TopicID) (*types.Topic, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	topic, exists := ms.topics[id]
	if !exists {
		return nil, errors.New("topic not found")
	}

	return topic, nil
}

// UpdateTopic updates a topic by ID
func (ms *MemoryStorage) UpdateTopic(id types.TopicID, name string) (*types.Topic, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	topic, exists := ms.topics[id]
	if !exists {
		return nil, errors.New("topic not found")
	}

	topic.Name = name

	return topic, nil
}

// DeleteTopic deletes a topic by ID
func (ms *MemoryStorage) DeprecateTopic(id types.TopicID) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, exists := ms.topics[id]
	if !exists {
		return errors.New("topic not found")
	}

	topic, err := ms.GetTopic(id)
	if err != nil {
		return err
	}

	topic.Deprecated = true

	return nil
}

func (ms *MemoryStorage) GetJob(jobID types.JobID) (*types.Job, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	job, exists := ms.jobs[jobID]
	if !exists {
		return nil, errors.New("job not found")
	}

	for _, v := range ms.tasks {
		if v.JobID == jobID {
			task, _ := ms.GetTask(v.Key)
			job.Tasks[v.Key] = task
		}
	}

	return job, nil
}

func (ms *MemoryStorage) GetTask(taskID types.TaskID) (*types.Task, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	task, exists := ms.tasks[taskID]
	if !exists {
		return nil, errors.New("task not found")
	}

	for _, v := range ms.units {
		if v.TaskID == taskID {
			unit, _ := ms.GetTaskUnit(v.Key)
			task.TaskUnits[v.Key] = unit
		}
	}

	return task, nil
}

func (ms *MemoryStorage) GetTaskUnit(taskUnitID types.TaskUnitID) (*types.TaskUnit, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	unit, exists := ms.units[taskUnitID]
	if !exists {
		return nil, errors.New("task unit not found")
	}

	return unit, nil
}

func (ms *MemoryStorage) UpdateJobStatus(jobID types.JobID, status types.StatusType) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	job, exists := ms.jobs[jobID]
	if !exists {
		return errors.New("job not found")
	}

	job.Status = status
	return nil
}

func (ms *MemoryStorage) UpdateTaskStatus(taskID types.TaskID, status types.StatusType) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	task, exists := ms.tasks[taskID]
	if !exists {
		return errors.New("task not found")
	}

	task.Status = status
	return nil
}

func (ms *MemoryStorage) UpdateTaskUnitStatus(taskUnitID types.TaskUnitID, status types.StatusType, err error) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	unit, exists := ms.units[taskUnitID]
	if !exists {
		return errors.New("task unit not found")
	}

	unit.Status = status
	unit.Error = err

	fmt.Println(unit.Key)

	return nil
}

// GetJobs retrieves all jobs for a topic.
func (ms *MemoryStorage) GetJobs(topicID types.TopicID) ([]types.Job, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	topic, exists := ms.topics[topicID]
	if !exists {
		return nil, fmt.Errorf("topic not found")
	}

	jobs := make([]types.Job, 0, len(topic.Jobs))
	for _, job := range topic.Jobs {
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// GetTaskUnits retrieves all task units for a task.
func (ms *MemoryStorage) GetTaskUnits(taskID types.TaskID) ([]types.TaskUnit, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	task, exists := ms.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	taskUnits := make([]types.TaskUnit, 0, len(task.TaskUnits))
	for _, taskUnit := range task.TaskUnits {
		taskUnits = append(taskUnits, *taskUnit)
	}

	return taskUnits, nil
}
