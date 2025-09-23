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
                checkout scm
                sh 'git log -1 --oneline'
            }
        }
        
        stage('Cleanup Duplicate Manifests') {
            steps {
                sh """
                    echo "=== Cleaning up duplicate manifests ==="
                    
                    # Удаляем дубликаты - оставляем только структурированные файлы в папках
                    if [ -f "k8s/go-api-deployment.yaml" ]; then
                        echo "🗑️  Removing duplicate: k8s/go-api-deployment.yaml"
                        rm -f k8s/go-api-deployment.yaml
                    fi
                    
                    if [ -f "k8s/nextjs-deployment.yaml" ]; then
                        echo "🗑️  Removing duplicate: k8s/nextjs-deployment.yaml"
                        rm -f k8s/nextjs-deployment.yaml
                    fi
                    
                    # Проверяем service.yaml файлы
                    if [ -f "k8s/go-api/service.yaml" ]; then
                        echo "🔍 Checking k8s/go-api/service.yaml"
                        cat k8s/go-api/service.yaml
                        if [ ! -s "k8s/go-api/service.yaml" ] || grep -q "no objects passed to apply" k8s/go-api/service.yaml; then
                            echo "🗑️  Removing empty service.yaml"
                            rm -f k8s/go-api/service.yaml
                        fi
                    fi
                    
                    echo "=== Final file structure ==="
                    find k8s/ -type f | sort
                """
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
                            echo "=== Docker Login ==="
                            echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                            
                            echo "=== Building Go API ==="
                            cd back
                            docker build -t $DOCKER_USER/go-api:${BUILD_NUMBER} .
                            docker push $DOCKER_USER/go-api:${BUILD_NUMBER}
                            
                            echo "=== Building Next.js ==="
                            cd ../front/FE_WeatherTime  
                            docker build -t $DOCKER_USER/nextjs:${BUILD_NUMBER} .
                            docker push $DOCKER_USER/nextjs:${BUILD_NUMBER}
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to K8s') {
            steps {
                script {
                    // Создаем namespace
                    sh """
                        echo "=== Creating namespace ==="
                        kubectl create namespace ${KUBE_NAMESPACE} --dry-run=client -o yaml | kubectl apply -f -
                    """
                    
                    // Применяем только нужные манифесты
                    sh """
                        echo "=== Applying Kubernetes manifests ==="
                        
                        # Применяем namespace отдельно
                        kubectl apply -f k8s/namespace.yaml
                        
                        # Применяем манифесты из папок (игнорируем корневые дубликаты)
                        for dir in k8s/go-api k8s/nextjs k8s/nginx; do
                            if [ -d "\$dir" ]; then
                                echo "📁 Applying manifests from: \$dir"
                                for file in \$dir/*.yaml \$dir/*.yml; do
                                    if [ -f "\$file" ] && [ -s "\$file" ]; then
                                        echo "📄 Applying: \$file"
                                        kubectl apply -f "\$file" -n ${KUBE_NAMESPACE} --validate=false
                                        if [ \$? -eq 0 ]; then
                                            echo "✅ Success: \$file"
                                        else
                                            echo "❌ Failed: \$file"
                                        fi
                                    fi
                                done
                            fi
                        done
                        
                        # Применяем ingress если есть
                        if [ -f "k8s/ingress.yaml" ] && [ -s "k8s/ingress.yaml" ]; then
                            echo "📄 Applying: k8s/ingress.yaml"
                            kubectl apply -f k8s/ingress.yaml -n ${KUBE_NAMESPACE}
                        fi
                    """
                    
                    // Ждем создания деплойментов
                    sh """
                        echo "=== Waiting for deployments ==="
                        sleep 10
                        
                        echo "=== Current deployments ==="
                        kubectl get deployments -n ${KUBE_NAMESPACE}
                        
                        # Ждем появления деплойментов
                        for i in {1..30}; do
                            if kubectl get deployment go-api -n ${KUBE_NAMESPACE} >/dev/null 2>&1 && \\
                               kubectl get deployment nextjs -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                                echo "✅ All deployments found"
                                break
                            fi
                            echo "⏳ Waiting for deployments... (\$i/30)"
                            sleep 2
                        done
                    """
                    
                    // Обновляем образы
                    sh """
                        echo "=== Updating images ==="
                        
                        if kubectl get deployment go-api -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            echo "🔄 Updating go-api image"
                            kubectl set image deployment/go-api go-api=${GO_API_IMAGE}:${env.BUILD_NUMBER} -n ${KUBE_NAMESPACE}
                        else
                            echo "❌ go-api deployment not found"
                        fi
                        
                        if kubectl get deployment nextjs -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            echo "🔄 Updating nextjs image"
                            kubectl set image deployment/nextjs nextjs=${NEXTJS_IMAGE}:${env.BUILD_NUMBER} -n ${KUBE_NAMESPACE}
                        else
                            echo "❌ nextjs deployment not found"
                        fi
                    """
                    
                    // Ждем rollout
                    sh """
                        echo "=== Waiting for rollout ==="
                        
                        if kubectl get deployment go-api -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            kubectl rollout status deployment/go-api -n ${KUBE_NAMESPACE} --timeout=300s
                        fi
                        
                        if kubectl get deployment nextjs -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            kubectl rollout status deployment/nextjs -n ${KUBE_NAMESPACE} --timeout=300s
                        fi
                    """
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                sh """
                    echo "=== Final status ==="
                    kubectl get all -n ${KUBE_NAMESPACE}
                    
                    echo "=== Pods details ==="
                    kubectl get pods -n ${KUBE_NAMESPACE} -o wide
                    
                    echo "=== Services ==="
                    kubectl get services -n ${KUBE_NAMESPACE}
                    
                    echo "=== Deployment status ==="
                    kubectl get deployments -n ${KUBE_NAMESPACE} -o wide
                """
            }
        }
    }
    
    post {
        always {
            sh 'docker system prune -f'
        }
        success {
            sh """
                echo "🚀 Deployment successful!"
                echo "📊 Application deployed to namespace: ${KUBE_NAMESPACE}"
            """
        }
        failure {
            sh """
                echo "❌ Deployment failed!"
                echo "🔍 Check the logs above for details"
            """
        }
    }
}
