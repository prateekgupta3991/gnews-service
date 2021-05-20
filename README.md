Simple News service in Go. Integrates with Google news api - https://newsapi.org/
Uses cassandra as storage layer and Go for backend. Docker deployable.

cassandra cluster
Docker deploy

How to create cassandra cluster in docker
Refer to - https://medium.com/@jayarajjg/setting-up-a-cassandra-cluster-on-your-laptop-using-docker-cf09b1bb651e
Assuming docker installed and cassandra image pulled.
NOTE : I have cassandra 3.11
Commands
Create one node - Exposed on port 9042
> docker run -p 9042:9042  --name my-cassandra-1 -m 2g -d cassandra:3.11

Check the IP
> docker inspect --format='{{ .NetworkSettings.IPAddress }}' my-cassandra-1

Create another node and link to prev node - Exposed on port 9043
> docker run --name my-cassandra-2 -m 2g -d -e CASSANDRA_SEEDS="$(docker inspect --format='{{ .NetworkSettings.IPAddress }}' my-cassandra-1)" cassandra:3.11

To use cassandra via client cqlsh
> docker run -it --link my-cassandra-1 --rm cassandra:3.11 bash -c 'exec cqlsh <<IP>>'

To check cassandra node status
> docker exec -i -t my-cassandra-1 bash -c 'nodetool status'

DB schema
Create keyspace "godemo"
> CREATE KEYSPACE "godemo" with replication = {'class' : 'SimpleStrategy', 'replication_factor' : 3};
> Use godemo

Create table
> create table user (uid int, name text, t_un text, chat_id int, primary key(uid, t_un));
> create table news_by_source (sid varchar, title_hash text, created_at timeuuid, sname varchar, sdesc text, surl text, scategory varchar, slang varchar, scountry varchar, nauthor varchar, ntitle text, ndesc text, nurl text, nurl_to_image text, npublished_at varchar, ncontent text, primary key ((sid), title_hash));
> CREATE TABLE news_sources (sid text, created_at timeuuid, scategory text, scountry text, sdesc text, slang text, sname text, surl text, PRIMARY KEY (scountry, sid));

Insert data
> insert into user(uid, name, t_un) values (1367340022, 'Prateek Gupta', 'Prtkgpt');

For running in your system, update the ip in cas.go file to your system ip via ifconfig.

Docker commands
docker build -f Dockerfile .
docker-compose -f docker-compose.yml build
docker-compose -f docker-compose.yml up

Debug the docker image
docker run -it --rm --entrypoint sh <image name>

works

Telegram bot
Poll for update
https://api.telegram.org/bot1853514787:AAHEi4brq8vXE39sYIqPTfFzfYNPvDDWmY0/getUpdates
Webhook
Set - https://api.telegram.org/bot{bot_token}/setWebhook?url={your_server_url}
Delete - https://api.telegram.org/bot{bot_token}/deleteWebhook
Send msg to user
https://api.telegram.org/bot{bot_token}/sendMessage?chat_id={chat_id}&text={text}