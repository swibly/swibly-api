<h1 align="center">👷 ARKHON API 👨‍💻</h1>

> [!CAUTION]
> O uso ou desenvolvimento **NÃO** oficial é desaconselhado devido às frequentes alterações de código, o que pode resultar em comportamentos imprevistos ou falhas críticas, tanto quanto mudança dos nomes, licenças e identificadores.

## ❓️ Sobre

**Arkhon** é uma plataforma dedicada à criação de plantas baixas digitais, com um foco específico em arquitetura. O projeto _arkhon-api_ constitui o Back-end da nossa plataforma, sendo conduzido por uma API escrita em _Golang_. Este componente desempenha um papel fundamental na infraestrutura técnica da **Arkhon**, facilitando a gestão eficiente dos dados e operações essenciais. O projeto foi concebido pela equipe **Swibly**, visando oferecer soluções tecnológicas acessíveis e eficazes para o campo da arquitetura digital.

## 🛠️ Instalação e Configuração

| Ferramenta | Versão (min) | Opcional? |
| ---------- | ------------ | --------- |
| Go         | `1.22`       | ❌        |
| Docker     | `26.1.0`     | ✔️        |
| Make       | `4.4.1`      | ✔️        |

**1. Clone este repositório:**

```bash
git clone https://github.com/swibly/arkhon-api.git
```

---

**2. Copie o arquivo `.env.example` para `.env`:**

```bash
cp .env.example .env
```

---

**3. (Opcional) Configure o projeto:**

Você pode configurar o projeto modificando o arquivo `.env` ou indo até a pasta `config/` e configurar os arquivos _YAML_.  
Este projeto está configurado para ser plug-and-play, o que significa que não é necessário fazer nenhuma configuração adicional para começar a usá-lo.

---

**4. (Ignore caso o banco não for local) Inicie o Banco de Dados:**

Inicie utilizando o comando `make`...

```bash
make up
```

...ou utilize o docker diretamente:

```bash
docker compose up -d
```

> Para desligar o banco: `make down` ou `docker compose down`

## 🚀 Uso

Utilize o comando `make` para iniciar a API...

```bash
make
```

...ou manualmente:

```bash
go build -race -o build/api -v ./cmd/api/main.go
./build/api
```

## 📝 Contribuição

Sinta-se à vontade para contribuir com melhorias! Antes de enviar uma solicitação de pull, certifique-se de que seu código está em conformidade com os padrões de código e de que você testou as alterações.

## 📃 Licença

Este projeto está licenciado sob a [ISC License](LICENSE)
