# ms-notificator

Este microsserviço tem por objetivo fornecer disparos e agendamentos de eventos de notificação. 

## Arquitetura



![Diagrama de Arquitetura](https://gcdnb.pbrd.co/images/U83QdQSIl8je.png?o=1)

*Diagrama de alto nível da solução*

### Componentes

Os componentes desenvolvidos e escolhidos foram baseados nos requisitos soliticitados.

- **notifier-client**: É o serviço ou cliente que irá enviar a notificação.
- **notificator-api**: Esta componente tem a responsabilidade de receber o evento de notificação, avaliar sua validade de acordo com as preferências do usuário e, com base no evento recebido, optar por publicar o evento de notificação de maneira instantânea ou fazer ao seu agendamento.
- **user-notification-preference**: É responsável por disponibilizar as preferências de notificação do usuário.
- **scheduledNotifications, notifications**: As duas coleções no banco de dados MongoDB responsáveis por armazenar as notificações e os agendamentos, o banco foi escolhido devido à escalabilidade e flexibilidade que o mesmo oferece. Especificamente, o MongoDB se torna interessante para gerenciar grandes volumes de dados, além de ter seu modelo de dado flexível para armazenar diferentes modelos de notificações.
- **notificator-worker**: Este componente assume a responsabilidade de consultar e encaminhar as notificações agendadas. Primeiramente, realiza uma consulta para verificar o status e a data dos agendamentos. Posteriormente, faz um processamento concorrente para o envio de todas as notificações disponíveis, garantindo uma execução eficiente e otimizada.
- **topic-SNS**: Escalabilidade: O SNS é altamente escalável, sendo capaz de lidar com um grande volume de mensagens e enviar para uma grande quantidade de destinatários. Oferece alta disponibilidade e durabilidade, garantindo que as mensagens de notificação sejam entregues de forma confiável, além de suportar a publicação de mensagens a partir de diversas fontes e para uma variedade de endpoints, incluindo Amazon SQS, HTTP/S web services, e-mails, e AWS Lambda, entre outros.

- **web-api**: É o serviço inscrito no tópico do SNS que recebe as noticiações.


### Casos de uso
O diagrama detalhado apresentado acima demostra a estrutura dos componentes e o fluxo de interações implicados no processo de envio ou agendamento de notificações. Essa estrutura é dividida em dois cenários de casos de uso:

1. O disparo de notificação instantâneo, no qual o cliente emite uma chamada HTTP para o serviço rotulado como "notificator-api". Nesse contexto, o evento é validado de acordo com as preferências definidas pelo usuário. Posteriormente, esse evento é publicado em um tópico SNS, onde uma aplicação web está inscrita para receber tais eventos. Além disso, o evento também é armazenado em uma coleção chamada "notifications" em um banco de dados MongoDB, proporcionando uma forma de persistência e registro de todas as notificações enviadas.
\
Abaixo segue um exemplo requisição.

```
curl --location 'http://localhost:8081/api/v1/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "rawMessage": "{\"message\": \"RawMessage\"}",
    "userId": "123"
}'
```

2. A agendamento de notificações é efetuada através de uma solicitação específica que inclui a data do agendamento no payload, enviada de um cliente ao serviço "notificator-api". Esse serviço, então, valida o pedido com base nas preferências do usuário. Em seguida, é gerado um registro na coleção "scheduledNotifications", no banco de dados Mango, com o status de "pending".

Posteriormente, o "notificator-worker", por meio de um cronjob, busca todas as notificações com o status de "pending" cuja data de agendamento é igual ou anterior à data atual. A notificação é então publicada em um tópico SNS, ao qual um cliente web está inscrito.

\
Abaixo segue um exemplo requisição.

```
curl --location 'http://localhost:8081/api/v1/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "rawMessage": "{\"message\": \"RawMessage\"}",
    "dateScheduling": "2023-02-26 11:22:22",
    "userId": "123"
}'
```
O campo "dateScheduling" é utilizado para identificar qual fluxo deve ser seguido, se a notificação é enviada imediatamente ou vai para o fluxo de agendamento.


## Dependências
- [Go](https://go.dev/learn/)
- [MongoDB](https://www.mongodb.com/docs/)
- [LocalStack](https://docs.localstack.cloud/getting-started/)
- [Docker](https://docs.docker.com/?_gl=1*owu4p0*_ga*MjExMzE3MDIxOS4xNjg0NDU3MTk4*_ga_XJWPQMJYHQ*MTY4ODMzMTc3My40LjEuMTY4ODMzMTc3NC41OS4wLjA.)

## Como executar o ambiente

### Como executar

1. Clone o repositório 


```cmd
$ git clone https://github.com/GeovaneCavalcante/ms-notificator
```
2. Dentro do diretório do repositório baixado construa o ambiente
```cmd
$ make setup-dev
```
3. Destruir todo ambiente 
```cmd
$ make down-all
```

### Utilizando a aplicação

1. Envio de notificação instantânea
```cmd
curl --location 'http://localhost:8081/api/v1/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "rawMessage": "{\"message\": \"RawMessage\"}",
    "userId": "123"
}'
```


2. Agendamento de notificação
```cmd
curl --location 'http://localhost:8081/api/v1/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "rawMessage": "{\"message\": \"RawMessage\"}",
    "dateScheduling": "2023-02-26 11:22:22",
    "userId": "123"
}'
```

3. Listagem de notificações recebidas no cliente web
```cmd
curl --location 'http://localhost:8083/notifications' \
--header 'accept: application/json'
```
### Docs Swagger
Documentações das APIs **notificator-api** e **client-web**:
- [**notificator-api**](http://localhost:8081/swagger/index.html#/)
- [**client-web**](http://localhost:8083/apidocs/#/default/get_notifications)
