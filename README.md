#### Учебный проект YandexLMS "Программирование на Go | 24". Спринт 5. Финальный проект 5.1.                           
---
### ОБЩАЯ ИНФОРМАЦИЯ
Веб-сервис позволяет пользователю асинхронно вычислять арифметическое выражение. Поддерживает 4 основные арифметические операции ```+``` ```-``` ```*``` ```/```, унарный ```-``` и скобки  ```(``` ```)``` любой глубины вложенности. Выражение не должно содержать знак ```=```. Все пробелы игнорируются. 
Сервис поддерживает многопользовательский режим (поддерживается JWT авторизация). 
Новый пользователь должен быть создан с помощью POST запроса на эндпоинт ```/api/v1/register  ```, Body запроса должно содержать JSON следующего вида:
```
{
  "login": "имя_пользователя",
  "password": "пароль"
}
```
Данные регистрируемого пользователя записываются в базу данных (SQLite) и сохраняются при перезагрузке ПО. Пароль хранится в хэшированном виде (не может быть получен в явном виде из БД). Так же в БД сохраняются данные вычисленных выражений в контексте того конкретного пользователя, который отправил запрос на вычисление этого вырадения. 
Если при запуске сервиса БД обнаружить не удалось, то она создается автоматически. Исходно в репозитарии присутствует БД с одним пользователем user01/psw_user_01 и зарегистрированным от его имени выражением.
Перед отправкой сервису запросов на вычисление выражений и получения результатов их вычисления пользователь должен авторизоваться с помощью POST запроса на эндпоинт ```/api/v1/login  ```, Body запроса должно содержать JSON следующего вида:
```
{
  "login": "имя_пользователя",
  "password": "пароль"
}
```
В случае успеха сервис возвращает код 200, а в Body ответа будет JSON следующего вида:
```
{"token":"eyJhbGciOiJIUzI1NiIs   ...   dWDbvsTt8o-u9hg"}
```
Данный JWT токен должен быть скопирован в Body во все последующие запросы к сервису.

Т.к. сервис позволяет распаралеливать процес вычисления каждого арифметичекого выражения между множеством потоков и асинхронно получать результаты вычислений, то POST запрос на эндпоинт ```/api/v1/calculate ``` передает вычисляемое выражение сервису, который разбивает его на группу взаимосвязанных задач, передает отдельные задачи на параллельное вычисление и отслеживает состояние процесса.
Body запроса должно содержать JSON следующего вида:
```
{
"token":"eyJhbGciOiJIUzI1NiIs   ...   dWDbvsTt8o-u9hg",
"expression": "Выражение"
}
```
Выражение может включать следующие символы ```()+-*/0123456789``` и должно быть корректным математическим выражением.
  В случае успеха сервис возвращает код ```200```, в качестве результата в body возвращает JSON следующего вида:
```json
{"id":  123}
```
Если входные данные не соответствуют указанным требованиям, сервис возвращает код ```422```, а в body JSON следующего вида:
```json
{"error": "Expression is not valid"}
```
С помошью полученного идентификатора выражения у сервиса можно запросить текущее состояние процесса вычислениния этого выражения. Для этого необходимо отправить GET запрос на эндпоинт ```/api/v1/expressions/{id} ``` , где {id} - уникальный идентификатор выражения (в нашем примере 123). 
В Body запроса должно содержать JWT токен полученный ранее при авторизации пользователя:
```
{
"token":"eyJhbGciOiJIUzI1NiIs   ...   dWDbvsTt8o-u9hg"
}
```

Если выражение с указанным id было ранее принято сервисом ОТ ДАННОГО ПОЛЬЗОВАТЕЛЯ, то сервис вернет код ```200```, в качестве результата в body отправит JSON  вида:
```json
{"expression":{"Expression":"(-1+2)*3 + (11+7)/2","id":123,"status":"success","result":12}}
```
Если выражения с указанным идентификатором не существует, то будет получен код ```422```. Полный список всех обрабатываемых выражений может быть получен по GET запросу  на эндпоинт ```/api/v1/expressions``` .
В Body запроса так же должен содержаться полученный при авторизации JWT токен:
```
{
"token":"eyJhbGciOiJIUzI1NiIs   ...   dWDbvsTt8o-u9hg"
}
```
В случае успеха в body получим JSON  вида:
```json
{"expressions":[
{"Expression":"-1 + 100","id":5,"status":"ready","result":0},
{"Expression":"2+2)*2","id":12,"status":"calculate","result":0},
{"Expression":"2/(1-1)","id":19,"status":"failed","result":0},
{"expression":{"Expression":"(-1+2)*3 + (11+7)/2","id":123,"status":"success","result":12}}
]}
```
Возможны состояния/статусы выражения:
* "ready" - выражение помещено в очередь на вычисление;
* "calculate" - выражение сейчас вычисляется;
* "failed" - в ходе вычисления произошли ошибки (например деление на 0);
* "success" - выражение успешно вычислено, результат содержится в поле "result".

Если в работе сервиса произошла какая-то ошибка, то будет возвращен код ```500``` и в body JSON следующего вида:
```json
{"error": "Internal server error"}
```
Архитектурно сервис состоит из 2 частей:
* Оркестратор - взаимодействует с пользователями с помощью приведенных выше запросов, выполняет предобработку поступивших на вычисление выражений (проверяет корректность, поводит декомпозицию на отдельные задачи), контролирует порядок исполнения отдельных задач агентами, собирает и хранит в базе данных результаты. В составе одного сервиса может работать только один оркестратор.
* Агенты - получают от оркестратора задачи и передают их одному из нескольких собственных процессов на исполнение. Результаты исполнения отдельных задач возвращаются оркестратору. В составе одного сервиса может работать любое число агентов. Каждый агент может содержать множество параллельных процессов - отдельных, не взаимодействующих между собой каналов исполнения задач. Агенты взаимодействуют с оркестратором (получают новые задачи и возвращают результаты) по gRPC (протоколу удаленного вызова проседур). Спецификация протокола представлена в файле .\proto\asyncCalculator.proto.
---
### НАСТРОЙКА КОМПОНЕНТОВ СЕРВИСА

Настройка сервиса осуществляется с помощью установки значений следующих переменных среды:
* ASYNC_CALCULATOR_HTTP_HOST - хост для отправки пользователями http запросов к оркестратору; 
* ASYNC_CALCULATOR_HTTP_PORT - порт для отправки пользователями http запросов к оркестратору;
* ASYNC_CALCULATOR_GRPC_HOST - хост для взаимодействия агентов с оркестратором по gRPC;
* ASYNC_CALCULATOR_GRPC_PORT - порт для взаимодействия агентов с оркестратором по gRPC;
* TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах;
* TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах;
* TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах;
* TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах;
* COMPUTING_POWER - количество каналов параллельного исполнения задач у агента.

Если переменные среды не установлены, то действуют следующие значения по умолчанию:
```
  ASYNC_CALCULATOR_HTTP_HOST = localhost
  ASYNC_CALCULATOR_HTTP_PORT = 8080
  ASYNC_CALCULATOR_GRPC_HOST = localhost
  ASYNC_CALCULATOR_GRPC_PORT = 5000
  TIME_ADDITION_MS = 5000
  TIME_SUBTRACTION_MS = 5000
  TIME_MULTIPLICATIONS_MS = 5000
  TIME_DIVISIONS_MS = 5000
  COMPUTING_POWER = 10
```

### ЗАПУСК КОМПОНЕНТОВ СЕРВИСА

Для запуска сервиса в ОС Windows необходимо:
1. Склонировать этот репозитарий к себе на компьютер например в каталог d:\Projects\yandexLMSGo. 
2. Запустить от имени администратора две командные строки (для этого можно например нажать Win + R, ввести cmd, нажать «OK»). 
3. В первой командной строке надо перейти в каталог п.1., затем, для запуска оркестратора, ввести
```sh
   go run cmd\orchestrator\main.go
```
   и нажать «Enter».

4. Во второй командной строке надо перейти в каталог п.1., затем, для запуска агента, ввести 
```sh
   go run cmd\agent\main.go
```
   и нажать «Enter».
   
Предварительно убедитесь, что Брандмауэр Windows и антивирус не блокируют обращения к используемым сервисом хосту и порту. Если это необходимо, выполните соответствующую настройку.

Обратите внимание, что для сборки проекта вам возможно понадобиться загрузить зависимости используя команду go get [package]@version, а для компиляции драйвера SQLite3 установить gcc компилятор.

---
### ЗАПУСК ЮНИТ-ТЕСТОВ
* оркестратора:
```bash
go test .\pkg\orchestrator\
```
* агента:
```bash
go test .\pkg\agent\
```

---
### ПРОВЕРКА РАБОТЫ СЕРВИСА С ПОМОЩЬЮ CURL 
Обратите внимание, что хотя в запросах использован JWT токен пользователя user01/psw_user_01 (уже зарегистрированного в БД, поставлямой в репозитории по умолчанию, повторная регистрация будет отклонена) с довольно большим сроком жизни, может понадобиться его обновить. В этом случае не забудте копировать в Body запросов новый токен, получаемый при логине. 

1. Регистрация пользователя user01 с паролем psw_user_01 (если пользователь с именем user01 в БД уже существует, то будет получен ответ "{"error":"Can't create user"}" - попробуйте создать пользователя с другим именем, например user_02)
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"login\": \"user01\",\"password\": \"psw_user_01\"}" localhost:8080/api/v1/register
```
##
2. Авторизация пользователя user01 с паролем psw_user_01
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"login\": \"user01\",\"password\": \"psw_user_01\"}" localhost:8080/api/v1/login
```
Ответ: ```{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc"}```
##
3. Простое выражение ```2+2*2```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"2+2*2\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
4. Простое выражение ```-2+2*2``` (первый унарный минус). Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-2+2*2\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
5. Простое выражение ```2+2*(-2)``` (унарный минус в скобках). Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"2+2*(-2)\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
6. Простое выражение со скобками ```-(2+2)*2``` (унарный минус перед скобками). Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"(2+2)*2\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
7. Сложное выражение ```-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
8. Сложное выражение ```-(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
9. Сложное выражение с пробелами ```-(3*( 12/ 4+(-2 + 8/2 ) ) +(7-20/5*(9-2* 2*2+1))*(  -10))```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-(3*( 12/ 4+(-2 + 8/2 ) ) +(7-20/5*(9-2* 2*2+1))*(  -10))\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"id":<число>}```
##
10. Получение полного списка состояний всех выражений. Наберите в командной строке
```sh
curl -X GET -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\"}" localhost:8080/api/v1/expressions
```
Ответ: ```{"expressions":[{"Expression":"1+2*3","id":1,"status":"success","result":7},{"Expression":"(1+2)*(3+4)","id":4,"status":"calculate","result":0}]}```
##
11. Получение состояния выражения с "id": 1. Наберите в командной строке
```sh
curl -X GET -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\"}" localhost:8080/api/v1/expressions/1
```
Ответ: ```{"expression":{"Expression":"1+2*3","id":1,"status":"success","result":7}}```
##
12. Получение состояния несуществующего выражения с "id": 2. Наберите в командной строке
```sh
curl -X GET -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\"}" localhost:8080/api/v1/expressions/2000
```
Ответ: ```404 not found.```
##
13. Ошибка в выражении (непарная скобка) ```-3*12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-3*12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10)\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"error": "Expression is not valid"}```
##
14. Ошибка в выражении (недопустимый символ ```=```) ```-(3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))=```. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-3*(12/4+(-2+8/2))+(7-20/5*(9-2*2*2+1))*(-10))=\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"error": "Expression is not valid"}```
##
15. Ошибка в выражении (недопустимый символы) -(3*(12/4+(-A+8/2))+(7-20/5*(9-2*2*2+b))*(-10))=. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"-3*(12/4+(-A+8/2))+(7-20/5*(9-2*2*2+b))*(-10))\"}" localhost:8080/api/v1/calculate
```
Ответ: ```{"error": "Expression is not valid"}```
##
16. Ошибка в хосте. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"2+2*2\"}" localhostttt:8080/api/v1/calculate
```
Ответ: ```curl: (6) Could not resolve host: localhostttt```
##
17. Ошибка в номере порта. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"2+2*2\"}" localhost:8081/api/v1/calculate
```
Ответ: ```curl: (7) Failed to connect to localhost port 8081 after 2255 ms: Could not connect to server```
##
18. Ошибка в пути. Наберите в командной строке
```sh
curl -X POST -H "Content-Type: application/json" -d "{\"token\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI5OTI4MjAsImlhdCI6MTc0Njk5MjgyMCwibG9naW4iOiJ1c2VyMDEiLCJuYmYiOjE3NDY5OTI4MjAsInVzZXJfaWQiOjF9.1gm_1vDApQpgv7uvTGPtqXxHzfJ1XV9lgydYZnRU3Lc\", \"expression\": \"2+2*2\"}" localhost:8080/api/v2/recalculate
```
Ответ: ```404 page not found```

---
### ПРОВЕРКА РАБОТЫ СЕРВИСА С ПОМОЩЬЮ POSTMAN 

1. Установите Postman (https://www.postman.com/).
2. Импортируйте в Postman коллекцию запросов для asyncCalculator. Файл коллекции ```authAsyncCalculator.postman_collection.json``` расположен в корне репозитория. Обратите внимание, что хотя в запросах использован JWT токен пользователя user01/psw_user_01 (уже зарегистрированного в БД, поставлямой в репозитории по умолчанию, повторная регистрация будет отклонена) с довольно большим сроком жизни, может понадобиться его обновить. В этом случае не забудте копировать в Body запросов новый токен, получаемый при логине.
3. Выберете необходимый вам запрос, при необходимости измените параметры, нажмите «Send».

---
В случае возникновения проблем с запуском и/или проверкой сервиса, а также при появлении вопросов по исходному коду прошу обращаться в телеграмм https://t.me/atrufmanov. Постараюсь оперативно ответить на все вопросы.
