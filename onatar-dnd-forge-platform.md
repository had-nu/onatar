# Onatar (D&D Character Builder based on 5.5e (2024))
Arquiteto de solução / product owner técnico
    1. Documento de Requisitos (funcionais e não-funcionais)
    2. Modelo de Ameaças (Threat Model) — porque mencionaste DevSecOps e tens background em AppSec
    3. Arquitetura de Referência — stack, deploy, dados
    4. Roadmap de MVP — o que entra na v1 vs. depois
    5. Definition of Done — critérios de aceitação por feature

## Fase 1: Levantamento de Requisitos
Responde às perguntas que conseguires. Se não tiveres certeza, diz "não sei ainda" — isso também é informação útil.

A. Âmbito & Público-alvo

    - Quem é o utilizador típico? Jogadores que querem gerir a ficha própria, mestres que querem ver fichas dos jogadores, ou ambos?
        >> Ambos. Para o jogador, a ideia é tornar D&D fácil de iniciar. Olhar apara uma ficha de personagem vazia é intimidador. Junte isso a quantidade de livros e suplementos para se ler e temos uma barreira de entrada absurda. Ao mesmo tempo, há uma abundância de informação disponível para ser coletada e organizada. Para o DM, gerenciar os diversos aspéctos de uma campanha pode virar um pesadelo se os recursos estiverem dispersos ou desorganizados. Ainda, um bom DM precisa considerar a historia de seus jogadores para fazer campanhas envolventes e conhecer a fundo a ficha dos personagens para não fazer metagame involuntário contra eles. Precisa ter visibilidade da história, do começo ao fim, como uma narrativa tradicional: o arco principal (longo), os arcos curtos, os atos, as side quests, etc. Vilões coerentes, propósito coerente, e respeito à cosmologia do cenário.

    > "Mais simples e acessível que D&D Beyond" — define "simples". Menos passos para criar personagem? Interface mais limpa? Menos regras automáticas (deixa o utilizador decidir)? Ou algo específico que odeias no D&D Beyond?
        >> D&D Beyond é uma ferramenta oficial que funciona muito bem para gerir campanha, personganes e npcs se voce investir uma pequena fortuna para comprar todos os livros. Criaçãod e homebrews, cominidades, gerador de mapas, etc. Tudo isso esta fora. Mas alguns recursos são chave para o Onatar: character builder com todas as possibilidades de escolha para o jogador personalizar seu personagem no detalhe (história, notas da trama, itens, classes, raças, etc), ou criar um personagem express no clique de um botão; que permita ter uma ficha de personagem dinamica e interativa onde ele pode evoluir aspectos no decorrer da trama, e gerenciar as estatisticas durante o combate; que permita ao DM criar fichas de monstros e NPCs relevantes com todas as possibilidades; ter visibilidade da ficha dos seus jogadores e status atual; que permita ao DM criar campanhas, gerenciar o ritmo, localização de plots e personagens importantes; quais memorias o npcs guardam dos pcs e porque, etc. 
    
    > Experiência técnica do utilizador médio: assume-se que conhece D&D 5e/ 5.5e (2024), ou precisas de tooltips explicando o que é "proficiência" ou "modificador de habilidade"?
        >> Consideramos que o jogador não sabe de absolutamente nada. Essa é a grande diferença dos builders comuns. A ideia é que a cada estapa, a cada escolha, seja indicado ao user o que ele está a escolher. Por exemplo, ao selecionar a classe sorecerer, que ela possa saber o que é um sorcerer, o como jogar com ele, quais os principais aspectos, como ele evolui, quais as subclasses, qual a lista de spells, melhores escolhas de background e species a depender do conceito do personagem; e assim para cada escolha (species, background, spells, itens, etc.). Se ele escolhe ser um Kalashtar, Sorcerer Aberrante de 6th nivel, ser destacado para ele melhores escolhas baseado no tipo de gameplay proposto para a classe, origem da specie, background, etc. Para DMs, idem. Embora, via de regra, quem quer ser DM de D&D sabe alguma coisa sobre o game, mas escolher cenário, encontros, dificuldades, e todos os aspectos de gerir uma campanha, pode ser desafidor. Assim, seria interessante a ferramenta ser inteligente o bastante para sugerir alguns caminhos usando as proprias indicações das fontes de dados (livros e suplementos).

B. Regras & Conteúdo

    Edição de D&D: SRD 5.1 (2014), D&D 2024 (nova edição), ou ambas? Isso muda drasticamente a estrutura de dados (ex: novo sistema de species/background).
        >> O escopo do projeto esta restrito ao D&D 2024 (nova edição) exclusivamente.

    Sources permitidas: apenas SRD (gratuito/open), ou pretendes incluir conteúdo de livros pagos (Xanathar, Tasha, etc.)? Nota: o SRD 5.1 é limitado (apenas 1 subclass por classe, etc.).
        >> ja temos isso na base de dados. Se não em SQL, em .md ou yaml

C. Funcionalidade Core (MVP)

    O que é absolutamente essencial na v1? Exemplo:
        >> Criador de personagem passo-a-passo (wizard)
        >> Visualizador de ficha web e pdf (character sheet)
        >> Lista de personagens salvos e a quais campanhas estão vinculandos.
        >> Ficha web dinâmica e interativa (para usar no tablet ou smartphone)
    O que é explicitamente fora do MVP? 
        >> campanhas, combate tracker, loja de conteúdo, compartilhamento entre jogadores, integração com VTT (Roll20, Foundry).

D. Modelo de Dados & Persistência

    - Contas de utilizador: registo obrigatório, ou modo guest (localStorage) com opção de registar depois? Isso impacta diretamente o modelo de auth e GDPR.
        >> para o MVP (localStorage). A ideia inicial é ser uma ferramenta containerizada de teste para uso no localhost com persistencia.

    - Partilha de fichas: um jogador pode mostrar a ficha ao mestre via link? Ou apenas o próprio dono vê?
        >> Ao vincular o ID de ficha ao ID de campanha o DM pode ver. Deve ser uma ação do user vincular o personagem a uma campanha.

    - Export/Import: queres exportar para PDF, JSON, ou importar do D&D Beyond?
        >> todas as opções. O personagem criado na plataforma estar disponivel na ficha dinamica web; poder ser exportada em pdf; e ser possível importar uma ficha do D&D Beyond (pdf) para usar ser salva no Onatar.

E. Stack Técnica & Infra

    - Backend: da tua experiência, prefere-se Go (como no projeto anterior)? Ou aberto a outra stack se houver justificação técnica? 
        >> em Go é mais facil para eu poder gerenciar o projeto. Mas estou aberto a outra stack que agregue de forma positiva ao projeto.

    - Frontend: React (como estava), ou aberto a alternativas (Vue, Svelte, HTMX+Go templates)? Qual é a tua preferência real para manutenção a longo prazo?
        >> não tenho expertise em frontend para essa escolha. React parece uma boa escolha para quando quiser tornar o projeto publico na web, mas se mostrou um desafio. Estou aberto a opções, desde que entregue os requisitos que espero.

    - Base de dados: PostgreSQL, SQLite (para self-hosted simples), ou outra?
        >> algo leve mas com boa performance. Maria DB talvez.

    - Deploy/Hosting: pretendes self-host (VPS próprio), ou usar serviço managed (Railway, Fly.io, Render, AWS)? Isso afeta decisões de infra (ex: SQLite vs PostgreSQL).
        >> para o MVP, VPS próprio. Vi ficar disponivel no gh para quem quiser montar sua propria mesa sem dificuldade.

F. DevSecOps & Qualidade

    - Autenticação: JWT simples, OAuth (Google/GitHub), ou ambos?
        >> Vamos ficar pelo gh mesmo.

    - Rate limiting / proteção de API: preocupa-te com abuse (spammers a criar 1000 personagens)?
        >> Não é o caso no MVP, pelas razões de E;

    - Testes: qual o nível de cobertura que consideras mínimo para dormir tranquilo? Unitários no backend, testes de integração na API, E2E no frontend?
        >> coerente com as escolhas anteriores.

    - CI/CD: usas GitHub Actions? Queres pipeline com tsc, go test, gosec, trivy scan de imagem Docker?
        >> Yep. Exceto pelo docker. Só usaremos docker no projeto quando ele estiver pelanemente funcional e pronto para o beta.

    - Observabilidade: logs estruturados, métricas (Prometheus), tracing? Ou apenas logs no stdout para começar?
        >> apenas logs no stdout para começar

G. Negócio & Legal

    - Modelo de negócio: 100% gratuito/open source, freemium (fichas ilimitadas vs limitadas), ou serviço pago?
        >> 100% gratuito/open source

    - Licenciamento: o projeto será open source (GitHub público)? Se sim, qual licença (MIT, AGPL, etc.)?
        >> AGPL-3

    - GDPR / privacidade: como é um serviço web com contas, precisamos de política de privacidade, consentimento, direito ao esquecimento. Queres lidar com isso desde o início ou adiar para pós-MVP?
        >> Não vamos lidar com dados de usuário no MVP