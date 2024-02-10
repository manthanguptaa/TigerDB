# TigerDB: A Scalable and Consistent Distributed Cache

**TigerDB** provides a horizontally scalable distributed cache, offering high performance and availability for key-value storage. It utilizes the Raft consensus algorithm for strong consistency and fault tolerance, ensuring data integrity even during node failures.

**Key Features:**

- **Scalability:** Effortlessly scale your cache by adding nodes to the cluster.
- **High Availability:** Experience uninterrupted data access despite node issues.
- **Consistency:** Robust Raft consensus guarantees consistent data across all nodes.
- **Performance:** Efficiently cache frequently accessed data for faster retrieval.
- **Open Source:** Freely use and modify TigerDB under the [insert license name].

**Getting Started:**

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/TigerDB
   ```

**Usage:**

**Client API:**

- Connect to the cache using the provided client library.
- Perform `SET`, `GET`, `Has`, and `Delete` operations on key-value pairs.
- Example usage will be included in the client library documentation.

**Cache Implementation:**

- The `Cacher` interface defines the core cache operations.
- The `Cache` struct implements `Cacher` with a thread-safe in-memory cache and optional TTL support.

**Roadmap:**

- The raft consensus isn't working properly right now. I still have to figure out why it isn't working. Any help is appreciated.