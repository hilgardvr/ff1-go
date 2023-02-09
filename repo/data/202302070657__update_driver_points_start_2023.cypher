//total points per 1st migration 2459/101(ppr)*3(races)
match (r:Race {season: 2023, race: 0})
match (d:Driver)-[hr:HAS_RACE]-(r)
set hr.points = hr.points/8