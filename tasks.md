### terminar nota de credito

### actualizar income sale con campo tiene factura

### guardar la factura en el modelo

### agregar columna de nuemro de punto de venta *




server {
  listen 80;
  server_name www.test.noagestion.com.ar test.noagestion.com.ar;

	client_max_body_size 50M;  # 10MB - ajusta a 50M, 100M, etc.

  location /api/ {
      proxy_pass http://localhost:3000;
      proxy_http_version 1.1;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 60s;
  }

  location /ecommerce {
      proxy_pass http://localhost:3001;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection 'upgrade';
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_read_timeout 60s;
      proxy_connect_timeout 60s;
      proxy_send_timeout 60s;

  }

  location /grafana/ {
    #rewrite ^/grafana/(.*)$ /$1 break;
    proxy_pass http://localhost:3030;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  location / {
    proxy_pass http://localhost:5000;
        proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";		
        proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
  }
}



