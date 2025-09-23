pipeline {
    agent any
    
    environment {
        DOCKER_USERNAME = 'debilyator'  // ⚠️ ЗАМЕНИТЕ на свой!
        KUBE_NAMESPACE = 'weather-app'
    }
    
    stages {
        stage('Checkout Code') {
            steps {
                checkout scm  // автоматически забирает код из GitHub
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
                            # Логинимся в Docker Hub
                            echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                            
                            # Собираем и пушим образы
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
                    echo "=== Using minikube kubectl ==="
                    
                    # Все команды через minikube kubectl (обходит SSL проблемы)
                    minikube kubectl -- create namespace $KUBE_NAMESPACE --dry-run=client -o yaml | minikube kubectl -- apply -f -
                    minikube kubectl -- apply -f k8s/ -n $KUBE_NAMESPACE
                    minikube kubectl -- set image deployment/go-api go-api=$DOCKER_USERNAME/go-api:latest -n $KUBE_NAMESPACE
                    minikube kubectl -- set image deployment/nextjs nextjs=$DOCKER_USERNAME/nextjs:latest -n $KUBE_NAMESPACE
                    
                    # Проверяем
                    minikube kubectl -- get all -n $KUBE_NAMESPACE
                '''
            }
        }
    }
    
    post {
        always {
            sh 'docker system prune -f'  // очистка
        }
    }
}
