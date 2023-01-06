# MMS TEST

Serviço de variações de médias móveis simples, de 20, 50 e 200 dias,
das moedas Bitcoin e Etherium.

## Descrição:
O serviço de media faz uso de uma api para carregar a tabela de médias, ao inicializar o serviço pela primeira vez um job ira fazer a carga inicial (caso esteja habilitado). Após a primeira execução é necessário que seja atualizado diariamente está tabela, isso é feito por um outro job que está schedulado para executar diariamente as 21:00 horas. Todas as configurações dos jobs pode ser consultada e/ou alterada no arquivo `resources/application.yml`.

```
extractor:
  enabled: true
  mocked: true
  url: https://mobile.mercadobitcoin.com.br/v4/%s/candle
  range: 365
  daily:
    maxPeriod: 200
    frequency: 1
    hour: "21:00"
    retry: 3
    retryTime: 30s
  pairs:
    - "BRLETH"
    - "BRLBTC"
```
O job diario tem a possibilidade de configurar a quantidade de retentativas e o tempo entre elas caso falhe o processamento. Mas para o caso de os retries não forem o suficiente, tambem disponibilizamos um endpoint para processamento manual, mais detalhes no swagger da aplicação.

## Como executar o serviço:

1. O serviço foi desenvolvido para utilizar o banco de dados PostgresSQL, que pode ser inicializado utilizando o arquivo `resources/docker/docker-compose.yml`. Ao inicializar o aquivo `resources/docker/init.sql` será executado para a criação do database e a tabela da aplicação.

2. É possivel utilizar o arquivo `Makefile` para buildar a aplicação (é necessário a instalação do make):
Comando:
`make run-install` (quando for a primeira execução)
`make run` (depois da primeira execução)


3. Com a aplicação rodando é possivel consultar a documentação acessando:
```
http://localhost:8080/swagger/index.html#/
```

Obs. A validação de não permitir consultas cuja data de início seja anterior a 365 dias está comentada para conseguir usar os dados mockados