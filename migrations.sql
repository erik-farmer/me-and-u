/* Migration place holder file until I figure out a way I like more.
Until then...it's writing things in the Turso UI
*/

-- Initial Migration
CREATE TABLE IF NOT EXISTS recipes (
    name VARCHAR UNIQUE NOT NULL,
    url VARCHAR,
    ingredients VARCHAR,
    steps VARCHAR,
    notes VARCHAR
);