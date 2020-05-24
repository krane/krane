## Developing Krane

We are actively seeking individuals who was contribute with the project. The below docs are a guide on getting started with developing on the krane ecosystem.

## PRs

The Krane Project accepts open pull request, the github actions pipeline is incharge of testing and building the project. If you break anything the pipeline will most likely catch it. Every PR merged to master needs an approval from an admin before merging.

## Bolt DB

Krane uses [bbolt](https://github.com/etcd-io/bbolt) as its datastore. Bbolt is a persisten db with extremely good performance.

To view bbolt in a `gui` like interface you can use the following: [here](https://github.com/br0xen/boltbrowser)
