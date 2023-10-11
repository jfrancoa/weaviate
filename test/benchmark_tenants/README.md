# Weaviate Tenant Benchmark

This Bash script is designed to automate the benchmarking process for tenant creation in Weaviate using a multi-node Kubernetes cluster. The script supports various options, including preparation, running benchmarks, cleanup, and an interactive mode.

## Usage

```bash
./benchmark_script.sh <options>
```

### Options:

- **prepare:** Prepares the benchmark environment by creating a multi-node cluster, adding required Helm repositories, installing Weaviate using Helm charts, and setting up Prometheus for monitoring.

- **run:** Runs the benchmark. Optionally, you can specify the number of tenants as an argument. If not specified, the default is 1000 tenants.

    ```bash
    ./benchmark_script.sh run          # Run with default number of tenants
    ./benchmark_script.sh run 5000     # Run with a specific number of tenants
    ```
    In this step, a dummy schema with multi tenancy is created in the Weaviate cluster and the desired number of tenants will then get created iteratively, one each time. 

- **cleanup:** Cleans up the benchmark environment by uninstalling Prometheus, Weaviate, and deleting the associated namespace. Ending up in the cluster removal using kind.

- **interactive:** Enters interactive mode, allowing you to choose options interactively. You can prepare, run, cleanup, or exit the interactive mode.

## Prerequisites

Before using the script, ensure the following prerequisites are met:

- [kind](https://kind.sigs.k8s.io/): Kubernetes in Docker
- [helm](https://helm.sh/): Package manager for Kubernetes
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/): Kubernetes command-line tool
- [Go](https://golang.org/): Programming language for running the benchmark
- Root permissions: The script requires root permissions for cluster creation. Kind doesn't support rootless creation. Check [this Github issue](https://github.com/kubernetes-sigs/kind/issues/2094) for more information 


## Running the Script

```bash
./benchmark_script.sh <options>
```

For example:

```bash
./benchmark_script.sh prepare
./benchmark_script.sh run 5000
./benchmark_script.sh run      # Same as running ./benchmark_script.sh run 1000
./benchmark_script.sh cleanup
./benchmark_script.sh interactive
```
