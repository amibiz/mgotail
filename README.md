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

## Changing Dependencies

### New Packages

When adding a new package, you can simply use `make vendor` to update your imports.
This should bring in the new dependency that was previously undeclared.
The change should be reflected in [Godeps.json](Godeps/Godeps.json) as well as [vendor/](vendor/).

### Existing Packages

First ensure that you have your desired version of the package checked out in your `$GOPATH`.

When to change the version of an existing package, you will need to use the godep tool.
You must specify the package with the `update` command, if you use multiple subpackages of a repo you will need to specify all of them.
So if you use package github.com/Clever/foo/a and github.com/Clever/foo/b, you will need to specify both a and b, not just foo.

```
# depending on github.com/Clever/foo
godep update github.com/Clever/foo

# depending on github.com/Clever/foo/a and github.com/Clever/foo/b
godep update github.com/Clever/foo/a github.com/Clever/foo/b
```

