## BaseGo
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
git clone https://github.com/filgueirasjulio/base.git
```

-   **Entrar no diretório**
    
```
cd base
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
./base migrate
```

-   **Popular banco de dados com dados de exemplo**
    - todas as seeds
```
./base seed 
```
    - seed de um metodo específico
```
./base seed --model=User
```

-   **Iniciar aplicação**
    
```
./base run
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
    
