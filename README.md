**Rate limiter**  

**Descrição**  
Este projeto criado em GO controla o total de requisições por segundo que podem ser recebidas. Esse total é definido em um arquivo .Env que pode ser diferente para requisições com ou sem TOKEN.  
O tempo de bloqueio e de expiração do TOKEN também são definidos no arquivo .Env.  
As informações são armazenadas no banco de dados REDIS.

**Configuração do arquivo .env**  
  
  * REQUISICOES_POR_SEGUNDO_IP: requisições por segundo sem TOKEN.
  * REQUISICOES_POR_SEGUNDO_TOKEN: requisições por segundo com TOKEN.
  * TEMPO_BLOQUEIO_EM_SEGUNDO_IP: tempo de bloqueio em segundos sem TOKEN.  
  * TEMPO_BLOQUEIO_EM_SEGUNDO_TOKEN: tempo de bloqueio em segundos com TOKEN.
  * TEMPO_EM_SEGUNDOS_EXPIRACAO_TOKEN: tempo de expiração em segundos do TOKEN.  
     
**Utilização**  
  
Tanto o REDIS quanto a aplicação estão executando em containers Docker. Portanto para iniciar a execução deve-se utilizar comando "docker-compose up -d".  
Para configurar o Handle que se quer que o rate limiter seja aplicado, deve-se passá-lo como parâmetro para a função rateLimiterMiddleware como exemplificado no arquivo main.go.  





