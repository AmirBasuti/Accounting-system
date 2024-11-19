CREATE TABLE voucher (
                         id SERIAL PRIMARY KEY,
                         number VARCHAR(64) NOT NULL UNIQUE,
                         version INT NOT NULL DEFAULT 1
);

CREATE TABLE voucher_item (
                              id SERIAL PRIMARY KEY,
                              voucher_id INT NOT NULL REFERENCES voucher(id) ON DELETE CASCADE,
                              sl_id INT NOT NULL REFERENCES sl(id),
                              dl_id INT REFERENCES dl(id),
                              debit INT CHECK (debit >= 0),
                              credit INT CHECK (credit >= 0)
);
