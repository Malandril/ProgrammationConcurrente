# Compte rendu tp chat elixir
### *Canava Thomas*

## Architecture

Pour ce tp j'ai choisis de refaire comme pour le tp de go, une fonction sur un thread elixir à part qui va s'occuper d'envoyer les messages à tous les utilisateurs.
Cette fonction garde grâce a des appels récursifs à chaque modification une liste d'utilisateurs.

Cette fonction `message_sender` grâce à receive reçoit des messages d'autres thread elixir et n'en traite qu'un à la fois dans le receive, ce qui permet d'éviter toute modification concurrente.

`message_sender` peut recevoir 4 types de messages:
- `:add_client` qui permet de rajouter un client à la liste de clients
- `:send_message` qui permet d'envoyer à tous les clients un message grâce à la liste de clients
- `:disconnect` qui permet de supprimer un client de la liste de clients

A chaque fois qu'un client se connecte un nouveau thread elixir avec la fonction `serve` est créé, et `message_sender` est notifié grace à la fonction `send`,
`serve` envoie d'abord `:add_client` puis `:send_message` avec un message qui indique la connexion de l'utilisateur.

Pour communiquer avec la fonction `message_sender` son *pid* est passé en paramètre de `serve`.

A la fin de `serve` c'est la méthode `listen_client` qui est appelée.

`listen_client` attend un message avec `:gen_tcp.recv` et grâce à un `case` on fait du pattern matching :
- si un message est bien reçu avec `:ok` on envoie `:send_message` à `message_sender`
- si il y a une erreur `:error` et un `:timeout` on envoie un message de timeout avec `:send_message`et on envoie `:disconnect` pour le supprimer de la liste
- si il y a une erreur `:error` et quelque chose autre que `:timeout` on envoie un message de déconnexion avec `:send_message`et on envoie `:disconnect` pour le supprimer de la liste