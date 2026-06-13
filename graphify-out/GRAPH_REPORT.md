# Graph Report - .  (2026-06-13)

## Corpus Check
- Corpus is ~7,693 words - fits in a single context window. You may not need a graph.

## Summary
- 224 nodes · 337 edges · 19 communities (12 shown, 7 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 48 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Core Queue Infrastructure|Core Queue Infrastructure]]
- [[_COMMUNITY_HTTP Request Context|HTTP Request Context]]
- [[_COMMUNITY_Job Data Model & Database|Job Data Model & Database]]
- [[_COMMUNITY_Job Lifecycle & Main Queue|Job Lifecycle & Main Queue]]
- [[_COMMUNITY_Delayed Queue Operations|Delayed Queue Operations]]
- [[_COMMUNITY_HTTP Handlers & Health|HTTP Handlers & Health]]
- [[_COMMUNITY_Job API Endpoints|Job API Endpoints]]
- [[_COMMUNITY_Base Queue Implementation|Base Queue Implementation]]
- [[_COMMUNITY_Logger Queue & Events|Logger Queue & Events]]
- [[_COMMUNITY_Job Repository Layer|Job Repository Layer]]
- [[_COMMUNITY_Configuration Management|Configuration Management]]
- [[_COMMUNITY_Email Sending|Email Sending]]
- [[_COMMUNITY_Email Validation|Email Validation]]
- [[_COMMUNITY_Commit Message Validation|Commit Message Validation]]
- [[_COMMUNITY_Project Documentation|Project Documentation]]
- [[_COMMUNITY_Dev Tools Configuration|Dev Tools Configuration]]
- [[_COMMUNITY_DLQ Delete Operations|DLQ Delete Operations]]
- [[_COMMUNITY_PR Template|PR Template]]
- [[_COMMUNITY_Validate Script|Validate Script]]

## God Nodes (most connected - your core abstractions)
1. `Context` - 36 edges
2. `QueueManager` - 32 edges
3. `New()` - 16 edges
4. `DelayedQueue` - 12 edges
5. `Config` - 11 edges
6. `LoggerQueue` - 11 edges
7. `MainQueue` - 10 edges
8. `run()` - 9 edges
9. `run()` - 9 edges
10. `Mailer` - 9 edges

## Surprising Connections (you probably didn't know these)
- `run()` --calls--> `LoadConfig()`  [INFERRED]
  cmd/mailer/main.go → internal/config/config.go
- `run()` --calls--> `NewMailer()`  [INFERRED]
  cmd/mailer/main.go → internal/emailqueue/mailer.go
- `run()` --calls--> `New()`  [INFERRED]
  cmd/mailer/main.go → internal/emailqueue/queue.go
- `run()` --calls--> `NewJobRepository()`  [INFERRED]
  cmd/mailer/main.go → internal/repository/job_repository.go
- `run()` --calls--> `LoadConfig()`  [INFERRED]
  cmd/server/main.go → internal/config/config.go

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Email Queue Lifecycle System** — basequeue_base_queue, emailqueue_delayed_queue, emailqueue_logger_queue [EXTRACTED 0.95]
- **Server Entry Point Wiring** — server_main, config_config, database_database, emailqueue_delayed_queue, emailqueue_logger_queue [EXTRACTED 0.90]
- **Jobs Lifecycle Schema and Code** — emailqueue_logger_queue, migrations_000002_create_jobs_table, migrations_000002_create_jobs_table_down [INFERRED 0.85]
- **Email Processing Flow** — handler_job_enqueue, emailqueue_queue_manager, emailqueue_main_queue, emailqueue_mailer, concept_worker_pool [INFERRED 0.85]
- **DLQ Management Flow** — handler_job_get_dlq, handler_job_replay_dlq, handler_job_delete_dlq, repository_job_repository, concept_dead_letter_queue [INFERRED 0.85]
- **Job Retry Mechanism** — concept_retry_with_backoff, concept_dead_letter_queue, concept_job_recovery, concept_real_time_scheduler [INFERRED 0.75]

## Communities (19 total, 7 thin omitted)

### Community 0 - "Core Queue Infrastructure"
Cohesion: 0.08
Nodes (25): Batch Logging, Dead Letter Queue (DLQ), Crash Recovery, Real-Time Timer Scheduler, Exponential Backoff Retry, Worker Pool, DelayedQueue, Config (+17 more)

### Community 2 - "Job Data Model & Database"
Cohesion: 0.16
Nodes (15): Job, Message, Config, NewPool(), Context, Pool, main(), printFailedJobs() (+7 more)

### Community 3 - "Job Lifecycle & Main Queue"
Cohesion: 0.14
Nodes (14): Job, JobAttempt, JobState, NewMainQueue(), MainQueue, Sender, JobHandler.GetDLQ, BaseQueue (+6 more)

### Community 4 - "Delayed Queue Operations"
Cohesion: 0.14
Nodes (5): DelayedQueue, jobHeap, BaseQueue, Job, Mutex

### Community 5 - "HTTP Handlers & Health"
Cohesion: 0.14
Nodes (13): NewContext(), Wrap(), Handler, Engine, HealthCheck(), HealthResponse, HandlerFunc, Context (+5 more)

### Community 6 - "Job API Endpoints"
Cohesion: 0.19
Nodes (10): DLQJobResponse, EnqueueRequest, EnqueueResponse, NewJobHandler(), JobHandler, Context, Job, JobAttempt (+2 more)

### Community 7 - "Base Queue Implementation"
Cohesion: 0.20
Nodes (7): NewBaseQueue(), BaseQueue, NewDelayedQueue(), TestDelayedQueueMinHeapOrdering(), TestDelayedQueueShutdownBehavior(), Mutex, T

### Community 8 - "Logger Queue & Events"
Cohesion: 0.21
Nodes (10): LogEvent, LogEventType, NewLoggerQueue(), LoggerQueue, BaseQueue, Context, Duration, Mutex (+2 more)

### Community 9 - "Job Repository Layer"
Cohesion: 0.29
Nodes (6): Context, Job, JobAttempt, Pool, NewJobRepository(), JobRepository

### Community 10 - "Configuration Management"
Cohesion: 0.47
Nodes (8): Config, getEnv(), getEnvAsDuration(), getEnvAsInt(), LoadConfig(), Duration, main(), usage()

### Community 11 - "Email Sending"
Cohesion: 0.47
Nodes (5): Mailer, NewMailer(), Duration, Message, Mutex

### Community 12 - "Email Validation"
Cohesion: 0.40
Nodes (3): TestValidateEmail(), validateEmail(), T

## Knowledge Gaps
- **48 isolated node(s):** `validate.sh script`, `Message`, `Job`, `Mutex`, `HandlerFunc` (+43 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `Core Queue Infrastructure` to `Job Data Model & Database`, `Job Lifecycle & Main Queue`, `Delayed Queue Operations`, `HTTP Handlers & Health`, `Base Queue Implementation`, `Logger Queue & Events`, `Job Repository Layer`?**
  _High betweenness centrality (0.421) - this node is a cross-community bridge._
- **Why does `QueueManager` connect `Core Queue Infrastructure` to `Job Lifecycle & Main Queue`, `Email Sending`, `Job API Endpoints`?**
  _High betweenness centrality (0.352) - this node is a cross-community bridge._
- **Why does `Context` connect `HTTP Request Context` to `Job Data Model & Database`, `HTTP Handlers & Health`?**
  _High betweenness centrality (0.241) - this node is a cross-community bridge._
- **Are the 6 inferred relationships involving `QueueManager` (e.g. with `Batch Logging` and `Dead Letter Queue (DLQ)`) actually correct?**
  _`QueueManager` has 6 INFERRED edges - model-reasoned connections that need verification._
- **Are the 11 inferred relationships involving `New()` (e.g. with `.Pop()` and `.Push()`) actually correct?**
  _`New()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **What connects `validate.sh script`, `Message`, `Job` to the rest of the system?**
  _54 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Core Queue Infrastructure` be split into smaller, more focused modules?**
  _Cohesion score 0.08377896613190731 - nodes in this community are weakly interconnected._