# Лабораторна робота №3.

### Тема: Асинхронний зв'язок між сервісами.

### Завдання:

1. Підключити базу даних в додаток.
2. Налаштувати *message broker*
3. Розширити додаток, додати асинхронну комунікацію між сервісами (мінімум 1 продюсер 2 консюмера).

### 1.1 *Persistent storage*

> Для команд, які працюють з хмарними провайдерами можна використовувати зовнішні бази даних, за межами кластера *k8s*.
> Для команд, які працюють з *minikube* потрібно налаштувати БД як частину кластера.

*Pod* в *k8s* не мають стану, це означає, що при перезавантаженні чи перестворенні *Pod* весь поточний стан буде видалено.
Для більшості мікро-сервісів, такий підхід доречний. Але при роботі зокрема, з базою даних необхідно зберігати стан. 
Для цього потрібно підключити диск в *Pod*.

Для роботи з диском в кластері *k8s* є кілька сутностей.  
**PersistentVolume** - сутність, яка зв'язує зовнішній диск з кластером.
Як правило явно не задається, а створюється за допомогою *StorageClass*.  
**PersistentVolumeClaim** - сутність, яка зв'язує *Pod* диск з *PersistentVolume*.  
**StorageClass** - конфігурація, що визначає правила для динамічного створення *PersistentVolume*.  

При роботі з *minikube* *StorageClass* створюється автоматично при створенні кластера.
Для хмарних провайдерів *StorageClass* потрібно визначати явно, наприклад для [gke](https://cloud.google.com/kubernetes-engine/docs/how-to/persistent-volumes/ssd-pd#ssd_persistent_disks)

Проста конфігурація *PersistentVolumeClaim* має містити об'єм та [accessModes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes)

```shell
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: demo-storage
spec:
  accessModes:
    - ReadWriteOnce  
  resources:
    requests:
      storage: 5Gi
```

### 1.2 База даних

З базою даних в *k8s* працювати так само як і з іншими сервісами, за виключенням того, що база даних вимагаю диск, для збереження даних 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:10.4
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
          env:
          - name: POSTGRES_DB
            value: demo
          - name: POSTGRES_USER
            value: demo
          - name: POSTGRES_PASSWORD # так працювати з конфігураціями в прод середовищі не варто, це буде розглянуто в наступних лабораторних
            value: demo
          volumeMounts: # прив'язка контейнеру, до *volume*
            - mountPath: /var/lib/postgresql/data
              name: postgredb 
      volumes: # визначення *volume*, прив'язка до PersistentVolumeClaim
        - name: postgredb
          persistentVolumeClaim:
            claimName: demo-storage
```

### 2. Налаштувати *message broker*

В прикладах буде використовуватись *Kafka*, у своїх проектах можна використовувати будь-який *message broker*.
Конфігурації для *kafka* знаходяться в директорії `k8s/kafka`. В цілому це стандартний *Deployment* та *Service*.

Після застосування налаштувань `kubectl apply -f k8s/kafka`. Перевіримо, що *кафка працює коректно*
Для цього можемо підключитись до терміналу в *Pod*.

```
kubectl exec -it kafka-86cf9f6968-dldlr -- /bin/bash
```

Для того, щоб надіслати повідомлення можна використати скрипт `kafka-console-producer.sh`. Так як скрипт запускається на тому ж хості, що і *kafka* адреса брокера буде `127.0.0.1:9092`

```shell
kafka-console-producer.sh --broker-list 127.0.0.1:9092 --topic test
```

Після цього зчитати повідомлення можна виконавши наступну команду

```shell
kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic test --from-beginning
```

Для перевірки які черги (*topic*) створено, можемо використати `kafka-topics.sh`

```shell
kafka-topics.sh --list --zookeeper zookeeper:2181
```

*kafka* зберігає всі повідомлення (на певний термін), якщо під час роботи повідомлення буде надіслано в *message broker*,
після чого відбудеться збій і *consumer* не встигне прочитати повідомлення, воно буде втрачено.
Якщо додати *Persistent storage* для *kafka*, всі повідомлення зберігатимуться а диск, тому після відмови системи,
вони будуть зчитані з диску і жодного повідомлення не буде втрачено

> *kafka* підтримує репліки, тобто, всі повідомлення будуть копіюватись на задану кількість *Pod*.
> При такому налаштуванні повідомлення можуть бути втрачені лише, якщо всі *Pod* впадуть.


### 3. Розширити додаток

Необхідно ускладнити логіку вашого додатку.  
Додати один сервіс (розширити існуючий), який буде надсилати повідомлення в *message broker*  
Додати 2 сервіси (розширити існуючі), які будуть читати з черги  

Варіанти асинхронних задач, які можуть підійти для різного роду додатків:  
1) нотифікації (браузер, пошта, смс)  
2) аудит (збереження історії змін сутностей)
