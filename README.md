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
> create table emp (eid uuid primary key, name text);
> create table news_by_source (sid varchar, title_hash text, created_at timeuuid, sname varchar, sdesc text, surl text, scategory varchar, slang varchar, scountry varchar, nauthor varchar, ntitle text, ndesc text, nurl text, nurl_to_image text, npublished_at varchar, ncontent text, primary key ((sid), title_hash));
> CREATE TABLE news_sources (sid text, created_at timeuuid, scategory text, scountry text, sdesc text, slang text, sname text, surl text, PRIMARY KEY (scountry, sid));

Insert data
> insert into emp(eid, name) values (uuid(), 'Prateek Gupta');

For running in your system, update the ip in cas.go file to your system ip via ifconfig.

Docker commands
docker build -f Dockerfile .
docker-compose -f docker-compose.yml build
docker-compose -f docker-compose.yml up

Debug the docker image
docker run -it --rm --entrypoint sh <image name>

works


id_rsa
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA7p5sE5f6k6h4g8cqtZ+fjxgUJJLxtHOEz9Vft+OWyOmEmwmhTNZD
Lr0lorQ1q/aSCzUfBPhyJQ8CHsDB7QGSXY9s3XEkhPHkjAcN0DguQwWWz7uUePwHif/EOw
6RMYJRVPA31UaoL4YpPixR7uWREpOJrmaxMTdyvCUSp5J5iEscglf3DoXQ08RRQz3SlSoB
xuxsfYh5Alq9G4T+Zs/QYTPtMXLeBANt0HdFAMm13xMDQHdV4p+Jd2giKz8N22i7aUbWxO
Pcby8IHjVk6lweY27TrfThEt99+YEholF3oSZpPq+Fnh8MRxj5QUvrg9uPKyb3PUHyhRvz
TaTPDiqjuwAAA8iO7Zo+ju2aPgAAAAdzc2gtcnNhAAABAQDunmwTl/qTqHiDxyq1n5+PGB
QkkvG0c4TP1V+345bI6YSbCaFM1kMuvSWitDWr9pILNR8E+HIlDwIewMHtAZJdj2zdcSSE
8eSMBw3QOC5DBZbPu5R4/AeJ/8Q7DpExglFU8DfVRqgvhik+LFHu5ZESk4muZrExN3K8JR
KnknmISxyCV/cOhdDTxFFDPdKVKgHG7Gx9iHkCWr0bhP5mz9BhM+0xct4EA23Qd0UAybXf
EwNAd1Xin4l3aCIrPw3baLtpRtbE49xvLwgeNWTqXB5jbtOt9OES3335gSGiUXehJmk+r4
WeHwxHGPlBS+uD248rJvc9QfKFG/NNpM8OKqO7AAAAAwEAAQAAAQAfZWrU/Sc0LHOG6zq8
YP9OhZ2I3mi9FIICEpIgkOpzDv6qo468nGiEitCb4tg5Ax1eKiQltEbYh2wA/d3GQHGwq8
FoNY3XjDhFEFyJ7ApyORcJyCWV8ZtQVf3Mw3LpL7th0KWYA073ydA9ZPl21G/NIOp5rvtb
fW80QDB/Ke+htfdE9P4KEm3QBOM0Ur5WE+mIhuz8tkrctsakH621mEbSF24au/e3pq4rVS
Sxhiabjz8BVk3EN/MunT8kX4HmD8lL0/DOB0CE7eoo5m3Ln7iaU6r+UV7WxNXAW+LQG+kZ
Ty67xgGGZyaEyKCGE4DQNWLBLgJJpuHXRBX9/qVSDRHpAAAAgB497Mp8VxQAXmNHPUzPxi
iOImARqyYo70vwED7+OAV4oQRup+yOY055LxeN183kaoGi/xN6Za9KOy8vNKLSl9HvduJM
X6KY8KQjzRzBMwgGbfX3KZ2OYZLGz8ANnH8K6g7OyoqauuVDYHvq1X2XvLkuXak2HNHK7W
igOo+w0yBBAAAAgQD4Fo/yCvl6kPh4O6B0t5VDQ9G0uJ2IlmJLfPjUBnpJjajwG5/6i0rW
WMNnVNqdssnfzRUYcoppGMVnGfqvNSIHuP7EO2L8EAipyMsGeGkoSQq9EKin1cOEm7CAr/
nuupUnNSJMiH7XWJR7+2MEojTiw1Ww6+03LwgJH/5JNhRJTwAAAIEA9jqM//cJHoXuCrEL
wcVoV6jiNMvbPy8XTxC/eK9sZ2AyisgJ+eW5R6LUfgbQTdPwPFd3wsBWdTslTlTY80fZDa
eiocDPgPqGMufRyAs+J4hCJTislhlJFZ+sfkVq02k+7oaxyi//mMZJMXVo64V2unX8VpEv
VwQU4wVnfFpYy9UAAAAOcGlAcmFzcGJlcnJ5cGkBAgMEBQ==
-----END OPENSSH PRIVATE KEY-----


id_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDunmwTl/qTqHiDxyq1n5+PGBQkkvG0c4TP1V+345bI6YSbCaFM1kMuvSWitDWr9pILNR8E+HIlDwIewMHtAZJdj2zdcSSE8eSMBw3QOC5DBZbPu5R4/AeJ/8Q7DpExglFU8DfVRqgvhik+LFHu5ZESk4muZrExN3K8JRKnknmISxyCV/cOhdDTxFFDPdKVKgHG7Gx9iHkCWr0bhP5mz9BhM+0xct4EA23Qd0UAybXfEwNAd1Xin4l3aCIrPw3baLtpRtbE49xvLwgeNWTqXB5jbtOt9OES3335gSGiUXehJmk+r4WeHwxHGPlBS+uD248rJvc9QfKFG/NNpM8OKqO7 pi@raspberrypi