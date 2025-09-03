CREATE TABLE validator
(
    consensus_address TEXT NOT NULL PRIMARY KEY, /* Validator consensus address */
    consensus_pubkey  TEXT NOT NULL UNIQUE /* Validator consensus public key */
);

CREATE TABLE pre_commit
(
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    voting_power      BIGINT                      NOT NULL,
    proposer_priority BIGINT                      NOT NULL,
    UNIQUE (validator_address, timestamp)
);
CREATE INDEX pre_commit_validator_address_index ON pre_commit (validator_address);
CREATE INDEX pre_commit_height_index ON pre_commit (height); 

CREATE TABLE block
(
    height           BIGINT  UNIQUE PRIMARY KEY,
    hash             TEXT    NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        BIGINT  DEFAULT 0,
    proposer_address TEXT REFERENCES validator (consensus_address),
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX block_height_index ON block (height);
CREATE INDEX block_hash_index ON block (hash);
CREATE INDEX block_proposer_address_index ON block (proposer_address);
ALTER TABLE block
    SET (
        autovacuum_vacuum_scale_factor = 0,
        autovacuum_analyze_scale_factor = 0,
        autovacuum_vacuum_threshold = 10000,
        autovacuum_analyze_threshold = 10000
        );

CREATE TABLE transaction
(
    hash         TEXT    NOT NULL,
    height       BIGINT  NOT NULL REFERENCES block (height),
    success      BOOLEAN NOT NULL,

    /* Body */
    messages     JSON   NOT NULL DEFAULT '[]'::JSON,
    memo         TEXT,
    signatures   TEXT[]  NOT NULL,

    /* AuthInfo */
    signer_infos JSONB   NOT NULL DEFAULT '[]'::JSONB,
    fee          JSONB   NOT NULL DEFAULT '{}'::JSONB,

    /* Tx response */
    gas_wanted   BIGINT           DEFAULT 0,
    gas_used     BIGINT           DEFAULT 0,
    raw_log      TEXT,
    logs         JSONB,

    events JSONB NOT NULL DEFAULT '[]'::JSONB,

    /* PSQL partition */
    partition_id BIGINT  NOT NULL DEFAULT 0,

    CONSTRAINT unique_tx UNIQUE (hash, partition_id)
)PARTITION BY LIST(partition_id);
CREATE INDEX transaction_hash_index ON transaction (hash);
CREATE INDEX transaction_height_index ON transaction (height);
CREATE INDEX transaction_partition_id_index ON transaction (partition_id);

CREATE TABLE message_type
(
    type      TEXT   NOT NULL UNIQUE,
    module    TEXT   NOT NULL,
    label     TEXT   NOT NULL,
    height    BIGINT NOT NULL
);
CREATE INDEX message_type_module_index ON message_type (module);
CREATE INDEX message_type_type_index ON message_type (type);

CREATE TABLE message
(
    transaction_hash            TEXT   NOT NULL,
    index                       BIGINT NOT NULL,
    type                        TEXT   NOT NULL REFERENCES message_type(type),
    value                       JSON  NOT NULL,
    involved_accounts_addresses TEXT[] NOT NULL,

    /* PSQL partition */
    partition_id                BIGINT NOT NULL DEFAULT 0,
    height                      BIGINT NOT NULL,
    FOREIGN KEY (transaction_hash, partition_id) REFERENCES transaction (hash, partition_id),
    CONSTRAINT unique_message_per_tx UNIQUE (transaction_hash, index, partition_id)
)PARTITION BY LIST(partition_id);
CREATE INDEX message_transaction_hash_index ON message (transaction_hash);
CREATE INDEX message_type_index ON message (type);
CREATE INDEX message_involved_accounts_index ON message USING GIN(involved_accounts_addresses);

/**
 * This function is used to find all the utils that involve any of the given addresses and have
 * type that is one of the specified types.
 */
CREATE FUNCTION messages_by_address(
addresses TEXT[],
types TEXT[] DEFAULT NULL::TEXT[],
"limit" BIGINT DEFAULT 100,
"offset" BIGINT DEFAULT 0)
 RETURNS SETOF message
 LANGUAGE plpgsql
 STABLE
AS $function$
DECLARE
    p regclass;
    current_count bigint := 0;
    current_offset bigint := "offset";
    partition_result_count bigint;
    has_match boolean;
    type_condition text;
BEGIN
    -- Check Type Condition
    IF types IS NULL OR array_length(types, 1) IS NULL OR array_length(types, 1) = 0 THEN
        type_condition := 'TRUE';
    ELSE
        type_condition := format('type = ANY(ARRAY[%s])',
            array_to_string(array(select quote_literal(t) from unnest(types) as t), ','));
    END IF;

    FOR p IN
        SELECT inhrelid
        FROM pg_inherits
        WHERE inhparent = 'message'::regclass
        ORDER BY inhrelid DESC
    LOOP
        -- Check EXISTS
        EXECUTE format(
            'SELECT EXISTS (
                SELECT 1 FROM %s
                WHERE %s AND involved_accounts_addresses && $1
                LIMIT 1
            )', p, type_condition
        )
        INTO has_match
        USING addresses;

        IF has_match THEN
            -- Check Counts
            EXECUTE format(
                'SELECT count(*) FROM %s
                WHERE %s AND involved_accounts_addresses && $1',
                p, type_condition
            )
            INTO partition_result_count
            USING addresses;

            -- Process offset
            IF current_offset > 0 THEN
                IF current_offset >= partition_result_count THEN
                    current_offset := current_offset - partition_result_count;
                    CONTINUE;
                ELSE
                    RETURN QUERY
                    EXECUTE format(
                        'SELECT * FROM %s
                        WHERE %s AND involved_accounts_addresses && $1
                        ORDER BY height DESC
                        LIMIT $2 OFFSET $3',
                        p, type_condition
                    )
                    USING addresses, "limit" - current_count, current_offset;

                    current_count := current_count + LEAST(partition_result_count - current_offset, "limit" - current_count);
                    current_offset := 0;
                END IF;
            ELSE
                RETURN QUERY
                EXECUTE format(
                    'SELECT * FROM %s
                    WHERE %s AND involved_accounts_addresses && $1
                    ORDER BY height DESC
                    LIMIT $2',
                    p, type_condition
                )
                USING addresses, "limit" - current_count;

                current_count := current_count + LEAST(partition_result_count, "limit" - current_count);
            END IF;

            IF current_count >= "limit" THEN
                EXIT;
            END IF;
        END IF;
    END LOOP;
END;
$function$

CREATE FUNCTION messages_by_type(
  types text [],
  "limit" bigint DEFAULT 100,
  "offset" bigint DEFAULT 0) 
  RETURNS SETOF message AS 
$$ 
SELECT * FROM message
WHERE (cardinality(types) = 0 OR type = ANY (types))
ORDER BY height DESC LIMIT "limit" OFFSET "offset" 
$$ LANGUAGE sql STABLE;

CREATE TABLE pruning
(
    last_pruned_height BIGINT NOT NULL
);