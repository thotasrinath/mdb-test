# MicroDB POC

### Steps To Start Test

1. Start Cassandra and Create below schema.

```sql
-- Create test keyspace
CREATE KEYSPACE test_db WITH replication = {'class': 'SimpleStrategy' ,'replication_factor':1};

-- Create teacher DB. Here details column is used for storage and retrival of MicroDB
CREATE TABLE teacher(
   id text PRIMARY KEY,
   details blob
   );
```
2. Configure conf.yaml with cassandra url and cache size.

```yaml
cassandra-url: localhost:9042
cache-size: 3
```

3. Start App in the application folder

```shell
go run .
```

### Available API's For Test

-- Here, Teacher table is in Cassandra and Students will be added as MicroDB to each teacher.

-- Add Student 
```cmd
curl --location --request POST 'http://localhost:8080/addstudent?teacherId=4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Code": "Engg",
    "Name": "Srinath",
    "Program":"Masters"
}'
```
-- Get all students for that teacher
```cmd
curl --location --request GET 'http://localhost:8080/getstudents?teacherId=1'
```