CREATE KEYSPACE IF NOT EXISTS "godemo" with replication = {'class' : 'SimpleStrategy', 'replication_factor' : 3};
USE godemo;

CREATE TABLE IF NOT EXISTS user (uid bigint, name text, t_un text, chat_id bigint, primary key(uid, t_un));
CREATE TABLE IF NOT EXISTS news_by_source (sid varchar, title_hash text, created_at timeuuid, sname varchar, sdesc text, surl text, scategory varchar, slang varchar, scountry varchar, nauthor varchar, ntitle text, ndesc text, nurl text, nurl_to_image text, npublished_at varchar, ncontent text, primary key ((sid), title_hash));
CREATE TABLE IF NOT EXISTS news_sources (sid text, created_at timeuuid, scategory text, scountry text, sdesc text, slang text, sname text, surl text, PRIMARY KEY (scountry, sid));

INSERT into user (uid, name, t_un) values (1367340022, 'Prateek Gupta', 'Prtkgpt');