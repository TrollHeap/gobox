## 📚 Explication détaillée : Pourquoi `string` et pas `[]int`

### Analyse de tes fichiers système

Parfait, tu as des cas réels  :[1][2][3]

```
index3: "0-23"   ← L3 partagé entre 24 threads (12 cores × 2 SMT)
index2: "0-1"    ← L2 partagé entre 2 threads (1 core × 2 SMT)
index1: "0-1"    ← L1i partagé entre 2 threads
```

***

## 🔴 Le problème fondamental : Les slices ne sont pas comparables

### Règle Go[4][5][6][7]

**Pour être clé de map, un type DOIT être "comparable"**, c'est-à-dire supporter les opérateurs `==` et `!=`.[5][6][4]

**Types comparables**  :[8][6][4][5]
- ✅ Types primitifs : `int`, `float64`, `string`, `bool`
- ✅ Pointeurs : `*int`, `*CPUInfo`
- ✅ Structs dont **tous les champs** sont comparables
- ✅ Arrays : `[36]int` (taille fixe)

**Types NON comparables**  :[9][10][11][5][8]
- ❌ **Slices** : `[]int`, `[]string`
- ❌ Maps : `map[string]int`
- ❌ Functions : `func()`

***

## 🧪 Démonstration pratique

### Scénario A : Si tu utilisais `[]int` (ne compile PAS)

```go
type CacheID struct {
    Level      int
    SharedCPUs []int  // ❌ slice = non comparable
}

func main() {
    seen := make(map[CacheID]bool)  // ❌ ERREUR DE COMPILATION
    // invalid map key type CacheID: Level contains []int which is not comparable
}
```


**Raison technique**  :[12][13][5]

Un slice contient **3 champs internes**  :[9][8]
```go
type slice struct {
    array *[...]int  // pointeur vers backing array
    len   int
    cap   int
}
```

**Comparer deux slices avec `==`** nécessiterait de décider  :[5][12]
- Compare-t-on les **pointeurs** (identité) ?
- Compare-t-on les **éléments** (égalité profonde) ?
- Et si les slices ont des capacités différentes mais mêmes éléments ?

Go refuse de choisir → **slices interdits comme clés**.[8][12][5]

[13][12][5][9][8]

---

### Scénario B : Avec `string` (compile et fonctionne)

```go
type CacheID struct {
    Level      int     // ✅ comparable
    SharedCPUs string  // ✅ comparable
}

func main() {
    seen := make(map[CacheID]bool)  // ✅ OK
    
    id1 := CacheID{Level: 3, SharedCPUs: "0-23"}
    id2 := CacheID{Level: 3, SharedCPUs: "0-23"}
    
    fmt.Println(id1 == id2)  // ✅ true (comparaison de struct)
    
    seen[id1] = true
    _, exists := seen[id2]  // ✅ true (même clé)
}
```


**Pourquoi `string` fonctionne**  :[5][8]

Les strings Go sont comparables **par valeur** (compare octet par octet)  :[8][5]
```go
"0-23" == "0-23"  // true
"0-1"  == "0-23"  // false
```


***

## 📊 Tableau comparatif : `string` vs `[]int`

| Aspect | `SharedCPUs string` | `SharedCPUs []int` |
|--------|---------------------|-------------------|
| **Comparable** | ✅ Oui (`==` défini) | ❌ Non (slice) |
| **Utilisable comme clé map** | ✅ Oui | ❌ Non (erreur compilation) |
| **Simplicité** | ✅ Pas de parsing nécessaire | ❌ Nécessite parsing de "0-23" |
| **Performance** | ✅ Comparaison directe | ❌ Nécessiterait boucle manuelle |
| **Mémoire** | ~8-24 octets (header string) | ~24 octets (header slice) + array |

[14][9][6][5][8]

***

## 🎯 Réponse en une phrase

**"On utilise `string` pour `SharedCPUs` car les slices ne sont pas comparables donc interdites comme clés de map, alors que les strings supportent `==` et permettent une déduplication simple."**[12][6][14][5]

---

## 🔍 Bonus : Exemple concret de déduplication

```go
type CacheID struct {
    Level      int
    SharedCPUs string
}

seen := make(map[CacheID]bool)

// Premier passage : cpu0/cache/index3
id1 := CacheID{Level: 3, SharedCPUs: "0-23"}
if !seen[id1] {
    seen[id1] = true
    totalSize += parseSize("16M")  // ✅ Compte ce cache
}

// Deuxième passage : cpu1/cache/index3
id2 := CacheID{Level: 3, SharedCPUs: "0-23"}  // ← MÊME string
if !seen[id2] {  // ← false (déjà vu)
    // ✅ Skip ce cache (déjà compté)
}

// Troisième passage : cpu2/cache/index2
id3 := CacheID{Level: 2, SharedCPUs: "0-1"}  // ← Différent (niveau + CPUs)
if !seen[id3] {
    seen[id3] = true
    totalSize += parseSize("512K")  // ✅ Compte ce cache
}
```


***

## ✅ Validation conceptuelle

**Tu as maintenant compris pourquoi `string` est nécessaire ?**

**Confirme en une phrase** : Pourquoi `[]int` ne peut pas être clé de map en Go ?[12][5]

🔒 **Réponse attendue → passage à l'Étape Bonus 2 (lecture de `shared_cpu_list` et implémentation map).**
