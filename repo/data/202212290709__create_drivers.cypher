merge (c:Constructor {name: "Red Bull")})
merge (:Driver {name: "Max", surname: "Verstappen"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Sergio", surname: "Perez"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Ferrari"})
merge (:Driver {name: "Charles", surname: "Leclerc"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Carlos", surname: "Sainz"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Mercedes"})
merge (:Driver {name: "Lewis", surname: "Hamilton"})-[:RACES_FOR]->(c)
merge (:Driver {name: "George", surname: "Russell"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Alpine"})
merge (:Driver {name: "Pierre", surname: "Gasly"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Esteban", surname: "Ocon"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "McLaren"})
merge (:Driver {name: "Lando", surname: "Norris"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Oscar", surname: "Piastri"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Alfa Romeo"})
merge (:Driver {name: "Zhou", surname: "Guanyu"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Valtteri", surname: "Bottas"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Aston Martin"})
merge (:Driver {name: "Fernando", surname: "Alonso"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Lance", surname: "Stroll"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Haas"})
merge (:Driver {name: "Kevin", surname: "Magnussen"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Nico", surname: "Hulkenberg"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "AlphaTauri"})
merge (:Driver {name: "Yuki", surname: "Tsunoda"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Nyck", surname: "de Vries"})-[:RACES_FOR]->(c);

merge (c:Constructor {name: "Williams"})
merge (:Driver {name: "Logan", surname: "Sargeant"})-[:RACES_FOR]->(c)
merge (:Driver {name: "Alex", surname: "Albon"})-[:RACES_FOR]->(c);