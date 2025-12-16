## Eliminar todo menos volumenes
~~~
docker system prune -a -f
~~~

## Consultas a la DB
~~~
docker exec -it mariadb-noa mariadb -u root -p
~~~

server

saltaget@gmail.com

SgProg12.

ssh -p5756 root@66.97.39.124

SaltaGet!NoaGestion25

services:
  db:
    image: mariadb:11
    container_name: mariadb-noa
    ports:
      - "3308:3306"
    env_file:
      - .env
    environment:
      - MYSQL_ROOT_PASSWORD=${DB_PASSWORD}
      - MYSQL_DATABASE=${DB_USER}
    volumes:
      - db-data:/var/lib/mysql
      - ./conf/mariadb/mariadb.cnf:/etc/mysql/mariadb.conf.d/99-custom.cnf:ro
    healthcheck:
      test: [ "CMD-SHELL", "mysqladmin ping -h localhost -u${DB_USER} -p${DB_PASSWORD} || exit 1" ]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s
    restart: unless-stopped
    networks:
      - monitoring

  api:
    image: danielmchachagua/noagestionback:0.3.0
    container_name: noa-gestion
    ports:
      - "3000:3000"
    env_file:
      - .env
    depends_on:
      - db
    restart: unless-stopped
    volumes:
      - ./backups:/app/backups
    networks:
      - default
      - monitoring
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
        - "9090:9090"
    volumes:
        - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
        - '--config.file=/etc/prometheus/prometheus.yml'
    networks:
        - monitoring

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    environment:
        - GF_SERVER_ROOT_URL=https://noa-gestion.saltaget.com/grafana/
        - GF_SERVER_SERVE_FROM_SUB_PATH=true
    ports:
        - "3030:3000"
    volumes:
        - grafana-storage:/var/lib/grafana
    depends_on:
        - prometheus
    networks:
        - monitoring

  loki:
    image: grafana/loki:latest
    container_name: loki
    ports:
        - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
    networks:
        - monitoring

  promtail:
     image: grafana/promtail:latest
     container_name: promtail
     user: root
     volumes:
        - /var/run/docker.sock:/var/run/docker.sock:ro
        - /var/lib/docker/containers:/var/lib/docker/containers:ro
        - ./promtail-config.yml:/etc/promtail/config.yml
     command: -config.file=/etc/promtail/config.yml
     depends_on:
        - loki
     networks:
        - monitoring
     restart: unless-stopped

  cadvisor:
     image: gcr.io/cadvisor/cadvisor:latest
     container_name: cadvisor
     ports:
        - "8080:8080"
     volumes:
        - /:/rootfs:ro
        - /var/run:/var/run:rw
        - /sys:/sys:ro
        - /var/lib/docker/:/var/lib/docker:ro
        - /sys/fs/cgroup:/sys/fs/cgroup:ro
        - /var/run/docker.sock:/var/run/docker.sock:ro
     restart: unless-stopped
     networks:
        - monitoring

  redis:
      image: redis:7-alpine
      container_name: noa_redis
      restart: unless-stopped
      ports:
         - "6379:6379"
      volumes:
         - redis_data:/data
      networks:
         - monitoring
      command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
      healthcheck:
         test: [ "CMD", "redis-cli", "ping" ]
         interval: 10s
         timeout: 3s
         retries: 5

  redis-commander:
      image: rediscommander/redis-commander:latest
      container_name: noa_redis_ui
      restart: unless-stopped
      environment:
         - REDIS_HOSTS=local:redis:6379
      ports:
         - "8081:8081"
      networks:
         - monitoring
      depends_on:
         - redis

networks:
  monitoring:


volumes:
  db-data:
  grafana-storage:
  redis_data:
    driver: local




container_memory_usage_bytes{id="/system.slice/docker-b419fda228ba25cd56e6df44a4947d6041960ca47455a703d3a13fb2d8070869.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-a4cfa9d4fa59c540cea451a8dbb3ece0c8ef41ffac9046422eec46627009c987.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-1591309e3a4cb88f0fa08f8a29671a3ccc36daebf905772c603337f57a93ba7e.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-318851d2c4c5c0a18110f5b436640b8aa5a404a1d37d67ae0794f6577bd96a4d.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-5677f4ca667b21648512e95911c0b6a83fa360c0d81625ba1c6a5a87804583a7.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-490f6d3ac492a0a00c2cdaed60c9d5f0de923d9f3fc219bfa115785cc7f6ca0d.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-20d779c3453a1cc5f122125aea088c875ce3b21f7b8b44fde1b8ba06ffcc34d6.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-8f1fa6be76baa544cfdd6867123792f22af330853f066b89e1bfe027a383afc5.scope", instance="cadvisor:8080", job="cadvisor"}
container_memory_usage_bytes{id="/system.slice/docker-2cf38438df6c5c53cdf9f0ac7633a41d005f57772e1cdffbc99d1884274acdf5.scope", instance="cadvisor:8080", job="cadvisor"}
