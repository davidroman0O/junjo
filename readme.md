# Junjo 

> WORK IN PROGRESS

> Hierarchical Workflow Management System with DAG-Structured Tasks in Go

Develop a robust workflow management system in Go, aimed at structuring and executing units of work in an organized, dependency-resolved manner. The system hierarchy encompasses Topics, Jobs, and Tasks, each level offering a layer of organization and control.

Core Components:

1. **Topics**: High-level categorizations of related work, e.g., "Provisioning" for provisioning bare metal machines.
2. **Jobs**: Individual scopes of work within a Topic, e.g., a job with ID `XXXX` under the "Provisioning" topic.
3. **Tasks**: Sets of related, finer-grained units of work associated with a parent Job, organized in a Directed Acyclic Graph (DAG) to dictate the execution order.
4. **Task Units**: The granular work units within a Task's DAG, to be executed in a defined sequence. Each Task Unit represents a specific action required to progress the parent Task towards completion.
5. **Commands**: A set of actions exposed via an API, enabling clients to report Task Unit progress (e.g., progress, success, error, pause, log) and mutate Task Unit states.
6. **API**: The interface for clients to interact with the system, fetch Task Units, and send Commands to report progress.

Workflow:

1. Define Topics and create Jobs under them.
2. For each Job, define a set of Tasks and organize them in a DAG to specify the execution order.
3. Clients fetch available Task Units from the API, based on their assigned Tasks.
4. Clients execute Task Units in accordance with the DAG's order and report progress using Commands via the API.
5. The system processes Commands, updates Task Unit and Task states, and advances the Job towards completion.
6. Upon completion of all Task Units within a Task's DAG, mark the Task as completed.
7. Once all Tasks within a Job are completed, mark the parent Job as completed.

That's it

No processing, job of the client
We just wait for all units to be completed, to valiate a task, which all tasks of a Job need to be validated to complete it.

Storage Implementations:
- in memory
- sqlite3

The developer should be allowed to provide a storage implementation of it's own, `junjo` only provide the mechanism and the rules.
