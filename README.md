# pathfinding-demo
Library to integrate all components of the server backend rendering for pathfinding.

To generate configs, run

```bash
go run simulation/generator/main.go \
  --output=demo/configs/ && go run demo/main.go \
  --configs=demo/configs/*json \
  --output=demo/output \
  --log=demo/output/logs
```
