#!/bin/bash

WEAVIATE_PORT=8383
PROMETHEUS_PORT=9091
GO_BIN_PATH="/usr/local/go/bin"

# Function to be invoked for "prepare" option
prepare_benchmark() {
    echo "benchmark_tenant # Preparing benchmark..."

    # Create a kind cluster using the yaml definition
    echo "benchmark_tenant # Creating multi-node cluster..."
    kind create cluster --config kind/cluster.yaml

    # Add required Helm repos and update the Helm repo
    echo "benchmark_tenant #  Adding required helm repositories..."
    helm repo add weaviate https://weaviate.github.io/weaviate-helm
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update

    # Create the Weaviate namespace
    echo "benchmark_tenant #  Installing weaviate..."
    kubectl create namespace weaviate

    # Install Weaviate using Helm charts
    helm upgrade --install weaviate weaviate/weaviate --namespace weaviate --values helm/values.yaml
    # Wait for Weaviate to be up and running
    kubectl wait sts/weaviate -n weaviate --for jsonpath='{.status.readyReplicas}'=3 --timeout=300s

    # Create a clusterIP service to expose weaviate's metrics on port 2112
    kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  labels:
    app: weaviate-metrics
  name: weaviate-metrics
  namespace: weaviate
spec:
  ports:
  - port: 2112
    protocol: TCP
    targetPort: 2112
  selector:
    app: weaviate
  type: ClusterIP
EOF

    echo "benchmark_tenant #  Installing prometheus..."
    # Install Prometheus using helm/prometheus.yaml values
    helm upgrade --install prometheus -f helm/prometheus.yaml prometheus-community/prometheus
    # Wait for Prometheus to be up
    kubectl wait deploy/prometheus-server --for=condition=available --timeout=300s

    echo "benchmark_tenant #  Exposing ports locally..."
    mkdir -p logs
    # Expose weaviate port locally on port defined in WEAVIATE_PORT
    nohup kubectl port-forward svc/weaviate $WEAVIATE_PORT:80 -n weaviate > logs/weaviate-port-forward.log 2>&1 &
    # Expose prometheus port in port defined in PROMETHEUS_PORT
    nohup kubectl port-forward svc/prometheus-server $PROMETHEUS_PORT:80 > logs/prometheus-port-forward.log 2>&1 &
    echo "##############################################"
    echo "# Monitor metrics in your browser by opening:#"
    echo "#        http://localhost:${PROMETHEUS_PORT}               #"
    echo "##############################################"
}

# Function to be invoked for "run" option
run_benchmark() {

    local number=${1:-""}
    echo "Running benchmark..."

    ${GO_BIN_PATH}/go run . --number ${number}

    # Store weaviate logs for debugging under logs directory
    timestamp=$(date +"%Y%m%d_%H%M%S")
    sudo kubectl logs sts/weaviate -n weaviate > logs/weaviate-${number}-${timestamp}.log
}

# Function to be invoked for "cleanup" option
cleanup_benchmark() {
    echo "Cleaning up benchmark..."
    
    #Kill kubectl port-forward processes running in the background
    pkill -f "kubectl port-forward"

    # Uninstall prometheus
    helm uninstall prometheus
    helm repo remove prometheus-community

    # Uninstall weaviate and delete the namespace
    helm uninstall weaviate --namespace weaviate
    kubectl delete namespace weaviate

    # Remove the cluster
    kind delete cluster --name benchmark-tenant
}

# Function to handle interactive mode
interactive_mode() {
    while true; do
        echo "Please, select one of the options:"
        echo "1. prepare"
        echo "2. run (default:1000 tenants)"
        echo "3. cleanup"
        echo "4. exit"

        read -p "Enter option number: " user_option

        case $user_option in
            1)
                prepare_benchmark
                ;;
            2)
                read -p "Enter the number of tenants (press Enter for default): " num_tenants
                if [[ -n $num_tenants && $num_tenants =~ ^[0-9]+$ ]]; then
                    run_benchmark "$num_tenants"
                else
                    run_benchmark
                fi
                ;;
            3)
                cleanup_benchmark
                ;;
            4)
                echo "Exiting interactive mode..."
                break
                ;;
            *)
                echo "Invalid option. Please try again."
                ;;
        esac
    done
}

# Main script

# Check if any options are passed
if [ $# -eq 0 ]; then
    echo "Usage: $0 <options>"
    echo "options:"
    echo "         prepare"
    echo "         run <number_of_tenants>"
    echo "         cleanup"
    echo "         interactive"
    exit 1
fi

# Check if required commands are available
if ! command -v kind &> /dev/null; then
    echo "Please, install the requirement for the benchmark: kind."
    exit 1
fi
if ! command -v helm &> /dev/null; then
    echo "Please, install the requirement for the benchmark: helm."
    exit 1
fi
if ! command -v kubectl &> /dev/null; then
    echo "Please, install the requirement for the benchmark: kubectl."
    exit 1
fi
if ! command -v ${GO_BIN_PATH}/go &> /dev/null; then
    echo "Please, install the requirement for the benchmark: go."
    echo "NOTE: If you are sure go binaries are installed, replace the GO_BIN_PATH variable bin the binaries directory."
    exit 1
fi

# Check if script is invoked with root permissions
if [ "$(id -u)" -ne 0 ]; then
    echo "This script requires root permissions. Please run with sudo."
    exit 1
fi

# Process command line options
case $1 in
    "prepare")
        prepare_benchmark
        ;;
    "run")
        if [[ -n $2 && $2 =~ ^[0-9]+$ ]]; then
            run_benchmark $2
        else
           run_benchmark
        fi
        ;;
    "cleanup")
        cleanup_benchmark
        ;;
    "interactive")
        interactive_mode
        ;;
    *)
        echo "Invalid option: $1. Use 'prepare', 'run', 'cleanup', or 'interactive'."
        exit 1
        ;;
esac

