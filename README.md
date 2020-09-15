# Optimizing communication between serverless functions(Master's Thesis)
- Model developed to show peer-to-peer communication between serverless functions

## Tech Stack
- Golang
- NATS server
- Kubernetes
- Docker
- Ubuntu linux server

## Problem
- Current serverless platforms donot have address aware function manager, i.e the developer doesnot have any idea of where the functions are launched on the serverless platform.

## Solution
- Development of a function aware manager on a custom serverless platform( using golang, docker and kubernetes).
- Implementing a scatter-gather communication that shows the master to distribute jobs to the workers.
- Collection of results from the workers by the master once the job is completed.

## Demo
- 



