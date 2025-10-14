#!/bin/bash

# Получаем IP Minikube
MINIKUBE_IP=$(minikube ip)

# Получаем NodePort сервиса (замените 'your-service-name' на имя вашего сервиса)
NODE_PORT=$(kubectl get svc ingress-nginx-controller -o jsonpath='{.spec.ports[0].nodePort}')

# Создаем новый конфиг
cat > /etc/nginx/conf.d/k8s.conf << EOF
map \$http_upgrade \$connection_upgrade { default upgrade; "" close; }

upstream k8s_ingress_http {
    server ${MINIKUBE_IP}:${NODE_PORT};
}

server {
    listen 80;
    server_name _;

    client_max_body_size 50m;

    location / {
        proxy_pass http://k8s_ingress_http;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;

        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \$connection_upgrade;

        proxy_connect_timeout 5s;
        proxy_read_timeout 60s;
        proxy_send_timeout 60s;
    }

    location ^~ /_next/ {
        proxy_pass http://k8s_ingress_http;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection \$connection_upgrade;
        proxy_buffering off;
        proxy_read_timeout 3600s;
        proxy_send_timeout 3600s;
    }
}
EOF

# Перезагружаем nginx
nginx -t && systemctl reload nginx

echo "✅ Nginx config updated: ${MINIKUBE_IP}:${NODE_PORT}"
