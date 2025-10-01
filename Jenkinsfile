pipeline {
  agent any

  environment {
    PROFILE = 'minikube'
  }

  stages {

    stage('Setup Kubernetes Access') {
      steps {
        sh '''#!/bin/bash
set -euxo pipefail

# 0) Проверим, что minikube запущен под этим пользователем
minikube -p "${PROFILE}" status

# 1) Получаем kubeconfig в temp-файл
export KUBECONFIG="$(mktemp)"
minikube -p "${PROFILE}" kubectl -- config view --raw > "$KUBECONFIG"

# 2) Узнаём адрес API сервера и настраиваем NO_PROXY (чтобы kubectl не шел через Jenkins proxy)
API="$(kubectl --kubeconfig="$KUBECONFIG" config view -o jsonpath='{.clusters[0].cluster.server}')"
HOSTPORT="${API#https://}"
HOST="${HOSTPORT%%:*}"
export NO_PROXY="${NO_PROXY},localhost,127.0.0.1,${HOST},${HOSTPORT},.svc,.cluster.local,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

# 3) Базовая проверка доступа (без прокси)
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" cluster-info || true
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" get nodes
'''
      }
    }

    stage('Checkout Code') {
      steps { checkout scm }
    }

    stage('Build images into Minikube') {
      steps {
        sh '''#!/bin/bash
set -euxo pipefail
# Строим образы прямо внутрь Minikube
minikube -p "${PROFILE}" image build -t app/go-api:latest ./back
minikube -p "${PROFILE}" image build -t app/nextjs:latest ./front/FE_WeatherTime
'''
      }
    }

    stage('Deploy to Kubernetes') {
      steps {
        sh '''#!/bin/bash
set -euxo pipefail

# kubeconfig ещё раз (новый temp на каждый шаг — безопаснее)
export KUBECONFIG="$(mktemp)"
minikube -p "${PROFILE}" kubectl -- config view --raw > "$KUBECONFIG"

API="$(kubectl --kubeconfig="$KUBECONFIG" config view -o jsonpath='{.clusters[0].cluster.server}')"
HOSTPORT="${API#https://}"
HOST="${HOSTPORT%%:*}"
export NO_PROXY="${NO_PROXY},localhost,127.0.0.1,${HOST},${HOSTPORT},.svc,.cluster.local,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

# Применяем манифесты (везде явно передаём kubeconfig и отключаем прокси)
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/00-ns.yaml
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/10-backend.yaml
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/20-frontend.yaml
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app apply -f k8s/30-ingress-web.yaml
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app apply -f k8s/31-ingress-api.yaml
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app get ingress -o wide

# Обновляем образы и ждём rollout
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app set image deploy/nextjs nextjs=app/nextjs:latest
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app set image deploy/go-api  go-api=app/go-api:latest

HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app rollout status deploy/nextjs
HTTP_PROXY= HTTPS_PROXY= http_proxy= https_proxy= kubectl --kubeconfig="$KUBECONFIG" -n app rollout status deploy/go-api
'''
      }
    }
  }
}
