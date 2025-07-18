# Task Manager Go

A high-performance, concurrent task management system built in Go that provides a scalable architecture for processing different types of tasks using a worker pool pattern.

---

## ✨ Features
- Thread-safe, buffered task queue
- Configurable worker pool for concurrent processing
- Dynamic processor registry for extensibility
- Modular, interface-driven architecture
- Graceful shutdown and task tracking
- Easy to add new task types and processors

---

## 🏗️ Architecture

### Core Components

The system is built around several key components that work together to provide a robust task processing pipeline:

#### 1. **Task Queue** (`pkg/queue/queue.go`)
- **Purpose**: Manages task distribution and synchronization
- **Features**:
  - Thread-safe task queuing with buffered channels
  - Graceful shutdown capabilities
  - Task completion tracking with `sync.WaitGroup`
  - Atomic shutdown state management
- **Key Methods**:
  - `AddTask(task any)`: Adds tasks to the queue
  - `GetTask()`: Retrieves tasks for processing
  - `WaitForTasks()`: Blocks until all tasks complete
  - `StartShutdown()`: Initiates graceful shutdown

#### 2. **Worker Pool** (`pkg/workers/worker.go`)
- **Purpose**: Concurrent task processing with multiple workers
- **Features**:
  - Configurable number of worker goroutines
  - Non-blocking task polling with sleep intervals
  - Graceful worker shutdown
  - Automatic task completion notification
- **Key Methods**:
  - `Start()`: Begins worker processing loop
  - `Stop()`: Gracefully stops the worker
  - `processTask()`: Handles individual task processing

#### 3. **Processor Registry** (`pkg/processors/registry/registry.go`)
- **Purpose**: Manages different task processors dynamically
- **Features**:
  - Type-based processor registration
  - Runtime processor lookup
  - Extensible processor system
- **Key Methods**:
  - `AddProcessor(type, processor)`: Registers new processors
  - `GetProcessor(type)`: Retrieves processor by type

#### 4. **Task Model** (`pkg/model/task.go`)
- **Purpose**: Defines the core task structure and interface
- **Features**:
  - Base `Task` struct with common fields
  - `TaskData` interface for type-specific implementations
- **Core Fields**:
  - `ID`: Unique task identifier
  - `Type`: Task type for processor routing
  - `Priority`: Task priority level
  - `CreatedAt`: Task creation timestamp

### Architecture Flow

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Task      │───▶│   Queue     │───▶│   Workers   │
│  Creation   │    │             │    │             │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                    │
                           ▼                    ▼
                   ┌─────────────┐    ┌─────────────┐
                   │  Registry   │◀───│  Processors │
                   │             │    │             │
                   └─────────────┘    └─────────────┘
```

---

## 🔧 Task Management

### How Tasks Are Managed

#### 1. **Task Creation**
Tasks are created by implementing the `TaskData` interface:

```go
type EmailTask struct {
    model.Task
    Receiver string
    Sender   string
    Subject  string
    Content  map[string]string
}

func (t EmailTask) GetTaskType() string {
    return "send_email"
}
```

#### 2. **Task Processing Pipeline**
1. **Submission**: Tasks are added to the queue via `queue.AddTask()`
2. **Distribution**: Workers continuously poll the queue for available tasks
3. **Routing**: Tasks are routed to appropriate processors based on their type
4. **Processing**: Processors handle the specific business logic
5. **Completion**: Workers notify the queue when tasks are complete

#### 3. **Processor System**
- **Interface-based**: All processors implement `TaskProcessor` interface
- **Type-driven**: Processors are registered by task type
- **Extensible**: New task types can be added without modifying core code

Example processor implementation:
```go
type EmailProcessor struct{}

func (ep EmailProcessor) CanProcess(taskType string) bool {
    return taskType == "send_email"
}

func (ep EmailProcessor) ProcessTask(t any) error {
    // Process email-specific logic
    return nil
}
```

#### 4. **Concurrency Management**
- **Worker Pool**: Multiple workers process tasks concurrently
- **Thread Safety**: Queue operations are thread-safe
- **Graceful Shutdown**: System can be stopped without losing tasks
- **Task Tracking**: WaitGroup ensures all tasks complete before shutdown

---

## 🚀 How to Execute

### Prerequisites
- Go 1.19 or higher
- Git (for cloning)

### Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd task-manager-go
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

### Configuration

The system is configured in `cmd/main.go`:

```go
// Queue configuration
queue := taskQueue.NewTaskQueue(10)  // Buffer size of 10

// Registry setup
registry := registry.NewRegistry()
registry.AddProcessor("send_email", &email.EmailProcessor{})

// Worker pool configuration
worker1 := workers.NewWorker(1, &queue, &registry)
worker2 := workers.NewWorker(2, &queue, &registry)
worker3 := workers.NewWorker(3, &queue, &registry)
```

### Adding New Task Types

1. **Create a new task struct**:
   ```go
   type NotificationTask struct {
       model.Task
       Message string
       Channel string
   }
   
   func (t NotificationTask) GetTaskType() string {
       return "send_notification"
   }
   ```

2. **Create a processor**:
   ```go
   type NotificationProcessor struct{}
   
   func (np NotificationProcessor) CanProcess(taskType string) bool {
       return taskType == "send_notification"
   }
   
   func (np NotificationProcessor) ProcessTask(t any) error {
       // Process notification logic
       return nil
   }
   ```

3. **Register the processor**:
   ```go
   registry.AddProcessor("send_notification", &NotificationProcessor{})
   ```

### Execution Flow Example

When you run the application, you'll see output like:
```
Starting worker 1...
Starting worker 2...
Starting worker 3...
Processed email task: ID=email-0, Receiver=user0@example.com, Sender=, Subject=
Processed email task: ID=email-1, Receiver=user1@example.com, Sender=, Subject=
...
Stop worker 1...
Stop worker 2...
Stop worker 3...
```

---

## 📊 Performance Characteristics

- **Concurrent Processing**: Multiple workers process tasks simultaneously
- **Non-blocking**: Workers don't block when no tasks are available
- **Memory Efficient**: Uses buffered channels to limit memory usage
- **Scalable**: Easy to add more workers or task types

---

## 🔍 Monitoring & Debugging

- **Task Completion**: System waits for all tasks to complete before shutdown
- **Error Handling**: Failed tasks are logged but don't crash the system
- **Worker Status**: Each worker logs its start/stop status
- **Queue Status**: Methods available to check queue state (`IsEmpty()`, `RemainingTasks()`)

---

## 🛠️ Extending the System

The modular architecture makes it easy to extend:

1. **New Task Types**: Implement `TaskData` interface
2. **New Processors**: Implement `TaskProcessor` interface
3. **Custom Queues**: Implement different queue strategies
4. **Monitoring**: Add metrics collection
5. **Persistence**: Add database storage for tasks

This system provides a solid foundation for building scalable, concurrent task processing applications in Go. 

---

## 🤝 Contributing

Contributions are welcome! To contribute:

1. **Fork the repository** and create your branch from `main`.
2. **Write clear, concise commit messages**.
3. **Add tests** for new features or bug fixes when possible.
4. **Open a pull request** with a clear description of your changes.
5. Ensure your code follows Go best practices and is formatted with `gofmt`.

For larger changes, please open an issue first to discuss your proposal.

---

## 📜 License

This project is licensed under the [MIT License](LICENSE). You are free to use, modify, and distribute this software. See the LICENSE file for details.

---

## 🧑‍💼 Code of Conduct

Please note that this project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

---

## 💬 Community & Support

- For questions, open an [issue](https://github.com/your-repo/task-manager-go/issues).
- For feature requests or discussions, use GitHub Discussions or open an issue.
- For security concerns, please contact the maintainers directly.

--- 