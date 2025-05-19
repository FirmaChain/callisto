/*######################### cosmos00.sql ##########################*/

/* change vacuum rules */
ALTER TABLE block
SET (
autovacuum_vacuum_scale_factor = 0,
autovacuum_analyze_scale_factor = 0,
autovacuum_vacuum_threshold = 10000,
autovacuum_analyze_threshold = 10000
);

/* transaction -> transaction partitioned */
ALTER TABLE transaction RENAME TO transaction_before_050;

CREATE TABLE transaction (
    hash         TEXT    NOT NULL,
    height       BIGINT  NOT NULL REFERENCES block (height),
    success      BOOLEAN NOT NULL,
    messages     JSON   NOT NULL DEFAULT '[]'::JSON,
    memo         TEXT,
    signatures   TEXT[]  NOT NULL,
    signer_infos JSONB   NOT NULL DEFAULT '[]'::JSONB,
    fee          JSONB   NOT NULL DEFAULT '{}'::JSONB,
    gas_wanted   BIGINT           DEFAULT 0,
    gas_used     BIGINT           DEFAULT 0,
    raw_log      TEXT,
    logs         JSONB,
    events       JSONB NOT NULL DEFAULT '[]'::JSONB,
    partition_id BIGINT  NOT NULL DEFAULT 0,

    CONSTRAINT unique_tx UNIQUE (hash, partition_id)
) PARTITION BY LIST(partition_id);

CREATE INDEX transaction_hash_index ON transaction (hash);
CREATE INDEX transaction_height_index ON transaction (height);
CREATE INDEX transaction_partition_id_index ON transaction (partition_id);

CREATE TABLE transaction_0 PARTITION OF transaction
FOR VALUES IN (0);

INSERT INTO transaction (
    hash, height, success, messages, memo, signatures,
    signer_infos, fee, gas_wanted, gas_used, raw_log, logs,
    events, partition_id
)
SELECT
    hash,
    height,
    success,
    messages::JSON,
    memo,
    signatures,
    signer_infos,
    fee,
    gas_wanted,
    gas_used,
    raw_log,
    logs,
    '[]'::JSONB, -- Check if it is a new column
    0            -- defualt value
FROM transaction_old;

/* Create message_type table MOVED IN migrateMsgTypes function
CREATE TABLE message_type
(
    type      TEXT   NOT NULL UNIQUE,
    module    TEXT   NOT NULL,
    label     TEXT   NOT NULL,
    height    BIGINT NOT NULL
);
CREATE INDEX message_type_module_index ON message_type (module);
CREATE INDEX message_type_type_index ON message_type (type); */

/* message -> message partitioned */
ALTER TABLE message RENAME TO message_old;

CREATE TABLE message (
    transaction_hash             TEXT   NOT NULL,
    index                        BIGINT NOT NULL,
    type                         TEXT   NOT NULL REFERENCES message_type(type),
    value                        JSON   NOT NULL,
    involved_accounts_addresses  TEXT[] NOT NULL,
    partition_id                 BIGINT NOT NULL DEFAULT 0,
    height                       BIGINT NOT NULL,
    FOREIGN KEY (transaction_hash, partition_id) REFERENCES transaction (hash, partition_id),
    CONSTRAINT unique_message_per_tx UNIQUE (transaction_hash, index, partition_id)
) PARTITION BY LIST(partition_id);

CREATE TABLE message_0 PARTITION OF message
FOR VALUES IN (0);

CREATE INDEX message_transaction_hash_index ON message (transaction_hash);
CREATE INDEX message_type_index ON message (type);
CREATE INDEX message_involved_accounts_index ON message USING GIN(involved_accounts_addresses);

INSERT INTO message (
    transaction_hash, index, type, value,
    involved_accounts_addresses, partition_id, height
)
SELECT
    m.transaction_hash,
    m.index,
    m.type,
    m.value::JSON,
    COALESCE(m.involved_accounts_addresses, ARRAY[]::TEXT[]),
    t.partition_id,
    t.height
FROM message_old m;

/* new message_by_address function */
DROP FUNCTION IF EXISTS messages_by_address(
    TEXT[],
    TEXT[],
    BIGINT,
    BIGINT
);

CREATE FUNCTION messages_by_address(
    addresses TEXT[],
    types TEXT[],
    "limit" BIGINT = 100,
    "offset" BIGINT = 0)
    RETURNS SETOF message AS
$$
SELECT * FROM message
WHERE (cardinality(types) = 0 OR type = ANY (types))
  AND addresses && involved_accounts_addresses
ORDER BY height DESC LIMIT "limit" OFFSET "offset"
$$ LANGUAGE sql STABLE;

/* Create message_by_address function */
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


/*######################### auth.sql ##########################*/
/* Create table vesting_account */
CREATE TABLE vesting_account
(
    id                  SERIAL                          PRIMARY KEY NOT NULL,
    type                TEXT                            NOT NULL,
    address             TEXT                            NOT NULL REFERENCES account (address),
    original_vesting    COIN[]                          NOT NULL DEFAULT '{}',
    end_time            TIMESTAMP WITHOUT TIME ZONE     NOT NULL,
    start_time          TIMESTAMP WITHOUT TIME ZONE
);

CREATE UNIQUE INDEX vesting_account_address_idx ON vesting_account (address);

/* Create table vesting_period */

CREATE TABLE vesting_period
(
    vesting_account_id  BIGINT  NOT NULL REFERENCES vesting_account (id),
    period_order        BIGINT  NOT NULL,
    length              BIGINT  NOT NULL,
    amount              COIN[]  NOT NULL DEFAULT '{}'
);

-- Can table ACCOUNT_BALANCE be eliminated?

/*######################### staking.sql ##########################*/
/* Add column to table vesting_account */
ALTER TABLE staking_pool 
ADD COLUMN    unbonding_tokens         TEXT    NOT NULL DEFAULT '',
ADD COLUMN    staked_not_bonded_tokens TEXT    NOT NULL DEFAULT '';

/* Drop column from table validator_status */
ALTER TABLE validator_status 
DROP COLUMN    tombstoned;

/*######################### gov.sql ##########################*/
/* Add column to table gov_params */
ALTER TABLE gov_params 
ADD COLUMN    params         JSONB    NOT NULL DEFAULT '[]'::JSONB;

/* Drop columns from table gov_params */
ALTER TABLE gov_params 
DROP COLUMN    deposit_params,
DROP COLUMN    voting_params,
DROP COLUMN    tally_params;

/* Drop columns from table proposal */
ALTER TABLE proposal
DROP COLUMN    proposal_route,
DROP COLUMN    proposal_type;

/* Add column to table proposal */
ALTER TABLE proposal
ADD COLUMN    metadata    TEXT    NOT NULL DEFAULT '';

/* Set default for content column */
ALTER TABLE proposal
ALTER COLUMN    content    SET DEFAULT '[]'::JSONB;

/* Drop constraint unique_deposit from table proposal_deposit */
ALTER TABLE proposal_deposit
DROP CONSTRAINT    unique_deposit;

/* Drop constraint proposal_deposit_height_fkey from table proposal_deposit */
ALTER TABLE proposal_deposit
DROP CONSTRAINT IF EXISTS    proposal_deposit_height_fkey;

/* Add columns to table proposal_deposit */
ALTER TABLE proposal_deposit
ADD COLUMN    timestamp    TIMESTAMP,
ADD COLUMN    transaction_hash    TEXT    NOT NULL DEFAULT '';

/* Add constraint to table proposal_deposit */
ALTER TABLE proposal_deposit
ADD CONSTRAINT    unique_deposit    UNIQUE    (proposal_id, depositor_address, transaction_hash);

/* Drop constraint unique_deposit from table proposal_vote */
ALTER TABLE    proposal_vote
DROP CONSTRAINT    unique_vote;

/* Drop constraint proposal_deposit_height_fkey from table proposal_deposit */
ALTER TABLE    proposal_vote
DROP CONSTRAINT IF EXISTS    proposal_vote_height_fkey;

/* Add columns to table proposal_deposit */
ALTER TABLE proposal_vote
ADD COLUMN    weight    TEXT    NOT NULL DEFAULT '',
ADD COLUMN    timestamp    TIMESTAMP;

/* Add constraint to table proposal_deposit */
ALTER TABLE    proposal_vote
ADD CONSTRAINT    unique_vote    UNIQUE    (proposal_id, voter_address, option);

/*######################### upgrade.sql ##########################*/
CREATE TABLE software_upgrade_plan
(
    proposal_id     INTEGER REFERENCES proposal (id) UNIQUE,
    plan_name       TEXT        NOT NULL,
    upgrade_height  BIGINT      NOT NULL,
    info            TEXT        NOT NULL,
    height          BIGINT      NOT NULL
);
CREATE INDEX software_upgrade_plan_proposal_id_index ON software_upgrade_plan (proposal_id);
CREATE INDEX software_upgrade_plan_height_index ON software_upgrade_plan (height);
