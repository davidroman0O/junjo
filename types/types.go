package types

import (
	"crypto/rand"
	"errors"
	"fmt"
)

// TODO: create better errors for entities
var (
	ErrTopicIDAlreadyExists   = errors.New("topic with same id already exists")
	ErrTopicNameAlreadyExists = errors.New("topic with same name already exists")

	ErrOwnerIDAlreadyExists   = errors.New("owner with same id already exists")
	ErrOwnerNameAlreadyExists = errors.New("owner with same name already exists")
)

// StorageInterface defines the methods required for managing data persistence.
// - creation: you never have to give an ID of the entitiy, it has to be managed by the type and your implementation
// - every unassigned Job, Task, TaskUnit are drafts, consider deleting them after a while
type StorageInterface interface {

	// Utility function to create a new `uuid` for each entities, bring your own
	NewUUID() (string, error)

	// Create an `Owner` that will handle `TaskUnit` based on `TaskDefinition`
	CreateOwner(name string, cfgs ...OwnerConfig) (*Owner, error) // [x]

	// Create a `TaskDefinition` which is the description of an abstract task defined for an `Owner`
	CreateTaskDefinition(name string, ownerID OwnerID, cfgs ...TaskDefinitionConfig) (*TaskDefinition, error) // [x]

	// Create new `Topic` which manage `Job` instances
	CreateTopic(name string, cfgs ...TopicConfig) (*Topic, error) // [x]

	// Create new `Task` without no JobID and no units
	CreateTask(cfgs ...TaskConfig) (*Task, error) // [x]

	// Create new `Job` without TaskID and no units
	CreateJob(cfgs ...JobConfig) (*Job, error) // [x]

	CreateTaskUnits(units []*TaskUnit) ([]TaskUnitID, error) // [ ]

	// Assign a drafted `Job` to a `Topic` for processing
	AssignJob(topicID TopicID, jobID JobID) error // [x]

	// Assign a drafted `Task` to a `Job`
	AssignTask(jobID JobID, taskID TaskID) error // [x]

	// Assign a drafted array `TaskUnit` to a `Task`
	AssignTaskUnits(taskID TaskID, ids []TaskUnitID) error // [x]

	// Owners should only see the `TaskUnit` they are supposed to accomplish
	GetInbox(ownerID OwnerID, params *QueryParams) ([]InboxTaskUnit, error) // [x]

	GetJob(jobID JobID) (*Job, error)
	UpdateTaskDefinition(id TaskDefinitionID, ownerID OwnerID, name string, description string, identifier string) (*TaskDefinition, error)
	DeprecateTaskDefinition(id TaskDefinitionID) error
	HasTaskDefinition(id TaskDefinitionID) (bool, error)
	GetTaskDefinition(id TaskDefinitionID) (*TaskDefinition, error)
	GetTaskDefinitions() ([]TaskDefinition, error)

	GetOwners() ([]Owner, error)
	GetOwner(ownerID OwnerID) (*Owner, error)
	UpdateOwner(ownerID OwnerID, name string) (*Owner, error)
	DeprecateOwner(ownerID OwnerID) (*Owner, error)
	HasOwner(id OwnerID) (bool, error)

	// A owner only own TaskUnit related to TaskDescription
	// GetInboxTopic(topic TopicID, owner OwnerID) ([]InboxTaskUnit, error)
	// GetInboxJob(jobID JobID, owner OwnerID) ([]InboxTaskUnit, error)
	// GetInboxTask(taskID TaskID, owner OwnerID) ([]InboxTaskUnit, error)

	GetTopics() ([]Topic, error)
	GetTopic(id TopicID) (*Topic, error)

	UpdateTopic(id TopicID, name string) (*Topic, error)
	HasTopic(id TopicID) (bool, error)
	DeprecateTopic(id TopicID) error

	GetJobs(id TopicID) ([]Job, error)
	HasJob(id JobID) (bool, error)
	CancelJob(jobID JobID) error

	GetTasks(jobID JobID) ([]Task, error)
	GetTask(taskID TaskID) (*Task, error)

	HasTask(id TaskID) (bool, error)
	CancelTask(taskID TaskID) error

	GetTaskUnits(taskID TaskID) ([]TaskUnit, error)
	GetTaskUnit(taskUnitID TaskUnitID) (*TaskUnit, error)

	HasTaskUnit(id TaskUnitID) (bool, error)

	UpdateJobStatus(jobID JobID, status StatusType) error
	UpdateTaskStatus(taskID TaskID, status StatusType) error
	UpdateTaskUnitStatus(taskUnitID TaskUnitID, status StatusType, error error) error
}

type TopicID string
type JobID string
type TaskID string
type TaskDefinitionID string
type TaskUnitID string
type OwnerID string

type CommandType string

var (
	ProgressCmd CommandType = "progress"
	SuccessCmd  CommandType = "success"
	ErrorCmd    CommandType = "error"
	PauseCmd    CommandType = "pause"
	LogCmd      CommandType = "log"
)

type StatusType string

var (
	NoneStatus     StatusType = "none"
	QueuedStatus   StatusType = "queued"
	ProgressStatus StatusType = "in-progress"
	SuccessStatus  StatusType = "success"
	ErrorStatus    StatusType = "error"
	PauseStatus    StatusType = "pause"
)

type OwnerConfig func(data *Owner)

func WithOwnerDescription(d string) OwnerConfig {
	return func(data *Owner) {
		data.Description = d
	}
}

// NewOwner creates a new Owner with a generated UUID as the ID.
func NewOwner(id OwnerID, name string, cfgs ...OwnerConfig) *Owner {
	owner := &Owner{
		Key:  id,
		Name: name,
	}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](owner)
	}
	return owner
}

type Owner struct {
	Key         OwnerID `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
}

// Command represents an action from a client to report progress or mutate task unit states.
type Command struct {
	Type    CommandType       `json:"type" db:"type"`
	Status  StatusType        `json:"status" db:"status"`
	Details string            `json:"details" db:"details"`
	Data    map[string]string `json:"data" db:"data"`
}

// TaskDefinition is the template of a TaskUnit (instance)
type TaskDefinition struct {
	Key         TaskDefinitionID `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Details     string           `json:"details" db:"details"`
	Identifier  string           `json:"identifier" db:"identifier"` // should be unique for the organization
	OwnerID     OwnerID          `json:"ownerID" db:"ownerID"`       // one owner can complete that unit
}

type TaskDefinitionConfig func(data *TaskDefinition)

func WithTaskDefDescription(d string) TaskDefinitionConfig {
	return func(data *TaskDefinition) {
		data.Description = d
	}
}

func WithTaskDefIdentifier(identifier string) TaskDefinitionConfig {
	return func(data *TaskDefinition) {
		data.Identifier = identifier
	}
}

// NewUnitDescription creates a new UnitDescription with a generated UUID as the ID.
func NewUnitDescription(id TaskDefinitionID, name string, ownerID OwnerID, cfgs ...TaskDefinitionConfig) *TaskDefinition {
	unitDescription := &TaskDefinition{
		Key:     id,
		Name:    name,
		OwnerID: ownerID,
	}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](unitDescription)
	}
	return unitDescription
}

type NodeTaskUnitConfig func(data *NodeTaskUnit)

func WithNodeWithTaskUnit(node TaskUnit) NodeTaskUnitConfig {
	return func(data *NodeTaskUnit) {
		data.Unit = &node
	}
}

func WithNodeWithTaskStatus(status StatusType) NodeTaskUnitConfig {
	return func(data *NodeTaskUnit) {
		if data.Unit == nil {
			data.Unit = &TaskUnit{}
		}
		data.Unit.Status = status
	}
}

func WithNodeWithTaskData(value map[string]string) NodeTaskUnitConfig {
	return func(data *NodeTaskUnit) {
		if data.Unit == nil {
			data.Unit = &TaskUnit{}
		}
		data.Unit.Data = value
	}
}

func WithNodeWithTaskKey(key TaskUnitID) NodeTaskUnitConfig {
	return func(data *NodeTaskUnit) {
		if data.Unit == nil {
			data.Unit = &TaskUnit{}
		}
		data.Unit.Key = key
	}
}

func WithNodeWithTaskDefinition(def *TaskDefinition) NodeTaskUnitConfig {
	return func(data *NodeTaskUnit) {
		data.Definition = def
	}
}

// Representation of a `TaskUnit` within a DAG
type NodeTaskUnit struct {
	Unit       *TaskUnit
	Definition *TaskDefinition
}

// NewNodeTaskUnit creates a new NodeTaskUnit for DAG
func NewNodeTaskUnit(cfgs ...NodeTaskUnitConfig) *NodeTaskUnit {
	nodeTaskUnit := &NodeTaskUnit{}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](nodeTaskUnit)
	}
	return nodeTaskUnit
}

func (u *NodeTaskUnit) ID() string {
	return string(u.Unit.Key)
}

func (u *NodeTaskUnit) Metadata() map[string]string {
	return u.Unit.Data
}

func (u *NodeTaskUnit) SetMetadata(metadata map[string]string) {
	if u.Unit.Data == nil {
		u.Unit.Data = map[string]string{}
	}
	for k, v := range metadata {
		u.Unit.Data[k] = v
	}
}

type TaskUnitConfig func(data *TaskUnit)

func WithTaskUnitDefinition(def *TaskDefinition) TaskUnitConfig {
	return func(data *TaskUnit) {
		data.TaskDefinitionID = def.Key
	}
}

func WithTaskUnitDefinitionKey(def TaskDefinitionID) TaskUnitConfig {
	return func(data *TaskUnit) {
		data.TaskDefinitionID = def
	}
}

func WithTaskUnitTaskID(id TaskID) TaskUnitConfig {
	return func(data *TaskUnit) {
		data.TaskID = id
	}
}

func WithTaskUnitStatus(status StatusType) TaskUnitConfig {
	return func(data *TaskUnit) {
		data.Status = status
	}
}

func WithTaskUnitData(value map[string]string) TaskUnitConfig {
	return func(data *TaskUnit) {
		data.Data = value
	}
}

func WithTaskUnitDepends(units []TaskUnit) TaskUnitConfig {
	return func(data *TaskUnit) {
		for i := 0; i < len(units); i++ {

			data.DependsOnIDs = append(data.DependsOnIDs, units[i].Key)
			data.DependsOn = append(data.DependsOn, &units[i])
		}

	}
}

// TaskUnit represents a granular work unit within a task's DAG.
// Using the taskDefinitionID, we know who is owner if it
type TaskUnit struct {
	Key              TaskUnitID        `json:"id" db:"id"`
	TaskDefinitionID TaskDefinitionID  `json:"taskDefinitionID" db:"taskDefinitionID"`
	DependsOnIDs     []TaskUnitID      `json:"dependsOnIds" db:"dependsOnIds"`
	DependsOn        []*TaskUnit       `json:"dependsOn,omitempty" db:"-"` // runtime, it will be useful to know what tasks was done before me
	Commands         []Command         `json:"commands" db:"commands"`     // tracking which commands was done on this taskunit
	Status           StatusType        `json:"status" db:"status"`
	Error            error             `json:"error" db:"error"`
	TaskID           TaskID            `json:"taskID" db:"taskID"`
	Data             map[string]string `json:"data" db:"data"` // original data, you have to run the commands to get the mutations of the data
}

func (j *TaskUnit) Mutate(cfgs ...TaskUnitConfig) {
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](j)
	}
}

// NewTaskUnit creates a new TaskUnit
func NewTaskUnit(id TaskUnitID, cfgs ...TaskUnitConfig) *TaskUnit {
	unit := &TaskUnit{
		Key:          id,
		Status:       NoneStatus,
		DependsOnIDs: []TaskUnitID{},
		DependsOn:    []*TaskUnit{},
		Commands:     make([]Command, 0),
	}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](unit)
	}
	return unit
}

type InboxTaskUnit struct {
	TopicID   TopicID    `json:"topicID" db:"topicID"`
	JobID     JobID      `json:"jobID" db:"jobID"`
	TaskID    TaskID     `json:"taskID" db:"taskID"`
	TaskUnits []TaskUnit `json:"taskUnits" db:"taskUnits"`
}

type TaskConfig func(data *Task)

func WithTaskJobID(jobID JobID) TaskConfig {
	return func(data *Task) {
		data.JobID = jobID
	}
}

func WithTaskUnitsIDs(id ...TaskUnitID) TaskConfig {
	return func(data *Task) {
		data.TaskUnitIDs = append(data.TaskUnitIDs, id...)
	}
}

func WithTaskUnits(units ...*TaskUnit) TaskConfig {
	return func(data *Task) {
		for i := 0; i < len(units); i++ {
			data.TaskUnits[units[i].Key] = units[i]
		}
	}
}

// Task represents a DAG of task units.
type Task struct {
	Key         TaskID                   `json:"id" db:"id"`
	JobID       JobID                    `json:"jobID" db:"jobID"`
	Status      StatusType               `json:"status" db:"status"`
	TaskUnitIDs []TaskUnitID             `json:"taskUnitIds" db:"taskUnitIds"` // instances of the nodes of the dag, those instances represent the dag
	TaskUnits   map[TaskUnitID]*TaskUnit `json:"taskUnits,omitempty" db:"-"`   // runtime
}

func (j *Task) Mutate(cfgs ...TaskConfig) {
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](j)
	}
}

// NewTask creates a new Task with a generated UUID as the ID.
func NewTask(id TaskID, cfgs ...TaskConfig) *Task {
	task := &Task{
		Key:         id,
		TaskUnitIDs: []TaskUnitID{},
		TaskUnits:   make(map[TaskUnitID]*TaskUnit),
		Status:      NoneStatus,
	}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](task)
	}
	// eventually assign the taskID to all units
	for i := 0; i < len(task.TaskUnitIDs); i++ {
		task.TaskUnits[task.TaskUnitIDs[i]].TaskID = task.Key
	}
	return task
}

type JobConfig func(data *Job)

func WithJobTaskID(topicID TopicID) JobConfig {
	return func(data *Job) {
		data.TopicID = topicID
	}
}

func WithJobTaskIDs(taskIDs ...TaskID) JobConfig {
	return func(data *Job) {
		data.TaskIDs = append(data.TaskIDs, taskIDs...)
	}
}

func WithJobTasks(tasks ...*Task) JobConfig {
	return func(data *Job) {
		for i := 0; i < len(tasks); i++ {
			data.Tasks[tasks[i].Key] = tasks[i]
		}
	}
}

func WithJobData(value map[string]string) JobConfig {
	return func(data *Job) {
		data.Data = value
	}
}

// Actual work that need to be done in that topic
type Job struct {
	Key     JobID             `json:"id" db:"id"`
	TaskIDs []TaskID          `json:"taskIds" db:"taskIds"`
	Tasks   map[TaskID]*Task  `json:"tasks,omitempty" db:"-"`
	Status  StatusType        `json:"status" db:"status"`
	Data    map[string]string `json:"data" db:"data"` // initial data to work with
	TopicID TopicID           `json:"topicID" db:"topicID"`
}

func (j *Job) Mutate(cfgs ...JobConfig) {
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](j)
	}
}

// NewJob creates a new Job with a generated UUID as the ID.
func NewJob(id JobID, cfgs ...JobConfig) *Job {
	job := &Job{
		Key:     JobID(GenerateUUID()),
		Tasks:   make(map[TaskID]*Task),
		TaskIDs: []TaskID{},
		Status:  NoneStatus, // Initialize with QueuedStatus
	}

	for i := 0; i < len(cfgs); i++ {
		cfgs[i](job)
	}

	return job
}

type TopicConfig func(data *Topic)

func WithTopicDescription(d string) TopicConfig {
	return func(data *Topic) {
		data.Description = d
	}
}

func WithTopicJobIDs(jobID ...JobID) TopicConfig {
	return func(data *Topic) {
		data.JobIDs = append(data.JobIDs, jobID...)
	}
}

func WithTopicJobs(job ...*Job) TopicConfig {
	return func(data *Topic) {
		for i := 0; i < len(job); i++ {
			data.Jobs[job[i].Key] = job[i]
		}
	}
}

// Topic represents a high-level category of related work.
type Topic struct {
	Key            TopicID        `json:"id" db:"id"`
	Name           string         `json:"name" db:"name"`
	Description    string         `json:"description" db:"description"`
	JobIDs         []JobID        `json:"jobIds" db:"jobIds"`
	Jobs           map[JobID]*Job `json:"jobs,omitempty" db:"-"` // runtime only
	TotalCompleted *int           `json:"totalCompleted" db:"-"` // runtime only
	TotalPending   *int           `json:"totalPending" db:"-"`   // runtime only
	TotalError     *int           `json:"totalError" db:"-"`     // runtime only
	TotalPause     *int           `json:"totalPause" db:"-"`     // runtime only
	Deprecated     bool           `json:"deprecated" db:"deprecated"`
}

func (j *Topic) Mutate(cfgs ...TopicConfig) {
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](j)
	}
}

// GenerateUUID generates a random UUID (version 4) and returns it as a string.
// The `StorageInterface` require a `NewUUID` function so you can change it yourself
func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}

	// Set version (4) and variant bits (2) as per UUID specification
	uuid[6] = (uuid[6] & 0x0F) | 0x40
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	// Format the UUID as a string
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// NewTopic creates a new Topic with a generated UUID as the ID.
func NewTopic(id TopicID, name string, cfgs ...TopicConfig) *Topic {
	topic := &Topic{
		Key:    id,
		Name:   name,
		JobIDs: []JobID{},
		Jobs:   make(map[JobID]*Job),
	}
	for i := 0; i < len(cfgs); i++ {
		cfgs[i](topic)
	}
	return topic
}

func TopologicalSort(taskUnits map[TaskUnitID]*TaskUnit) ([]TaskUnitID, error) {
	var order []TaskUnitID
	tempMark := make(map[TaskUnitID]bool)
	permMark := make(map[TaskUnitID]bool)
	var visit func(u TaskUnitID) error

	visit = func(u TaskUnitID) error {
		if _, found := tempMark[u]; found {
			return fmt.Errorf("cyclic dependency detected")
		}

		if _, found := permMark[u]; !found {
			tempMark[u] = true

			for _, m := range taskUnits[u].DependsOnIDs {
				err := visit(m)
				if err != nil {
					return err
				}
			}

			delete(tempMark, u)
			permMark[u] = true
			order = append(order, u)
		}

		return nil
	}

	for id := range taskUnits {
		if _, found := permMark[id]; !found {
			err := visit(id)
			if err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}

type QueryConfig func(p *QueryParams)

func WithQueryOffset(offset int) QueryConfig {
	return func(p *QueryParams) {
		p.Offset = &offset
	}
}

func WithQuerySize(size int) QueryConfig {
	return func(p *QueryParams) {
		p.Size = &size
	}
}

// Simple Query
type QueryParams struct {
	Offset *int `json:"offset"`
	Size   *int `json:"size"`
}

func NewQuery(cfgs ...QueryConfig) *QueryParams {
	p := &QueryParams{}

	for i := 0; i < len(cfgs); i++ {
		cfgs[i](p)
	}

	return p
}
