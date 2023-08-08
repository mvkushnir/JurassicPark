How to run:
* Set up DB tables with `go build setup.go`
* Run API with `go build main.go`
* Given more time, or preparing for an actual prod environment, I'd want to containerize this to make everything self contained, and allow for easy local dev. Currently this assumes you already have a working MySQL server.

Requirements / functionality:
* Create a cage:

`curl -X POST -H 'Content-Type: application/json' -d '{"number": 1, "powerStatus": "ACTIVE", "capacity": 10}' 'localhost:8080/cages'`

`curl -X POST -H 'Content-Type: application/json' -d '{"number": 2, "powerStatus": "ACTIVE", "capacity": 2}' 'localhost:8080/cages'`

* View all cages, optionally filtering by power status:

`curl 'localhost:8080/cages'`

`curl 'localhost:8080/cages?powerStatus=ACTIVE'`

* View details of a specific cage

`curl 'localhost:8080/cages/1'`

* View dinosaurs in a specific cage

`curl 'localhost:8080/cages/1/dinosaurs'`

* Create a dinosaur, and put it in a specific cage (will error if the cage does not exist, another dinosaur already exists with the given name, or the cage conditions are not satisfied (trying to place a carnivore with another species, trying to place an herbivore with a carnivore, or cage is full / not active))

`curl -v -X POST -H 'Content-Type: application/json' -d '{"name": "Jeff", "cage": 1, "species": "Velociraptor"}' 'localhost:8080/dinosaurs'`

* View a dinosaur (by name)

`curl 'localhost:8080/dinosaurs/Jeff'`

* View all dinosaurs (optionally filtering by species)

`curl 'localhost:8080/dinosaurs?species=Velociraptor'`

* Move dinosaur to another cage (by dinosaur name and cage number)

`curl -X PUT 'localhost:8080/dinosaurs/Jeff/1'`

* Change power status of cage (fails if dinosaurs are in cage)

`curl -X PUT 'localhost:8080/cages/2/DOWN'`

Assumptions:
* All dinosaurs can be uniquely identified by name
* The only operations we might want to do are to move a dinosaur to a different cage, and change the status of a cage

Future work:
* Input validation - currently the API just assumes that valid values are being given for dinosaur species and cage power statuses
* Unit & integration testing - ideally would have full unit test coverage for all functionality, both on the API and core layers
* Better error handling & messaging - I'd want to do a better job of detecting specific types of errors and returning appropriate HTTP responses, rather than just passing the exact error along.
* I picked this simple SQL driver for a quick solution, but in a production codebase I'd probably want to use some driver that does things in a more object-oriented way and lets us hide the SQL details.
* If this were to be built for a concurrent environment, first I'd ensure that locking is implemented so that only one instance is accessing the database at a time. If we were to use caching then we might have to add some sort of messaging queues or some other inter-instance communication to ensure updates are correctly propagated, but this implementation is currently simple enough that there's no need; the database is the sole source of truth.