pipeline {
    agent any
    
    environment {
        DOCKER_USERNAME = 'debilyator'
        KUBE_NAMESPACE = 'weather-app'
    }
    
    stages {
        stage('Setup Kubernetes Access') {
            steps {
                sh '''
                    echo "=== Fixing Kubernetes access ==="
                    
                    # Обновляем kubeconfig
                    minikube update-context || true
                    
                    # Проверяем доступ
                    kubectl cluster-info || echo "Cluster info failed, continuing..."
                '''
            }
        }
        
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }
        
        stage('Deploy to Kubernetes') {
            steps {
                sh '''
                    minikube image build -t app/go-api:latest ./back
                    minikube image build -t app/nextjs:latest ./front/FE_WeatherTime
                    
                    # применить манифесты
                    minikube kubectl -- apply -f k8s/00-ns.yaml --validate=false --insecure-skip-tls-verify=true
                    minikube kubectl -- apply -f k8s/10-backend.yaml --validate=false --insecure-skip-tls-verify=true
                    minikube kubectl -- apply -f k8s/20-frontend.yaml --validate=false --insecure-skip-tls-verify=true
                    minikube kubectl -- apply -f k8s/30-ingress.yaml --validate=false --insecure-skip-tls-verify=true
                    
                    # переключить deployments на тег latest (если в yaml другой)
                    minikube kubectl -- -n app set image deploy/nextjs nextjs=app/nextjs:latest --validate=false --insecure-skip-tls-verify=true
                    minikube kubectl -- -n app set image deploy/go-api go-api=app/go-api:latest --validate=false --insecure-skip-tls-verify=true
                '''
            }
        }
    }
}
