Simple Login service in Go
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
> create table emp(eid uuid primarykey, name text);

Insert data
> insert into emp(id, name) values (uuid(), 'Prateek Gupta');