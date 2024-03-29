match (alex:Driver {name: "Alex", surname: "Albon"})
match (carlos:Driver {name: "Carlos", surname: "Sainz"})
match (charles:Driver {name: "Charles", surname: "Leclerc"})
match (esteban:Driver {name: "Esteban", surname: "Ocon"})
match (fernando:Driver {name: "Fernando", surname: "Alonso"})
match (george:Driver {name: "George", surname: "Russell"})
match (kevin:Driver {name: "Kevin", surname: "Magnussen"})
match (lance:Driver {name: "Lance", surname: "Stroll"})
match (lando:Driver {name: "Lando", surname: "Norris"})
match (lewis:Driver {name: "Lewis", surname: "Hamilton"})
match (logan:Driver {name: "Logan", surname: "Sargeant"})
match (max:Driver {name: "Max", surname: "Verstappen"})
match (nico:Driver {name: "Nico", surname: "Hulkenberg"})
match (nyck:Driver {name: "Nyck", surname: "de Vries"})
match (oscar:Driver {name: "Oscar", surname: "Piastri"})
match (pierre:Driver {name: "Pierre", surname: "Gasly"})
match (sergio:Driver {name: "Sergio", surname: "Perez"})
match (valtteri:Driver {name: "Valtteri", surname: "Bottas"})
match (yuki:Driver {name: "Yuki", surname: "Tsunoda"})
match (zhou:Driver {name: "Zhou", surname: "Guanyu"})

merge (r:Race {season: 2023, race: 0})
merge (alex)-[:HAS_RACE {points: 0}]-(r)
merge (carlos)-[:HAS_RACE {points: 30}]-(r)
merge (charles)-[:HAS_RACE {points: 38}]-(r)
merge (esteban)-[:HAS_RACE {points: 11}]-(r)
merge (fernando)-[:HAS_RACE {points: 4}]-(r)
merge (george)-[:HAS_RACE {points: 34}]-(r)
merge (kevin)-[:HAS_RACE {points: 3}]-(r)
merge (lance)-[:HAS_RACE {points: 2}]-(r)
merge (lando)-[:HAS_RACE {points: 15}]-(r)
merge (lewis)-[:HAS_RACE {points: 30}]-(r)
merge (logan)-[:HAS_RACE {points: 0}]-(r)
merge (max)-[:HAS_RACE {points: 56}]-(r)
merge (nico)-[:HAS_RACE {points: 3}]-(r)
merge (nyck)-[:HAS_RACE {points: 2}]-(r)
merge (oscar)-[:HAS_RACE {points: 15}]-(r)
merge (pierre)-[:HAS_RACE {points: 11}]-(r)
merge (sergio)-[:HAS_RACE {points: 38}]-(r)
merge (valtteri)-[:HAS_RACE {points: 6}]-(r)
merge (yuki)-[:HAS_RACE {points: 1}]-(r)
merge (zhou)-[:HAS_RACE {points: 0}]-(r)

merge (r2:Race {season: 2023, race: 1})