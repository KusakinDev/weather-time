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
        stage('Build Images') {
            parallel {
                stage('Build Go API') {
                    steps {
                        dir('back') {
                            script {
                                docker.build("${GO_API_IMAGE}:${env.BUILD_ID}")
                            }
                        }
                    }
                }
                stage('Build Next.js') {
                    steps {
                        dir('front/FE_WeatherTime') {
                            script {
                                docker.build("${NEXTJS_IMAGE}:${env.BUILD_ID}")
                            }
                        }
                    }
                }
            }
        }
        
        stage('Test Images') {
            parallel {
                stage('Test Go API') {
                    steps {
                        script {
                            docker.image("${GO_API_IMAGE}:${env.BUILD_ID}").inside {
                                sh 'echo "Running tests for Go API"'
                                // Добавьте ваши тесты здесь
                            }
                        }
                    }
                }
                stage('Test Next.js') {
                    steps {
                        script {
                            docker.image("${NEXTJS_IMAGE}:${env.BUILD_ID}").inside {
                                sh 'echo "Running tests for Next.js"'
                                // Добавьте ваши тесты здесь
                            }
                        }
                    }
                }
            }
        }
        
        stage('Push Images') {
            steps {
                script {
                    docker.withRegistry("https://${DOCKER_REGISTRY}", 'docker-hub-credentials') {
                        docker.image("${GO_API_IMAGE}:${env.BUILD_ID}").push()
                        docker.image("${NEXTJS_IMAGE}:${env.BUILD_ID}").push()
                        
                        // Также пушим как latest
                        docker.image("${GO_API_IMAGE}:${env.BUILD_ID}").push('latest')
                        docker.image("${NEXTJS_IMAGE}:${env.BUILD_ID}").push('latest')
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
