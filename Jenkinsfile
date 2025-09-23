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
        
        stage('Build and Push Docker Images') {
            steps {
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'docker-hub-credentials',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh '''
                            echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                            
                            cd back
                            docker build -t $DOCKER_USER/go-api:latest .
                            docker push $DOCKER_USER/go-api:latest
                            
                            cd ../front/FE_WeatherTime
                            docker build -t $DOCKER_USER/nextjs:latest .
                            docker push $DOCKER_USER/nextjs:latest
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to Kubernetes') {
            steps {
                sh '''
                    echo "=== Deploying with SSL fix ==="
                    
                    # Вариант 1: Отключаем валидацию SSL
                    kubectl create namespace $KUBE_NAMESPACE --dry-run=client -o yaml | kubectl apply -f - --validate=false
                    kubectl apply -f k8s/ -n $KUBE_NAMESPACE --validate=false
                    
                    # Обновляем образы
                    kubectl set image deployment/go-api go-api=$DOCKER_USERNAME/go-api:latest -n $KUBE_NAMESPACE --validate=false
                    kubectl set image deployment/nextjs nextjs=$DOCKER_USERNAME/nextjs:latest -n $KUBE_NAMESPACE --validate=false
                    
                    echo "=== Deployment status ==="
                    kubectl get all -n $KUBE_NAMESPACE --validate=false
                '''
            }
        }
    }
}
