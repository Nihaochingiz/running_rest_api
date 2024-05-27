
1. docker volume create roach
2. docker network create -d bridge mynet
3. docker compose up --build
4. docker exec -it roach ./cockroach sql --insecure
5. CREATE DATABASE runningdb;
6. CREATE USER totoro;
7. GRANT ALL ON DATABASE runningdb TO totoro;
8. docker compose up --build

CREATE TABLE IF NOT EXISTS running_statistics (
			id SERIAL PRIMARY KEY,
			date VARCHAR(10),
			distance VARCHAR(10),
			time VARCHAR(10),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(date, distance, time)
		)


# Replace the placeholders with your actual data
curl -X POST -H "Content-Type: application/json" -d '{
  "date": "28.05.2024",
  "distance": "1 km",
  "time": "30 minutes"
}' http://localhost/running_statistics