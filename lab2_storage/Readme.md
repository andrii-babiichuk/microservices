# Лабораторна робота №2.

### Тема: робота з Базою даних

### Завдання:

1. Підключити базу даних в додаток.
2. Написати міграцію.
3. Додати автоматичний запуск міграцій.
4. Розширити сервіс, додати роботу з базою.

### 1.1 *Persistent storage*

> Для команд, які працюють з хмарними провайдерами можна використовувати зовнішні бази даних, за межами кластера *k8s*.
> Для команд, які працюють з *minikube* потрібно налаштувати БД як частину кластера.

*Pod* в *k8s* не мають стану, це означає, що при перезавантаженні чи перестворенні *Pod* весь поточний стан буде видалено.
Для більшості мікро-сервісів, такий підхід доречний. Але при роботі, зокрема, з базою даних необхідно зберігати стан. 
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

З базою даних в *k8s* працювати так само як і з іншими сервісами, за виключенням того, що база даних має мати диск, для збереження даних 

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
          image: postgres:14.1-alpine3.15
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
          volumeMounts:
            - mountPath: /var/lib/postgresql/data
              name: postgredb
      volumes:
        - name: postgredb
          persistentVolumeClaim:
            claimName: postgres-volume-claim
```


### 1.3 Конфігурація

Для коректної роботи `postgres` необхідно задати деякі змінні середовища.

> Які саме значення потрібно визначити можна подивитись в [офіційному образі Докер](https://hub.docker.com/_/postgres) в розділі *Environment Variables*

Для початку визначимо `POSTGRES_DB`.

Найпростіший спосіб це зробити - описати властивість контейнера `env`:

```yaml
...
      containers:
        - name: postgres
          image: postgres:14.1-alpine3.15
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
          env:
            POSTGRES_DB: "demo"
...
```

В такому підході є недолік. `POSTGRES_DB` може відрізнятись для різних середовищ. Для того, щоб розділити
статичні конфігурації контейнера від конфігурацій специфічних до середовища потрібно створити окремий об'єкт *ConfigMap*.
*ConfigMap* являє собою просту мапу ключ:значення.

`k8s/postgres/config-map.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-config
data:
  POSTGRES_DB: demo
```

Після цього потрібно прив'язати значення *ConfigMap* до контейнера:

1) Прив'язка до конкретного значення:
```yaml
...
      containers:
        - name: postgres
          image: postgres:14.1-alpine3.15
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
           env:
            - name: POSTGRES_DB
              valueFrom:
                configMapKeyRef:
                  name: postgres-config
                  key: POSTGRES_DB
...
```

2) Прив'язка всієї *ConfigMap*
```yaml
...
      containers:
        - name: postgres
          image: postgres:14.1-alpine3.15
          imagePullPolicy: "IfNotPresent"
          ports:
            - containerPort: 5432
          envFrom:
            - configMapRef:
                name: postgres-config
...
```

> Є інші варіанти використання *ConfigMap*, детальніше на [офіційному сайті](https://kubernetes.io/docs/concepts/configuration/configmap/)

Тепер потрібно визначити:

`POSTGRES_USER` і `POSTGRES_PASSWORD`. Це чутливі дані та зберігати їх у відкритому вигляді не можна,
на жаль *Kubernetes* не має власного рішення для зберігання чутливих даних,
як саме працювати з чутливими даними буде розглянути пізніше, а поки можна використати *Secret*

*Secret* це майже теж саме, що і *ConfigMap*, але *Secret* зберігає дані закодовані в *base64*.
Крім цього для деяких даних є окремо визначені [типи *Secret*](https://kubernetes.io/docs/concepts/configuration/secret/#secret-types).
Для довільних даних використовується тип *Opaque*:

Таким чином маємо `2_storage/k8s/service1/secret.yaml`
```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: service1-secret
data:
  POSTGRES_USER: "ZGVtbw=="
  POSTGRES_PASSWORD: "ZGVtbw=="
```

Так само як і *ConfigMap* підключаємо його в *Pod*

```yaml
...
containers:
  - name: postgres
    image: postgres:14.1-alpine3.15
    imagePullPolicy: "IfNotPresent"
    ports:
      - containerPort: 5432
          envFrom:
            - configMapRef:
                name: postgres-config
            - secretRef:
                name: postgres-secret
...
```

> Зверніть увагу, що при зміні конфігурації - сервіс не отримає зміни, способи, як цього досягти буде розглянуто пізніше.
> Поки для того, щоб сервіс отримав оновлені конфігурації можна створювати нові *ConfigMap* (переіменовувати) 

Тепер виконавши команду `kubectl apply -f k8s/postgres` отримаємо робочу БД розгорнуту в кластері з визначеною базовою конфігурацією. 

### 2 Написати міграцію.

Для творення необхідної структури БД необхідно написати міграцію за допомогою будь-яких технологій, (для *go* можу порекомендувати [golang-migrate](https://github.com/golang-migrate/migrate), хоча її можна використовувати з будь-якою мовою програмування).  
Міграція обов'язково має мати зворотню міграцію (*rollback*), що поверне базу в стан до міграції.

Приклад елементарних міграйцій можна знайти в папці `services/service1/migrations`.

`20211128095755_spending_init.up.sql` В міграції створюється таблиця з 1 колонкою, `counter`, та вставляється значення `0`:

```postgresql
CREATE TABLE IF NOT EXISTS visits (counter INTEGER);
INSERT INTO visits VALUES (0);
```

`20211128095755_spending_init.down.sql` Зворотня міграція просто видаляє таблицю:
```postgresql
DROP TABLE visits;
```

> Зверніть увагу на ім'я файлів, вони починаются з `20211128095755` це таймстемп створення міграції, і використовується для версіонування міграцій.
> Рекомендую використовувати саме таймстемп як версію, а не простий інкремент. 

### 3. Додати автоматичний запуск міграцій.

Є два варіанти, як можна запускати міграцію в *kubernetes*:

1) Cтворити [*Kubernetes Job*](https://kubernetes.io/docs/concepts/workloads/controllers/job/) та запускати міграцію в рамках цієї *Job*.
2) Визначити [*initContainer*](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) в сервісі, який відповідає за міграцію.

Альтернативно можна використати фреймворк для запуску міграцій, наприклад [*liquibase*](https://liquibase.org/get-started/quickstart) (для *Java*) чи інші.

#### Міграція за допомогою *Kubernetes Job*

*Kubernetes Job* це, по суті, *Pod*, який завершує своє виконання і після завершення очищує ресурси.
Такі міграції незалежні від сервісу і можуть запускатись в будь-який час, як правило, як окремі кроки *CI/CD*.

#### Міграція за допомогою *initContainer*
*initContainer* - це додаткова конфігурація сервісу, що запускає 1 або кілька контейнерів до запуску основного.
Як і в *Kubernetes Job* контейнери описані в *initContainer* мають завершити своє виконання, і лише після успішного завершення може запуститись основний контейнер сервісу.
Цей підхід зв'язує сервіс і міграцію, таким чином, що кожен сервіс має подбати про те, щоб привести БД в стан, який йому необхідний.
Основний недолік цього підходу, що *initContainer* присутній в кожному екземплярі сервісу і відповідно при запуску кількох екземплярів сервісу,
скрипт міграції буде запускатись разом з кожним екземпляром сервісу, тому треба слідкувати,
щоб БД була заблокована на запис під час міграції (але це треба робити в будь-якому випадку, для надійної роботи тому це не проблема).

> Принципова різниця в підходах у тому, що в першому випадку міграцію можна запускати будь-коли, і зв'язок з сервісом неявний. 
> На мою думку, підхід з *Kubernetes Job* кращий, коли є 1 БД для всіх сервісів, або коли відбувається перехід з монолітної архітектури на мікросервісну (зв'язок БД і сервісу не однозначний). 
> Тоді в 1 *Kubernetes Job* можна запускати міграції одразу для кількох або всіх сервісів.
> Підхід з *initContainer* більш доречний для роботи з мікросервісами, так як змушує одразу розділяти базу між сервісами,
> і краще розділяє відповідальності між сервісами.

В прикладі розглянемо другий варіант.

Спершу потрібно додати *Dokerfile* з інструментом, що запускатиме міграцію, та міститеме файли міграції

```Dockerfile
FROM migrate/migrate:v4.15.1
COPY services/service1/migrations migrations
```

> CLI для `golang-migrate` може працювати з *git*, як хранилищем міграцій, тобто не обов'язково створювати `Dockerfile`.
> Можна використати `migrate/migrate:v4.15.1` та передати шлях до директорії через git.

Тепер треба створити образ міграцій 

`docker build -t service1-migrations:0.1 -f services/service1/migrations/Dockerfile .`

І додати *initContainers* в deployment

```yaml
      initContainers:
        - name: run-migrations
          image: service1-migrations:0.1
          command: ["migrate", "-path", "/migrations", "-database",  "$(POSTGRESQL_URL)", "goto", "$(VERSION)"] # CLI команда запуску міграцій, де $(POSTGRESQL_URL) і $(VERSION) змінні середовища
          env:
            - name: POSTGRESQL_URL
              value: "postgres://demo:demo@postgres:5432/demo?sslmode=disable"
            - name: VERSION
              value: "20211128095755"
```

Тепер запустивши команду `kubectl apply -f k8s/service1` міграція виконається.

### 4. Розширити сервіс, додати роботу з базою.

Реалізувати, як мінімум *CRUD* для сутностей, з якими працює сервіс.
Елементарний приклад роботи з БД для *golang* можна знайти в першому сервісі `services/service1/main.go`
