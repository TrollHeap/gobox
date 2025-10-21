# TODO CPU infos
##  Étape 1 — Définir les structures

Dans `types.go`, Définir deux structures :

1. **`CPUCore`** : pour représenter un cœur logique.

   * Contiendra les infos brutes extraites par bloc (`id`, `core_id`, `mhz`, `flags`, etc.)
   * Pas besoin de tout remplir tout de suite, do it simple.

2. **`CPUInfo`** : synthèse globale de la machine.

   * C’est ce que `GetCPUInfo()` retournera.
   * Mettre dedans : `Vendor`, `ModelName`, `Cores`, `Threads`, `FrequencyMHz`, `CacheL3KB`.

> But : Pourvoir dire “équivalent structuré de `lscpu`”.

---

## Étape 2 — Lire `/proc/cpuinfo`

Dans `probe.go`, crée une fonction interne du type `parseProcCPUInfo()`.

Ta mission :

1. Lire le fichier avec `os.ReadFile`.
2. Découper le texte en **blocs** séparés par des lignes vides (`\n\n`).
3. Parcourir chaque bloc (chaque processeur logique).
4. Pour chaque ligne :

   * si elle commence par `"vendor_id"`, Noter e nom du constructeur.
   * si c’est `"model name"`, Noter le modèle.
   * si c’est `"cpu MHz"`, Additionner pour calculer une moyenne à la fin.
   * si c’est `"cpu cores"`, Enregistrer la valeur (souvent répétée).
   * si c’est `"physical id"`, garder dans un `map[string]bool` pour compter les sockets uniques.
   * si c’est `"cache size"`, noter la taille une seule fois (en KB).

> Objectif ici : remplir progressivement les variables pour former `CPUInfo`.

---

## Étape 3 — Calculs de synthèse

Une fois la boucle terminée :

* la moyenne de fréquence = somme des `cpu MHz` / nombre de processeurs logiques.
* le nombre de threads = nombre total de blocs `processor :`.
* le nombre de sockets = taille du `map`.
* le nombre de threads par cœur = `threads / cores`.

Vérifications ces relations :

```
threads >= cores
threadsPerCore ≥ 1
```

---

## Étape 4 — Retour du résultat

Dans `detect.go`, écrire `GetCPUInfo()` qui :

* appeller `parseProcCPUInfo()`
* renvoyer un `CPUInfo` complet
* gerer les erreurs s’il n’y a pas de fichier ou si `/proc/cpuinfo` est vide.

Fonction pour appeler depuis CLI / TUI.

---

## Étape 5 — Test console rapide

Faire un petit `main.go` temporaire ou add un test rendu

```go
info, _ := cpu.GetCPUInfo()
fmt.Printf("%+v\n", info)
```

Obtenir un affichage brut des valeurs.
Ensuite, formater

---

## Étape 6 — Vérifications

Teste sur plusieurs systèmes (si tu peux) :

* Laptop Intel
* Machine virtuelle
* Serveur avec CPU AMD

Compare la sortie avec `lscpu` pour vérifier la cohérence :

```
lscpu | grep -E 'Model name|Socket|Thread|Core|MHz|Cache'
```

Trouver les mêmes valeurs, ±1 % pour les fréquences (elles fluctuent).

---

## Étape 7 — Extensions possibles (plus tard)

* Lire `/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq` pour obtenir la fréquence max.
* Lire `/sys/class/thermal/thermal_zone*/temp` pour ajouter la température.
* Ajouter une section “Vulnérabilités connues” depuis `/sys/devices/system/cpu/vulnerabilities`.

