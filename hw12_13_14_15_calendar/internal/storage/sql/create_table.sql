-- 1
DROP TABLE IF EXISTS events;
-- 2
create table events (
    id serial primary key,
    owner_id varchar(20),
    title text,
    descr text,
    start_time timestamp not null,
    end_time timestamp not null
);
-- 3
create index owner_idx on events (owner_id);
create index start_idx on events using btree (start_time);
--4.1.
GRANT SELECT,
    INSERT ON TABLE events TO otus_hw_user;
-- 4.2.    
GRANT USAGE,
    SELECT ON SEQUENCE events_id_seq TO otus_hw_user;