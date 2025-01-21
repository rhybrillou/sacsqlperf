# SAC SQL Performance analysis

This repo aims at providing tooling to troubleshoot
possible performance issues around scoped access control 
in the [stackrox](github.com/stackrox/stackrox) repository.

The current state of the tool allows to replay specific database queries,
injecting scoped access control filters in the statements, and to extract
database execution plans for these queries.

## Pre-requisites

The tool comes in the form of a built container image that is run
in an existing stackrox deployment with a large scale dataset 
(e.g. 5000 namespace, 20000 deployments, 80000 alerts).

See the [stackrox scale](github.com/stackrox/stackrox/tree/master/scale)
directory for details on how to generate a stackrox deployment
with large scale dataset.

## Running the tool

The tool is a go program. See the go.mod file for go requirements.

The tool can be built using the `make image` command.

To deploy the built image, an image pull secret is needed
in the `stackrox` namespace of the stackrox cluster. The test itself
can be run by applying the `sqltest-deploy.yaml` file on the stackrox cluster.

## Profiling specific queries

At the moment, the queries being profiled are described in structured form
in the `main.go` file in the `testedQueries` variable.

A query object needs knowledge of the target table, as well as the columns
that contain the cluster ID and namespace name information on the scoping table.