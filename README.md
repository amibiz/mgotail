# mgotail

mgotail is a simple library to tail mongodb [oplogs](http://docs.mongodb.org/manual/core/replica-set-oplog/) in Go.

## Documentation

[![GoDoc](https://godoc.org/github.com/Clever/mgotail?status.png)](https://godoc.org/github.com/Clever/mgotail).

## Tests

The tests assume an env variable, MONGO_URL, that defines a mongo connection string.
This mongo instance must be running a replicaset db named "TailTest".
You can do this by running mongod from the command line:

```
mongod --noprealloc --nojournal --smallfiles --oplogSize 10 --replSet TailTest
```

Or via docker

```
docker run -d -p 27017 rgarcia/mongodb mongod --noprealloc --nojournal --smallfiles --oplogSize 10 --replSet TailTest
```

Once mongodb is running, you can run the tests:

```
MONGO_URL=... go test
```

## Vendoring

Please view the [dev-handbook for instructions](https://github.com/Clever/dev-handbook/blob/master/golang/godep.md).
