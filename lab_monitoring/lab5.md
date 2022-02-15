# Лабораторна робота №5.

### Тема: Моніторинг стану додатку

### Завдання:

1. Додати систему централізованого логування
2. Додати систему моніторингу
3. Додати метрики для моніторингу стану (мінімум 4 графіки)


### 2. Додати систему моніторингу

#### 1. Перш ніж додати систему моніторингу необхідно увімкнути `metrics-server` в `minikube`

```shell
minikube addons enable metrics-server
```

Після цього можна отримати поточний стан використання ресурсів (доступний через кілька секунд після активації): 

```shell
kubectl top node
kubectl top pod
```

#### 2. Для моніторингу потрібно встановити *Prometheus* та *Grafana*   

Спершу створимо *Namespace*, щоб не мішати компоненти моніторингу з основним додатком

```shell
kubectl create namespace monitoring
```

Та встановимо *Prometheus* та *Grafana* за допомогою *helm*

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

helm install --namespace monitoring prometheus prometheus-community/kube-prometheus-stack
```

#### 3. Відкрити доступ до *Grafana*

Можна за допомогою *Ingress*. Для цього пропишемо окремий хост для *Grafana*

```shell
sudo vim /etc/hosts
```
```shell
192.168.39.76   grafana // де 192.168.39.76 ip мінікуба
```

Поміняємо контекст, та встановимо *Ingress* в *Namespace* `monitoring`.
```shell
kubectl config set-context --current --namespace=monitoring

kubectl apply -f grafana_ingress.yaml
```

> На демо я замість зміни *Namespace* змінив *Context* `kubectl config use-context monitoring`.
> Контекст `monitoring` в мене не налаштований, тому були проблеми.
> В більшості є лише один контекст *minikube* і змінювати його не треба

Інший варіант це просто прокинути порт
```shell
kubectl port-forward --namespace monitoring service/prometheus-grafana 3000:80
```

Пароль до *Grafana* можна отримати виконавши:
```shell
kubectl get secret --namespace monitoring prometheus-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

### Додаткові посилання

- [Базові функції *Grafana*](https://grafana.com/docs/grafana/latest/datasources/prometheus/#query-variable)
- [Функції *Prometheus*](https://prometheus.io/docs/prometheus/latest/querying/functions)
- Основні метрики:
    - [cadvisor](https://github.com/google/cadvisor/blob/master/docs/storage/prometheus.md)
    - [pod-metrics](https://github.com/kubernetes/kube-state-metrics/blob/release-1.9/docs/pod-metrics.md) (version 1.9)
