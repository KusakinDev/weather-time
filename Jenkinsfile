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
        
        stage('Debug Structure') {
            steps {
                sh """
                    echo "=== Current Directory Structure ==="
                    pwd
                    ls -la
                    
                    echo "=== k8s Directory Contents ==="
                    find . -name "*.yaml" -o -name "*.yml" | sort || echo "No YAML files found"
                    
                    if [ -d "k8s" ]; then
                        echo "=== k8s Folder Details ==="
                        find k8s/ -type f | sort
                        echo "=== k8s File Contents ==="
                        find k8s/ -name "*.yaml" -exec echo "=== File: {} ===" \\; -exec cat {} \\;
                    else
                        echo "❌ k8s directory not found!"
                    fi
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
        
        stage('Prepare Manifests') {
            steps {
                script {
                    sh """
                        echo "=== Preparing Kubernetes Manifests ==="
                        
                        # Создаем k8s директорию если её нет
                        mkdir -p k8s
                        
                        # Проверяем существование манифестов
                        if [ ! -f "k8s/go-api-deployment.yaml" ]; then
                            echo "⚠️  go-api-deployment.yaml not found, creating basic one..."
                            cat > k8s/go-api-deployment.yaml << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-api
  namespace: weather-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-api
  template:
    metadata:
      labels:
        app: go-api
    spec:
      containers:
      - name: go-api
        image: ${DOCKER_USERNAME}/go-api:${env.BUILD_NUMBER}
        ports:
        - containerPort: 8000
EOF
                        fi
                        
                        if [ ! -f "k8s/nextjs-deployment.yaml" ]; then
                            echo "⚠️  nextjs-deployment.yaml not found, creating basic one..."
                            cat > k8s/nextjs-deployment.yaml << EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nextjs
  namespace: weather-app
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nextjs
  template:
    metadata:
      labels:
        app: nextjs
    spec:
      containers:
      - name: nextjs
        image: ${DOCKER_USERNAME}/nextjs:${env.BUILD_NUMBER}
        ports:
        - containerPort: 3000
EOF
                        fi
                        
                        echo "=== Final k8s structure ==="
                        find k8s/ -type f | sort
                    """
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
                    
                    // Применяем манифесты с подробным выводом
                    sh """
                        echo "=== Applying Kubernetes manifests ==="
                        
                        # Применяем каждый файл отдельно с проверкой
                        for file in k8s/*.yaml k8s/*.yml k8s/*/*.yaml k8s/*/*.yml; do
                            if [ -f "\$file" ]; then
                                echo "📄 Applying: \$file"
                                kubectl apply -f "\$file" --validate=false
                                if [ \$? -eq 0 ]; then
                                    echo "✅ Success: \$file"
                                else
                                    echo "❌ Failed: \$file"
                                fi
                            fi
                        done
                        
                        echo "=== Waiting for deployments to be created ==="
                        # Ждем появления деплойментов
                        for i in {1..30}; do
                            if kubectl get deployment go-api -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                                echo "✅ go-api deployment found"
                                break
                            fi
                            echo "⏳ Waiting for go-api deployment... (\$i/30)"
                            sleep 2
                        done
                        
                        for i in {1..30}; do
                            if kubectl get deployment nextjs -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                                echo "✅ nextjs deployment found"
                                break
                            fi
                            echo "⏳ Waiting for nextjs deployment... (\$i/30)"
                            sleep 2
                        done
                    """
                    
                    // Проверяем что создалось
                    sh """
                        echo "=== Current Kubernetes state ==="
                        kubectl get all -n ${KUBE_NAMESPACE} || echo "No resources found in namespace"
                        
                        echo "=== Deployments list ==="
                        kubectl get deployments -n ${KUBE_NAMESPACE} || echo "No deployments found"
                        
                        echo "=== Pods list ==="
                        kubectl get pods -n ${KUBE_NAMESPACE} || echo "No pods found"
                    """
                    
                    // Пробуем обновить образы если деплойменты существуют
                    sh """
                        if kubectl get deployment go-api -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            echo "=== Updating go-api image ==="
                            kubectl set image deployment/go-api go-api=${GO_API_IMAGE}:${env.BUILD_NUMBER} -n ${KUBE_NAMESPACE}
                            kubectl rollout status deployment/go-api -n ${KUBE_NAMESPACE} --timeout=300s
                        else
                            echo "❌ go-api deployment not found, skipping rollout"
                        fi
                        
                        if kubectl get deployment nextjs -n ${KUBE_NAMESPACE} >/dev/null 2>&1; then
                            echo "=== Updating nextjs image ==="
                            kubectl set image deployment/nextjs nextjs=${NEXTJS_IMAGE}:${env.BUILD_NUMBER} -n ${KUBE_NAMESPACE}
                            kubectl rollout status deployment/nextjs -n ${KUBE_NAMESPACE} --timeout=300s
                        else
                            echo "❌ nextjs deployment not found, skipping rollout"
                        fi
                    """
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                sh """
                    echo "=== Final deployment status ==="
                    kubectl get all -n ${KUBE_NAMESPACE}
                    
                    echo "=== Pods details ==="
                    kubectl get pods -n ${KUBE_NAMESPACE} -o wide
                    
                    echo "=== Services ==="
                    kubectl get services -n ${KUBE_NAMESPACE}
                    
                    echo "=== Application URLs ==="
                    minikube service list -n ${KUBE_NAMESPACE} || echo "Minikube service command not available"
                """
            }
        }
    }
    
    post {
        always {
            sh '''
                echo "=== Cleaning up ==="
                docker system prune -f || true
            '''
            script {
                // Сохраняем логи деплоймента
                sh """
                    kubectl get events -n ${KUBE_NAMESPACE} --sort-by='.lastTimestamp' > deployment-events.log || true
                    kubectl describe namespace ${KUBE_NAMESPACE} > namespace-describe.log || true
                """
                archiveArtifacts artifacts: '*.log', fingerprint: true
            }
        }
        success {
            sh """
                echo "🚀 Deployment completed successfully!"
                echo "📊 Check application with: minikube service nginx-service -n ${KUBE_NAMESPACE} --url"
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
