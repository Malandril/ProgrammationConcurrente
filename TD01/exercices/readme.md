# Programmation Concurrente

## Canava Thomas

## Exercice 9 Création d'un serveur de chat simple

### Architecture

Pour le serveur de chat j'ai decidé de créér trois canaux de communication principaux.
* un canal sur lequel sont envoyés les messages: `messageChan`
* un canal sur lequel sont envoyés les connexions: `connexion`
* un canal sur lequel sont envoyés les déconnexions: `deconnexion`|

Ainsi qu'une map qui permet de lier un nom d'utilisateur a une connexion et de recomuniquer avec les utilisateurs.

Ces trois cannaux sont "lus" dans une go routine `messageSender`. Cette go routine s'occupe de propager l'information à tous les clients grâce à la map de clients.

Il y a un *select* sur ces trois canaux et qui permet d'éviter d'éxecuter ces trois évenements differents de manière concurrente. Cela évite à la map qui contient les utilisateurs d'être modifiée ou lu en même temps.

Une go routine `handleConnection` est créé par utilisateur qui se connecte.

Quand l'utilisateur a renseigné son nom un message est envoyé sur le cannal de connexion pour notifier de l'arrivé de l'utilisateur.
Ensuite un canal interne à la go routine est créé `receivedMessage`.

Ensuite dans la go routine `handleConnection` une go routine va être créé qui va lire le réseau et envoyer sur le canal interne `receivedMessage`quand un message arrive, et sur le cannal `messageChan` le message reçu. 
Cette go routine interne et le canal `receivedMessage` permettent grâce à un *select* de savoir si un message n'a pas été envoyé depuis longtemps par le client et de le deconnecter si c'est le cas, ou si le client choisis de terminer la connection.

Il y a donc en tout 2 goroutines par client et une goroutine sur le serveur.

Les goroutines des clients communiquent avec la goroutine du serveur grâce à 3 canaux.
les deux goroutines clients communiquent grâce à un canal.