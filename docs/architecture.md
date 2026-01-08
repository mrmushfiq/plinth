# Architecture Overview

Plinth is a distributed object storage system designed for reliability, cost-awareness, and ML workload optimization.

## System Components

### 1. API Gateway

The API gateway is the entry point for all client requests.

**Responsibilities:**
- S3-compatible HTTP API
- Request routing
- Authentication and authorization
- Request/response translation

**Technology:**
- Go HTTP server
- S3 protocol compatibility
- Optional AWS SigV4 authentication

### 2. Metadata Service

Central metadata store using PostgreSQL.

**Responsibilities:**
- Object metadata (keys, versions, etags)
- Bucket management
- Placement information (which nodes store replicas)
- Cost tracking
- Access patterns

**Schema:**
```
buckets
  - id, name, versioning_enabled, created_at

objects
  - id, bucket_name, object_key, version_id
  - size_bytes, etag, content_type
  - placement (jsonb - node IDs)
  - state (pending/committed/tombstoned)

cost_tracking
  - bucket_name, total_bytes, estimated_monthly_cost

node_health
  - node_id, status, disk_usage, last_heartbeat
```

### 3. Placement Controller

Determines where objects should be stored using consistent hashing.

**Responsibilities:**
- Node selection for new objects
- Rebalancing on node addition/removal
- Consistent hashing (minimizes rebalancing)
- Tier-aware placement (hot/warm/cold)

**Algorithm:**
```
1. Hash object key
2. Find nodes on the hash ring
3. Select RF (replication factor) nodes
4. Consider node health and capacity
5. Return node list
```

### 4. Data Nodes

Storage nodes that hold actual object data.

**Responsibilities:**
- Store objects on local disk
- Serve read requests
- Calculate and verify checksums (xxHash)
- Health reporting
- gRPC API for internal communication

**Storage Layout:**
```
data/
├── node1/
│   ├── objects/
│   │   ├── ab/cd/abcd1234...  # Object data
│   │   └── ef/gh/efgh5678...
│   └── metadata/              # Local metadata cache
```

### 5. Repair Worker

Background service that maintains data durability.

**Responsibilities:**
- Detect under-replicated objects
- Rebuild replicas from healthy nodes
- Verify checksums (scrubbing)
- Handle node failures

**Repair Cycles:**
- Repair: Every 60 seconds
- Scrub: Every 5 minutes

## Data Flow

### Write Path (PUT Object)

```
1. Client → Gateway (HTTP PUT)
2. Gateway validates request
3. Gateway queries Placement Controller
   → Returns N node IDs (where N = replication factor)
4. Gateway writes to Metadata Service
   → Creates object record (state: pending)
5. Gateway writes to Data Nodes in parallel
   → Waits for W successes (write quorum)
   → Each node calculates checksum
6. Gateway updates Metadata Service
   → Updates placement, state: committed
7. Gateway responds to client with ETag
```

**Quorum Example:**
- Replication Factor (RF) = 3
- Write Quorum (W) = 2
- System waits for 2 successful writes out of 3

### Read Path (GET Object)

```
1. Client → Gateway (HTTP GET)
2. Gateway queries Metadata Service
   → Gets object metadata + placement (node IDs)
3. Gateway reads from Data Nodes
   → Tries first available node
   → Falls back to other replicas on error
4. Gateway verifies checksum
5. Gateway streams data to client
```

**Read Optimization:**
- Prefer local/nearby nodes
- Parallel reads for range requests
- Read repair on checksum mismatch

### Delete Path

```
1. Client → Gateway (HTTP DELETE)
2. Gateway creates delete marker in Metadata
   → state: tombstoned
3. Background garbage collection
   → Eventually deletes from data nodes
```

## Consistency Model

### Write Consistency

- **Quorum Writes**: Must succeed on W out of RF nodes
- **Idempotent**: Same PUT can be retried safely
- **Eventual Consistency**: Replicas converge over time

### Read Consistency

- **Read-Your-Writes**: Guaranteed via metadata
- **Monotonic Reads**: Same object returns consistent data
- **Eventual Consistency**: Different replicas may lag briefly

### Conflict Resolution

- **Last-Write-Wins**: Based on timestamp
- **Version Vectors**: Optional for versioned buckets

## Failure Scenarios

### Data Node Failure

1. Node stops responding to health checks
2. Metadata marks node as "offline"
3. Gateway routes requests to other replicas
4. Repair worker detects under-replicated objects
5. Repair worker copies from healthy replicas to new node

**Recovery Time:**
- Detection: ~30 seconds (health check interval)
- Repair: Depends on object count (target: <10min for 1000 objects)

### Metadata Database Failure

- **Single Point of Failure** (for now)
- Mitigation: PostgreSQL replication (future)
- Mitigation: Regular backups

### Network Partition

- Gateway continues serving if it can reach metadata + quorum nodes
- Partitioned nodes cannot serve writes
- Reads continue from accessible nodes
- Repair runs after partition heals

### Checksum Mismatch

1. Gateway reads object from node
2. Verifies checksum
3. On mismatch: marks replica as corrupt
4. Reads from different replica
5. Repair worker rebuilds corrupt replica

## Scaling

### Horizontal Scaling

**Add Data Node:**
```bash
docker-compose scale datanode=5
```

Placement controller automatically includes new nodes.
Rebalancing happens gradually via repair worker.

**Add Gateway:**
```bash
docker-compose scale gateway=3
```

Put load balancer in front of gateways.

### Vertical Scaling

- Data nodes: More disk, more objects
- Gateway: More CPU/memory, more concurrent requests
- Metadata: Larger PostgreSQL instance

## Performance Characteristics

### Latency

- **PUT**: 1 network RTT + disk write (target: <100ms p95)
- **GET**: 1 network RTT + disk read (target: <50ms p95)
- **DELETE**: Fast (metadata only, ~10ms)

### Throughput

- **Bottleneck**: Typically network or disk I/O
- **Parallelism**: Multiple gateways + data nodes scale linearly
- **Optimization**: Batch operations for ML workloads

### Capacity

- **Per Node**: Limited by disk size
- **Total**: Sum of all node capacities × (1 / RF)
- **Example**: 3 nodes × 1TB × (1/3) = 1TB usable

## Security

### Authentication

- AWS SigV4 signatures (planned)
- API keys (simple mode)
- Pre-signed URLs for temporary access

### Authorization

- Bucket-level policies (planned)
- Object-level ACLs (future)

### Encryption

- In-transit: TLS for HTTP and gRPC
- At-rest: Optional (future)

## Cost Model

### Storage Costs

```
Cost = (Total Bytes / 1GB) × Cost Per GB × RF
```

**Tiering:**
- Hot tier (SSD): $0.10/GB/month
- Warm tier (HDD): $0.05/GB/month
- Cold tier (Archive): $0.01/GB/month

### Egress Costs

```
Cost = (Bytes Read / 1GB) × Egress Cost Per GB
```

Default: $0.09/GB

### Tracking

- Per-bucket cost calculation
- Access pattern tracking
- Admin API for cost reports

## Observability

### Metrics (Prometheus)

```
plinth_objects_total
plinth_objects_under_replicated
plinth_repair_operations_total
plinth_node_disk_usage_bytes
plinth_http_requests_total
plinth_http_request_duration_seconds
```

### Logging

- Structured JSON logs
- Trace IDs for request correlation
- Log levels: debug, info, warn, error

### Health Checks

```
GET /health
→ {"status": "healthy", "nodes": [...]}
```

## Future Enhancements

1. **Erasure Coding**: Reduce storage overhead
2. **Multi-DC Replication**: Geographic redundancy
3. **S3 Select**: Query objects without downloading
4. **Lambda Functions**: Compute on objects
5. **Cross-Region Replication**: Disaster recovery

## References

- [Amazon DynamoDB Paper](https://www.allthingsdistributed.com/files/amazon-dynamo-sosp2007.pdf)
- [Consistent Hashing](https://en.wikipedia.org/wiki/Consistent_hashing)
- [Quorum Consensus](https://en.wikipedia.org/wiki/Quorum_(distributed_computing))

