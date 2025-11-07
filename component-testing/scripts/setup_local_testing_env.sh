# ==================================================================================
#
#       Copyright (c) 2025 Samsung Electronics Co., Ltd. All Rights Reserved.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#          http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# ==================================================================================

echo -e "The following script will setup the local testing environment which includes: \nCreating kind Cluster, deploying modelmgmtservice component & mocks"
# Prerequisties
# 1. Docker

CLUSTER_NAME="testing"

echo "Step-1: Installing & Creating Kind Cluster"
# For AMD64 / x86_64
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-amd64
# For ARM64
[ $(uname -m) = aarch64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.30.0/kind-linux-arm64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

kind create cluster --config kind.yaml --name ${CLUSTER_NAME}
kubectl wait --for=condition=Ready nodes --all --timeout=120s --context kind-${CLUSTER_NAME}

echo "Step-2: Building Components-Image"
docker build -t modelmgmtservice:latest ../../.
kind load docker-image modelmgmtservice:latest --name ${CLUSTER_NAME}

echo "Step-3: Deploying Patches, Hacks and mocks"
kubectl create namespace traininghost
echo "Installing PostgreSQL (tm-db) in namespace traininghost..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

helm install tm-db bitnami/postgresql \
--set image.repository="bitnamilegacy/postgresql" \
--set image.tag="17.6.0" \
--set global.security.allowInsecureImages=true \
--set auth.postgresPassword=postgres \
--set primary.persistence.enabled=false \
--namespace traininghost

echo "Waiting for PostgreSQL pod to be ready..."
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=postgresql -n traininghost --timeout=300s

echo "PostgreSQL installation complete."
kubectl get pods -n traininghost -l app.kubernetes.io/name=postgresql

echo "Creating mock secret: leofs-secret"
          kubectl create secret generic leofs-secret \
            --from-literal=password="dummy-s3-secret-key" \
            -n traininghost --dry-run=client -o yaml | kubectl apply -f -


echo "Step-4: Setup Helm & Deploy dataextraction component"
PREV_WORK_DIR=$(pwd)
cd /tmp/
git clone https://github.com/o-ran-sc/aiml-fw-aimlfw-dep.git
cd aiml-fw-aimlfw-dep
./bin/install_common_templates_to_helm.sh

helm dep up helm/modelmgmtservice
helm install mme helm/modelmgmtservice -f RECIPE_EXAMPLE/example_recipe_local_images_oran_latest.yaml --kube-context kind-${CLUSTER_NAME}

cd $PREV_WORK_DIR