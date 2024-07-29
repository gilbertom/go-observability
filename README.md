# Desafio Prático Tracing Distribuído e Span
Projeto do Desafio Prático Observabilidade & Open Telemetry para conclusão da Pós Graduação em Go Expert da Full Cycle.

<p align="center">
  <img src="https://blog.golang.org/gopher/gopher.png" alt="">
</p>

Esta aplicação, desenvolvida em Go, recebe um POST com um JSON contendo um CEP e valida se ele possui 8 caracteres. Em seguida, realiza um GET para consultar a API ViaCEP e obter a localização. Posteriormente, consulta a WeatherAPI para obter a temperatura em Celsius. Por fim, a API retorna a cidade e a temperatura em Celsius, Fahrenheit e Kelvin.

<br>

## Índice


- [Instalação e Start da aplicação](#instalação-e-start-da-aplicação)
- [Pré requisitos](#pré-requisitos)
- [Como Usar](#como-usar)
- [Tracing](#tracing)
- [Contato](#contato)
- [Agradecimentos](#agradecimentos)

<br>


## Instalação e Start da aplicação

```sh
$ git clone https://github.com/gilbertom/go-observability.git
$ docker-compose up -d --build
```
<br>

## Pré requisitos
Esta aplicação utiliza a WeatherAPI e é obrigatório ter uma chave de acesso.  
Necessário também ter o Go e Docker instalado.

1. Crie uma conta no site <a href="https://www.weatherapi.com">WeatherAPI</a>.
2. Na página <a href="https://www.weatherapi.com/my/">Home</a>, copie a chave de acesso (API Key).
3. Adicione a chave de acesso no arquivo ".env" localizado na pasta /cmd/app, utilizando a variável API_KEY_WEATHER.

<br>


## Como Usar

No diretório /api temos o arquivo 'post_temperature.http' que envia uma requisição POST para a aplicação. 

Obs.: Imprescindível instalar a Extensão 'HTTP Client' no seu Visual Studio Code.  

Exemplo
```sh
POST http://localhost:8081/serviceA HTTP/1.1
Host: localhost:8081
Content-Type: application/json

{
  "cep": "28951620"
}
```

Response
  ```sh
  HTTP/1.1 200 OK
  Content-Type: application/json
  Date: Sun, 28 Jul 2024 23:27:59 GMT
  Content-Length: 86
  Connection: close

  {
    "city": "Armação dos Búzios",
    "tempC": 32.3,
    "tempF": 90.13999999999999,
    "tempK": 305.3
  }
  ```
<br>  

___

  Request de um CEP inválido
  ```sh
  POST http://localhost:8081/serviceA HTTP/1.1
  Host: localhost:8081
  Content-Type: application/json

  {
    "cep": "28951620A"
  }
  ```

  Response
  ```sh
  HTTP/1.1 422 Unprocessable Entity
  Content-Type: text/plain; charset=utf-8
  X-Content-Type-Options: nosniff
  Date: Sun, 28 Jul 2024 23:28:35 GMT
  Content-Length: 16
  Connection: close

  invalid zipcode
  ```
<br>  


## Tracing
Esta aplicação possui instrumentação que possibilita o rastreamento das requisições. Para acessar o Zipkin, visite: http://localhost:9411

Abaixo um exemplo de requisição:

![Exemplo de Imagem do Zipkin](https://github.com/gilbertom/go-observability/blob/master/assets/Zipkin.jpg)

<br>



## Contato
Para entrar em contato com o desenvolvedor deste projeto:
[gilbertomakiyama@gmail.com](mailto:gilbertomakiyama@gmail.com)

<br>


## Agradecimentos
Gostaria de expressar minha sincera gratidão a todo o time do curso de Pós-Graduação em Go Avançado da FullCycle pelo empenho, dedicação e excelência no ensino. Suas contribuições foram fundamentais para o meu desenvolvimento e sucesso. Muito obrigado!
