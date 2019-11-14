# Description
L'API doit relier chacune des impressions et clics présents dans le fichier events.csv au point
d'intérêt le plus proche.

# Objectifs de ce projet :
Le but est de construire une API REST pour relier des impressions publicitaires et clics à une liste de
points d'intérêts.

# Fonctionnalités
L'API doit exposer une route sur laquelle nous enverrons un objet JSON décrivant les points
d'intérêts.

Pour tester :
1. Effectuer une POST request à l'adresse : http://localhost:5000/impressionsAndClicks/
> Exemple de body pour la POST Request :
```json
[
	{
		"lat": 48.86,
		"lon": 2.35,
		"name": "Chatelet"
	},
	{
		"lat": 48.8759992,
		"lon": 2.3481253,
		"name": "Arc de triomphe"
	}
]
```
# Déploiement
## Local (nécessite go)
1. Build le projet et lancer l'application : `go build && ./audion`
## Container (nécessite docker)
1. Construire le container : `docker build -t audion .`
1. Lancer le container : `docker run -it --rm --name test -p 5000:5000 audion`