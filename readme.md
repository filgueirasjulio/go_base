## TradeApi
================

## Como Usar
-------------

## Pré-requisitos

-   Git
-   Docker
-   Go (opcional, para build manual)
    
## Passo a Passo

-   **Clonar repositório**
    
```
git clone https://github.com/Tfos-Software/tradeapi.git
```

-   **Entrar no diretório**
    
```
cd tradeapi
```

-   **Configurar ambiente**
    
```
cp .env.example .env
```

-   **Subir containers Docker**

```
docker-compose up -d
```

-   **Buildar aplicação**

```
go build
```

-   **Migrar banco de dados**   

```
./tradeapi migrate
```

-   **Popular banco de dados com dados de exemplo**
    - todas as seeds
```
./tradeapi seed 
```
    - seed de um metodo específico
```
./tradeapi seed --model=User
```

-   **Iniciar aplicação**
    
```
./tradeapi run
```

## Testar API
-------------

-   para testar a API.
    ```
    http://localhost:8082/api/users
    ```

-   para visualizar documentação Swagger.
    ```
    http://localhost:8082/api/swagger
    ```
    gerar nova documentação
    ```
    swag init -o utils/docs
    ```
    
