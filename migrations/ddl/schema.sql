CREATE SEQUENCE MySequence OPTIONS (
    sequence_kind="bit_reversed_positive"
);

CREATE TABLE album (
    id INT64 NOT NULL DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE MySequence)),
    title STRING(255) NOT NULL,
    artist STRING(255) NOT NULL,
    price FLOAT64 
) PRIMARY KEY (id);
