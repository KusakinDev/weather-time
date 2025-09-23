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
                    # Создаем namespace если нет
                    kubectl create namespace $KUBE_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
                    
                    # Применяем манифесты
                    kubectl apply -f k8s/ -n $KUBE_NAMESPACE
                    
                    # Обновляем образы в деплойментах
                    kubectl set image deployment/go-api go-api=$DOCKER_USERNAME/go-api:latest -n $KUBE_NAMESPACE
                    kubectl set image deployment/nextjs nextjs=$DOCKER_USERNAME/nextjs:latest -n $KUBE_NAMESPACE
                    
                    # Ждем готовности подов
                    sleep 30
                    kubectl get pods -n $KUBE_NAMESPACE
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
