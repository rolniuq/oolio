-- Create order_queue table for batch processing
CREATE TABLE IF NOT EXISTS order_queue (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_req JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    error TEXT,
    order_data JSONB,
    retry_count INTEGER NOT NULL DEFAULT 0
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_order_queue_status ON order_queue(status);
CREATE INDEX IF NOT EXISTS idx_order_queue_created_at ON order_queue(created_at);
CREATE INDEX IF NOT EXISTS idx_order_queue_retry_count ON order_queue(retry_count);

-- Create index for complex queries used by worker
CREATE INDEX IF NOT EXISTS idx_order_queue_worker_fetch ON order_queue(status, retry_count, created_at);