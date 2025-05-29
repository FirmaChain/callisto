/*######################### cosmos00.sql ##########################*/

/* change vacuum rules */
ALTER TABLE block
SET (
autovacuum_vacuum_scale_factor = 0,
autovacuum_analyze_scale_factor = 0,
autovacuum_vacuum_threshold = 10000,
autovacuum_analyze_threshold = 10000
);

/* transaction */
-- Add new column to the parent transaction table
ALTER TABLE transaction
ADD COLUMN events JSONB NOT NULL DEFAULT '[]'::JSONB;

-- Change messages column type to JSON in the parent transaction table
ALTER TABLE transaction
ALTER COLUMN messages DROP DEFAULT,
ALTER COLUMN messages TYPE JSON USING messages::JSON,
ALTER COLUMN messages SET DEFAULT '[]'::JSON;

-- Add new column to the child transaction tables
SELECT format(
  'ALTER TABLE %I ADD COLUMN events JSONB NOT NULL DEFAULT ''[]''::JSONB;',
  relname
)
FROM pg_class
WHERE relname LIKE 'transaction_%' AND relkind = 'r';

-- Change messages column type to JSON in the child transaction tables
SELECT format(
  $$
  ALTER TABLE %I
  ALTER COLUMN messages DROP DEFAULT,
  ALTER COLUMN messages TYPE JSON USING messages::JSON,
  ALTER COLUMN messages SET DEFAULT '[]'::JSON;
  $$,
  relname
)
FROM pg_class
WHERE relname LIKE 'transaction_%' AND relkind = 'r';

/* message */
-- Change value column from JSONB to JSON
ALTER TABLE message
ALTER COLUMN value DROP DEFAULT,
ALTER COLUMN value TYPE JSON USING value::JSON;

-- Add foreign key to message_type(type)
ALTER TABLE message
ADD CONSTRAINT message_type_fk
FOREIGN KEY (type) REFERENCES message_type(type);

-- Alter child message tables
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN
        SELECT inhrelid::regclass AS partition_name
        FROM pg_inherits
        WHERE inhparent = 'message'::regclass
    LOOP
        RAISE NOTICE 'Altering partition: %', r.partition_name;

        -- Add FK constraint on type column
        EXECUTE format(
            'ALTER TABLE %I
             ADD CONSTRAINT %I
             FOREIGN KEY (type) REFERENCES message_type(type);',
            r.partition_name,
            r.partition_name || '_type_fk'
        );
    END LOOP;
END
$$;


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
