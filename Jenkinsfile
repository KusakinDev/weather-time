pipeline {
    agent any
    
    environment {
        DOCKER_REGISTRY = 'docker.io'
        DOCKER_USERNAME = 'debilyator'
        GO_API_IMAGE = "${DOCKER_USERNAME}/go-api"
        NEXTJS_IMAGE = "${DOCKER_USERNAME}/nextjs"
        KUBE_NAMESPACE = 'weather-app'
    }
    
    stages {
        stage('Checkout Code') {
            steps {
                checkout scm  // забирает код из GitHub
                sh 'git log -1 --oneline'  // показываем последний коммит
            }
        }
        stage('Docker Operations') {
            steps {
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'docker-hub-credentials', 
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh '''
                            # Логинимся один раз в начале
                            echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                            
                            # Все docker команды в одной сессии
                            cd back
                            docker build -t debilyator/go-api:${BUILD_NUMBER} .
                            docker push debilyator/go-api:${BUILD_NUMBER}
                            
                            cd ../front/FE_WeatherTime  
                            docker build -t debilyator/nextjs:${BUILD_NUMBER} .
                            docker push debilyator/nextjs:${BUILD_NUMBER}
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to K8s') {
            steps {
                script {
                    // Создаем namespace если не существует
                    sh """
                        kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
                    """
                    
                    // Обновляем образы в манифестах
                    sh """
                        sed -i 's|your-dockerhub-username|${DOCKER_USERNAME}|g' k8s/go-api/deployment.yaml
                        sed -i 's|your-dockerhub-username|${DOCKER_USERNAME}|g' k8s/nextjs/deployment.yaml
                        sed -i 's|latest|${env.BUILD_ID}|g' k8s/go-api/deployment.yaml
                        sed -i 's|latest|${env.BUILD_ID}|g' k8s/nextjs/deployment.yaml
                    """
                    
                    // Применяем все манифесты
                    sh """
                        kubectl apply -f k8s/ -n ${KUBE_NAMESPACE}
                    """
                    
                    // Ждем развертывания
                    sh """
                        kubectl rollout status deployment/go-api -n ${KUBE_NAMESPACE} --timeout=300s
                        kubectl rollout status deployment/nextjs -n ${KUBE_NAMESPACE} --timeout=300s
                        kubectl rollout status deployment/nginx -n ${KUBE_NAMESPACE} --timeout=300s
                    """
                    
                    // Показываем информацию о развертывании
                    sh """
                        echo "=== Deployment Status ==="
                        kubectl get all -n ${KUBE_NAMESPACE}
                        echo "=== Services ==="
                        kubectl get services -n ${KUBE_NAMESPACE}
                    """
                }
            }
        }
    }
    
    post {
        always {
            // Очистка
            sh 'docker system prune -f'
        }
        success {
            // Уведомление об успешном деплое
            sh """
                echo "🚀 Deployment successful!"
                echo "📊 Check application: minikube service nginx-service -n ${KUBE_NAMESPACE} --url"
            """
        }
        failure {
            // Уведомление о неудаче
            sh 'echo "❌ Deployment failed!"'
        }
    }
}
