# Diretrizes de Contribuição

Bem-vindo ao projeto `go-data-api-microservices`! Agradecemos o seu interesse em contribuir. Para garantir um processo de colaboração suave e eficaz, por favor, siga estas diretrizes.

## Como Contribuir

1.  **Fork o Repositório**: Comece fazendo um fork do repositório para a sua conta GitHub.
2.  **Clone o Repositório**: Clone o seu fork para a sua máquina local:
    ```bash
    git clone https://github.com/YOUR_USERNAME/go-data-api-microservices.git
    cd go-data-api-microservices
    ```
3.  **Crie uma Nova Branch**: Crie uma branch para a sua feature ou correção de bug. Use nomes descritivos para as branches (ex: `feature/nova-autenticacao`, `bugfix/correcao-login`).
    ```bash
    git checkout -b feature/sua-feature
    ```
4.  **Faça Suas Alterações**: Implemente suas alterações, garantindo que o código siga os padrões de estilo e as melhores práticas do projeto.
5.  **Testes**: Certifique-se de que todos os testes existentes passem e adicione novos testes para cobrir suas alterações, se aplicável.
    ```bash
    cd src
    go test ./...
    ```
6.  **Commit Suas Alterações**: Escreva mensagens de commit claras e concisas. Use o formato `Tipo(Escopo): Descrição` (ex: `feat(auth): Adiciona autenticação JWT`).
    ```bash
    git commit -m "feat(sua-feature): Adiciona sua nova funcionalidade"
    ```
7.  **Push para o seu Fork**: Envie suas alterações para o seu fork no GitHub.
    ```bash
    git push origin feature/sua-feature
    ```
8.  **Abra um Pull Request**: Vá para o repositório original no GitHub e abra um Pull Request da sua branch para a branch `main` (ou `master`, dependendo da configuração do projeto). Descreva suas alterações detalhadamente e referencie quaisquer issues relevantes.

## Padrões de Código

*   **Formatação**: Use `go fmt` para formatar seu código Go.
*   **Linting**: Use `golangci-lint` para verificar problemas de linting.
*   **Comentários**: Comente seu código onde for necessário para explicar a lógica complexa.

## Relatando Bugs

Se você encontrar um bug, por favor, abra uma issue no GitHub com as seguintes informações:

*   Uma descrição clara e concisa do bug.
*   Passos para reproduzir o comportamento.
*   O comportamento esperado.
*   O comportamento real.
*   Sua versão do Go e sistema operacional.

## Sugerindo Melhorias

Para sugerir uma nova funcionalidade ou melhoria, abra uma issue no GitHub com as seguintes informações:

*   Uma descrição clara e concisa da funcionalidade proposta.
*   O problema que ela resolve ou a melhoria que ela traz.
*   Exemplos de uso, se aplicável.

Obrigado por sua contribuição!
