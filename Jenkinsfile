pipeline {
  agent any

environment {
  PROFILE = 'jenkins'
  MINIKUBE_HOME = '/var/lib/jenkins/.minikube'
  KUBECONFIG    = '/var/lib/jenkins/.kube/config'
}


  stages {
    stage('Who am I') {
  steps {
    sh '''#!/bin/bash
whoami
id
echo "HOME=$HOME"
echo "MINIKUBE_HOME=${MINIKUBE_HOME:-<empty>}"
echo "KUBECONFIG=${KUBECONFIG:-<empty>}"
'''
  }
}


    stage('Setup Kubernetes Access') {
      steps {
        sh '''#!/bin/bash
set -Eeuo pipefail

# Убедимся, что кластер поднят под пользователем jenkins
minikube -p "${PROFILE}" status || minikube -p "${PROFILE}" start --driver=docker

# Сформируем временный kubeconfig (сырой), чтобы не трогать глобальный
TMP_KUBECONFIG="$(mktemp)"
minikube -p "${PROFILE}" kubectl -- config view --raw > "$TMP_KUBECONFIG"

# Разберём адрес API и подготовим NO_PROXY (защитимся, если NO_PROXY не задан)
API="$(kubectl --kubeconfig="$TMP_KUBECONFIG" config view -o jsonpath='{.clusters[0].cluster.server}' || true)"
HOSTPORT="${API#https://}"
HOST="${HOSTPORT%%:*}"
NO_PROXY="${NO_PROXY:-}"
export NO_PROXY="${NO_PROXY},localhost,127.0.0.1,${HOST:-},${HOSTPORT:-},.svc,.cluster.local,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

# Проверка доступа (обнуляем прокси только на вызов kubectl)
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$TMP_KUBECONFIG" get nodes
'''
      }
    }

    stage('Checkout Code') {
      steps { checkout scm }
    }

    stage('Build images into Minikube') {
      steps {
        sh '''#!/bin/bash
set -Eeuo pipefail
minikube -p "${PROFILE}" image build -t app/go-api:latest ./back
minikube -p "${PROFILE}" image build -t app/nextjs:latest ./front/FE_WeatherTime
'''
      }
    }

    stage('Deploy to Kubernetes') {
      steps {
        sh '''#!/bin/bash
set -Eeuo pipefail

TMP_KUBECONFIG="$(mktemp)"
minikube -p "${PROFILE}" kubectl -- config view --raw > "$TMP_KUBECONFIG"

API="$(kubectl --kubeconfig="$TMP_KUBECONFIG" config view -o jsonpath='{.clusters[0].cluster.server}' || true)"
HOSTPORT="${API#https://}"
HOST="${HOSTPORT%%:*}"
NO_PROXY="${NO_PROXY:-}"
export NO_PROXY="${NO_PROXY},localhost,127.0.0.1,${HOST:-},${HOSTPORT:-},.svc,.cluster.local,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

# Применяем манифесты (везде без прокси)
for f in k8s/00-ns.yaml k8s/10-backend.yaml k8s/20-frontend.yaml k8s/30-ingress-web.yaml k8s/31-ingress-api.yaml; do
  HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$TMP_KUBECONFIG" apply -f "$f"
done

HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$TMP_KUBECONFIG" -n app set image deploy/nextjs nextjs=app/nextjs:latest
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$TMP_KUBECONFIG" -n app set image deploy/go-api  go-api=app/go-api:latest
'''
      }
    }
  }
}
