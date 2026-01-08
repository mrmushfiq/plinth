-- Plinth Object Storage Metadata Schema

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Buckets table
CREATE TABLE IF NOT EXISTS buckets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Versioning
    versioning_enabled BOOLEAN DEFAULT FALSE,
    
    -- Metadata
    region VARCHAR(50) DEFAULT 'us-east-1',
    
    CONSTRAINT bucket_name_valid CHECK (name ~ '^[a-z0-9][a-z0-9-]*[a-z0-9]$')
);

-- Object state enum
CREATE TYPE object_state AS ENUM ('pending', 'committed', 'tombstoned');

-- Objects table
CREATE TABLE IF NOT EXISTS objects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bucket_name VARCHAR(255) NOT NULL REFERENCES buckets(name) ON DELETE CASCADE,
    object_key TEXT NOT NULL,
    
    -- Version tracking
    version_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    is_latest BOOLEAN DEFAULT TRUE,
    is_delete_marker BOOLEAN DEFAULT FALSE,
    
    -- Object properties
    size_bytes BIGINT NOT NULL,
    etag VARCHAR(255) NOT NULL,
    content_type VARCHAR(255) DEFAULT 'application/octet-stream',
    
    -- Placement info (stores node IDs where replicas exist)
    placement JSONB NOT NULL,
    
    -- State management
    state object_state DEFAULT 'committed',
    
    -- User metadata
    metadata JSONB,
    
    -- Tags
    tags JSONB,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Composite unique constraint
    UNIQUE (bucket_name, object_key, version_id)
);

-- Index for efficient lookups
CREATE INDEX idx_objects_bucket_key ON objects(bucket_name, object_key);
CREATE INDEX idx_objects_bucket_key_latest ON objects(bucket_name, object_key) WHERE is_latest = TRUE;
CREATE INDEX idx_objects_state ON objects(state);
CREATE INDEX idx_objects_created_at ON objects(created_at);

-- Multipart uploads table
CREATE TABLE IF NOT EXISTS multipart_uploads (
    upload_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bucket_name VARCHAR(255) NOT NULL REFERENCES buckets(name) ON DELETE CASCADE,
    object_key TEXT NOT NULL,
    
    -- Upload metadata
    initiated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    content_type VARCHAR(255),
    metadata JSONB,
    
    -- State
    state VARCHAR(20) DEFAULT 'active' CHECK (state IN ('active', 'completed', 'aborted')),
    
    UNIQUE (bucket_name, object_key, upload_id)
);

-- Multipart upload parts
CREATE TABLE IF NOT EXISTS multipart_parts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    upload_id UUID NOT NULL REFERENCES multipart_uploads(upload_id) ON DELETE CASCADE,
    part_number INTEGER NOT NULL,
    
    -- Part properties
    size_bytes BIGINT NOT NULL,
    etag VARCHAR(255) NOT NULL,
    
    -- Placement
    placement JSONB NOT NULL,
    
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE (upload_id, part_number)
);

-- Cost tracking table
CREATE TABLE IF NOT EXISTS cost_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bucket_name VARCHAR(255) NOT NULL REFERENCES buckets(name) ON DELETE CASCADE,
    
    -- Storage costs
    total_bytes BIGINT DEFAULT 0,
    object_count INTEGER DEFAULT 0,
    
    -- Egress tracking
    bytes_read_last_month BIGINT DEFAULT 0,
    bytes_read_current_month BIGINT DEFAULT 0,
    
    -- Cost configuration (dollars per GB)
    storage_cost_per_gb DECIMAL(10, 6) DEFAULT 0.023,
    egress_cost_per_gb DECIMAL(10, 6) DEFAULT 0.09,
    
    -- Calculated costs
    estimated_monthly_storage_cost DECIMAL(10, 2) GENERATED ALWAYS AS 
        ((total_bytes::DECIMAL / 1073741824) * storage_cost_per_gb) STORED,
    
    last_calculated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE (bucket_name)
);

-- Access patterns table (for tiering decisions)
CREATE TABLE IF NOT EXISTS access_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bucket_name VARCHAR(255) NOT NULL,
    object_key TEXT NOT NULL,
    
    -- Access tracking
    last_accessed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    access_count INTEGER DEFAULT 0,
    total_bytes_read BIGINT DEFAULT 0,
    
    -- Tiering info
    current_tier VARCHAR(20) DEFAULT 'hot' CHECK (current_tier IN ('hot', 'warm', 'cold')),
    
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE (bucket_name, object_key)
);

-- Repair/health tracking
CREATE TABLE IF NOT EXISTS repair_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    object_id UUID NOT NULL REFERENCES objects(id) ON DELETE CASCADE,
    
    -- Issue details
    issue_type VARCHAR(50) NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Resolution
    repaired_at TIMESTAMP WITH TIME ZONE,
    repair_status VARCHAR(20) DEFAULT 'pending' CHECK (repair_status IN ('pending', 'in_progress', 'completed', 'failed')),
    
    -- Details
    source_node VARCHAR(100),
    target_nodes TEXT[],
    error_message TEXT
);

CREATE INDEX idx_repair_log_status ON repair_log(repair_status);
CREATE INDEX idx_repair_log_detected_at ON repair_log(detected_at);

-- Node health table
CREATE TABLE IF NOT EXISTS node_health (
    node_id VARCHAR(100) PRIMARY KEY,
    
    -- Status
    status VARCHAR(20) DEFAULT 'healthy' CHECK (status IN ('healthy', 'degraded', 'offline')),
    
    -- Capacity
    total_disk_bytes BIGINT,
    used_disk_bytes BIGINT,
    available_disk_bytes BIGINT,
    disk_usage_percent DECIMAL(5, 2) GENERATED ALWAYS AS 
        (CASE WHEN total_disk_bytes > 0 THEN (used_disk_bytes::DECIMAL / total_disk_bytes * 100) ELSE 0 END) STORED,
    
    -- Health metrics
    object_count INTEGER DEFAULT 0,
    last_heartbeat_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Node metadata
    tier VARCHAR(20) DEFAULT 'hot',
    labels JSONB,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_buckets_updated_at BEFORE UPDATE ON buckets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_objects_updated_at BEFORE UPDATE ON objects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cost_tracking_updated_at BEFORE UPDATE ON cost_tracking
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_access_patterns_updated_at BEFORE UPDATE ON access_patterns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_node_health_updated_at BEFORE UPDATE ON node_health
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default cost tracking for new buckets
CREATE OR REPLACE FUNCTION create_default_cost_tracking()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO cost_tracking (bucket_name) VALUES (NEW.name);
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER create_bucket_cost_tracking AFTER INSERT ON buckets
    FOR EACH ROW EXECUTE FUNCTION create_default_cost_tracking();

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO plinth;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO plinth;

