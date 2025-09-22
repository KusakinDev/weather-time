pipeline {
    agent any // Запускать на любом доступном агенте (ноде)

    triggers {
        pollSCM('H/5 * * * *') // Опционально: периодическая проверка репозитория каждые 5 минут
    }

    stages {
        stage('Checkout') {
            steps {
                // Шаг 1: Получаем код из GitHub
                checkout scm
            }
        }

        stage('Test') {
            steps {
                // Шаг 2: (Опционально) Запускаем тесты
                sh 'echo "Running tests..."'
                // sh 'npm test' // для Node.js
                // sh 'pytest'   // для Python
            }
        }

        stage('Deploy to Production') {
            steps {
                // Шаг 3: Деплой на продакшн-сервер через SSH
                sshPublisher(
                    publishers: [
                        sshPublisherDesc(
                            configName: 'weather-time', // Имя сервера, настроенное в Jenkins
                            transfers: [
                                sshTransfer(
                                    // Команды, которые выполнятся на удаленном сервере
                                    execCommand: """
                                        cd /root/weather-time
                                        docker-compose down
                                        git pull origin master
                                        docker-compose up -d --build
                                    """
                                )
                            ],
                            usePromotionTimestamp: false,
                            useWorkspaceInPromotion: false,
                            verbose: true
                        )
                    ]
                )
            }
        }
    }

    post {
        success {
            // Действия при успешном завершении
            slackSend channel: '#deploys', message: "SUCCESS: Deployment of ${env.JOB_NAME} #${env.BUILD_NUMBER} is complete!"
        }
        failure {
            // Действия при ошибке
            slackSend channel: '#deploys', message: "FAILED: Deployment of ${env.JOB_NAME} #${env.BUILD_NUMBER} failed!"
        }
    }
}
