package junjo

import (
	"github.com/davidroman0O/junjo/types"
)

/// I DON'T LIKE IT ANYMOOOOOOOOOOOOOOOOOOOOOOORRRREEEEEEEEEEEEEEEEEEEEEE
/// aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah

///
/// Junjo is a hierarchical workflow management system with DAG-Structured tasks intented to orchestrate distributed systems around a set of controlled task queues.
/// Each `Owner` of `TaskDefinition` will only see the tasks they need to process and report on, while the overall progression of the workflow is managed by `junjo`.
/// Sometimes you don't want your distributed task system to do the work, you want to outsource the work and let other systems do what they know best.
///
/// The main motto of `junjo` is "mind your own business"
///
/// You have the possibility to create new `StorageInterface` for your own usage, follow the guidelines to have all the rules (TODO) you need to implement.
///

// `Junjoold` is your main interface with it's systems which provide an API to facilitate your interaction with it
// Most of the API leverage the storage implemenation, except few exeptions.
type Junjoold struct {
	storageImplementation types.StorageInterface
}

// New `Junjo` Api
func NewJ(implt types.StorageInterface) *Junjoold {
	return &Junjoold{
		storageImplementation: implt,
	}
}

// Create a new `Topic`
func (j *Junjoold) CreateTopic(name string, cfgs ...types.TopicConfig) (*types.Topic, error) {
	return j.
		storageImplementation.
		CreateTopic(name, cfgs...)
}

// Create a new `Owner`
func (j *Junjoold) CreateOwner(name string, cfgs ...types.OwnerConfig) (*types.Owner, error) {
	return j.
		storageImplementation.
		CreateOwner(name, cfgs...)
}

// Create a new `TaskDefinition` for a `Owner`
func (j *Junjoold) CreateTaskDefinition(name string, ownerID types.OwnerID, cfgs ...types.TaskDefinitionConfig) (*types.TaskDefinition, error) {
	return j.
		storageImplementation.
		CreateTaskDefinition(name, ownerID, cfgs...)
}

// Create a new `WorkUnitDag` to define your own DAG based on your `TaskDefinition`
func (j *Junjoold) CreateDagTaskUnits() *types.WorkUnitDag {
	// execption, we don't need storage implementation here, just `types`
	return types.NewWorkUnitDag(j.storageImplementation)
}

// Create new `Task` with it's TaskUnit, probably from a WorkUnitDag
// By default that task as `Status == none` with no JobID
func (j *Junjoold) CreateTask(cfgs ...types.TaskConfig) (*types.Task, error) {
	return j.
		storageImplementation.
		CreateTask(cfgs...)
}

// Create new `Job` with it's Tasks
// By default that `Job` will have `Status == none` with no TopicID AND no initial `Data`
func (j *Junjoold) CreateJob(cfgs ...types.JobConfig) (*types.Job, error) {
	return j.
		storageImplementation.
		CreateJob(cfgs...)
}

// Create new `TaskUnit` as an array
// By default those `TaskUnit` has `Status == none` with not JobID
func (j *Junjoold) CreateTaskUnits(units []*types.TaskUnit) ([]types.TaskUnitID, error) {
	return j.
		storageImplementation.
		CreateTaskUnits(units)
}

// Assign a drafted `Job` to a `Topic` for processing
// Consider every orphan `Job` as a draft (that you might take in charge for deletion)
func (j *Junjoold) AssignJob(topicID types.TopicID, jobID types.JobID) error {
	return j.storageImplementation.AssignJob(topicID, jobID)
}

// Assign a drafted `Task` to a `Job`
// Consider every orphan `Task` as a draft (that you might take in charge for deletion)
func (j *Junjoold) AssignTask(jobID types.JobID, taskID types.TaskID) error {
	return j.storageImplementation.AssignTask(jobID, taskID)
}

// Assign a drafted `TaskUnits` to a `Task`
// Consider every orphan `TaskUnit` as a draft (that you might take in charge for deletion)
func (j *Junjoold) AssignTaskUnits(taskID types.TaskID, ids []types.TaskUnitID) error {
	return j.storageImplementation.AssignTaskUnits(taskID, ids)
}

// Workers/Owners will only see the tasks their need to accomplish
func (j *Junjoold) GetInbox(ownerID types.OwnerID, cfgs ...types.QueryConfig) ([]types.InboxAllTaskUnit, error) {
	params := types.NewQuery(cfgs...)
	return j.storageImplementation.GetInbox(ownerID, params)
}
