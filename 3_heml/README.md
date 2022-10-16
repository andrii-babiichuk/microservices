# Лабораторна робота №3.

### Тема: робота з менеджером пакетів для *Kubernetes* - *Helm*

### Завдання:

1. Описати додаток за допомогою *Helm*.
   - Винести всі необхідні параметри в змінні.
   - Використати зовнішні залежності для бази даних та інших сторонніх систем.
   - Створити іменовані шаблони та використати їх для різних сутностей.
   - Винести шаблони в бібліотеку.

## 1 *Helm*

**Helm** - менеджер пакетів для *Kubernetes*.
Це набір інструментів, що автоматизує процес встановлення, конфігурацію, оновлення та видалення додатків *Kubernetes*.

Основні переваги використання *Helm*:

1) Робота зі сторонніми пакетами, наприклад для розгортання бази даних, меседж брокера, системи моніторингу та іншими.
2) Потужний інструмент шаблонізації, що дозволяє відділити змінювані параметри від статичного шаблону, використовувати елементарні конструкції, як цикли та умовні оператори, створювати універсальні шаблони, які можна перевикористовувати. 
3) Інструмент для управління життєвим циклом додатку, за допомогою *Helm* можна виконувати оновлення версій, відкат до конкретної версії, управління різними середовищами. 

### 1.1 Встановлення сторонніх пакетів

Для початку можна встановити один зі сторонніх додатків, наприклад `nginx`.
При роботі з *k8s*, без *Helm*, необхідно було б описати, як мінімум *Service* та *Deployment*.
При роботі з *Helm* потрібно: 
1) Додати репозиторій де міститься інформація про застосунок `helm repo add bitnami https://charts.bitnami.com`.

> Є різні репозиторії, що містять різні реалізації `nginx`. Їх можна пошукати на сайті [https://artifacthub.io](https://artifacthub.io) або виконавши команду `helm search hub nginx --list-repo-url`

2) Встановити додаток `helm install nginx bitnami/nginx`, де *bitnami/nginx* додаток з репозиторію *bitnami*, що встановлюється, а `nginx` - імʼя релізу. Реліз - екземпляр додатку. Один і той самий додаток може мати кілька релізів, наприклад, можна розгортати різні екземпляри для різного середовища `nginx-stage`, `nginx-prod`.

Після виконання цієї команди в кластері буде розгорнутий сервер `nginx` і встановлено всі необхідні обʼєкти *k8s*.

В браузері перейшовши на адресу `minikube ip` (або `http://localhost/`, якщо працюєте через `minikube tunnel`) можна побачити початкову сторінку `nginx`.

Переглянути список встановлених релізів, можна за допомогою команди `helm list`

```shell
NAME 	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART       	APP VERSION
nginx	default  	1       	2022-06-04 13:51:54.358577 +0300 EEST	deployed	nginx-12.0.0	1.22.0
```

Виконавши команду `kubectl get all` можна побачити всі компоненти, що було встановлено 
```
NAME                        READY   STATUS    RESTARTS   AGE
pod/nginx-8f97b57c7-ldqr7   1/1     Running   0          13m

NAME                 TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
...
service/nginx        LoadBalancer   10.109.7.11   127.0.0.1     80:31387/TCP   13m

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/nginx   1/1     1            1           13m

NAME                              DESIRED   CURRENT   READY   AGE
replicaset.apps/nginx-8f97b57c7   1         1         1       13m
```

Тепер *nginx* можна видалити виконавши `helm uninstall nginx`, далі він не знадобиться.
Детальніше про роботу з репозиторієм можна познайомитись на [офіційному сайті](https://helm.sh/docs/helm/helm_repo/#see-also)

### 1.2 Опис власного сервісу за допомогою *Helm*.

Для опису власних додатків потрібно описати *Helm Chart*. *Chart* - базовий елемент *Helm*,
що являє собою набір файлів, що описують всі дані, необхідні *k8s* для створення своїх об'єктів.
Простий *Helm Chart* містить наступну структуру:

```
demo-chart/
    Chart.yaml          # загальні дані про Chart
    values.yaml         # значення змінних Chart
    charts/             # Chart від якийх залежить поточний Chart
    templates/          # шаблони, які разом зі значеннями сформують об'єкти kubernetes
```

> Тут показано основні елементи, з повною структурою *Chart* можна ознайомитсь на [офіційному сайті](https://helm.sh/docs/topics/charts/#the-chart-file-structure)

`Chart.yaml` виглядає наступним чином:

```yaml
apiVersion: v2           # версія helm API, для Helm 3 - v2, для Helm 2 - v1 (обов'язково)
type: application        # Тип Chart, може бути application та library (library неможливо встановити, такий чарт має використовуватись в чарті типу application) (опціонально, за замовчуванням application)
version: 0.1.0           # Версія чарту по формату "SemVer 2" (обов'язково)
name: service2           # Довільне ім'я Chart, або додатку (обов'язково)
appVersion: "0.1"        # Версія додатку, не обов'язково "SemVer 2", бажано в лапках (опціонально)
description: "demo"      # Довільний опис (опціонально)
dependencies: []         # Сторонні залежності
```

> Коли змінюється шаблон чи *Chart.yaml* - потрібно змінювати `version`, коли змінюється код додатку, потрібно змінювати `appVersion`.

З повним списком властивостей, можна ознайомитись на [офіційному сайті](https://helm.sh/docs/topics/charts/#the-chartyaml-file)

У найпростішому вигляді `templates` буде містити такі ж самі обʼєкти, що і *k8s* без *Helm*.
Тому для того, щоб перевести поточний проект на *Helm* можна просто скопіювати обʼєкти *k8s* в директорію *templates* та описати для кожного сервісу *Chart.yaml*.
Приклад такого перенесення - [helm/v1](helm/v1). Тут містяться такі ж обʼєкти, що і в попередній роботі. 

Тепер для встановлення додатку потрібно перейти в папку *helm* та виконати команду `helm install local v1`.
Виконаємо `helm list`, щоб перевірити встановлені релізи. 

```shell
NAME 	NAMESPACE	REVISION	UPDATED                              	STATUS  	CHART     	APP VERSION
local	default  	1       	2022-08-27 17:30:20.049683 +0300 EEST	deployed	demo-0.0.1	0.1
```

### 2.1 Змінні шаблонів

В шаблонах *Helm* є можливість використовувати змінні. За замовчуванням в *Helm* є ряд змінних, основні це - `Chart`, що містить поля визначені в `Chart.yaml`
та `Release`, що містить інформацію про реліз, зокрема імʼя.
Ці змінні, як правило використовуються для ідентифікації обʼєктів *Helm*, наприклад імені ресурсу. 
Звернутись до цих змінних можна наступним чином `{{ .Chart.Name }}` 
`{{ }}` - визначаються динамічні конструкції шаблону. `.` на початку позначає контекст виконання,
за замовчуванням він глобальний, як його можна змінювати, розглянемо пізніше. Також зверніть увагу, що `Name` з великої букви.
В *Helm* стандартні змінні іменуються з великої літери, користувацькі рекомендується іменувати з малої літери.

Перша частина шаблону тепер буде мати такий вигляд

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}
...
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-{{ .Chart.Name }}
...
```

### 2.2 Користувацькі змінні

Змінні визначені користувачем описуються у файлі `values.yaml` та доступні через стандартну змінну *Helm* - `Values`.
Це можуть бути, як примітиви: строки, числа, булеві значення; так і вкладені структури та масиви.
Звернутись до цих змінних можна через обʼєкт `Values`, наприклад `{{ .Values.replicaCount }}`

Тоді `values.yaml` буде мати такий вигляд:

```yaml
replicaCount: 2
```

*Helm* надає можливість перевіряти згенеровані шаблони, за допомогою команди `helm template v2` (з директорії `3_heml/helm`).

### 2.3 Цикли

Деякі змінні можуть бути перелічуваними, наприклад сервіс може мати кілька відкритих портів

Спершу запишемо їх у `values.yaml`

> Хоч, в даному випадку, сервіс слухає 1 порт, часто порти для сервісу описуються саме масивом.
> Так як це дозволяє уніфікувати шаблони (як саме розглянемо пізніше)

```yaml
service:
  ports:
    - name: http
      port: 80
      containerPort: 8080
```

В шаблоні для отримання значень масиву використовується наступна конструкція:

```yaml
{{ range .Values.service.ports }}
  
{{ end }}
```

`range` підміняється контекст, (`.`) тепер буде означати значення елементу масиву.
Тому виведення масиву буде виглядати наступним чином:

```yaml
    {{ range .Values.service.ports }}
    - name: {{ .name }}
      port: {{ .port }}
      targetPort: {{ .name }}
    {{ end }}
```

Альтернативно можна використовувати іменовані змінні:

```yaml
    {{ range $key, $value := .Values.service.ports }}
    - name: {{ $value.name }}
      port: {{ $value.port }}
      targetPort: {{ $value.name }}
    {{ end }}
```

Такий варіант особливо корисно використовувати коли виводимо мапу значень, а не масив і потрібно отримати значення ключа.
Див приклад [ingress.yaml](helm/v2/templates/ingress.yaml)

### 2.4 Функції

Для додаткового форматування чи обробки виводу в *Helm* існують функції. Простий приклад такої функції - "default".
Функція повертає або значення змінної, або значення за замовчуванням, якщо змінна містить порожнє значення.
Синтаксис використання функцій - наступний:

`{{ default 1 .Values.replicaCount }}`, де `default` - назва функції, `1` та `.Values.replicaCount` - аргументи

Альтернативно можна використати "конвеєрний" підхід (*pipelines*)

`{{ .Values.replicaCount | default 1 }}`

В такому випадку `.Values.replicaCount` - останній аргумент, `default 1` - імʼя функції з усіма аргументами крім останнього.
Такий підхід більш поширений, оскільки дозволяє будували ланцюжок з викликів функції.
Наприклад:

`{{ password | b64enc | quotes}}`

`password` спершу закодовується в `base64`, а потім додаються лапки (для коректної обробки строк в *yaml*).

Детальніше про функції та конвеєри на [офіційному сайті](https://helm.sh/docs/chart_template_guide/functions_and_pipelines).
В *Helm* немає можливості визначати користувацькі функції, рекомендую переглянути [список доступних функцій](https://helm.sh/docs/chart_template_guide/function_list/).

### 2.5 Інші конструкції

Зверніть увагу на файл [helm/v2/charts/service2/templates/service.yaml](helm/v2/charts/service2/templates/service.yaml)
Перед `range` та `end` стоїть символ `-`.

```yaml
    {{- range .Values.service.ports }}
    - name: {{ .name }}
      port: {{ .port }}
      targetPort: {{ .name }}
    {{- end }}
```

Для демонстрації, приберемо їх, та запустимо `helm template v2`, сервіс для `service2` буде виглядати так:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: release-name-service2
spec:
  ports:
    
    - name: http
      port: 80
      targetPort: http
    
  selector:
    app: service2
```

Як бачимо, згенерований шаблон містить два порожні рядки. Це відбувається тому,
що коли шаблонізатор *Helm* генерує кінцеві файли,
він видаляє вміст між `{{` і `}}`. Але не видаляє сам рядок. 
Для контролю такого перенесення рядка в *Helm* використовується символ "-".
Одразу після відкриття дужки видаляє перенос з попередньго рядка.
"-" перед закриттям дужок видаляє перенос на наступний рядок.
Потрібно буде обережним з видаленням пробілів,
якщо, в цьому випадку, поставити дефіс перед закриттям дужок, шаблон буде з помилкою.

Також важливий функціонал це [умовні оператори](https://helm.sh/docs/chart_template_guide/control_structures/#ifelse) та   
[зміна контексту](https://helm.sh/docs/chart_template_guide/control_structures/#modifying-scope-using-with), з якими можна ознайомитись на офіційному сайті.

Тепер коли наші шаблони відредаговані можна оновити додаток: `helm uninstall local` потім `helm install local v2`.
Так як імена ресурсів змінились краще перестворити додаток, а не оновити.

### 3. Залежності

В поточній версії додатку в системі є одна стороння залежність - це база даних. Для опису якої використовується кілька обʼєктів *k8s*.
За допомогою *Helm* ми можемо встановити залежність на основі стороннього *Helm Chart*.
В даному випадку ми використаємо репозиторій *bitnami*, який містить ряд *Helm Chart* для різних інструментів.

Якщо репозиторій ще не доданий потрібно його додати

`helm repo add bitnami https://charts.bitnami.com`

Для пошуку у встановлених репозиторіях виконуємо команду:

`helm search repo postgresql`

```yaml
NAME                 	CHART VERSION	APP VERSION	DESCRIPTION                                       
bitnami/postgresql   	10.16.2      	11.14.0    	Chart for PostgreSQL, an object-relational data...
bitnami/postgresql-ha	8.2.8        	11.14.0    	This PostgreSQL cluster solution includes the P...
```

Для встановлення залежного *Chart*, додамо `dependencies` в `Chart.yaml`

```yaml
dependencies:
- name: postgresql
  version: 10.16.2
  repository: https://charts.bitnami.com/bitnami
```

та виконаємо команду 

`helm dependency build v3`

Після чого в папці *charts* зʼявиться архів з *Helm Chart* для бази даних.

Цей *Helm Chart* містить значення підключення до БД за замовчуванням, які нам не підходять.
Для конфігурування різних параметрів у файлі values.yaml додамо обʼєкт `postgresql`, що відповідає імені чарту.
Список всіх параметрів можна подивитись в [офіційному репозиторії](https://github.com/bitnami/charts/tree/master/bitnami/postgresql)
Обираємо ті, що нам потрібні та заповнюємо:

```yaml
postgresql:
  fullnameOverride: postgres # використовується як хост 
  postgresqlDatabase: "demo"
  postgresqlUsername: "demo"
  postgresqlPassword: "demo"
  persistence:
    size: 5Gi
    mountPath: /var/lib/postgresql/data
```

Тепер можна видалити папку `postgres` з папки `charts` та встановити додаток.

### 4 Іменовані шаблони

Під час створення шаблонів, зʼявляються частини, які повторюються в різних обʼєктах.
Для того, щоб оптимізувати код, можна виносити такі частини в іменовані шаблони.

За замовчуванням файл, що міститиме такі шаблони буде називатись `_helpers.tpl` та розміщений в директорії `templates`. 

Для того, щоб додати такий шаблон потрібно використати наступну конструкцію:

```yaml
{{- define "імʼя шаблону" -}}}
{{ ... вміст шаблону}}
{{- end -}}
```

Для прикладу можна створити шаблон, що визначає імʼя ресурсу

```yaml
{{- define "common.fullname" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}
```

Тепер у всіх обʼєктах замість використання:

```yaml
{{ .Release.Name }}-{{ .Chart.Name }}
```

потрібно вставити ```{{ include "common.fullname" . }}```

`.` в кінці позначає контекст відносно якого беруться значення

Це найпростіша версія, як правило виділяється окремо імʼя *chart* `name` та імʼя сервісу `fullname`, а також деякі інші приклад можна переглянути тут: [helm/v4/templates/_helper.tpl](helm/v4/templates/_helper.tpl)

> Зверніть увагу, що іменовані шаблони доступні в самому чарті та у вкладених чартах, але не навпаки,
> якщо описати шаблони `_helpers.tpl` наприклад в сервісі 1, то вони будуть доступні лише в межах цього чарту.

> В цій версії, було змінено `selector` через це оновлення чарту викличе помилку.
> Тому цю версію потрібно встановити окремо, спершу `helm delete local` і потім `helm install local v4`.

### 3. Бібліотека

Якщо подивитись на шаблони сервісів, то вони майже однакові (за виключенням складнішої конфігурації в `service1`).
Для зменшення дублювання коду, і простішої підтримки великої кількості сервісів, такі шаблони можна виділити в бібліотеку.
В директорії `charts` створимо *Chart* `common` типу *library*. У `values.yaml` пропишемо властивості за замовчуванням.

Розглянемо зміни в шаблоні, при винесенні в бібліотеці:

```yaml
{{- define "common.service" -}}

{{- $common := dict "Values" .Values.common -}} # Values з common Chart
{{- $noCommon := omit .Values "common" -}} # Values з основного Chart
{{- $overrides := dict "Values" $noCommon -}} # записуємо common Values в змінну Values
{{- $noValues := omit . "Values" -}} # містить усе крім Values, зокрема Chart, Release
{{- with merge $noValues $overrides $common -}}
---
  ...
```

- Перший рядок визначає назву шаблону
- Наступні рядки об'єднують `Values` з бібліотеки `Values` з сервісу, та значення не `Values` (`Relesase`,`Chart`). Таким чином в common Chart можна задавати значення за замовчуванням.


Після визначення шаблонів в бібліотеці в самих сервісах просто підключаємо створені шаблони:
```yaml
{{- template "common.deployment" . -}}
```

Недолік такого підходу в тому, що для коректного встановлення необхідно зібрати залежності для всіх чартів
```
helm dep build v5/charts/client
helm dep build v5/charts/service1
helm dep build v5/charts/service2
helm dep build v5
```

Це потрібно робити для того, щоб значення за замовчуванням з `common` були доступні в сервісі

Після цього `helm install local v5` встановить додаток.

### Управління версіями

Ми розглянули 2 важливі можливості *Helm* - встановлення сторонніх застосунків, та шаблонізатор для власного застосунку.
Але є ще один вадливий функціонал - це контроль версій. *Helm* дозволяє легко повертатись до попередньої версії.

Після того як застосунок встановлено за допомогою `helm install local v5` можна перевірити історії зміни застосунку.

`helm history local`

```shell
REVISION	UPDATED                 	STATUS  	CHART     	APP VERSION	DESCRIPTION     
1       	Tue Aug 30 20:29:53 2022	deployed	demo-0.0.5	0.1        	Install complete
```

Так як було виконано лише `install` в історії 1 запис (`uninstall` видаляє реліз і всю історію)
Змінимо версію чарта на `0.0.6` і виконаємо `helm upgrade local v5`

В історії зʼявиться ще один запис:

```shell
1       	Tue Aug 30 20:29:53 2022	superseded	demo-0.0.5	0.1        	Install complete
2       	Tue Aug 30 20:33:28 2022	deployed  	demo-0.0.6	0.1        	Upgrade complete
```

Тепер зробимо зміну (з помилкою) наприклад в файлі `values` змінимо `postgres.persistence.size` з 5 на `10Gi` (розмір диску не можна змінювати).
Знову виконаємо `helm upgrade local v5`

```shell
REVISION	UPDATED                 	STATUS    	CHART     	APP VERSION	DESCRIPTION                                                                                                                                                                                                                                                                                           
1       	Tue Aug 30 20:29:53 2022	superseded	demo-0.0.5	0.1        	Install complete                                                                                                                                                                                                                                                                                      
2       	Tue Aug 30 20:33:28 2022	deployed  	demo-0.0.6	0.1        	Upgrade complete                                                                                                                                                                                                                                                                                      
3       	Tue Aug 30 20:34:29 2022	failed    	demo-0.0.6	0.1        	Upgrade "local" failed: cannot patch "postgres" with kind StatefulSet: StatefulSet.apps "postgres" is invalid: spec: Forbidden: updates to statefulset spec for fields other than 'replicas', 'template', 'updateStrategy', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden
```

Тепер статус `failed` і також видно помилку через яку не вдалось встановити чарт.
Виправлення помилки в чарті можне зайняти час, а додаток потрібно відновити як омога швидше,
тому спершу потрібно відкатити версію до останньої успішної, а потім зайнятись виправленням помилки.

`helm rollback local 2` - тут 2 це номер `REVISION`.

```shell
REVISION	UPDATED                 	STATUS    	CHART     	APP VERSION	DESCRIPTION                                                                                                                                                                                                                                                                                           
1       	Tue Aug 30 20:29:53 2022	superseded	demo-0.0.5	0.1        	Install complete                                                                                                                                                                                                                                                                                      
2       	Tue Aug 30 20:33:28 2022	superseded	demo-0.0.6	0.1        	Upgrade complete                                                                                                                                                                                                                                                                                      
3       	Tue Aug 30 20:34:29 2022	failed    	demo-0.0.6	0.1        	Upgrade "local" failed: cannot patch "postgres" with kind StatefulSet: StatefulSet.apps "postgres" is invalid: spec: Forbidden: updates to statefulset spec for fields other than 'replicas', 'template', 'updateStrategy', 'persistentVolumeClaimRetentionPolicy' and 'minReadySeconds' are forbidden
4       	Tue Aug 30 20:42:42 2022	deployed  	demo-0.0.6	0.1        	Rollback to 2
```

Знову в полі `DESCRIPTION` чітко описано, що відбулось - `Rollback to 2`.

### Додаткові посилання

[An Introduction to Helm](https://www.youtube.com/watch?v=Zzwq9FmZdsU&t=2s&ab_channel=CNCF%5BCloudNativeComputingFoundation%5D)
