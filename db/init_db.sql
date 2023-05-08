CREATE TABLE blocks (
    block_num BIGINT PRIMARY KEY,
    block_hash VARCHAR(255),
    block_time BIGINT,
    parent_hash VARCHAR(255),
    tx_hashes TEXT[]
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    tx_hash VARCHAR(255),
    block_num BIGINT,
    from_address VARCHAR(255),
    to_address VARCHAR(255),
    nonce INT,
    data TEXT,
    value VARCHAR(255),
    logs JSONB,
    FOREIGN KEY (block_num) REFERENCES blocks(block_num)
);
