pipeline {
  agent any

  environment {
    PROFILE = 'minikube'
  }

  stages {
    stage('Setup Kubernetes Access') {
      steps {
        sh '''

          # 2) Получим kubeconfig для этого профиля в отдельный временный файл
          export KUBECONFIG=$(mktemp)
          minikube -p minikube kubectl -- config view --raw > "$KUBECONFIG"

          # 3) Базовая проверка доступа
          kubectl --kubeconfig="$KUBECONFIG" cluster-info || true
          kubectl --kubeconfig="$KUBECONFIG" get nodes
        '''
      }
    }

    stage('Checkout Code') {
      steps { checkout scm }
    }

    stage('Build images into Minikube') {
      steps {
        sh '''
          set -euxo pipefail
          # Строим образы прямо внутрь Minikube
          minikube -p ${PROFILE} image build -t app/go-api:latest ./back
          minikube -p ${PROFILE} image build -t app/nextjs:latest ./front/FE_WeatherTime
        '''
      }
    }

    stage('Deploy to Kubernetes') {
      steps {
        sh '''
          set -euxo pipefail
          export KUBECONFIG=$(mktemp)
          minikube -p ${PROFILE} kubeconfig > "$KUBECONFIG"

          # Применяем манифесты обычным kubectl с явным kubeconfig
          kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/00-ns.yaml
          kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/10-backend.yaml
          kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/20-frontend.yaml
          kubectl --kubeconfig="$KUBECONFIG" apply -f k8s/30-ingress.yaml

          # Обновляем образы в деплойментах
          kubectl --kubeconfig="$KUBECONFIG" -n app set image deploy/nextjs nextjs=app/nextjs:latest
          kubectl --kubeconfig="$KUBECONFIG" -n app set image deploy/go-api  go-api=app/go-api:latest

          # (опц.) дождаться роллаута
          kubectl --kubeconfig="$KUBECONFIG" -n app rollout status deploy/nextjs
          kubectl --kubeconfig="$KUBECONFIG" -n app rollout status deploy/go-api
        '''
      }
    }
  }
}
