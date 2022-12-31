merge (red:Constructor {name: "Red Bull"})
merge (:Driver {name: "Max", surname: "Verstappen"})-[:RACES_FOR]->(red)
merge (:Driver {name: "Sergio", surname: "Perez"})-[:RACES_FOR]->(red)

merge (ferrari:Constructor {name: "Ferrari"})
merge (:Driver {name: "Charles", surname: "Leclerc"})-[:RACES_FOR]->(ferrari)
merge (:Driver {name: "Carlos", surname: "Sainz"})-[:RACES_FOR]->(ferrari)

merge (mercedes:Constructor {name: "Mercedes"})
merge (:Driver {name: "Lewis", surname: "Hamilton"})-[:RACES_FOR]->(mercedes)
merge (:Driver {name: "George", surname: "Russell"})-[:RACES_FOR]->(mercedes)

merge (alpine:Constructor {name: "Alpine"})
merge (:Driver {name: "Pierre", surname: "Gasly"})-[:RACES_FOR]->(alpine)
merge (:Driver {name: "Esteban", surname: "Ocon"})-[:RACES_FOR]->(alpine)

merge (mclaren:Constructor {name: "McLaren"})
merge (:Driver {name: "Lando", surname: "Norris"})-[:RACES_FOR]->(mclaren)
merge (:Driver {name: "Oscar", surname: "Piastri"})-[:RACES_FOR]->(mclaren)

merge (alfa:Constructor {name: "Alfa Romeo"})
merge (:Driver {name: "Zhou", surname: "Guanyu"})-[:RACES_FOR]->(alfa)
merge (:Driver {name: "Valtteri", surname: "Bottas"})-[:RACES_FOR]->(alfa)

merge (aston:Constructor {name: "Aston Martin"})
merge (:Driver {name: "Fernando", surname: "Alonso"})-[:RACES_FOR]->(aston)
merge (:Driver {name: "Lance", surname: "Stroll"})-[:RACES_FOR]->(aston)

merge (haas:Constructor {name: "Haas"})
merge (:Driver {name: "Kevin", surname: "Magnussen"})-[:RACES_FOR]->(haas)
merge (:Driver {name: "Nico", surname: "Hulkenberg"})-[:RACES_FOR]->(haas)

merge (alphatauri:Constructor {name: "AlphaTauri"})
merge (:Driver {name: "Yuki", surname: "Tsunoda"})-[:RACES_FOR]->(alphatauri)
merge (:Driver {name: "Nyck", surname: "de Vries"})-[:RACES_FOR]->(alphatauri)

merge (williams:Constructor {name: "Williams"})
merge (:Driver {name: "Logan", surname: "Sargeant"})-[:RACES_FOR]->(williams)
merge (:Driver {name: "Alex", surname: "Albon"})-[:RACES_FOR]->(williams)