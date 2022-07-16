start nats-server.exe  -m 8222 -p 4222 -cluster nats://localhost:4232 -routes nats://localhost:4233,nats://localhost:4234
start nats-server.exe  -m 8223 -p 4223 -cluster nats://localhost:4233 -routes nats://localhost:4232,nats://localhost:4234
start nats-server.exe  -m 8224 -p 4224 -cluster nats://localhost:4234 -routes nats://localhost:4232,nats://localhost:4233